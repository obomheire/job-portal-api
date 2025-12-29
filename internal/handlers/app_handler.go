package handlers

import (
	"job-portal-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	service *services.AppService
}

func NewAppHandler(service *services.AppService) *AppHandler {
	return &AppHandler{service: service}
}

func (h *AppHandler) HealthCheck(c *gin.Context) {
	status := h.service.HealthCheck(c.Request.Context())

	if status["database"] == "down" {
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}
