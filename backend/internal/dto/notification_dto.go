package dto

import (
	"time"

	"github.com/google/uuid"
)

// @Description Request to create a notification rule
type CreateNotificationRequest struct {
	Name        string                 `json:"name" binding:"required" example:"High CPU Alert"`
	Description string                 `json:"description" example:"Alert when CPU usage is high"`
	Enabled     bool                   `json:"enabled" example:"true"`
	Conditions  []NotificationCondition `json:"conditions" binding:"required"`
	DeviceIDs   []uuid.UUID            `json:"device_ids" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// @Description Notification condition
type NotificationCondition struct {
	Parameter string      `json:"parameter" example:"cpu"`
	Operator  string      `json:"operator" example:">"`
	Value     interface{} `json:"value" example:"70.0"`
}

// @Description Response for notification rule
type NotificationResponse struct {
	ID          uuid.UUID               `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID      uuid.UUID               `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string                  `json:"name" example:"High CPU Alert"`
	Description string                  `json:"description" example:"Alert when CPU usage is high"`
	Enabled     bool                    `json:"enabled" example:"true"`
	Conditions  []NotificationCondition `json:"conditions"`
	DeviceIDs   []uuid.UUID             `json:"device_ids"`
	CreatedAt   time.Time               `json:"created_at" example:"2023-01-01T12:00:00Z"`
	UpdatedAt   time.Time               `json:"updated_at" example:"2023-01-01T12:00:00Z"`
}