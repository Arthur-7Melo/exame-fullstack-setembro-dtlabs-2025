package services

import (
	"errors"
	"testing"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	custom_errors "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDeviceRepository struct {
	mock.Mock
}

func (m *MockDeviceRepository) FindBySN(sn string) (*models.Device, error) {
	args := m.Called(sn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockDeviceRepository) Create(device *models.Device) error {
	args := m.Called(device)
	return args.Error(0)
}

func (m *MockDeviceRepository) FindByID(id uuid.UUID) (*models.Device, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Device), args.Error(1)
}

func (m *MockDeviceRepository) FindByUserID(userID uuid.UUID) ([]models.Device, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Device), args.Error(1)
}

func (m *MockDeviceRepository) Update(device *models.Device) error {
	args := m.Called(device)
	return args.Error(0)
}

func (m *MockDeviceRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestDeviceService_CreateDevice(t *testing.T) {
	userID := uuid.New()
	validSN := "123456789012"

	t.Run("Success - Valid device creation", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindBySN", validSN).Return((*models.Device)(nil), gorm.ErrRecordNotFound)
		mockRepo.On("Create", mock.AnythingOfType("*models.Device")).Return(nil)

		device, err := service.CreateDevice(userID, "Test Device", "Test Location", validSN, "Test Description")

		assert.NoError(t, err)
		assert.NotNil(t, device)
		assert.Equal(t, "Test Device", device.Name)
		assert.Equal(t, "Test Location", device.Location)
		assert.Equal(t, validSN, device.SN)
		assert.Equal(t, "Test Description", device.Description)
		assert.Equal(t, userID, device.UserID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Empty device name", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		device, err := service.CreateDevice(userID, "", "Test Location", validSN, "Test Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.NewValidationError("Device name is required"), err)
		assert.Nil(t, device)
	})

	t.Run("Error - Empty location", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		device, err := service.CreateDevice(userID, "Test Device", "", validSN, "Test Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.NewValidationError("Device location is required"), err)
		assert.Nil(t, device)
	})

	t.Run("Error - Empty SN", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		device, err := service.CreateDevice(userID, "Test Device", "Test Location", "", "Test Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.NewValidationError("Device serial number is required"), err)
		assert.Nil(t, device)
	})

	t.Run("Error - Invalid SN format", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		device, err := service.CreateDevice(userID, "Test Device", "Test Location", "123", "Test Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.NewValidationError("Serial number must be exactly 12 digits"), err)
		assert.Nil(t, device)
	})

	t.Run("Error - SN already exists", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		existingDevice := &models.Device{SN: validSN}
		mockRepo.On("FindBySN", validSN).Return(existingDevice, nil)

		device, err := service.CreateDevice(userID, "Test Device", "Test Location", validSN, "Test Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDeviceAlreadyExists, err)
		assert.Nil(t, device)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on FindBySN", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindBySN", validSN).Return((*models.Device)(nil), errors.New("database error"))

		device, err := service.CreateDevice(userID, "Test Device", "Test Location", validSN, "Test Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, device)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on Create", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindBySN", validSN).Return((*models.Device)(nil), gorm.ErrRecordNotFound)
		mockRepo.On("Create", mock.AnythingOfType("*models.Device")).Return(errors.New("database error"))

		device, err := service.CreateDevice(userID, "Test Device", "Test Location", validSN, "Test Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, device)

		mockRepo.AssertExpectations(t)
	})
}

func TestDeviceService_GetDevice(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	deviceID := uuid.New()
	device := &models.Device{
		UUID:   deviceID,
		Name:   "Test Device",
		UserID: userID,
	}

	t.Run("Success - Get device", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return(device, nil)

		result, err := service.GetDevice(userID, deviceID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, deviceID, result.UUID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return((*models.Device)(nil), gorm.ErrRecordNotFound)

		result, err := service.GetDevice(userID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDeviceNotFound, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return((*models.Device)(nil), errors.New("database error"))

		result, err := service.GetDevice(userID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Forbidden (different user)", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return(device, nil)

		result, err := service.GetDevice(otherUserID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrForbidden, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})
}

func TestDeviceService_ListDevices(t *testing.T) {
	userID := uuid.New()
	devices := []models.Device{
		{UUID: uuid.New(), Name: "Device 1", UserID: userID},
		{UUID: uuid.New(), Name: "Device 2", UserID: userID},
	}

	t.Run("Success - List devices", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByUserID", userID).Return(devices, nil)

		result, err := service.ListDevices(userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByUserID", userID).Return(([]models.Device)(nil), errors.New("database error"))

		result, err := service.ListDevices(userID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})
}

func TestDeviceService_UpdateDevice(t *testing.T) {
	userID := uuid.New()
	deviceID := uuid.New()
	device := &models.Device{
		UUID:        deviceID,
		Name:        "Old Name",
		Location:    "Old Location",
		Description: "Old Description",
		UserID:      userID,
	}

	t.Run("Success - Update device", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return(device, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.Device")).Return(nil)

		result, err := service.UpdateDevice(userID, deviceID, "New Name", "New Location", "New Description")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, "New Location", result.Location)
		assert.Equal(t, "New Description", result.Description)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Empty device name", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return(device, nil)

		result, err := service.UpdateDevice(userID, deviceID, "", "New Location", "New Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.NewValidationError("Device name is required"), err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Empty location", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return(device, nil)

		result, err := service.UpdateDevice(userID, deviceID, "New Name", "", "New Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.NewValidationError("Device location is required"), err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return((*models.Device)(nil), gorm.ErrRecordNotFound)

		result, err := service.UpdateDevice(userID, deviceID, "New Name", "New Location", "New Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDeviceNotFound, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on update", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return(device, nil)
		mockRepo.On("Update", mock.AnythingOfType("*models.Device")).Return(errors.New("database error"))

		result, err := service.UpdateDevice(userID, deviceID, "New Name", "New Location", "New Description")

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)
		assert.Nil(t, result)

		mockRepo.AssertExpectations(t)
	})
}

func TestDeviceService_DeleteDevice(t *testing.T) {
	userID := uuid.New()
	deviceID := uuid.New()
	device := &models.Device{
		UUID:   deviceID,
		Name:   "Test Device",
		UserID: userID,
	}

	t.Run("Success - Delete device", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return(device, nil)
		mockRepo.On("Delete", deviceID).Return(nil)

		err := service.DeleteDevice(userID, deviceID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Device not found", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return((*models.Device)(nil), gorm.ErrRecordNotFound)

		err := service.DeleteDevice(userID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDeviceNotFound, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - Database error on delete", func(t *testing.T) {
		mockRepo := new(MockDeviceRepository)
		service := NewDeviceService(mockRepo)

		mockRepo.On("FindByID", deviceID).Return(device, nil)
		mockRepo.On("Delete", deviceID).Return(errors.New("database error"))

		err := service.DeleteDevice(userID, deviceID)

		assert.Error(t, err)
		assert.Equal(t, custom_errors.ErrDatabaseError, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestIsValidSN(t *testing.T) {
	t.Run("Valid SN", func(t *testing.T) {
		assert.True(t, isValidSN("123456789012"))
	})

	t.Run("Too short", func(t *testing.T) {
		assert.False(t, isValidSN("123"))
	})

	t.Run("Too long", func(t *testing.T) {
		assert.False(t, isValidSN("1234567890123"))
	})

	t.Run("Non-digit characters", func(t *testing.T) {
		assert.False(t, isValidSN("12345678901a"))
	})

	t.Run("Empty string", func(t *testing.T) {
		assert.False(t, isValidSN(""))
	})
}