package handlers

import (
	"net/http"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) CreateNotification(c *gin.Context) {
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

	var req dto.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.DetailedErrorResponse{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	notification, err := h.notificationService.CreateNotification(uuidUserID, req)
	if err != nil {
		if customErr, ok := err.(errors.CustomError); ok {
			c.JSON(customErr.StatusCode(), dto.DetailedErrorResponse{
				Code:    dto.ErrorCodeFromStatusCode(customErr.StatusCode()),
				Message: customErr.Message(),
				Details: "Please check your input and try again",
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

	c.JSON(http.StatusCreated, notification)
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
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

	notifications, err := h.notificationService.GetUserNotifications(uuidUserID)
	if err != nil {
		if customErr, ok := err.(errors.CustomError); ok {
			c.JSON(customErr.StatusCode(), dto.DetailedErrorResponse{
				Code:    dto.ErrorCodeFromStatusCode(customErr.StatusCode()),
				Message: customErr.Message(),
				Details: "Failed to get notifications",
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

	c.JSON(http.StatusOK, notifications)
}