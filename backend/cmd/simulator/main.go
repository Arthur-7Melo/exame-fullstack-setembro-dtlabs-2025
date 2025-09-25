package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/dto"
	"github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025/internal/utils/logger"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
    amqpURL := os.Getenv("AMQP_URL")
    if amqpURL == "" {
        amqpURL = "amqp://guest:guest@rabbitmq:5672/" 
    }
    
    conn, err := amqp091.Dial(amqpURL)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()
    
    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()
    
    q, err := ch.QueueDeclare(
        "heartbeats",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        logger.Logger.Error("Failed to declare a queue", "error", err.Error())
        os.Exit(1)
    }
    
    deviceIDs := os.Getenv("DEVICE_IDS")
    if deviceIDs == "" {
        logger.Logger.Warn("DEVICE_IDS environment variable is not set. Simulator will start but won't send any heartbeats.")
        deviceIDs = ""
    }
    
    var deviceUUIDs []uuid.UUID
    if deviceIDs != "" {
        for _, id := range strings.Split(deviceIDs, ",") {
            deviceID, err := uuid.Parse(strings.TrimSpace(id))
            if err != nil {
                logger.Logger.Error("Invalid device ID", "id", id)
                os.Exit(1)
            }
            deviceUUIDs = append(deviceUUIDs, deviceID)
        }
    }
    
    logger.Logger.Info("Starting heartbeat simulator", "device_count", len(deviceUUIDs))

    if len(deviceUUIDs) == 0 {
        logger.Logger.Warn("No devices configured. Simulator is running but won't send heartbeats.")
    }
    
    for {
        if len(deviceUUIDs) == 0 {
            logger.Logger.Info("Waiting for devices to be configured...")
            time.Sleep(1 * time.Minute)
            continue
        }
        
        for _, deviceID := range deviceUUIDs {
            msg := dto.HeartbeatMessage{
                DeviceID:     deviceID.String(),
                CPU:          rand.Float64() * 100,
                RAM:          rand.Float64() * 100,
                DiskFree:     rand.Float64() * 100,
                Temperature:  20 + rand.Float64()*60,
                Latency:      rand.Intn(500),
                Connectivity: rand.Intn(2),
                BootTime:     time.Now().UTC().Add(-time.Duration(rand.Intn(86400)) * time.Second),
            }
            
            body, err := json.Marshal(msg)
            if err != nil {
                logger.Logger.Error("Error marshaling message", "error", err)
                continue
            }
            
            err = ch.Publish(
                "",
                q.Name,
                false,
                false,
                amqp091.Publishing{
                    ContentType: "application/json",
                    Body:        body,
                })
            if err != nil {
                logger.Logger.Error("Failed to publish a message", "error", err)
            } else {
                logger.Logger.Debug("Published heartbeat for device", "device", deviceID)
            }
        }
        
        logger.Logger.Info("Completed heartbeat cycle", "device_count", len(deviceUUIDs))
        time.Sleep(1 * time.Minute)
    }
}