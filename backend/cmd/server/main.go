package main

import (
	"os"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/database"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/handlers"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/repository"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/routers"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	logger "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Logger.Warn(".env file not found", "warn", err.Error())
	}

	// Initialize database connection
	dbConfig := database.NewDBConfig()
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		logger.Logger.Error("Failed to connect to database", "error", err.Error())
		os.Exit(1)
	}

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		logger.Logger.Error("failed to auto migrate", "error", err.Error())
		os.Exit(1)
	}

	// Initialize services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "62774aa06a16f84f7acefe1c0be66aca07b665743eb459f90db56afd4deace4b" 
		logger.Logger.Warn("Warning: JWT_SECRET not set, using fallback key")
	}

	jwtService := services.NewJWTService(jwtSecret)
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(*userRepo, jwtService)

	router := gin.Default()

	authHandler := handlers.NewAuthHandler(authService)
	routers.SetupAuthRouter(router, authHandler, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	if err := router.Run(":" + port); err != nil {
		logger.Logger.Error("Failed to start server", "error", err.Error())
		os.Exit(1)
	}
}