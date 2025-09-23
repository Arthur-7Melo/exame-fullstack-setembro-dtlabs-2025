package mq

import (
	"encoding/json"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/services"
	logger "github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/logger"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type HeartbeatConsumer struct {
    heartbeatService services.HeartbeatService
    conn             *amqp091.Connection 
    channel          *amqp091.Channel    
    queueName        string
}

func NewHeartbeatConsumer(amqpURL, queueName string, heartbeatService services.HeartbeatService) (*HeartbeatConsumer, error) {
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

    logger.Logger.Info("Declaring queue", "queue", queueName)
    
    _, err = ch.QueueDeclare(
        queueName,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        logger.Logger.Error("Failed to declare queue", "error", err)
        return nil, err
    }

    logger.Logger.Info("RabbitMQ consumer initialized successfully")
    
    return &HeartbeatConsumer{
        heartbeatService: heartbeatService,
        conn:             conn,
        channel:          ch,
        queueName:        queueName,
    }, nil
}

func (c *HeartbeatConsumer) Start() error {
    logger.Logger.Info("Starting to consume messages", "queue", c.queueName)
    
    msgs, err := c.channel.Consume(
        c.queueName,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        logger.Logger.Error("Failed to start consuming", "error", err)
        return err
    }

    logger.Logger.Info("Successfully started consuming messages")

    go func() {
        for d := range msgs {
            logger.Logger.Debug("Received message", "body_length", len(d.Body))
            
            var msg dto.HeartbeatMessage // Using DTO now
            if err := json.Unmarshal(d.Body, &msg); err != nil {
                logger.Logger.Error("Error decoding message", "error", err)
                continue
            }

            logger.Logger.Debug("Parsed heartbeat message", 
                "device_id", msg.DeviceID,
                "cpu", msg.CPU)

            deviceID, err := uuid.Parse(msg.DeviceID)
            if err != nil {
                logger.Logger.Error("Invalid device ID", "device_id", msg.DeviceID, "error", err)
                continue
            }

            heartbeat, err := c.heartbeatService.CreateHeartbeat(
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
                logger.Logger.Error("Error saving heartbeat", 
                    "device_id", deviceID, 
                    "error", err)
            } else {
                logger.Logger.Info("Heartbeat saved successfully", 
                    "device_id", deviceID,
                    "heartbeat_id", heartbeat.ID)
            }
        }
    }()

    return nil
}

func (c *HeartbeatConsumer) Close() {
    logger.Logger.Info("Closing RabbitMQ connection")
    c.channel.Close()
    c.conn.Close()
}