package controller

import (
	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {
	// Đăng ký endpoint health
	healthController := NewHealthController(cfg)
	r.GET("/healthz", healthController.HealthCheck)
}