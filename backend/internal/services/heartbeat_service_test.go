package services

import (
	"errors"
	"testing"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	custom_errors "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockHeartbeatRepository struct {
	mock.Mock
}

func (m *MockHeartbeatRepository) Create(heartbeat *models.Heartbeat) error {
	args := m.Called(heartbeat)
	return args.Error(0)
}

func (m *MockHeartbeatRepository) FindByDeviceID(deviceID uuid.UUID, startTime, endTime time.Time) ([]models.Heartbeat, error) {
	args := m.Called(deviceID, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Heartbeat), args.Error(1)
}

func (m *MockHeartbeatRepository) FindLatestByDeviceID(deviceID uuid.UUID) (*models.Heartbeat, error) {
	args := m.Called(deviceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Heartbeat), args.Error(1)
}

func TestHeartbeatService_CreateHeartbeat(t *testing.T) {
	deviceID := uuid.New()
	bootTime := time.Now().UTC().Add(-time.Hour * 24)
	cpu := 45.67
	ram := 67.89
	diskFree := 23.45
	temperature := 35.67
	latency := 150
	connectivity := 1

	t.Run("Success - Create heartbeat", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockHeartbeatRepo.On("Create", mock.AnythingOfType("*models.Heartbeat")).Return(nil)

		heartbeat, err := service.CreateHeartbeat(deviceID, cpu, ram, diskFree, temperature, latency, connectivity, bootTime)

		assert.NoError(t, err)
		assert.NotNil(t, heartbeat)
		assert.Equal(t, deviceID, heartbeat.DeviceID)
		assert.Equal(t, cpu, heartbeat.CPU)
		assert.Equal(t, ram, heartbeat.RAM)
		assert.Equal(t, diskFree, heartbeat.DiskFree)
		assert.Equal(t, temperature, heartbeat.Temperature)
		assert.Equal(t, latency, heartbeat.Latency)
		assert.Equal(t, connectivity, heartbeat.Connectivity)
		assert.Equal(t, bootTime, heartbeat.BootTime)

		mockHeartbeatRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on Create", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockHeartbeatRepo.On("Create", mock.AnythingOfType("*models.Heartbeat")).Return(errors.New("database error"))

		heartbeat, err := service.CreateHeartbeat(deviceID, cpu, ram, diskFree, temperature, latency, connectivity, bootTime)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, heartbeat)

		mockHeartbeatRepo.AssertExpectations(t)
	})
}

func TestHeartbeatService_GetDeviceHeartbeats(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	deviceID := uuid.New()
	startTime := time.Now().UTC().Add(-time.Hour * 24)
	endTime := time.Now().UTC()

	device := &models.Device{
		UUID:   deviceID,
		UserID: userID,
	}

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
			BootTime:     startTime.Add(time.Hour),
			CreatedAt:    startTime.Add(time.Hour),
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
			BootTime:     startTime.Add(2 * time.Hour),
			CreatedAt:    startTime.Add(2 * time.Hour),
		},
	}

	t.Run("Success - Get device heartbeats", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockHeartbeatRepo.On("FindByDeviceID", deviceID, startTime, endTime).Return(heartbeats, nil)

		result, err := service.GetDeviceHeartbeats(userID, deviceID, startTime, endTime)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, heartbeats[0].ID, result[0].ID)
		assert.Equal(t, heartbeats[1].ID, result[1].ID)

		mockDeviceRepo.AssertExpectations(t)
		mockHeartbeatRepo.AssertExpectations(t)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return((*models.Device)(nil), gorm.ErrRecordNotFound)

		result, err := service.GetDeviceHeartbeats(userID, deviceID, startTime, endTime)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDeviceNotFound, err)
		assert.Nil(t, result)

		mockDeviceRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on FindByID", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return((*models.Device)(nil), errors.New("database error"))

		result, err := service.GetDeviceHeartbeats(userID, deviceID, startTime, endTime)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, result)

		mockDeviceRepo.AssertExpectations(t)
	})

	t.Run("Error - Forbidden (different user)", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)

		result, err := service.GetDeviceHeartbeats(otherUserID, deviceID, startTime, endTime)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrForbidden, err)
		assert.Nil(t, result)

		mockDeviceRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on FindByDeviceID", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockHeartbeatRepo.On("FindByDeviceID", deviceID, startTime, endTime).Return(([]models.Heartbeat)(nil), errors.New("database error"))

		result, err := service.GetDeviceHeartbeats(userID, deviceID, startTime, endTime)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, result)

		mockDeviceRepo.AssertExpectations(t)
		mockHeartbeatRepo.AssertExpectations(t)
	})
}

