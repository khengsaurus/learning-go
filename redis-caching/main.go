package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"redis-caching/middleware"
)

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
	port         = ":8080"
	nominatimUrl = "https://nominatim.openstreetmap.org/search?q=%s&format=json"
)

func main() {
	fmt.Printf("Starting server on port %s", port)
	mux := http.NewServeMux()

	handler := http.HandlerFunc(Handler)
	mux.Handle("/api", middleware.UseJson(handler))

	http.ListenAndServe(port, mux)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	data, err1 := GetData(q)

	if err1 != nil {
		w.Write([]byte("Error calling Nominatim API"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := APIResponse{
		Cache: false,
		Data:  data,
	}

	if err2 := json.NewEncoder(w).Encode(resp); err2 != nil {
		fmt.Printf("Error encoding response: %v", err2)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetData(q string) ([]NominatinLocation, error) {
	escapedQ := url.PathEscape(q)
	url := fmt.Sprintf(nominatimUrl, escapedQ)
	resp, err1 := http.Get(url)
	if err1 != nil {
		return nil, err1
	}

	data := make([]NominatinLocation, 0)
	err2 := json.NewDecoder(resp.Body).Decode(&data)
	if err2 != nil {
		return nil, err2
	}

	return data, nil
}
