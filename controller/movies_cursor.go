package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/Dungsenpai-ux/Practice_Go/service"
)

// ListMoviesWithPaging handles GET /movies/list with offset or cursor paging
func ListMoviesWithPaging(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")
	cursor := r.URL.Query().Get("cursor")

	// Default size is 10 if not provided
	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 10
	}

	var movies []model.Movie
	var nextCursor string

	if cursor != "" {
		// Cursor-based paging
		movies, nextCursor, err = service.GetMoviesByCursor(r.Context(), cursor, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Measure performance for cursor-based paging
		// Assume cursor is in the format "id,year"
		var cursorID int
		var cursorYear int
		_, err := fmt.Sscanf(cursor, "%d,%d", &cursorID, &cursorYear)
		if err != nil {
			log.Printf("Invalid cursor format: %v", err)
			cursorID, cursorYear = 0, 0
		}
		explain, err := service.MeasureQuery(r.Context(), "SELECT id, title, year, genres FROM movies WHERE (id, year) > ($1, $2) ORDER BY id, year LIMIT $3", cursorID, cursorYear, size)
		if err != nil {
			log.Printf("Error measuring cursor query: %v", err)
		} else {
			log.Printf("Cursor-based query (cursor=%s, size=%d): %s", cursor, size, explain)
		}
	} else {
		// Offset-based paging
		page, err := strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = 1
		}
		offset := (page - 1) * size
		movies, err = service.GetMoviesByOffset(r.Context(), offset, size)
		if err != nil {
			http.Error(w, "Internal Server Error: Unable to fetch movies", http.StatusInternalServerError)
			log.Printf("Error fetching movies: %v", err)
			return
		}
		// Measure performance for offset-based paging
		explain, err := service.MeasureQuery(r.Context(), "SELECT id, title, year, genres FROM movies ORDER BY id LIMIT $1 OFFSET $2", size, offset)
		if err != nil {
			log.Printf("Error measuring offset query: %v", err)
		} else {
			log.Printf("Offset-based query (page=%d, size=%d): %s", page, size, explain)
		}
		// Generate next cursor for offset-based paging
		if len(movies) > 0 {
			lastMovie := movies[len(movies)-1]
			nextCursor = fmt.Sprintf("%d,%d", lastMovie.ID, lastMovie.Year)
		}
	}

	// Return response with movies and next cursor
	response := map[string]interface{}{
		"movies":      movies,
		"next_cursor": nextCursor,
	}
	json.NewEncoder(w).Encode(response)
}