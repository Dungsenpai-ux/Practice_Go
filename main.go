package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/Dungsenpai-ux/Practice_Go/controller"
)

func main() {
	// Khởi tạo router Gin
	r := gin.Default()

	// Khởi tạo cấu hình
	cfg := config.NewConfig()

	// Đăng ký endpoint health
	healthController := controller.NewHealthController(cfg)
	r.GET("/healthz", healthController.HealthCheck)

	// Khởi chạy server
	r.Run(":8080")
}