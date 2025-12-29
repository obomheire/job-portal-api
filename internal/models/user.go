package models

import (
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"-"` // Never send password in JSON response
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsAdmin        bool      `json:"is_admin"`
	ProfilePicture *string   `json:"profile_picture,omitempty"` // Pointer to allow null values
}
