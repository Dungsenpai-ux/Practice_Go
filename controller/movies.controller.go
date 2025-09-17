package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/Dungsenpai-ux/Practice_Go/service"
)

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := service.InsertMovie(r.Context(), movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/movies/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}
	movie, err := service.GetMovieByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(movie)
}

func SearchMovies(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	yearStr := r.URL.Query().Get("year")
	year := 0
	if yearStr != "" {
		year, _ = strconv.Atoi(yearStr)
	}
	movies, err := service.SearchMovies(r.Context(), q, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(movies)
}
