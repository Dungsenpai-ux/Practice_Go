package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Dungsenpai-ux/Practice_Go/config"
	"github.com/gin-gonic/gin"
)

type HealthController struct {
	config *config.Config
}

func NewHealthController(config *config.Config) *HealthController {
	return &HealthController{config: config}
}

func (h *HealthController) HealthCheck(c *gin.Context) {
	resp := h.build()
	c.JSON(http.StatusOK, resp)
}

// ServeHTTP allows usage without gin (net/http)
func (h *HealthController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := h.build()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *HealthController) build() HealthResponseDTO {
	return HealthResponseDTO{Status: "ok", Version: h.config.Version, Time: time.Now().UTC().Format(time.RFC3339)}
}
