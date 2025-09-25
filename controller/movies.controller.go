package controller

import (
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/bradfitz/gomemcache/memcache"
)

var (
	cacheHits        = expvar.NewInt("cache_hits")
	cacheMisses      = expvar.NewInt("cache_misses")
	cacheErrors      = expvar.NewInt("cache_errors")
	cacheNegatives   = expvar.NewInt("cache_negative_writes")
	cacheWriteErrors = expvar.NewInt("cache_write_errors")
)

// MovieHandler aggregates dependencies for movie endpoints
type MovieHandler struct {
	Service interface {
		InsertMovie(ctx context.Context, movie model.Movie) (int, error)
		GetMovieByID(ctx context.Context, id int) (model.Movie, error)
		SearchMovies(ctx context.Context, q string, year int) ([]model.Movie, error)
	}
	Cache *memcache.Client
}

// NewMovieHandler builds a handler instance
func NewMovieHandler(s interface {
	InsertMovie(ctx context.Context, movie model.Movie) (int, error)
	GetMovieByID(ctx context.Context, id int) (model.Movie, error)
	SearchMovies(ctx context.Context, q string, year int) ([]model.Movie, error)
}, cache *memcache.Client) *MovieHandler {
	return &MovieHandler{Service: s, Cache: cache}
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var req MovieCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	movieModel := req.ToModel()
	id, err := h.Service.InsertMovie(r.Context(), movieModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if h.Cache != nil {
		cacheKey := fmt.Sprintf("movie:%d", id)
		if err := h.Cache.Delete(cacheKey); err != nil && err != memcache.ErrCacheMiss {
			log.Printf("Error invalidating cache for %s: %v", cacheKey, err)
		}
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/movies/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID không hợp lệ", http.StatusBadRequest)
		return
	}
	if h.Cache != nil {
		cacheKey := fmt.Sprintf("movie:%d", id)
		item, err := h.Cache.Get(cacheKey)
		switch err {
		case nil:
			cacheHits.Add(1)
			log.Printf("Cache hit for %s", cacheKey)
			var movie model.Movie
			if err := json.Unmarshal(item.Value, &movie); err == nil {
				json.NewEncoder(w).Encode(FromModel(movie))
				return
			}
			if string(item.Value) == `{"error":"không tìm thấy phim"}` {
				http.Error(w, "không tìm thấy phim", http.StatusNotFound)
				return
			}
			log.Printf("Error unmarshaling cached movie: %v", err)
			cacheErrors.Add(1)
		case memcache.ErrCacheMiss:
			cacheMisses.Add(1)
			log.Printf("Cache miss for %s", cacheKey)
		default:
			cacheErrors.Add(1)
			log.Printf("Error checking cache for %s: %v", cacheKey, err)
		}
	}
	movie, err := h.Service.GetMovieByID(r.Context(), id)
	if err != nil {
		cacheKey := fmt.Sprintf("movie:%d", id)
		item := &memcache.Item{Key: cacheKey, Value: []byte(`{"error":"không tìm thấy phim"}`), Expiration: 30}
		if h.Cache != nil {
			if err := h.Cache.Set(item); err != nil {
				cacheWriteErrors.Add(1)
				log.Printf("Error setting negative cache for %s: %v", cacheKey, err)
			} else {
				cacheNegatives.Add(1)
			}
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if h.Cache != nil {
		data, err := json.Marshal(movie)
		if err != nil {
			cacheErrors.Add(1)
			log.Printf("Error marshaling movie for cache: %v", err)
		} else {
			cacheKey := fmt.Sprintf("movie:%d", id)
			item := &memcache.Item{Key: cacheKey, Value: data, Expiration: int32(5 * 60)}
			if err := h.Cache.Set(item); err != nil {
				cacheWriteErrors.Add(1)
				log.Printf("Error setting cache for %s: %v", cacheKey, err)
			}
		}
	}
	json.NewEncoder(w).Encode(FromModel(movie))
}

func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	yearStr := r.URL.Query().Get("year")
	year := 0
	if yearStr != "" {
		year, _ = strconv.Atoi(yearStr)
	}
	movies, err := h.Service.SearchMovies(r.Context(), q, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(FromModelSlice(movies))
}
