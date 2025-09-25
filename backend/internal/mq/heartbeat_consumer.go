// internal/mq/heartbeat_consumer.go
package mq

import (
	"encoding/json"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/models"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/repository"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	logger "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/logger"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type HeartbeatConsumer struct {
	heartbeatService   services.HeartbeatService
	notificationService services.NotificationService
	deviceRepo         repository.DeviceRepository
	conn               *amqp091.Connection
	channel            *amqp091.Channel
	queueName          string
}

func NewHeartbeatConsumer(amqpURL, queueName string, heartbeatService services.HeartbeatService, notificationService services.NotificationService, deviceRepo repository.DeviceRepository) (*HeartbeatConsumer, error) {
	logger.Logger.Info("Connecting to RabbitMQ", "url", amqpURL)

	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		logger.Logger.Error("Failed to connect to RabbitMQ", "error", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Logger.Error("Failed to open channel", "error", err)
		return nil, err
	}

	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		logger.Logger.Error("Failed to declare queue", "error", err)
		return nil, err
	}

	return &HeartbeatConsumer{
		heartbeatService:   heartbeatService,
		notificationService: notificationService,
		deviceRepo:         deviceRepo,
		conn:               conn,
		channel:            ch,
		queueName:          queueName,
	}, nil
}

func (c *HeartbeatConsumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		logger.Logger.Error("Failed to start consuming", "error", err)
		return err
	}

	go func() {
		for d := range msgs {
			var msg dto.HeartbeatMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				logger.Logger.Error("Error decoding heartbeat", "error", err)
				continue
			}

			deviceID, err := uuid.Parse(msg.DeviceID)
			if err != nil {
				logger.Logger.Error("Invalid device ID", "device_id", msg.DeviceID, "error", err)
				continue
			}

			device, err := c.deviceRepo.FindByID(deviceID)
			if err != nil {
				logger.Logger.Error("Error getting device", "device_id", msg.DeviceID, "error", err)
				continue
			}

			_, err = c.heartbeatService.CreateHeartbeat(
				deviceID,
				msg.CPU,
				msg.RAM,
				msg.DiskFree,
				msg.Temperature,
				msg.Latency,
				msg.Connectivity,
				msg.BootTime,
			)
			if err != nil {
				logger.Logger.Error("Error saving heartbeat", "error", err)
				continue
			}

			heartbeat := &models.Heartbeat{
				DeviceID:     deviceID,
				CPU:          msg.CPU,
				RAM:          msg.RAM,
				DiskFree:     msg.DiskFree,
				Temperature:  msg.Temperature,
				Latency:      msg.Latency,
				Connectivity: msg.Connectivity,
				BootTime:     msg.BootTime,
				CreatedAt:    time.Now(),
			}

			if err := c.notificationService.CheckHeartbeat(heartbeat); err != nil {
				logger.Logger.Error("Error checking notifications", "error", err)
			}

			logger.Logger.Info("Processed heartbeat", "device_id", deviceID, "device_sn", device.SN)
		}
	}()

	return nil
}

func (c *HeartbeatConsumer) Close() {
	c.channel.Close()
	c.conn.Close()
}