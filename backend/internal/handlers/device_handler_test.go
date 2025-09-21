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

type MockDeviceService struct {
	mock.Mock
}

func (m *MockDeviceService) CreateDevice(userID uuid.UUID, name, location, sn, description string) (*models.Device, error) {
	args := m.Called(userID, name, location, sn, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockDeviceService) GetDevice(userID, deviceID uuid.UUID) (*models.Device, error) {
	args := m.Called(userID, deviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockDeviceService) ListDevices(userID uuid.UUID) ([]models.Device, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Device), args.Error(1)
}

func (m *MockDeviceService) UpdateDevice(userID, deviceID uuid.UUID, name, location, description string) (*models.Device, error) {
	args := m.Called(userID, deviceID, name, location, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockDeviceService) DeleteDevice(userID, deviceID uuid.UUID) error {
	args := m.Called(userID, deviceID)
	return args.Error(0)
}

func TestDeviceHandler_ListDevices(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()

	t.Run("Success - List devices", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		devices := []models.Device{
			{UUID: uuid.New(), Name: "Device 1", UserID: userID},
			{UUID: uuid.New(), Name: "Device 2", UserID: userID},
		}

		mockDeviceService.On("ListDevices", userID).Return(devices, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)

		handler.ListDevices(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Device
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Len(t, response, 2)
		mockDeviceService.AssertExpectations(t)
	})

	t.Run("Error - User ID not found in context", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.ListDevices(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Unauthorized", response.Message)
	})

	t.Run("Error - Database error", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		mockDeviceService.On("ListDevices", userID).Return(([]models.Device)(nil), custom_errors.ErrDatabaseError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)

		handler.ListDevices(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "database error", response.Message)
		mockDeviceService.AssertExpectations(t)
	})
}

func TestDeviceHandler_CreateDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	deviceID := uuid.New()
	validSN := "123456789012"

	t.Run("Success - Create device", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		device := &models.Device{
			UUID:        deviceID,
			Name:        "Test Device",
			Location:    "Test Location",
			SN:          validSN,
			Description: "Test Description",
			UserID:      userID,
		}

		createReq := dto.CreateDeviceRequest{
			Name:        "Test Device",
			Location:    "Test Location",
			SN:          validSN,
			Description: "Test Description",
		}

		mockDeviceService.On("CreateDevice", userID, createReq.Name, createReq.Location, createReq.SN, createReq.Description).Return(device, nil)

		jsonData, _ := json.Marshal(createReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Request, _ = http.NewRequest("POST", "/devices", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateDevice(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Device
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Test Device", response.Name)
		assert.Equal(t, validSN, response.SN)
		mockDeviceService.AssertExpectations(t)
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Request, _ = http.NewRequest("POST", "/devices", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateDevice(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Invalid request body", response.Message)
	})

	t.Run("Error - Device already exists", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		createReq := dto.CreateDeviceRequest{
			Name:        "Test Device",
			Location:    "Test Location",
			SN:          validSN,
			Description: "Test Description",
		}

		mockDeviceService.On("CreateDevice", userID, createReq.Name, createReq.Location, createReq.SN, createReq.Description).Return((*models.Device)(nil), custom_errors.ErrDeviceAlreadyExists)

		jsonData, _ := json.Marshal(createReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Request, _ = http.NewRequest("POST", "/devices", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateDevice(c)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "device with this serial number already exists", response.Message)
		mockDeviceService.AssertExpectations(t)
	})
}

func TestDeviceHandler_GetDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	deviceID := uuid.New()

	t.Run("Success - Get device", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		device := &models.Device{
			UUID:   deviceID,
			Name:   "Test Device",
			UserID: userID,
		}

		mockDeviceService.On("GetDevice", userID, deviceID).Return(device, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}

		handler.GetDevice(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Device
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, deviceID, response.UUID)
		mockDeviceService.AssertExpectations(t)
	})

	t.Run("Error - Invalid device ID", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid-uuid"}}

		handler.GetDevice(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Invalid device ID", response.Message)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		mockDeviceService.On("GetDevice", userID, deviceID).Return((*models.Device)(nil), custom_errors.ErrDeviceNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}

		handler.GetDevice(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "device not found", response.Message)
		mockDeviceService.AssertExpectations(t)
	})
}

func TestDeviceHandler_UpdateDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	deviceID := uuid.New()

	t.Run("Success - Update device", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		device := &models.Device{
			UUID:        deviceID,
			Name:        "Updated Device",
			Location:    "Updated Location",
			Description: "Updated Description",
			UserID:      userID,
		}

		updateReq := dto.UpdateDeviceRequest{
			Name:        "Updated Device",
			Location:    "Updated Location",
			Description: "Updated Description",
		}

		mockDeviceService.On("UpdateDevice", userID, deviceID, updateReq.Name, updateReq.Location, updateReq.Description).Return(device, nil)

		jsonData, _ := json.Marshal(updateReq)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("PUT", "/devices/"+deviceID.String(), bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateDevice(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Device
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Updated Device", response.Name)
		assert.Equal(t, "Updated Description", response.Description)
		mockDeviceService.AssertExpectations(t)
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("PUT", "/devices/"+deviceID.String(), bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateDevice(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Invalid request body", response.Message)
	})
}
func TestDeviceHandler_DeleteDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	deviceID := uuid.New()

	t.Run("Success - Delete device", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		mockDeviceService.On("DeleteDevice", userID, deviceID).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}

		handler.DeleteDevice(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String()) 
		mockDeviceService.AssertExpectations(t)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		mockDeviceService.On("DeleteDevice", userID, deviceID).Return(custom_errors.ErrDeviceNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}

		handler.DeleteDevice(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "device not found", response.Message)
		mockDeviceService.AssertExpectations(t)
	})

	t.Run("Error - Database error on delete", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		mockDeviceService.On("DeleteDevice", userID, deviceID).Return(custom_errors.ErrDatabaseError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}

		handler.DeleteDevice(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "database error", response.Message)
		mockDeviceService.AssertExpectations(t)
	})

	t.Run("Error - Invalid device ID", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid-uuid"}}

		handler.DeleteDevice(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Invalid device ID", response.Message)
	})

	t.Run("Error - User ID not found in context", func(t *testing.T) {
		mockDeviceService := new(MockDeviceService)
		handler := NewDeviceHandler(mockDeviceService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}

		handler.DeleteDevice(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response dto.DetailedErrorResponse
		json.Unmarshal(w.Body.Bytes(), &response)

		assert.Equal(t, "Unauthorized", response.Message)
	})
}