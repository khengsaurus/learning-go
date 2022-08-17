package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var Director1 = Director{Firstname: "John", Lastname: "Doe"}
var Director2 = Director{Firstname: "Jane", Lastname: "Smith"}
var Movie1 = Movie{ID: "1", Isbn: "1", Title: "One", Director: &Director1}
var Movie2 = Movie{ID: "2", Isbn: "2", Title: "One", Director: &Director1}
var Movie3 = Movie{ID: "3", Isbn: "3", Title: "One", Director: &Director2}

var movies = []Movie{}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode((movies))
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, "400 missing parameter", http.StatusBadRequest)
		return
	}
	for _, item := range movies {
		if item.ID == id {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Movie not found"))
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie // HELP: pass it to the memory address ?
	json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000))
	movies = append(movies, movie)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(movie)
}

// HELP: how to update in place?
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	for index, item := range movies {
		if item.ID == id {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = id
			movies = append(movies, movie)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Movie not found"))
}

func deleteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, "400 missing parameter", http.StatusBadRequest)
		return
	}
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Movie deleted"))
}

func main() {
	r := mux.NewRouter()
	movies = []Movie{Movie1, Movie2, Movie3}

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovies).Methods("DELETE")

	fmt.Printf("Server listening at port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
