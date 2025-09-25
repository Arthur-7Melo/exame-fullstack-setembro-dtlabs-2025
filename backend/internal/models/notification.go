package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Notification struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID      `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Enabled     bool           `json:"enabled"`
	Conditions  datatypes.JSON `json:"conditions" gorm:"type:jsonb"`
	DeviceIDs   datatypes.JSON `json:"device_ids" gorm:"type:jsonb"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}