package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/Dungsenpai-ux/Practice_Go/config"
	"net/http"
	"time"
)

type HealthController struct {
	config *config.Config
}

func NewHealthController(config *config.Config) *HealthController {
	return &HealthController{config: config}
}

func (h *HealthController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": h.config.Version,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}