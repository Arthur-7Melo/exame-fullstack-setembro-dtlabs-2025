package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	custom_errors "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MockRedisPublisher struct {
	mock.Mock
}

func (m *MockRedisPublisher) Publish(ctx context.Context, channel string, message interface{}) error {
	args := m.Called(ctx, channel, message)
	return args.Error(0)
}

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(notification *models.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) FindByUserID(userID uuid.UUID) ([]models.Notification, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Notification), args.Error(1)
}

func (m *MockNotificationRepository) FindActiveByUserID(userID uuid.UUID) ([]models.Notification, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Notification), args.Error(1)
}

func TestNotificationService_CreateNotification(t *testing.T) {
	userID := uuid.New()
	validConditions := []dto.NotificationCondition{
		{Parameter: "cpu", Operator: ">", Value: 80.0},
	}
	validDeviceIDs := []uuid.UUID{uuid.New()}

	t.Run("Success - Valid notification creation", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		mockNotifRepo.On("Create", mock.AnythingOfType("*models.Notification")).Return(nil)

		req := dto.CreateNotificationRequest{
			Name:        "Test Notification",
			Description: "Test Description",
			Enabled:     true,
			Conditions:  validConditions,
			DeviceIDs:   validDeviceIDs,
		}

		notification, err := service.CreateNotification(userID, req)

		assert.NoError(t, err)
		assert.NotNil(t, notification)
		assert.Equal(t, "Test Notification", notification.Name)
		assert.Equal(t, "Test Description", notification.Description)
		assert.True(t, notification.Enabled)
		assert.Equal(t, userID, notification.UserID)

		mockNotifRepo.AssertExpectations(t)
	})

	t.Run("Error - Empty notification name", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		req := dto.CreateNotificationRequest{
			Name:        "",
			Description: "Test Description",
			Conditions:  validConditions,
			DeviceIDs:   validDeviceIDs,
		}

		notification, err := service.CreateNotification(userID, req)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.NewValidationError("Notification name is required"), err)
		assert.Nil(t, notification)
	})

	t.Run("Error - Invalid parameter", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		invalidConditions := []dto.NotificationCondition{
			{Parameter: "invalid_param", Operator: ">", Value: 80.0},
		}

		req := dto.CreateNotificationRequest{
			Name:        "Test Notification",
			Description: "Test Description",
			Conditions:  invalidConditions,
			DeviceIDs:   validDeviceIDs,
		}

		notification, err := service.CreateNotification(userID, req)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.NewValidationError("Invalid parameter: invalid_param"), err)
		assert.Nil(t, notification)
	})

	t.Run("Error - Invalid operator", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		invalidConditions := []dto.NotificationCondition{
			{Parameter: "cpu", Operator: "invalid_op", Value: 80.0},
		}

		req := dto.CreateNotificationRequest{
			Name:        "Test Notification",
			Description: "Test Description",
			Conditions:  invalidConditions,
			DeviceIDs:   validDeviceIDs,
		}

		notification, err := service.CreateNotification(userID, req)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.NewValidationError("Invalid operator: invalid_op"), err)
		assert.Nil(t, notification)
	})

	t.Run("Error - Database error on Create", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		mockNotifRepo.On("Create", mock.AnythingOfType("*models.Notification")).Return(errors.New("database error"))

		req := dto.CreateNotificationRequest{
			Name:        "Test Notification",
			Description: "Test Description",
			Conditions:  validConditions,
			DeviceIDs:   validDeviceIDs,
		}

		notification, err := service.CreateNotification(userID, req)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, notification)

		mockNotifRepo.AssertExpectations(t)
	})
}

func TestNotificationService_GetUserNotifications(t *testing.T) {
	userID := uuid.New()
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

	t.Run("Success - Get user notifications", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		mockNotifRepo.On("FindByUserID", userID).Return(notifications, nil)

		result, err := service.GetUserNotifications(userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		mockNotifRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		mockNotifRepo.On("FindByUserID", userID).Return(([]models.Notification)(nil), errors.New("database error"))

		result, err := service.GetUserNotifications(userID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, result)

		mockNotifRepo.AssertExpectations(t)
	})
}

