package controller

import (
	"net/http"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// Endpoint Health
	healthController := NewHealthController(cfg)
	r.GET("/healthz", healthController.HealthCheck)
}

func RegisterHTTPRoutes() {
	// Endpoint Movies
	config.Init()
	http.Handle("POST /movies", config.InstrumentHandlerWithPath("/movies", http.HandlerFunc(CreateMovie)))
	http.Handle("GET /movies/", config.InstrumentHandlerWithPath("/movies/", http.HandlerFunc(GetMovie)))
	http.Handle("GET /movies/search", config.InstrumentHandlerWithPath("/movies/search", http.HandlerFunc(SearchMovies)))

	// Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
}
