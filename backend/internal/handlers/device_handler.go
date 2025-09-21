package handlers

import (
	"net/http"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeviceHandler struct {
	deviceService services.DeviceService
}

func NewDeviceHandler(deviceService services.DeviceService) *DeviceHandler {
	return &DeviceHandler{deviceService: deviceService}
}

func (h *DeviceHandler) ListDevices(c *gin.Context) {
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

	devices, err := h.deviceService.ListDevices(uuidUserID)
	if err != nil {
		if customErr, ok := err.(errors.CustomError); ok {
			c.JSON(customErr.StatusCode(), dto.DetailedErrorResponse{
				Code:    dto.ErrorCodeFromStatusCode(customErr.StatusCode()),
				Message: customErr.Message(),
				Details: "Failed to list devices",
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

	c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) CreateDevice(c *gin.Context) {
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

	var req dto.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.DetailedErrorResponse{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	device, err := h.deviceService.CreateDevice(uuidUserID, req.Name, req.Location, req.SN, req.Description)
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

	c.JSON(http.StatusCreated, device)
}

func (h *DeviceHandler) GetDevice(c *gin.Context) {
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

	device, err := h.deviceService.GetDevice(uuidUserID, deviceID)
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

	c.JSON(http.StatusOK, device)
}

func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
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

	var req dto.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.DetailedErrorResponse{
			Code:    dto.ErrorCodeInvalidRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	device, err := h.deviceService.UpdateDevice(uuidUserID, deviceID, req.Name, req.Location, req.Description)
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

	c.JSON(http.StatusOK, device)
}

func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
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

	if err := h.deviceService.DeleteDevice(uuidUserID, deviceID); err != nil {
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

	c.AbortWithStatus(http.StatusNoContent)
}