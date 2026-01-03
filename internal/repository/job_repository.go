package repository

import (
	"context"
	"fmt"
	"job-portal-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository struct {
	pool *pgxpool.Pool
}

func NewJobRepository(pool *pgxpool.Pool) *JobRepository {
	return &JobRepository{pool: pool}
}

func (r *JobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	query := `
		INSERT INTO jobs (title, description, location, salary, experience_level, skills, job_type, company, company_logo, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		job.Title, job.Description, job.Location, job.Salary, job.ExperienceLevel, job.Skills, job.JobType, job.Company, job.CompanyLogo, job.UserID,
	).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}
	return nil
}

func (r *JobRepository) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	query := `SELECT id, title, description, location, salary, experience_level, skills, job_type, company, company_logo, created_at, updated_at, user_id FROM jobs`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all jobs: %w", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		if err := rows.Scan(
			&job.ID, &job.Title, &job.Description, &job.Location, &job.Salary, &job.ExperienceLevel, &job.Skills, &job.JobType, &job.Company, &job.CompanyLogo, &job.CreatedAt, &job.UpdatedAt, &job.UserID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *JobRepository) GetJobsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Job, error) {
	query := `SELECT id, title, description, location, salary, experience_level, skills, job_type, company, company_logo, created_at, updated_at, user_id FROM jobs WHERE user_id = $1`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs by user id: %w", err)
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		if err := rows.Scan(
			&job.ID, &job.Title, &job.Description, &job.Location, &job.Salary, &job.ExperienceLevel, &job.Skills, &job.JobType, &job.Company, &job.CompanyLogo, &job.CreatedAt, &job.UpdatedAt, &job.UserID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *JobRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	query := `SELECT id, title, description, location, salary, experience_level, skills, job_type, company, company_logo, created_at, updated_at, user_id FROM jobs WHERE id = $1`
	var job models.Job
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&job.ID, &job.Title, &job.Description, &job.Location, &job.Salary, &job.ExperienceLevel, &job.Skills, &job.JobType, &job.Company, &job.CompanyLogo, &job.CreatedAt, &job.UpdatedAt, &job.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get job by id: %w", err)
	}
	return &job, nil
}

func (r *JobRepository) UpdateJob(ctx context.Context, job *models.Job) error {
	query := `
		UPDATE jobs
		SET title = $1, description = $2, location = $3, salary = $4, experience_level = $5, skills = $6, job_type = $7, company = $8, company_logo = $9, updated_at = NOW()
		WHERE id = $10
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		job.Title, job.Description, job.Location, job.Salary, job.ExperienceLevel, job.Skills, job.JobType, job.Company, job.CompanyLogo, job.ID,
	).Scan(&job.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}
	return nil
}

func (r *JobRepository) DeleteJob(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM jobs WHERE id = $1`
	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("job not found")
	}
	return nil
}
