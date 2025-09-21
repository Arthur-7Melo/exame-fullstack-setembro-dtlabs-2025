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

	_ "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title DT Labs Exam API
// @version 1.0
// @description API for DT Labs Fullstack Exam - Device Management System

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	if err := godotenv.Load(); err != nil {
		logger.Logger.Warn(".env file not found", "warn", err.Error())
	}

	dbConfig := database.NewDBConfig()
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		logger.Logger.Error("Failed to connect to database", "error", err.Error())
		os.Exit(1)
	}

	if err := database.AutoMigrate(db); err != nil {
		logger.Logger.Error("failed to auto migrate", "error", err.Error())
		os.Exit(1)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "62774aa06a16f84f7acefe1c0be66aca07b665743eb459f90db56afd4deace4b" 
		logger.Logger.Warn("Warning: JWT_SECRET not set, using fallback key")
	}

	jwtService := services.NewJWTService(jwtSecret)
	
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	
	// Initialize services
	authService := services.NewAuthService(userRepo, jwtService)
	deviceService := services.NewDeviceService(deviceRepo)
	
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	router := gin.Default()

	url := ginSwagger.URL("/swagger/doc.json") 
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	routers.SetupAuthRouter(router, authHandler, jwtService)
		routers.SetupDeviceRoutes(router, deviceHandler, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	if err := router.Run(":" + port); err != nil {
		logger.Logger.Error("Failed to start server", "error", err.Error())
		os.Exit(1)
	}
}