func TestHeartbeatService_GetLatestDeviceHeartbeat(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	deviceID := uuid.New()

	device := &models.Device{
		UUID:   deviceID,
		UserID: userID,
	}

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

	t.Run("Success - Get latest device heartbeat", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockHeartbeatRepo.On("FindLatestByDeviceID", deviceID).Return(heartbeat, nil)

		result, err := service.GetLatestDeviceHeartbeat(userID, deviceID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, heartbeat.ID, result.ID)
		assert.Equal(t, deviceID, result.DeviceID)

		mockDeviceRepo.AssertExpectations(t)
		mockHeartbeatRepo.AssertExpectations(t)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return((*models.Device)(nil), gorm.ErrRecordNotFound)

		result, err := service.GetLatestDeviceHeartbeat(userID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDeviceNotFound, err)
		assert.Nil(t, result)

		mockDeviceRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on FindByID", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return((*models.Device)(nil), errors.New("database error"))

		result, err := service.GetLatestDeviceHeartbeat(userID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, result)

		mockDeviceRepo.AssertExpectations(t)
	})

	t.Run("Error - Forbidden (different user)", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)

		result, err := service.GetLatestDeviceHeartbeat(otherUserID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrForbidden, err)
		assert.Nil(t, result)

		mockDeviceRepo.AssertExpectations(t)
	})

	t.Run("Error - Heartbeat not found", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockHeartbeatRepo.On("FindLatestByDeviceID", deviceID).Return((*models.Heartbeat)(nil), gorm.ErrRecordNotFound)

		result, err := service.GetLatestDeviceHeartbeat(userID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDeviceNotFound, err)
		assert.Nil(t, result)

		mockDeviceRepo.AssertExpectations(t)
		mockHeartbeatRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on FindLatestByDeviceID", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockHeartbeatRepo.On("FindLatestByDeviceID", deviceID).Return((*models.Heartbeat)(nil), errors.New("database error"))

		result, err := service.GetLatestDeviceHeartbeat(userID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, result)

		mockDeviceRepo.AssertExpectations(t)
		mockHeartbeatRepo.AssertExpectations(t)
	})
}

func TestHeartbeatService_EdgeCases(t *testing.T) {
	userID := uuid.New()
	deviceID := uuid.New()
	startTime := time.Now().UTC().Add(-time.Hour * 24)
	endTime := time.Now().UTC()

	device := &models.Device{
		UUID:   deviceID,
		UserID: userID,
	}

	t.Run("Success - Empty heartbeats list", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		emptyHeartbeats := []models.Heartbeat{}

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockHeartbeatRepo.On("FindByDeviceID", deviceID, startTime, endTime).Return(emptyHeartbeats, nil)

		result, err := service.GetDeviceHeartbeats(userID, deviceID, startTime, endTime)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)

		mockDeviceRepo.AssertExpectations(t)
		mockHeartbeatRepo.AssertExpectations(t)
	})

	t.Run("Success - Single heartbeat", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		singleHeartbeat := []models.Heartbeat{
			{
				ID:           uuid.New(),
				DeviceID:     deviceID,
				CPU:          45.67,
				RAM:          67.89,
				CreatedAt:    time.Now().UTC().Add(-time.Hour),
			},
		}

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockHeartbeatRepo.On("FindByDeviceID", deviceID, startTime, endTime).Return(singleHeartbeat, nil)

		result, err := service.GetDeviceHeartbeats(userID, deviceID, startTime, endTime)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)

		mockDeviceRepo.AssertExpectations(t)
		mockHeartbeatRepo.AssertExpectations(t)
	})

	t.Run("Success - Extreme values in heartbeat", func(t *testing.T) {
		mockHeartbeatRepo := new(MockHeartbeatRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		service := NewHeartbeatService(mockHeartbeatRepo, mockDeviceRepo)

		mockHeartbeatRepo.On("Create", mock.AnythingOfType("*models.Heartbeat")).Return(nil)

		// Test with extreme values
		heartbeat, err := service.CreateHeartbeat(
			deviceID,
			100.0,    // CPU at 100%
			100.0,    // RAM at 100%
			0.0,      // Disk free at 0%
			100.0,    // Temperature at 100°C
			1000,     // High latency
			0,        // No connectivity
			time.Now().UTC().Add(-time.Hour*720), // 30 days ago
		)

		assert.NoError(t, err)
		assert.NotNil(t, heartbeat)
		assert.Equal(t, 100.0, heartbeat.CPU)
		assert.Equal(t, 100.0, heartbeat.RAM)
		assert.Equal(t, 0.0, heartbeat.DiskFree)
		assert.Equal(t, 100.0, heartbeat.Temperature)
		assert.Equal(t, 1000, heartbeat.Latency)
		assert.Equal(t, 0, heartbeat.Connectivity)

		mockHeartbeatRepo.AssertExpectations(t)
	})
}