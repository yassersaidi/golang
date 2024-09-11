package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type Movie struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Startin server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))

}

func getMovies(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	moviesData := readData("./movies.json", &movies)

	json.NewEncoder(response).Encode(moviesData)
}

func deleteMovie(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	moviesData := readData("./movies.json", &movies)

	for index, item := range moviesData {
		if item.Id == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			writeData("./movies.json", movies)
			break
		}
	}

	json.NewEncoder(response).Encode(movies)
}

func getMovie(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	moviesData := readData("./movies.json", &movies)

	for _, item := range moviesData {
		if item.Id == params["id"] {
			json.NewEncoder(response).Encode(item)
			return
		}
	}
}

func createMovie(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(request.Body).Decode(&movie)

	movie.Id = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	writeData("./movies.json", movies)
	json.NewEncoder(response).Encode(movie)
}

func updateMovie(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)
	moviesData := readData("./movies.json", &movies)

	for index, item := range moviesData {
		if item.Id == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(request.Body).Decode(&movie)
			movie.Id = params["id"]
			movies = append(movies, movie)
			writeData("./movies.json", movies)
			json.NewEncoder(response).Encode(movie)
			return
		}
	}

}

func readData[T any](path string, parsedTo *[]T) []T {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, parsedTo)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return *parsedTo
}

func writeData[T any](path string, data []T) {
	dataToWrite, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	writeErr := os.WriteFile(path, dataToWrite, os.ModeAppend)
	if writeErr != nil {
		log.Fatal(writeErr)
	}
}
