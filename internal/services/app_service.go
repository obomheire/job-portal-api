package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppService struct {
	db *pgxpool.Pool
}

func NewAppService(db *pgxpool.Pool) *AppService {
	return &AppService{db: db}
}

func (s *AppService) HealthCheck(ctx context.Context) map[string]string {
	status := make(map[string]string)
	status["server"] = "running"

	// Create a context with timeout for the ping
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.db.Ping(ctx); err != nil {
		status["database"] = "down"
		status["error"] = err.Error()
	} else {
		status["database"] = "up"
	}

	return status
}
