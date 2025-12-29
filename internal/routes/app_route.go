package routes

import (
	"job-portal-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAppRoutes(r *gin.Engine, handler *handlers.AppHandler) {
	api := r.Group("/api")
	{
		api.GET("/health", handler.HealthCheck)
	}
}
