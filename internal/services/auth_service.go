package services

import (
	"context"
	"errors"
	"net/mail"
	"time"

	"job-portal-api/internal/models"
	"job-portal-api/internal/repository"
	"job-portal-api/pkg/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("all fields are required")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email format")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	tokenString, err := utils.GenerateAccessToken(user)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil

}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Generate 6 digit random number
	// For simplicity, using a fixed randomizer or simple math. In production, use crypto/rand
	// Using utils? Or just implement here.
	resetToken := utils.GenerateRandomNumericString(6)
	expiresAt := time.Now().Add(60 * time.Minute)

	if err := s.userRepo.UpdatePasswordResetToken(ctx, user.ID, resetToken, expiresAt); err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, email, newPassword, token string) error {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.PasswordResetToken == nil || *user.PasswordResetToken != token {
		return errors.New("invalid token")
	}

	if user.PasswordResetExpires == nil || time.Now().After(*user.PasswordResetExpires) {
		return errors.New("token expired")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, user.ID, string(hashedPassword))
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("incorrect current password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}

func (s *AuthService) ChangeUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Admin check is done in handler/middleware. Here we just update.
	// Verify user exists
	_, err := s.userRepo.GetUserById(ctx, userID)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}
