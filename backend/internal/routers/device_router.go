package routers

import (
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/handlers"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/middlewares"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupDeviceRoutes(router *gin.Engine, deviceHandler *handlers.DeviceHandler, jwtService services.JWTService) {
    authMiddleware := middlewares.AuthMiddleware(jwtService)
    deviceRoutes := router.Group("/api/v1/devices")
    deviceRoutes.Use(authMiddleware)
    {
        deviceRoutes.GET("", deviceHandler.ListDevices)
        deviceRoutes.POST("", deviceHandler.CreateDevice)
        deviceRoutes.GET("/:id", deviceHandler.GetDevice)
        deviceRoutes.PUT("/:id", deviceHandler.UpdateDevice)
        deviceRoutes.DELETE("/:id", deviceHandler.DeleteDevice)
    }
}