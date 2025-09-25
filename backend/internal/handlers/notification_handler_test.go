package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	custom_errors "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) CreateNotification(userID uuid.UUID, req dto.CreateNotificationRequest) (*models.Notification, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Notification), args.Error(1)
}

func (m *MockNotificationService) GetUserNotifications(userID uuid.UUID) ([]models.Notification, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Notification), args.Error(1)
}

func (m *MockNotificationService) CheckHeartbeat(heartbeat *models.Heartbeat) error {
	args := m.Called(heartbeat)
	return args.Error(0)
}

func TestNotificationHandler_CreateNotification(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	notificationID := uuid.New()

	t.Run("Success - Create notification", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		createReq := dto.CreateNotificationRequest{
			Name:        "Test Notification",
			Description: "Test Description",
			Enabled:     true,
			Conditions: []dto.NotificationCondition{
				{Parameter: "cpu", Operator: ">", Value: 80.0},
			},
			DeviceIDs: []uuid.UUID{uuid.New()},
		}

		notification := &models.Notification{
			ID:          notificationID,
			UserID:      userID,
			Name:        "Test Notification",
			Description: "Test Description",
			Enabled:     true,
		}

		mockNotificationService.On("CreateNotification", userID, createReq).Return(notification, nil)

		jsonData, _ := json.Marshal(createReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Request, _ = http.NewRequest("POST", "/notifications", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNotification(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Notification
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Test Notification", response.Name)
		assert.Equal(t, "Test Description", response.Description)
		assert.True(t, response.Enabled)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Request, _ = http.NewRequest("POST", "/notifications", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNotification(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("Error - Gin validation error (empty name)", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		createReq := dto.CreateNotificationRequest{
			Name:        "",
			Description: "Test Description",
			Enabled:     true,
			Conditions: []dto.NotificationCondition{
				{Parameter: "cpu", Operator: ">", Value: 80.0},
			},
			DeviceIDs: []uuid.UUID{uuid.New()},
		}

		jsonData, _ := json.Marshal(createReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Request, _ = http.NewRequest("POST", "/notifications", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNotification(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Invalid request body", response.Message)
		mockNotificationService.AssertNotCalled(t, "CreateNotification")
	})

	t.Run("Error - Validation error from service", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		createReq := dto.CreateNotificationRequest{
			Name:        "Valid Notification Name",
			Description: "Test Description",
			Enabled:     true,
			Conditions: []dto.NotificationCondition{
				{Parameter: "invalid_param", Operator: ">", Value: 80.0},
			},
			DeviceIDs: []uuid.UUID{uuid.New()},
		}

		validationError := custom_errors.NewValidationError("Invalid parameter: invalid_param")

		mockNotificationService.On("CreateNotification", userID, createReq).Return((*models.Notification)(nil), validationError)

		jsonData, _ := json.Marshal(createReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Request, _ = http.NewRequest("POST", "/notifications", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNotification(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Invalid parameter: invalid_param", response.Message)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("Error - User ID not found in context", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		createReq := dto.CreateNotificationRequest{
			Name:        "Test Notification",
			Description: "Test Description",
			Enabled:     true,
			Conditions: []dto.NotificationCondition{
				{Parameter: "cpu", Operator: ">", Value: 80.0},
			},
			DeviceIDs: []uuid.UUID{uuid.New()},
		}

		jsonData, _ := json.Marshal(createReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/notifications", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNotification(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Unauthorized", response.Message)
	})

	t.Run("Error - Invalid user ID type", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		createReq := dto.CreateNotificationRequest{
			Name:        "Test Notification",
			Description: "Test Description",
			Enabled:     true,
			Conditions: []dto.NotificationCondition{
				{Parameter: "cpu", Operator: ">", Value: 80.0},
			},
			DeviceIDs: []uuid.UUID{uuid.New()},
		}

		jsonData, _ := json.Marshal(createReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", "invalid-uuid-type")
		c.Request, _ = http.NewRequest("POST", "/notifications", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNotification(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Internal server error", response.Message)
	})

	t.Run("Error - Database error", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		createReq := dto.CreateNotificationRequest{
			Name:        "Test Notification",
			Description: "Test Description",
			Enabled:     true,
			Conditions: []dto.NotificationCondition{
				{Parameter: "cpu", Operator: ">", Value: 80.0},
			},
			DeviceIDs: []uuid.UUID{uuid.New()},
		}

		mockNotificationService.On("CreateNotification", userID, createReq).Return((*models.Notification)(nil), custom_errors.ErrDatabaseError)

		jsonData, _ := json.Marshal(createReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Request, _ = http.NewRequest("POST", "/notifications", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateNotification(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "database error", response.Message)
		mockNotificationService.AssertExpectations(t)
	})
}

func TestNotificationHandler_GetNotifications(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()

	t.Run("Success - Get notifications", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		notifications := []models.Notification{
			{
				ID:          uuid.New(),
				UserID:      userID,
				Name:        "Notification 1",
				Description: "Description 1",
				Enabled:     true,
			},
			{
				ID:          uuid.New(),
				UserID:      userID,
				Name:        "Notification 2",
				Description: "Description 2",
				Enabled:     false,
			},
		}

		mockNotificationService.On("GetUserNotifications", userID).Return(notifications, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Notification
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Len(t, response, 2)
		assert.Equal(t, "Notification 1", response[0].Name)
		assert.Equal(t, "Notification 2", response[1].Name)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("Success - Empty notifications list", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		emptyNotifications := []models.Notification{}

		mockNotificationService.On("GetUserNotifications", userID).Return(emptyNotifications, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Notification
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Len(t, response, 0)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("Error - User ID not found in context", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Unauthorized", response.Message)
	})

	t.Run("Error - Invalid user ID type", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", "invalid-uuid-type")

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Internal server error", response.Message)
	})

	t.Run("Error - Database error", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		mockNotificationService.On("GetUserNotifications", userID).Return(([]models.Notification)(nil), custom_errors.ErrDatabaseError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "database error", response.Message)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("Error - Validation error", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		validationError := custom_errors.NewValidationError("validation error")

		mockNotificationService.On("GetUserNotifications", userID).Return(([]models.Notification)(nil), validationError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "validation error", response.Message)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("Error - Forbidden error", func(t *testing.T) {
		mockNotificationService := new(MockNotificationService)
		handler := NewNotificationHandler(mockNotificationService)

		mockNotificationService.On("GetUserNotifications", userID).Return(([]models.Notification)(nil), custom_errors.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)

		handler.GetNotifications(c)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "access to this resource is forbidden", response.Message)
		mockNotificationService.AssertExpectations(t)
	})
}