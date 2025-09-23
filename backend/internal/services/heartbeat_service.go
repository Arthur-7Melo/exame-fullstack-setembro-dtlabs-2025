// internal/services/heartbeat_service.go
package services

import (
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/repository"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HeartbeatService interface {
    CreateHeartbeat(deviceID uuid.UUID, cpu, ram, diskFree, temperature float64, latency, connectivity int, bootTime time.Time) (*models.Heartbeat, error)
    GetDeviceHeartbeats(userID, deviceID uuid.UUID, startTime, endTime time.Time) ([]models.Heartbeat, error)
    GetLatestDeviceHeartbeat(userID, deviceID uuid.UUID) (*models.Heartbeat, error)
}

type heartbeatService struct {
    heartbeatRepo repository.HeartbeatRepository
    deviceRepo    repository.DeviceRepository
}

func NewHeartbeatService(heartbeatRepo repository.HeartbeatRepository, deviceRepo repository.DeviceRepository) HeartbeatService {
    return &heartbeatService{
        heartbeatRepo: heartbeatRepo,
        deviceRepo:    deviceRepo,
    }
}

func (s *heartbeatService) CreateHeartbeat(deviceID uuid.UUID, cpu, ram, diskFree, temperature float64, latency, connectivity int, bootTime time.Time) (*models.Heartbeat, error) {
    heartbeat := &models.Heartbeat{
        DeviceID:     deviceID,
        CPU:          cpu,
        RAM:          ram,
        DiskFree:     diskFree,
        Temperature:  temperature,
        Latency:      latency,
        Connectivity: connectivity,
        BootTime:     bootTime,
    }

    if err := s.heartbeatRepo.Create(heartbeat); err != nil {
        return nil, errors.ErrDatabaseError
    }

    return heartbeat, nil
}

func (s *heartbeatService) GetDeviceHeartbeats(userID, deviceID uuid.UUID, startTime, endTime time.Time) ([]models.Heartbeat, error) {
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

    heartbeats, err := s.heartbeatRepo.FindByDeviceID(deviceID, startTime, endTime)
    if err != nil {
        return nil, errors.ErrDatabaseError
    }

    return heartbeats, nil
}

func (s *heartbeatService) GetLatestDeviceHeartbeat(userID, deviceID uuid.UUID) (*models.Heartbeat, error) {
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

    heartbeat, err := s.heartbeatRepo.FindLatestByDeviceID(deviceID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.ErrDeviceNotFound
        }
        return nil, errors.ErrDatabaseError
    }

    return heartbeat, nil
}