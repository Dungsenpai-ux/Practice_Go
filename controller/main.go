package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, health *HealthController, movie *MovieHandler) {
	r.GET("/healthz", health.HealthCheck)
	// If you later move movie endpoints to gin instead of net/http:
	// r.POST("/movies", movie.CreateMovie)
	// r.GET("/movies/:id", wrapper...) etc.
}

func RegisterHTTPRoutes(movie *MovieHandler) {
	http.HandleFunc("POST /movies", movie.CreateMovie)
	http.HandleFunc("GET /movies/", movie.GetMovie)
	http.HandleFunc("GET /movies/search", movie.SearchMovies)
}
