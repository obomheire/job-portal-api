package main

import (
	"log"
	"os"

	"job-portal-api/internal/handlers"
	"job-portal-api/internal/repository"
	"job-portal-api/internal/routes"
	"job-portal-api/internal/services"
	"job-portal-api/pkg/cloudinary"

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

	// Initialize Cloudinary
	cldService, err := cloudinary.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary service: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(pool)
	jobRepo := repository.NewJobRepository(pool)

	// Initialize services
	appService := services.NewAppService(pool)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo, jobRepo, cldService)
	jobService := services.NewJobService(jobRepo, cldService)

	// Initialize handlers
	appHandler := handlers.NewAppHandler(appService)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	jobHandler := handlers.NewJobHandler(jobService)

	// Setup routes
	api := r.Group("/api")
	routes.RegisterAppRoutes(api, appHandler)
	routes.RegisterAuthRoutes(api, authHandler)
	routes.RegisterUserRoutes(api, userHandler)
	routes.RegisterJobRoutes(api, jobHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
