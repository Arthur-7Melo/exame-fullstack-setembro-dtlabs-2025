package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDeviceRequest struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	SN          string `json:"sn"`
	Description string `json:"description"`
}

type UpdateDeviceRequest struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type DeviceResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	SN          string    `json:"sn"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}