package routes

import (
	"job-portal-api/internal/handlers"
	"job-portal-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup, handler *handlers.UserHandler) {
	user := r.Group("/users")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/:id", handler.GetUserById)
		user.PUT("/:id", handler.UpdateUser)
	}
}
