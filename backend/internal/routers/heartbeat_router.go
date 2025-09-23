package routers

import (
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/handlers"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/middlewares"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupHeartbeatRoutes(router *gin.Engine, heartbeatHandler *handlers.HeartbeatHandler, jwtService services.JWTService) {
    authMiddleware := middlewares.AuthMiddleware(jwtService)
    heartbeatRoutes := router.Group("/api/v1/devices/:id/heartbeats")
    heartbeatRoutes.Use(authMiddleware)
    {
        heartbeatRoutes.GET("", heartbeatHandler.GetDeviceHeartbeats)
        heartbeatRoutes.GET("/latest", heartbeatHandler.GetLatestDeviceHeartbeat)
    }
}