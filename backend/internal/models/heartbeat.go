package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Heartbeat struct {
    ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    DeviceID     uuid.UUID `json:"device_id" gorm:"type:uuid;not null"`
    CPU          float64   `json:"cpu" gorm:"not null"`
    RAM          float64   `json:"ram" gorm:"not null"`                 
    DiskFree     float64   `json:"disk_free" gorm:"not null"`              
    Temperature  float64   `json:"temperature" gorm:"not null"`           
    Latency      int       `json:"latency" gorm:"not null"`              
    Connectivity int       `json:"connectivity" gorm:"not null"`           
    BootTime     time.Time `json:"boot_time" gorm:"not null"`              
    CreatedAt    time.Time `json:"created_at" gorm:"not null"`
}

func (h *Heartbeat) BeforeCreate(tx *gorm.DB) error {
    h.ID = uuid.New()
    h.CreatedAt = time.Now().UTC()
    return nil
}