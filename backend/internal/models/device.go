package models

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	UUID        uuid.UUID `json:"uuid" db:"uuid"`
	Name        string    `json:"name" db:"name"`
	Location    string    `json:"location" db:"location"`
	SN          string    `json:"sn" db:"sn"`
	Description string    `json:"description" db:"description"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}