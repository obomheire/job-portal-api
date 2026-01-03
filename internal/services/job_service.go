package services

import (
	"context"
	"errors"
	"mime/multipart"

	"job-portal-api/internal/models"
	"job-portal-api/internal/repository"
	"job-portal-api/pkg/cloudinary"

	"github.com/google/uuid"
)

type JobService struct {
	repo       *repository.JobRepository
	cldService *cloudinary.Service
}

func NewJobService(repo *repository.JobRepository, cldService *cloudinary.Service) *JobService {
	return &JobService{
		repo:       repo,
		cldService: cldService,
	}
}

func (s *JobService) CreateJob(ctx context.Context, job *models.Job, file multipart.File, filename string) (*models.Job, error) {
	if file != nil {
		imageUrl, err := s.cldService.UploadImage(ctx, file, filename)
		if err != nil {
			return nil, err
		}
		job.CompanyLogo = imageUrl
	}

	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *JobService) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	return s.repo.GetAllJobs(ctx)
}

func (s *JobService) GetJobsByUser(ctx context.Context, userID uuid.UUID) ([]models.Job, error) {
	return s.repo.GetJobsByUserID(ctx, userID)
}

func (s *JobService) GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	return s.repo.GetJobByID(ctx, id)
}

func (s *JobService) UpdateJob(ctx context.Context, jobID uuid.UUID, updateData *models.Job, file multipart.File, filename string, requestUser *models.User) (*models.Job, error) {
	existingJob, err := s.repo.GetJobByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	// Authorization check
	if !requestUser.IsAdmin && existingJob.UserID != requestUser.ID {
		return nil, errors.New("unauthorized to update this job")
	}

	// Update fields
	existingJob.Title = updateData.Title
	existingJob.Description = updateData.Description
	existingJob.Location = updateData.Location
	existingJob.Salary = updateData.Salary
	existingJob.ExperienceLevel = updateData.ExperienceLevel
	existingJob.Skills = updateData.Skills
	existingJob.JobType = updateData.JobType
	existingJob.Company = updateData.Company

	if file != nil {
		imageUrl, err := s.cldService.UploadImage(ctx, file, filename)
		if err != nil {
			return nil, err
		}
		existingJob.CompanyLogo = imageUrl
	}

	if err := s.repo.UpdateJob(ctx, existingJob); err != nil {
		return nil, err
	}

	return existingJob, nil
}

func (s *JobService) DeleteJob(ctx context.Context, id uuid.UUID, requestUser *models.User) error {
	existingJob, err := s.repo.GetJobByID(ctx, id)
	if err != nil {
		return err
	}

	if !requestUser.IsAdmin && existingJob.UserID != requestUser.ID {
		return errors.New("unauthorized to delete this job")
	}

	return s.repo.DeleteJob(ctx, id)
}