func TestNotificationService_CheckHeartbeat(t *testing.T) {
	userID := uuid.New()
	deviceID := uuid.New()
	device := &models.Device{
		UUID:   deviceID,
		UserID: userID,
		SN:     "123456789012",
	}

	heartbeat := &models.Heartbeat{
		DeviceID:     deviceID,
		CPU:          90.0,
		RAM:          85.0,
		DiskFree:     25.0,
		Temperature:  75.0,
		Latency:      100,
		Connectivity: 1,
		CreatedAt:    time.Now(),
	}

	conditionsJSON, _ := json.Marshal([]dto.NotificationCondition{
		{Parameter: "cpu", Operator: ">", Value: 80.0},
	})
	deviceIDsJSON, _ := json.Marshal([]uuid.UUID{deviceID})

	notification := models.Notification{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        "High CPU Alert",
		Description: "CPU usage is too high",
		Enabled:     true,
		Conditions:  datatypes.JSON(conditionsJSON),
		DeviceIDs:   datatypes.JSON(deviceIDsJSON),
	}

	t.Run("Success - Conditions met, notification sent", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockNotifRepo.On("FindActiveByUserID", userID).Return([]models.Notification{notification}, nil)
		mockRedis.On("Publish", mock.Anything, "notifications:"+userID.String(), mock.AnythingOfType("[]uint8")).Return(nil)

		err := service.CheckHeartbeat(heartbeat)

		assert.NoError(t, err)

		mockDeviceRepo.AssertExpectations(t)
		mockNotifRepo.AssertExpectations(t)
		mockRedis.AssertExpectations(t)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		mockDeviceRepo.On("FindByID", deviceID).Return((*models.Device)(nil), gorm.ErrRecordNotFound)

		err := service.CheckHeartbeat(heartbeat)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDeviceNotFound, err)

		mockDeviceRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on device lookup", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		mockDeviceRepo.On("FindByID", deviceID).Return((*models.Device)(nil), errors.New("database error"))

		err := service.CheckHeartbeat(heartbeat)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)

		mockDeviceRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on notifications lookup", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockNotifRepo.On("FindActiveByUserID", userID).Return(([]models.Notification)(nil), errors.New("database error"))

		err := service.CheckHeartbeat(heartbeat)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)

		mockDeviceRepo.AssertExpectations(t)
		mockNotifRepo.AssertExpectations(t)
	})

	t.Run("Success - Conditions not met, no notification sent", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		conditionsJSON, _ := json.Marshal([]dto.NotificationCondition{
			{Parameter: "cpu", Operator: "<", Value: 50.0},
		})
		notification.Conditions = datatypes.JSON(conditionsJSON)

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockNotifRepo.On("FindActiveByUserID", userID).Return([]models.Notification{notification}, nil)

		err := service.CheckHeartbeat(heartbeat)

		assert.NoError(t, err)

		mockRedis.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)

		mockDeviceRepo.AssertExpectations(t)
		mockNotifRepo.AssertExpectations(t)
	})

	t.Run("Success - Notification not for this device", func(t *testing.T) {
		mockNotifRepo := new(MockNotificationRepository)
		mockDeviceRepo := new(MockDeviceRepository)
		mockRedis := new(MockRedisPublisher)
		service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

		otherDeviceID := uuid.New()
		deviceIDsJSON, _ := json.Marshal([]uuid.UUID{otherDeviceID})
		notification.DeviceIDs = datatypes.JSON(deviceIDsJSON)
		notification.Conditions = datatypes.JSON(conditionsJSON) // Restaurar condições

		mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
		mockNotifRepo.On("FindActiveByUserID", userID).Return([]models.Notification{notification}, nil)

		err := service.CheckHeartbeat(heartbeat)

		assert.NoError(t, err)

		mockRedis.AssertNotCalled(t, "Publish", mock.Anything, mock.Anything, mock.Anything)

		mockDeviceRepo.AssertExpectations(t)
		mockNotifRepo.AssertExpectations(t)
	})

t.Run("Error - Redis publish error", func(t *testing.T) {
    mockNotifRepo := new(MockNotificationRepository)
    mockDeviceRepo := new(MockDeviceRepository)
    mockRedis := new(MockRedisPublisher)
    service := NewNotificationService(mockNotifRepo, mockDeviceRepo, mockRedis)

    conditionsJSON, _ := json.Marshal([]dto.NotificationCondition{
        {Parameter: "cpu", Operator: ">", Value: 80.0},
    })
    deviceIDsJSON, _ := json.Marshal([]uuid.UUID{deviceID})
    
    notification := models.Notification{
        ID:          uuid.New(),
        UserID:      userID,
        Name:        "High CPU Alert",
        Description: "CPU usage is too high",
        Enabled:     true,
        Conditions:  datatypes.JSON(conditionsJSON),
        DeviceIDs:   datatypes.JSON(deviceIDsJSON),
    }

    mockDeviceRepo.On("FindByID", deviceID).Return(device, nil)
    mockNotifRepo.On("FindActiveByUserID", userID).Return([]models.Notification{notification}, nil)
    mockRedis.On("Publish", mock.Anything, "notifications:"+userID.String(), mock.AnythingOfType("[]uint8")).Return(errors.New("redis error"))
    err := service.CheckHeartbeat(heartbeat)

    assert.NoError(t, err, "CheckHeartbeat should return nil even when Redis fails")
    mockRedis.AssertCalled(t, "Publish", mock.Anything, "notifications:"+userID.String(), mock.AnythingOfType("[]uint8"))

    mockDeviceRepo.AssertExpectations(t)
    mockNotifRepo.AssertExpectations(t)
    mockRedis.AssertExpectations(t)
})
}

