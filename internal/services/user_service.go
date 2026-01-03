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
	cld      *cloudinary.Service
}

func NewUserService(userRepo *repository.UserRepository, cld *cloudinary.Service) *UserService {
	return &UserService{userRepo: userRepo, cld: cld}
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
