package repository

import (
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HeartbeatRepository interface {
    Create(heartbeat *models.Heartbeat) error
    FindByDeviceID(deviceID uuid.UUID, startTime, endTime time.Time) ([]models.Heartbeat, error)
    FindLatestByDeviceID(deviceID uuid.UUID) (*models.Heartbeat, error)
}

type heartbeatRepository struct {
    db *gorm.DB
}

func NewHeartbeatRepository(db *gorm.DB) HeartbeatRepository {
    return &heartbeatRepository{db: db}
}

func (r *heartbeatRepository) Create(heartbeat *models.Heartbeat) error {
    return r.db.Create(heartbeat).Error
}

func (r *heartbeatRepository) FindByDeviceID(deviceID uuid.UUID, startTime, endTime time.Time) ([]models.Heartbeat, error) {
    var heartbeats []models.Heartbeat
    err := r.db.Where("device_id = ? AND created_at BETWEEN ? AND ?", deviceID, startTime, endTime).
        Order("created_at DESC").
        Find(&heartbeats).Error
    return heartbeats, err
}

func (r *heartbeatRepository) FindLatestByDeviceID(deviceID uuid.UUID) (*models.Heartbeat, error) {
    var heartbeat models.Heartbeat
    err := r.db.Where("device_id = ?", deviceID).
        Order("created_at DESC").
        First(&heartbeat).Error
    if err != nil {
        return nil, err
    }
    return &heartbeat, nil
}