func TestAppliesToDevice(t *testing.T) {
	service := &notificationService{}
	deviceID := uuid.New()
	otherDeviceID := uuid.New()

	t.Run("Applies to all devices (empty device IDs)", func(t *testing.T) {
		deviceIDsJSON, _ := json.Marshal([]uuid.UUID{})
		notification := models.Notification{
			DeviceIDs: datatypes.JSON(deviceIDsJSON),
		}
		result := service.appliesToDevice(notification, deviceID)
		assert.True(t, result)
	})

	t.Run("Applies to specific device", func(t *testing.T) {
		deviceIDs := []uuid.UUID{deviceID}
		deviceIDsJSON, _ := json.Marshal(deviceIDs)
		notification := models.Notification{
			DeviceIDs: datatypes.JSON(deviceIDsJSON),
		}
		result := service.appliesToDevice(notification, deviceID)
		assert.True(t, result)
	})

	t.Run("Does not apply to other device", func(t *testing.T) {
		deviceIDs := []uuid.UUID{otherDeviceID}
		deviceIDsJSON, _ := json.Marshal(deviceIDs)
		notification := models.Notification{
			DeviceIDs: datatypes.JSON(deviceIDsJSON),
		}
		result := service.appliesToDevice(notification, deviceID)
		assert.False(t, result)
	})

	t.Run("Invalid device IDs JSON", func(t *testing.T) {
		notification := models.Notification{
			DeviceIDs: datatypes.JSON("invalid json"),
		}
		result := service.appliesToDevice(notification, deviceID)
		assert.False(t, result)
	})
}

func TestCheckCondition(t *testing.T) {
	service := &notificationService{}
	heartbeat := &models.Heartbeat{
		CPU:          85.0,
		RAM:          75.0,
		DiskFree:     30.0,
		Temperature:  40.0,
		Latency:      100,
		Connectivity: 1,
	}

	t.Run("CPU condition - greater than", func(t *testing.T) {
		condition := dto.NotificationCondition{
			Parameter: "cpu",
			Operator:  ">",
			Value:     80.0,
		}
		result := service.checkCondition(condition, heartbeat)
		assert.True(t, result)
	})

	t.Run("CPU condition - less than", func(t *testing.T) {
		condition := dto.NotificationCondition{
			Parameter: "cpu",
			Operator:  "<",
			Value:     90.0,
		}
		result := service.checkCondition(condition, heartbeat)
		assert.True(t, result)
	})

	t.Run("RAM condition - equals", func(t *testing.T) {
		condition := dto.NotificationCondition{
			Parameter: "ram",
			Operator:  "==",
			Value:     75.0,
		}
		result := service.checkCondition(condition, heartbeat)
		assert.True(t, result)
	})

	t.Run("Invalid parameter", func(t *testing.T) {
		condition := dto.NotificationCondition{
			Parameter: "invalid",
			Operator:  ">",
			Value:     50.0,
		}
		result := service.checkCondition(condition, heartbeat)
		assert.False(t, result)
	})

	t.Run("Invalid operator", func(t *testing.T) {
		condition := dto.NotificationCondition{
			Parameter: "cpu",
			Operator:  "invalid",
			Value:     50.0,
		}
		result := service.checkCondition(condition, heartbeat)
		assert.False(t, result)
	})
}

func TestIsValidParameter(t *testing.T) {
	t.Run("Valid parameters", func(t *testing.T) {
		validParams := []string{"cpu", "ram", "disk_free", "temperature", "latency", "connectivity"}
		for _, param := range validParams {
			assert.True(t, isValidParameter(param), "Parameter %s should be valid", param)
		}
	})

	t.Run("Invalid parameters", func(t *testing.T) {
		invalidParams := []string{"", "invalid", "cpu_usage", "memory"}
		for _, param := range invalidParams {
			assert.False(t, isValidParameter(param), "Parameter %s should be invalid", param)
		}
	})
}

func TestIsValidOperator(t *testing.T) {
	t.Run("Valid operators", func(t *testing.T) {
		validOps := []string{">", "<", ">=", "<=", "==", "!="}
		for _, op := range validOps {
			assert.True(t, isValidOperator(op), "Operator %s should be valid", op)
		}
	})

	t.Run("Invalid operators", func(t *testing.T) {
		invalidOps := []string{"", "===", "=>", "=<", "!", "equals"}
		for _, op := range invalidOps {
			assert.False(t, isValidOperator(op), "Operator %s should be invalid", op)
		}
	})
}