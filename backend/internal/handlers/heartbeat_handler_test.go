package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	custom_errors "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHeartbeatService struct {
	mock.Mock
}

func (m *MockHeartbeatService) CreateHeartbeat(deviceID uuid.UUID, cpu, ram, diskFree, temperature float64, latency, connectivity int, bootTime time.Time) (*models.Heartbeat, error) {
	args := m.Called(deviceID, cpu, ram, diskFree, temperature, latency, connectivity, bootTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Heartbeat), args.Error(1)
}

func (m *MockHeartbeatService) GetDeviceHeartbeats(userID, deviceID uuid.UUID, startTime, endTime time.Time) ([]models.Heartbeat, error) {
	args := m.Called(userID, deviceID, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Heartbeat), args.Error(1)
}

func (m *MockHeartbeatService) GetLatestDeviceHeartbeat(userID, deviceID uuid.UUID) (*models.Heartbeat, error) {
	args := m.Called(userID, deviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Heartbeat), args.Error(1)
}

func TestHeartbeatHandler_GetDeviceHeartbeats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	deviceID := uuid.New()

	t.Run("Success - Get device heartbeats", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		heartbeats := []models.Heartbeat{
			{
				ID:           uuid.New(),
				DeviceID:     deviceID,
				CPU:          45.67,
				RAM:          67.89,
				DiskFree:     23.45,
				Temperature:  35.67,
				Latency:      150,
				Connectivity: 1,
				BootTime:     time.Now().UTC().Add(-time.Hour),
				CreatedAt:    time.Now().UTC().Add(-time.Hour),
			},
			{
				ID:           uuid.New(),
				DeviceID:     deviceID,
				CPU:          55.12,
				RAM:          72.34,
				DiskFree:     18.90,
				Temperature:  38.12,
				Latency:      120,
				Connectivity: 1,
				BootTime:     time.Now().UTC().Add(-2 * time.Hour),
				CreatedAt:    time.Now().UTC().Add(-2 * time.Hour),
			},
		}

		mockHeartbeatService.On("GetDeviceHeartbeats", userID, deviceID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(heartbeats, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)
		
		startTime := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second)
		endTime := time.Now().UTC().Truncate(time.Second)
		
		q := c.Request.URL.Query()
		q.Add("start", startTime.Format(time.RFC3339))
		q.Add("end", endTime.Format(time.RFC3339))
		c.Request.URL.RawQuery = q.Encode()

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Heartbeat
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, heartbeats[0].ID, response[0].ID)
		assert.Equal(t, heartbeats[1].ID, response[1].ID)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Success - Get device heartbeats with default time range", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		heartbeats := []models.Heartbeat{
			{
				ID:        uuid.New(),
				DeviceID:  deviceID,
				CPU:       45.67,
				RAM:       67.89,
				CreatedAt: time.Now().UTC().Add(-time.Hour),
			},
		}

		mockHeartbeatService.On("GetDeviceHeartbeats", userID, deviceID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(heartbeats, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Heartbeat
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Error - User ID not found in context", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized", response.Message)
		assert.Equal(t, "User ID not found in context", response.Details)
	})

	t.Run("Error - Invalid user ID type in context", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", "invalid-uuid-type")
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response.Message)
		assert.Equal(t, "Invalid user ID type", response.Details)
	})

	t.Run("Error - Invalid device ID", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid-uuid"}}
		c.Request, _ = http.NewRequest("GET", "/devices/invalid-uuid/heartbeats", nil)

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid device ID", response.Message)
	})

	t.Run("Error - Invalid start time format", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)
		q := c.Request.URL.Query()
		q.Add("start", "invalid-time-format")
		q.Add("end", time.Now().UTC().Format(time.RFC3339))
		c.Request.URL.RawQuery = q.Encode()

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid start time format", response.Message)
	})

	t.Run("Error - Invalid end time format", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)
		q := c.Request.URL.Query()
		q.Add("start", time.Now().UTC().Add(-24*time.Hour).Format(time.RFC3339))
		q.Add("end", "invalid-time-format")
		c.Request.URL.RawQuery = q.Encode()

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid end time format", response.Message)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		mockHeartbeatService.On("GetDeviceHeartbeats", userID, deviceID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(([]models.Heartbeat)(nil), custom_errors.ErrDeviceNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)
		
		startTime := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second)
		endTime := time.Now().UTC().Truncate(time.Second)
		
		q := c.Request.URL.Query()
		q.Add("start", startTime.Format(time.RFC3339))
		q.Add("end", endTime.Format(time.RFC3339))
		c.Request.URL.RawQuery = q.Encode()

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "device not found", response.Message)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Error - Forbidden access", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		mockHeartbeatService.On("GetDeviceHeartbeats", userID, deviceID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(([]models.Heartbeat)(nil), custom_errors.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)
		
		startTime := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second)
		endTime := time.Now().UTC().Truncate(time.Second)
		
		q := c.Request.URL.Query()
		q.Add("start", startTime.Format(time.RFC3339))
		q.Add("end", endTime.Format(time.RFC3339))
		c.Request.URL.RawQuery = q.Encode()

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "access to this resource is forbidden", response.Message)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Error - Database error", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		mockHeartbeatService.On("GetDeviceHeartbeats", userID, deviceID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(([]models.Heartbeat)(nil), custom_errors.ErrDatabaseError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)
		
		startTime := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second)
		endTime := time.Now().UTC().Truncate(time.Second)
		
		q := c.Request.URL.Query()
		q.Add("start", startTime.Format(time.RFC3339))
		q.Add("end", endTime.Format(time.RFC3339))
		c.Request.URL.RawQuery = q.Encode()

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "database error", response.Message)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Error - Generic error", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		mockHeartbeatService.On("GetDeviceHeartbeats", userID, deviceID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(([]models.Heartbeat)(nil), assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)
		
		startTime := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second)
		endTime := time.Now().UTC().Truncate(time.Second)
		
		q := c.Request.URL.Query()
		q.Add("start", startTime.Format(time.RFC3339))
		q.Add("end", endTime.Format(time.RFC3339))
		c.Request.URL.RawQuery = q.Encode()

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response.Message)

		mockHeartbeatService.AssertExpectations(t)
	})
}

