package controller

import (
	"expvar"
	"net/http"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/service/middleware"
	"github.com/gin-gonic/gin"
)

// (Legacy) SetupRoutes – giữ lại nếu cần tương thích cũ
func SetupRoutes(r *gin.Engine, health *HealthController) {
	r.GET("/healthz", health.HealthCheck)
}

// RegisterHTTPRoutes registers movie endpoints onto provided mux, applying logging/metrics middleware once.
func RegisterHTTPRoutes(mux *http.ServeMux, movie *MovieHandler) {
	mux.Handle("POST /movies", middleware.LoggingAndMetrics(http.HandlerFunc(movie.CreateMovie)))
	mux.Handle("GET /movies/", middleware.LoggingAndMetrics(http.HandlerFunc(movie.GetMovie)))
	mux.Handle("GET /movies/search", middleware.LoggingAndMetrics(http.HandlerFunc(movie.SearchMovies)))
}

// BuildRouter: tạo *gin.Engine và đăng ký tất cả endpoint (Gin mode Option B)
func BuildRouter(cfg *config.Config, health *HealthController, movie *MovieHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/healthz", health.HealthCheck)
	r.POST("/movies", gin.WrapF(movie.CreateMovie))
	r.GET("/movies/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Request.URL.Path = "/movies/" + id
		gin.WrapF(movie.GetMovie)(c)
	})
	r.GET("/movies/search", gin.WrapF(movie.SearchMovies))
	// Serve expvar metrics directly (no redirect) so /metrics returns JSON
	handler := expvar.Handler()
	r.GET("/metrics", gin.WrapH(handler))
	r.GET("/debug/vars", gin.WrapH(handler))
	return r
}
