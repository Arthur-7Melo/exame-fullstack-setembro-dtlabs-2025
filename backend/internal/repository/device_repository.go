package repository

import (
	"errors"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	FindBySN(sn string) (*models.Device, error)
	Create(device *models.Device) error
	FindByID(id uuid.UUID) (*models.Device, error)
	FindByUserID(userID uuid.UUID) ([]models.Device, error)
	Update(device *models.Device) error
	Delete(id uuid.UUID) error
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) FindBySN(sn string) (*models.Device, error) {
	var device models.Device
	err := r.db.Where("sn = ?", sn).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) Create(device *models.Device) error {
	return r.db.Create(device).Error
}

func (r *deviceRepository) FindByID(id uuid.UUID) (*models.Device, error) {
	var device models.Device
	err := r.db.First(&device, "uuid = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &device, nil
}

func (r *deviceRepository) FindByUserID(userID uuid.UUID) ([]models.Device, error) {
	var devices []models.Device
	err := r.db.Where("user_id = ?", userID).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *deviceRepository) Update(device *models.Device) error {
	result := r.db.Model(&models.Device{}).Where("uuid = ?", device.UUID).Updates(device)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *deviceRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("uuid = ?", id).Delete(&models.Device{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
