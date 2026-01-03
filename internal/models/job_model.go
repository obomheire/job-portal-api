package models

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Location        string     `json:"location"`
	Salary          string     `json:"salary"`
	ExperienceLevel string     `json:"experience_level"`
	Skills          []string   `json:"skills"`
	JobType         string     `json:"job_type"`
	Company         string     `json:"company"`
	CompanyLogo     FileUpload `json:"company_logo"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	UserID          uuid.UUID  `json:"user_id"`
}
