package routes

import (
	"job-portal-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAppRoutes(r *gin.RouterGroup, handler *handlers.AppHandler) {
	api := r.Group("/app")
	{
		api.GET("/health", handler.HealthCheck)
	}
}
