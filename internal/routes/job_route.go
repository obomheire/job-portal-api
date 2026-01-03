package routes

import (
	"job-portal-api/internal/handlers"
	"job-portal-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterJobRoutes(r *gin.RouterGroup, handler *handlers.JobHandler) {
	jobs := r.Group("/jobs")
	jobs.Use(middleware.AuthMiddleware())
	{
		jobs.POST("/", handler.CreateJob)
		jobs.GET("/", handler.GetAllJobs)
		jobs.GET("/me", handler.GetJobsByUser)
		jobs.GET("/:id", handler.GetJobByID)
		jobs.PUT("/:id", handler.UpdateJob)
		jobs.DELETE("/:id", handler.DeleteJob)
	}
}
