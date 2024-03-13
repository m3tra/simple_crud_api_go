package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       int       `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

var currMaxID int = 0
var movies map[int]Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("getMovies")

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Println(err)
	}
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	log.Println("getMovie")
	fmt.Printf("Params:\n%v\n", params)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println(err)
		return
	}

	movie, exists := movies[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println(movie)

	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Println(err)
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("createMovie")

	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(movie)

	if m, exists := movies[movie.ID]; exists {
		log.Printf("Already exists:\n%v\n", m)
		w.WriteHeader(http.StatusConflict)
		return
	}

	currMaxID++
	movie.ID = currMaxID
	movies[movie.ID] = movie
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(movie); err != nil {
		log.Println(err)
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	log.Println("updateMovie")

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad id"))
		return
	}

	var changes Movie
	if err := json.NewDecoder(r.Body).Decode(&changes); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Changes:\n%v\n", changes)

	if changes == *new(Movie) {
		log.Println("Error: Empty object received")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movie, exists := movies[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("Movie to update:\n%v\n", movie)

	changed := false
	if len(changes.Isbn) > 0 && movie.Isbn != changes.Isbn {
		movie.Isbn = changes.Isbn
		changed = true
	}
	if len(changes.Title) > 0 && movie.Title != changes.Title {
		movie.Title = changes.Title
		changed = true
	}
	if len(changes.Director.First_name) > 0 && movie.Director.First_name != changes.Director.First_name {
		movie.Director.First_name = changes.Director.First_name
		changed = true
	}
	if len(changes.Director.Last_name) > 0 && movie.Director.Last_name != changes.Director.Last_name {
		movie.Director.Last_name = changes.Director.Last_name
		changed = true
	}

	if changed {
		if err := json.NewEncoder(w).Encode(movie); err != nil {
			log.Println(err)
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	log.Println("deleteMovie")

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad id"))
		return
	}

	movie, exists := movies[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Println(movie)
	delete(movies, id)
	log.Println("Deleted")

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.Println(err)
	}
}

func main() {
	const port int16 = 8080

	r := mux.NewRouter()

	movies = make(map[int]Movie)

	movies[1] = Movie{
		ID:    1,
		Isbn:  "12345",
		Title: "Movie One",
		Director: &Director{
			First_name: "John",
			Last_name:  "Doe",
		},
	}
	movies[2] = Movie{
		ID:    2,
		Isbn:  "12346",
		Title: "Django",
		Director: &Director{
			First_name: "Quentin",
			Last_name:  "Tarantino",
		},
	}
	currMaxID = 2

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PATCH")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port %d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), r))
}
