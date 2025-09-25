package services

import (
	"regexp"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/repository"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceService interface {
    CreateDevice(userID uuid.UUID, name, location, sn, description string) (*models.Device, error)
    GetDevice(userID, deviceID uuid.UUID) (*models.Device, error)
    ListDevices(userID uuid.UUID) ([]models.Device, error)
    UpdateDevice(userID, deviceID uuid.UUID, name, location, description string) (*models.Device, error)
    DeleteDevice(userID, deviceID uuid.UUID) error
}

type deviceService struct {
    deviceRepo repository.DeviceRepository
}

func NewDeviceService(deviceRepo repository.DeviceRepository) DeviceService {
    return &deviceService{deviceRepo: deviceRepo}
}

func (s *deviceService) CreateDevice(userID uuid.UUID, name, location, sn, description string) (*models.Device, error) {
    if name == "" {
        return nil, errors.NewValidationError("Device name is required")
    }
    if location == "" {
        return nil, errors.NewValidationError("Device location is required")
    }
    if sn == "" {
        return nil, errors.NewValidationError("Device serial number is required")
    }

    if !isValidSN(sn) {
        return nil, errors.NewValidationError("Serial number must be exactly 12 digits")
    }

    existing, err := s.deviceRepo.FindBySN(sn)
    if err != nil && err != gorm.ErrRecordNotFound {
        return nil, errors.ErrDatabaseError
    }
    if existing != nil {
        return nil, errors.ErrDeviceAlreadyExists
    }

    device := &models.Device{
        UUID:        uuid.New(),
        Name:        name,
        Location:    location,
        SN:          sn,
        Description: description,
        UserID:      userID,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if err := s.deviceRepo.Create(device); err != nil {
        return nil, errors.ErrDatabaseError
    }

    return device, nil
}

func (s *deviceService) GetDevice(userID, deviceID uuid.UUID) (*models.Device, error) {
    device, err := s.deviceRepo.FindByID(deviceID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.ErrDeviceNotFound
        }
        return nil, errors.ErrDatabaseError
    }

    if device.UserID != userID {
        return nil, errors.ErrForbidden
    }

    return device, nil
}

func (s *deviceService) ListDevices(userID uuid.UUID) ([]models.Device, error) {
    devices, err := s.deviceRepo.FindByUserID(userID)
    if err != nil {
        return nil, errors.ErrDatabaseError
    }
    return devices, nil
}

func (s *deviceService) UpdateDevice(userID, deviceID uuid.UUID, name, location, description string) (*models.Device, error) {
    device, err := s.GetDevice(userID, deviceID)
    if err != nil {
        return nil, err
    }

    if name == "" {
        return nil, errors.NewValidationError("Device name is required")
    }
    if location == "" {
        return nil, errors.NewValidationError("Device location is required")
    }

    device.Name = name
    device.Location = location
    device.Description = description
    device.UpdatedAt = time.Now()

    if err := s.deviceRepo.Update(device); err != nil {
        return nil, errors.ErrDatabaseError
    }

    return device, nil
}

func (s *deviceService) DeleteDevice(userID, deviceID uuid.UUID) error {
    _, err := s.GetDevice(userID, deviceID)
    if err != nil {
        return err
    }

    if err := s.deviceRepo.Delete(deviceID); err != nil {
        return errors.ErrDatabaseError
    }

    return nil
}

func isValidSN(sn string) bool {
    match, _ := regexp.MatchString(`^\d{12}$`, sn)
    return match
}