package repository

import (
	"context"
	"errors"
	"fmt"
	"job-portal-api/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, email, password) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at, is_admin, profile_picture
	`
	err := r.pool.QueryRow(ctx, query, user.Username, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.IsAdmin, &user.ProfilePicture)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, password, is_admin, profile_picture, password_reset_token, password_reset_expires FROM users WHERE email = $1`
	var user models.User
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.ProfilePicture, &user.PasswordResetToken, &user.PasswordResetExpires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, username, email, password, is_admin, profile_picture, password_reset_token, password_reset_expires FROM users WHERE id = $1`
	var user models.User
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.ProfilePicture, &user.PasswordResetToken, &user.PasswordResetExpires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, is_admin = $3, profile_picture = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query, user.Username, user.Email, user.IsAdmin, user.ProfilePicture, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, email, is_admin, profile_picture, created_at, updated_at FROM users`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.ProfilePicture, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete jobs created by the user
	_, err = tx.Exec(ctx, "DELETE FROM jobs WHERE user_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user jobs: %w", err)
	}

	// Delete user
	_, err = tx.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdatePasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `UPDATE users SET password_reset_token = $1, password_reset_expires = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, token, expiresAt, userID)
	if err != nil {
		return fmt.Errorf("failed to update password reset token: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error {
	// Also clear the reset token
	query := `UPDATE users SET password = $1, password_reset_token = NULL, password_reset_expires = NULL, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, password, userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}
