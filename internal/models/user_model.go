package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                   uuid.UUID  `json:"id"`
	Username             string     `json:"username"`
	Password             string     `json:"-"` // Never send password in JSON response
	Email                string     `json:"email"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	IsAdmin              bool       `json:"is_admin"`
	ProfilePicture       FileUpload `json:"profile_picture"` // Default empty object
	PasswordResetToken   *string    `json:"-"`
	PasswordResetExpires *time.Time `json:"-"`
}