func TestHeartbeatHandler_GetLatestDeviceHeartbeat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	deviceID := uuid.New()

	t.Run("Success - Get latest device heartbeat", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		heartbeat := &models.Heartbeat{
			ID:           uuid.New(),
			DeviceID:     deviceID,
			CPU:          45.67,
			RAM:          67.89,
			DiskFree:     23.45,
			Temperature:  35.67,
			Latency:      150,
			Connectivity: 1,
			BootTime:     time.Now().UTC().Add(-time.Hour * 12),
			CreatedAt:    time.Now().UTC().Add(-time.Minute * 5),
		}

		mockHeartbeatService.On("GetLatestDeviceHeartbeat", userID, deviceID).Return(heartbeat, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats/latest", nil)

		handler.GetLatestDeviceHeartbeat(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Heartbeat
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, heartbeat.ID, response.ID)
		assert.Equal(t, deviceID, response.DeviceID)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Error - User ID not found in context", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats/latest", nil)

		handler.GetLatestDeviceHeartbeat(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized", response.Message)
		assert.Equal(t, "User ID not found in context", response.Details)
	})

	t.Run("Error - Invalid user ID type in context", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", "invalid-uuid-type")
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats/latest", nil)

		handler.GetLatestDeviceHeartbeat(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response.Message)
		assert.Equal(t, "Invalid user ID type", response.Details)
	})

	t.Run("Error - Invalid device ID", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid-uuid"}}
		c.Request, _ = http.NewRequest("GET", "/devices/invalid-uuid/heartbeats/latest", nil)

		handler.GetLatestDeviceHeartbeat(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid device ID", response.Message)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		mockHeartbeatService.On("GetLatestDeviceHeartbeat", userID, deviceID).Return((*models.Heartbeat)(nil), custom_errors.ErrDeviceNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats/latest", nil)

		handler.GetLatestDeviceHeartbeat(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "device not found", response.Message)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Error - Forbidden access", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		mockHeartbeatService.On("GetLatestDeviceHeartbeat", userID, deviceID).Return((*models.Heartbeat)(nil), custom_errors.ErrForbidden)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats/latest", nil)

		handler.GetLatestDeviceHeartbeat(c)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "access to this resource is forbidden", response.Message)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Error - Database error", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		mockHeartbeatService.On("GetLatestDeviceHeartbeat", userID, deviceID).Return((*models.Heartbeat)(nil), custom_errors.ErrDatabaseError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats/latest", nil)

		handler.GetLatestDeviceHeartbeat(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "database error", response.Message)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Error - Generic error", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		mockHeartbeatService.On("GetLatestDeviceHeartbeat", userID, deviceID).Return((*models.Heartbeat)(nil), assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats/latest", nil)

		handler.GetLatestDeviceHeartbeat(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.DetailedErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Internal server error", response.Message)

		mockHeartbeatService.AssertExpectations(t)
	})
}

func TestHeartbeatHandler_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	deviceID := uuid.New()

	t.Run("Success - Empty heartbeats list", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		emptyHeartbeats := []models.Heartbeat{}

		mockHeartbeatService.On("GetDeviceHeartbeats", userID, deviceID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(emptyHeartbeats, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Heartbeat
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)

		mockHeartbeatService.AssertExpectations(t)
	})

	t.Run("Success - Single heartbeat", func(t *testing.T) {
		mockHeartbeatService := new(MockHeartbeatService)
		handler := NewHeartbeatHandler(mockHeartbeatService)

		singleHeartbeat := []models.Heartbeat{
			{
				ID:        uuid.New(),
				DeviceID:  deviceID,
				CPU:       45.67,
				RAM:       67.89,
				CreatedAt: time.Now().UTC().Add(-time.Hour),
			},
		}

		mockHeartbeatService.On("GetDeviceHeartbeats", userID, deviceID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(singleHeartbeat, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", userID)
		c.Params = gin.Params{gin.Param{Key: "id", Value: deviceID.String()}}
		c.Request, _ = http.NewRequest("GET", "/devices/"+deviceID.String()+"/heartbeats", nil)

		handler.GetDeviceHeartbeats(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Heartbeat
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)

		mockHeartbeatService.AssertExpectations(t)
	})
}