package handlers

import (
	"net/http"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HeartbeatHandler struct {
    heartbeatService services.HeartbeatService
}

func NewHeartbeatHandler(heartbeatService services.HeartbeatService) *HeartbeatHandler {
    return &HeartbeatHandler{heartbeatService: heartbeatService}
}

func (h *HeartbeatHandler) GetDeviceHeartbeats(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInvalidCredentials,
            Message: "Unauthorized",
            Details: "User ID not found in context",
        })
        return
    }

    uuidUserID, ok := userID.(uuid.UUID)
    if !ok {
        c.JSON(http.StatusInternalServerError, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInternalError,
            Message: "Internal server error",
            Details: "Invalid user ID type",
        })
        return
    }

    deviceID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInvalidRequest,
            Message: "Invalid device ID",
            Details: err.Error(),
        })
        return
    }

    // Parse query parameters for time range
    startTimeStr := c.DefaultQuery("start", time.Now().Add(-24*time.Hour).Format(time.RFC3339))
    endTimeStr := c.DefaultQuery("end", time.Now().Format(time.RFC3339))

    startTime, err := time.Parse(time.RFC3339, startTimeStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInvalidRequest,
            Message: "Invalid start time format",
            Details: "Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)",
        })
        return
    }

    endTime, err := time.Parse(time.RFC3339, endTimeStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInvalidRequest,
            Message: "Invalid end time format",
            Details: "Use RFC3339 format (e.g., 2023-01-01T00:00:00Z)",
        })
        return
    }

    heartbeats, err := h.heartbeatService.GetDeviceHeartbeats(uuidUserID, deviceID, startTime, endTime)
    if err != nil {
        if customErr, ok := err.(errors.CustomError); ok {
            c.JSON(customErr.StatusCode(), dto.DetailedErrorResponse{
                Code:    dto.ErrorCodeFromStatusCode(customErr.StatusCode()),
                Message: customErr.Message(),
                Details: "Failed to get device heartbeats",
            })
        } else {
            c.JSON(http.StatusInternalServerError, dto.DetailedErrorResponse{
                Code:    dto.ErrorCodeInternalError,
                Message: "Internal server error",
                Details: err.Error(),
            })
        }
        return
    }

    c.JSON(http.StatusOK, heartbeats)
}

func (h *HeartbeatHandler) GetLatestDeviceHeartbeat(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInvalidCredentials,
            Message: "Unauthorized",
            Details: "User ID not found in context",
        })
        return
    }

    uuidUserID, ok := userID.(uuid.UUID)
    if !ok {
        c.JSON(http.StatusInternalServerError, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInternalError,
            Message: "Internal server error",
            Details: "Invalid user ID type",
        })
        return
    }

    deviceID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.DetailedErrorResponse{
            Code:    dto.ErrorCodeInvalidRequest,
            Message: "Invalid device ID",
            Details: err.Error(),
        })
        return
    }

    heartbeat, err := h.heartbeatService.GetLatestDeviceHeartbeat(uuidUserID, deviceID)
    if err != nil {
        if customErr, ok := err.(errors.CustomError); ok {
            c.JSON(customErr.StatusCode(), dto.DetailedErrorResponse{
                Code:    dto.ErrorCodeFromStatusCode(customErr.StatusCode()),
                Message: customErr.Message(),
                Details: "Failed to get latest device heartbeat",
            })
        } else {
            c.JSON(http.StatusInternalServerError, dto.DetailedErrorResponse{
                Code:    dto.ErrorCodeInternalError,
                Message: "Internal server error",
                Details: err.Error(),
            })
        }
        return
    }

    c.JSON(http.StatusOK, heartbeat)
}