package services

import (
	"context"
	"job-portal-api/internal/models"
	"job-portal-api/internal/repository"
	"job-portal-api/pkg/cloudinary"
	"mime/multipart"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
	jobRepo  *repository.JobRepository
	cld      *cloudinary.Service
}

func NewUserService(userRepo *repository.UserRepository, jobRepo *repository.JobRepository, cld *cloudinary.Service) *UserService {
	return &UserService{userRepo: userRepo, jobRepo: jobRepo, cld: cld}
}

func (s *UserService) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserById(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *UserService) UploadProfilePicture(ctx context.Context, userID uuid.UUID, file multipart.File) (string, error) {
	// Check if user exists
	user, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return "", err
	}

	// Upload to Cloudinary using userID as filename to overwrite existing
	imageURL, publicID, err := s.cld.UploadImage(ctx, file, userID.String())
	if err != nil {
		return "", err
	}

	// Update user record
	user.ProfilePicture = models.FileUpload{
		URL:      imageURL,
		PublicID: publicID,
	}
	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return imageURL, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	// Delete profile picture if exists
	if user.ProfilePicture.PublicID != "" {
		if err := s.cld.DeleteAsset(ctx, user.ProfilePicture.PublicID); err != nil {
			return err
		}
	}

	// Get all jobs created by the user
	jobs, err := s.jobRepo.GetJobsByUserID(ctx, id)
	if err != nil {
		return err
	}

	// Delete company logos for each job
	for _, job := range jobs {
		if job.CompanyLogo.PublicID != "" {
			if err := s.cld.DeleteAsset(ctx, job.CompanyLogo.PublicID); err != nil {
				// We log or simply return error. Returning error seems safer to ensure consistency,
				// though it might block deletion if one image fails.
				// Given the previous pattern, let's return error.
				return err
			}
		}
	}

	return s.userRepo.DeleteUser(ctx, id)
}
