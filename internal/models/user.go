package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"-"` // Never send password in JSON response
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsAdmin        bool      `json:"is_admin"`
	ProfilePicture *string   `json:"profile_picture,omitempty"` // Pointer to allow null values
}
