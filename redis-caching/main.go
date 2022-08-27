package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"redis-caching/middleware"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type API struct {
	redisClient *redis.Client
}

type APIResponse struct {
	Cache bool                `json:"cache"`
	Data  []NominatinLocation `json:"data"`
}

type NominatinLocation struct {
	PlaceID     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmID       int      `json:"osm_id"`
	Boundingbox []string `json:"boundingbox"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
	Icon        string   `json:"icon"`
}

var (
	container    = os.Getenv("CONTAINER") == "true"
	nominatimUrl = "https://nominatim.openstreetmap.org/search?q=%s&format=json"
	defaultTTL   = time.Second * 300
)

func NewAPI() *API {
	var opts *redis.Options

	if container {
		fmt.Printf("Using local redis service\n")
		redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_SERVICE"))
		opts = (&redis.Options{
			Addr:     redisAddress,
			Password: "",
			DB:       0,
		})
	} else {
		fmt.Printf("Using remote redis service\n")
		opts = (&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PW"),
			DB:       0,
		})
	}
	rdb := redis.NewClient(opts)

	return &API{redisClient: rdb}
}

func main() {
	if !container {
		envErr := godotenv.Load("local.env")
		if envErr != nil {
			panic(envErr)
		}
	}
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Printf("Starting server on port %s\n", port)

	api := NewAPI()
	mux := http.NewServeMux()

	handler := http.HandlerFunc(api.Handler)
	mux.Handle("/api", middleware.UseJson(handler))

	http.ListenAndServe(port, mux)
}

func (api *API) Handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	data, cacheHit, err1 := api.GetData(r.Context(), q)

	if err1 != nil {
		w.Write([]byte("Error calling Nominatim API"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := APIResponse{
		Cache: cacheHit,
		Data:  data,
	}

	if err2 := json.NewEncoder(w).Encode(resp); err2 != nil {
		fmt.Printf("Error encoding response: %v\n", err2)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) GetData(ctx context.Context, q string) (
	[]NominatinLocation, bool, error) {
	escapedQ := url.PathEscape(strings.ToLower(q))
	cacheValue, redisErr := api.redisClient.Get(ctx, escapedQ).Result()
	if redisErr == redis.Nil {
		// Cached value does not exist. Call API -> get value -> cache
		url := fmt.Sprintf(nominatimUrl, escapedQ)
		resp, getErr := http.Get(url)
		if getErr != nil {
			return nil, false, getErr
		}

		data := make([]NominatinLocation, 0)
		decodeErr := json.NewDecoder(resp.Body).Decode(&data)
		if decodeErr != nil {
			return nil, false, decodeErr
		}

		b, marshallErr := json.Marshal(data)
		if marshallErr != nil {
			return nil, false, marshallErr
		}

		redisSetErr := api.redisClient.Set(
			ctx,
			escapedQ,
			bytes.NewBuffer(b).Bytes(),
			defaultTTL,
		).Err()
		if redisSetErr != nil {
			fmt.Printf("Failed to cache value for key %s\n", escapedQ)
		}

		return data, false, nil
	} else if redisErr != nil {
		return nil, false, redisErr
	} else {
		// Cached value exists. Extend TTL & build response
		api.redisClient.Expire(ctx, escapedQ, defaultTTL).Result()
		data := make([]NominatinLocation, 0)
		unmarshalErr := json.Unmarshal(
			bytes.NewBufferString(cacheValue).Bytes(), &data)
		if unmarshalErr != nil {
			return nil, false, unmarshalErr
		}
		return data, true, nil
	}
}
