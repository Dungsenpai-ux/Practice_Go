package controller

import (
	"expvar"
	"net/http"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/service/middleware"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware(cfg.OtelService))
	r.Use(func(c *gin.Context) { c.Next() })
	r.GET("/healthz", health.HealthCheck)
	r.POST("/movies", gin.WrapF(movie.CreateMovie))
	r.GET("/movies/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.Request.URL.Path = "/movies/" + id
		gin.WrapF(movie.GetMovie)(c)
	})
	r.GET("/movies/search", gin.WrapF(movie.SearchMovies))
	handler := expvar.Handler()
	r.GET("/metrics", gin.WrapH(handler))
	r.GET("/debug/vars", gin.WrapH(handler))
	return r
}
