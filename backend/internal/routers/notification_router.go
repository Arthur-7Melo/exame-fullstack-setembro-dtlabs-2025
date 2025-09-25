package routers

import (
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/handlers"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/middlewares"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupNotificationRoutes(router *gin.Engine, notificationHandler *handlers.NotificationHandler, jwtService services.JWTService) {
	authMiddleware := middlewares.AuthMiddleware(jwtService)
	notificationRoutes := router.Group("/api/v1/notifications")
	notificationRoutes.Use(authMiddleware)
	{
		notificationRoutes.GET("", notificationHandler.GetNotifications)
		notificationRoutes.POST("", notificationHandler.CreateNotification)
	}
}