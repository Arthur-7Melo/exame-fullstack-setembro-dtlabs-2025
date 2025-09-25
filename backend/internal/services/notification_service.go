package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/repository"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/errors"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/logger"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RedisPublisher interface {
	Publish(ctx context.Context, channel string, message interface{}) error
}

type NotificationService interface {
	CreateNotification(userID uuid.UUID, req dto.CreateNotificationRequest) (*models.Notification, error)
	GetUserNotifications(userID uuid.UUID) ([]models.Notification, error)
	CheckHeartbeat(heartbeat *models.Heartbeat) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	deviceRepo       repository.DeviceRepository
	redisClient      RedisPublisher // Usando interface em vez do tipo concreto
}

func NewNotificationService(notificationRepo repository.NotificationRepository, deviceRepo repository.DeviceRepository, redisClient RedisPublisher) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		deviceRepo:       deviceRepo,
		redisClient:      redisClient,
	}
}

func (s *notificationService) CreateNotification(userID uuid.UUID, req dto.CreateNotificationRequest) (*models.Notification, error) {
	if req.Name == "" {
		return nil, errors.NewValidationError("Notification name is required")
	}

	for _, condition := range req.Conditions {
		if !isValidParameter(condition.Parameter) {
			return nil, errors.NewValidationError("Invalid parameter: " + condition.Parameter)
		}
		if !isValidOperator(condition.Operator) {
			return nil, errors.NewValidationError("Invalid operator: " + condition.Operator)
		}
	}

	conditionsJSON, err := json.Marshal(req.Conditions)
	if err != nil {
		return nil, errors.NewValidationError("Invalid conditions format")
	}

	deviceIDsJSON, err := json.Marshal(req.DeviceIDs)
	if err != nil {
		return nil, errors.NewValidationError("Invalid device IDs format")
	}

	notification := &models.Notification{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Conditions:  datatypes.JSON(conditionsJSON),
		DeviceIDs:   datatypes.JSON(deviceIDsJSON),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return nil, errors.ErrDatabaseError
	}

	return notification, nil
}

func (s *notificationService) GetUserNotifications(userID uuid.UUID) ([]models.Notification, error) {
	notifications, err := s.notificationRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.ErrDatabaseError
	}
	return notifications, nil
}

func (s *notificationService) CheckHeartbeat(heartbeat *models.Heartbeat) error {
	device, err := s.deviceRepo.FindByID(heartbeat.DeviceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrDeviceNotFound
		}
		return errors.ErrDatabaseError
	}

	notifications, err := s.notificationRepo.FindActiveByUserID(device.UserID)
	if err != nil {
		return errors.ErrDatabaseError
	}

	for _, notification := range notifications {
		if !s.appliesToDevice(notification, device.UUID) {
			continue
		}

		if s.checkConditions(notification.Conditions, heartbeat) {
			if err := s.sendNotification(device.UserID, notification, device, heartbeat); err != nil {
				logger.Logger.Error("Error sending notification", "error", err)
			}
		}
	}

	return nil
}

func (s *notificationService) appliesToDevice(notification models.Notification, deviceID uuid.UUID) bool {
	var deviceIDs []uuid.UUID
	if err := json.Unmarshal(notification.DeviceIDs, &deviceIDs); err != nil {
		return false
	}

	if len(deviceIDs) == 0 {
		return true
	}

	for _, id := range deviceIDs {
		if id == deviceID {
			return true
		}
	}
	return false
}

func (s *notificationService) checkConditions(conditionsJSON datatypes.JSON, heartbeat *models.Heartbeat) bool {
	var conditions []dto.NotificationCondition
	if err := json.Unmarshal(conditionsJSON, &conditions); err != nil {
		return false
	}

	for _, condition := range conditions {
		if !s.checkCondition(condition, heartbeat) {
			return false
		}
	}
	return true
}

func (s *notificationService) checkCondition(condition dto.NotificationCondition, heartbeat *models.Heartbeat) bool {
	var value float64

	switch condition.Parameter {
	case "cpu":
		value = heartbeat.CPU
	case "ram":
		value = heartbeat.RAM
	case "disk_free":
		value = heartbeat.DiskFree
	case "temperature":
		value = heartbeat.Temperature
	case "latency":
		value = float64(heartbeat.Latency)
	case "connectivity":
		value = float64(heartbeat.Connectivity)
	default:
		return false
	}

	var conditionValue float64
	switch v := condition.Value.(type) {
	case float64:
		conditionValue = v
	case int:
		conditionValue = float64(v)
	case float32:
		conditionValue = float64(v)
	default:
		return false
	}

	switch condition.Operator {
	case ">":
		return value > conditionValue
	case "<":
		return value < conditionValue
	case ">=":
		return value >= conditionValue
	case "<=":
		return value <= conditionValue
	case "==":
		return value == conditionValue
	case "!=":
		return value != conditionValue
	default:
		return false
	}
}

func (s *notificationService) sendNotification(userID uuid.UUID, notification models.Notification, device *models.Device, heartbeat *models.Heartbeat) error {
	notificationMessage := map[string]interface{}{
		"id":              notification.ID.String(),
		"user_id":         userID.String(),
		"name":            notification.Name,
		"description":     notification.Description,
		"device_sn":       device.SN,
		"triggered_value": s.getTriggeredValue(notification.Conditions, heartbeat),
		"timestamp":       time.Now().Format(time.RFC3339),
		"heartbeat_data": map[string]interface{}{
			"cpu":          heartbeat.CPU,
			"ram":          heartbeat.RAM,
			"disk_free":    heartbeat.DiskFree,
			"temperature":  heartbeat.Temperature,
			"latency":      heartbeat.Latency,
			"connectivity": heartbeat.Connectivity,
		},
	}

	messageJSON, err := json.Marshal(notificationMessage)
	if err != nil {
		return errors.ErrDatabaseError
	}

	ctx := context.Background()
	if err := s.redisClient.Publish(ctx, "notifications:"+userID.String(), messageJSON); err != nil {
		return errors.ErrDatabaseError
	}

	logger.Logger.Info("Notification sent via Redis", 
		"user_id", userID.String(),
		"notification_id", notification.ID.String(),
		"device_sn", device.SN)

	return nil
}

func (s *notificationService) getTriggeredValue(conditionsJSON datatypes.JSON, heartbeat *models.Heartbeat) float64 {
	var conditions []dto.NotificationCondition
	if err := json.Unmarshal(conditionsJSON, &conditions); err != nil {
		return 0
	}

	for _, condition := range conditions {
		switch condition.Parameter {
		case "cpu":
			return heartbeat.CPU
		case "ram":
			return heartbeat.RAM
		case "disk_free":
			return heartbeat.DiskFree
		case "temperature":
			return heartbeat.Temperature
		case "latency":
			return float64(heartbeat.Latency)
		case "connectivity":
			return float64(heartbeat.Connectivity)
		}
	}
	return 0
}

func isValidParameter(parameter string) bool {
	validParameters := []string{"cpu", "ram", "disk_free", "temperature", "latency", "connectivity"}
	for _, p := range validParameters {
		if p == parameter {
			return true
		}
	}
	return false
}

func isValidOperator(operator string) bool {
	validOperators := []string{">", "<", ">=", "<=", "==", "!="}
	for _, op := range validOperators {
		if op == operator {
			return true
		}
	}
	return false
}