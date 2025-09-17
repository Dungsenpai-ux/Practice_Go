package controller

import (
	"net/http"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// Endpoint Health
	healthController := NewHealthController(cfg)
	r.GET("/healthz", healthController.HealthCheck)
}

func RegisterHTTPRoutes() {
	// Endpoint Movies
	http.HandleFunc("POST /movies", CreateMovie)
	http.HandleFunc("GET /movies/", GetMovie)
	http.HandleFunc("GET /movies/search", SearchMovies)
}
