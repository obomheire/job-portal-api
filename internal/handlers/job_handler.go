package handlers

import (
	"job-portal-api/internal/models"
	"job-portal-api/internal/services"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JobHandler struct {
	service *services.JobService
}

func NewJobHandler(service *services.JobService) *JobHandler {
	return &JobHandler{service: service}
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	userIdStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	var job models.Job
	job.Title = c.PostForm("title")
	job.Description = c.PostForm("description")
	job.Location = c.PostForm("location")
	job.Salary = c.PostForm("salary")
	job.ExperienceLevel = c.PostForm("experience_level")
	job.JobType = c.PostForm("job_type")
	job.Company = c.PostForm("company")
	job.Skills = c.PostFormArray("skills")
	job.UserID = userID

	if job.Title == "" || job.Description == "" || job.Location == "" || job.Salary == "" || job.ExperienceLevel == "" || job.Company == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All required fields must be provided"})
		return
	}

	var file multipart.File
	var filename string

	fileHeader, err := c.FormFile("company_logo")
	if err == nil {
		f, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer f.Close()
		file = f
		filename = fileHeader.Filename
	}

	createdJob, err := h.service.CreateJob(c.Request.Context(), &job, file, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdJob)
}

func (h *JobHandler) GetAllJobs(c *gin.Context) {
	jobs, err := h.service.GetAllJobs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetJobsByUser(c *gin.Context) {
	userIdStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	jobs, err := h.service.GetJobsByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetJobByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := h.service.GetJobByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	userIdStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIdStr)
	isAdmin := c.GetBool("is_admin")

	requestUser := &models.User{
		ID:      userID,
		IsAdmin: isAdmin,
	}

	var job models.Job
	job.Title = c.PostForm("title")
	job.Description = c.PostForm("description")
	job.Location = c.PostForm("location")
	job.Salary = c.PostForm("salary")
	job.ExperienceLevel = c.PostForm("experience_level")
	job.JobType = c.PostForm("job_type")
	job.Company = c.PostForm("company")
	job.Skills = c.PostFormArray("skills")

	var file multipart.File
	var filename string

	fileHeader, err := c.FormFile("company_logo")
	if err == nil {
		f, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer f.Close()
		file = f
		filename = fileHeader.Filename
	}

	updatedJob, err := h.service.UpdateJob(c.Request.Context(), id, &job, file, filename, requestUser)
	if err != nil {
		if err.Error() == "unauthorized to update this job" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedJob)
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	userIdStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIdStr)
	isAdmin := c.GetBool("is_admin")

	requestUser := &models.User{
		ID:      userID,
		IsAdmin: isAdmin,
	}

	if err := h.service.DeleteJob(c.Request.Context(), id, requestUser); err != nil {
		if err.Error() == "unauthorized to delete this job" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully"})
}
