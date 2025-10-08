package controller

import (
    "expvar"
    "net/http"

    "github.com/Dungsenpai-ux/Practice_Go/config"
    "github.com/Dungsenpai-ux/Practice_Go/service/middleware"
    "github.com/gin-gonic/gin"
    "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)


func SetupRoutes(r *gin.Engine, health *HealthController) {
    r.GET("/healthz", health.HealthCheck)
}

func RegisterHTTPRoutes(mux *http.ServeMux, movie *MovieHandler) {
    mux.Handle("POST /movies", middleware.LoggingAndMetrics(http.HandlerFunc(movie.CreateMovie)))
    mux.Handle("GET /movies/", middleware.LoggingAndMetrics(http.HandlerFunc(movie.GetMovie)))
    mux.Handle("GET /movies/search", middleware.LoggingAndMetrics(http.HandlerFunc(movie.SearchMovies)))
}

func BuildRouter(cfg *config.Config, health *HealthController, movie *MovieHandler) *gin.Engine {
    r := gin.New()
    
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    
    // Tracing middleware
    r.Use(otelgin.Middleware(cfg.OtelService))
    
    // Routes
    r.GET("/healthz", health.HealthCheck)
    r.POST("/movies", gin.WrapF(movie.CreateMovie))
    r.GET("/movies/:id", func(c *gin.Context) {
        id := c.Param("id")
        c.Request.URL.Path = "/movies/" + id
        movie.GetMovie(c.Writer, c.Request)
    })
    r.GET("/movies/search", gin.WrapF(movie.SearchMovies))
    
    // Metrics endpoints
    handler := expvar.Handler()
    r.GET("/metrics", gin.WrapH(handler))
    r.GET("/debug/vars", gin.WrapH(handler))
    
    return r
}