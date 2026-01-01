package main

import (
	"log"
	"os"

	"job-portal-api/internal/handlers"
	"job-portal-api/internal/repository"
	"job-portal-api/internal/routes"
	"job-portal-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dsn := os.Getenv("POSTGRES_DB")
	if dsn == "" {
		log.Fatal("POSTGRES_DB environment variable is not set")
	}

	// Run migrations
	if err := repository.RunMigrations(dsn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	pool, err := repository.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	log.Println("Database connected successfully")

	r := gin.Default()

	// Initialize services and handlers
	appService := services.NewAppService(pool)
	appHandler := handlers.NewAppHandler(appService)

	userRepo := repository.NewUserRepository(pool)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup routes
	routes.RegisterAppRoutes(r, appHandler)
	routes.RegisterAuthRoutes(r, authHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
