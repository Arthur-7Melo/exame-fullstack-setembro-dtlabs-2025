package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/database"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/handlers"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/mq"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/repository"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/routers"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	logger "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/logger"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/pkg/redis"
	"github.com/gin-contrib/cors"
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

	redisClient := redis.NewRedisClient()
	logger.Logger.Info("Connected to Redis")

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
	heartbeatRepo := repository.NewHeartbeatRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	
	// Initialize services
	authService := services.NewAuthService(userRepo, jwtService)
	deviceService := services.NewDeviceService(deviceRepo)
	heartbeatService := services.NewHeartbeatService(heartbeatRepo, deviceRepo)
	notificationService := services.NewNotificationService(notificationRepo, deviceRepo, redisClient)
	
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)
	heartbeatHandler := handlers.NewHeartbeatHandler(heartbeatService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	url := ginSwagger.URL("/swagger/doc.json") 
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Health check endpoint including Redis
	router.GET("/health", func(c *gin.Context) {
		ctx := context.Background()
		if err := redisClient.Ping(ctx); err != nil {
			c.JSON(503, gin.H{"status": "unhealthy", "redis": "disconnected"})
			return
		}

		if err := db.Exec("SELECT 1").Error; err != nil {
			c.JSON(503, gin.H{"status": "unhealthy", "database": "disconnected"})
			return
		}

		c.JSON(200, gin.H{"status": "healthy", "redis": "connected", "database": "connected"})
	})

	routers.SetupAuthRouter(router, authHandler, jwtService)
	routers.SetupDeviceRoutes(router, deviceHandler, jwtService)
	routers.SetupHeartbeatRoutes(router, heartbeatHandler, jwtService)
	routers.SetupNotificationRoutes(router, notificationHandler, jwtService)

	amqpURL := os.Getenv("AMQP_URL")
    if amqpURL == "" {
      amqpURL = "amqp://guest:guest@rabbitmq:5672/"
    }
        
   heartbeatConsumer, err := mq.NewHeartbeatConsumer(
		amqpURL,
		"heartbeats",
		heartbeatService,
		notificationService,
		deviceRepo,
	)
    if err != nil {
			logger.Logger.Error("Failed to create heartbeat consumer", "error", err.Error())
			os.Exit(1)
    }
    defer heartbeatConsumer.Close()
        
    if err := heartbeatConsumer.Start(); err != nil {
			logger.Logger.Error("Failed to start heartbeat consumer", "error", err.Error())
			os.Exit(1)
    }
	logger.Logger.Info("Heartbeat consumer started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	go func() {
		logger.Logger.Info("Starting HTTP server", "port", port)
		if err := router.Run(":" + port); err != nil {
			logger.Logger.Error("Failed to start server", "error", err.Error())
			os.Exit(1)
		}
	}()

	logger.Logger.Info("Application started successfully. Press Ctrl+C to shutdown.")
	<-quit
	logger.Logger.Info("Shutting down application...")
}