package routes

import (
	"job-portal-api/internal/handlers"
	"job-portal-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, handler *handlers.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/forgot-password", handler.ForgotPassword)
		auth.POST("/reset-password", handler.ResetPassword)
		auth.POST("/change-password", middleware.AuthMiddleware(), handler.ChangePassword)
		// Admin only
		auth.POST("/users/:id/change-password", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handler.ChangeUserPassword)
	}
}
