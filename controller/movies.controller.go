package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/model"
	"github.com/Dungsenpai-ux/Practice_Go/service"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/rs/zerolog/log"
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
	// Invalidate cache for the new movie
	if config.Memcached != nil {
		cacheKey := fmt.Sprintf("movie:%d", id)
		if err := config.Memcached.Delete(cacheKey); err != nil && err != memcache.ErrCacheMiss {
			log.Error().Err(err).Str("cache_key", cacheKey).Msg("invalidate cache failed")
		}
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

	// Check cache first
	if config.Memcached != nil {
		cacheKey := fmt.Sprintf("movie:%d", id)
		item, err := config.Memcached.Get(cacheKey)
		switch err {
		case nil:
			// Cache hit
			config.CacheHit.Inc()
			log.Info().Str("cache_key", cacheKey).Msg("cache hit")
			var movie model.Movie
			if err := json.Unmarshal(item.Value, &movie); err == nil {
				json.NewEncoder(w).Encode(movie)
				return
			}
			// Negative cache hit
			if string(item.Value) == `{"error":"không tìm thấy phim"}` {
				http.Error(w, "không tìm thấy phim", http.StatusNotFound)
				return
			}
			log.Error().Err(err).Msg("unmarshal cached movie failed")
		case memcache.ErrCacheMiss:
			config.CacheMiss.Inc()
			log.Info().Str("cache_key", cacheKey).Msg("cache miss")
		default:
			log.Error().Err(err).Str("cache_key", cacheKey).Msg("check cache failed")
		}
	}

	// Cache miss, query database
	movie, err := service.GetMovieByID(r.Context(), id)
	if err != nil {
		// Negative cache for not found
		cacheKey := fmt.Sprintf("movie:%d", id)
		item := &memcache.Item{
			Key:        cacheKey,
			Value:      []byte(`{"error":"không tìm thấy phim"}`),
			Expiration: 30, // TTL 30 seconds
		}
		if err := config.Memcached.Set(item); err != nil {
			log.Printf("Error setting negative cache for %s: %v", cacheKey, err)
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Cache the movie
	if config.Memcached != nil {
		data, err := json.Marshal(movie)
		if err != nil {
			log.Error().Err(err).Msg("marshal movie for cache failed")
		} else {
			cacheKey := fmt.Sprintf("movie:%d", id)
			item := &memcache.Item{
				Key:        cacheKey,
				Value:      data,
				Expiration: int32(5 * 60), // TTL 5 minutes
			}
			if err := config.Memcached.Set(item); err != nil {
				log.Error().Err(err).Str("cache_key", cacheKey).Msg("set cache failed")
			}
		}
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
