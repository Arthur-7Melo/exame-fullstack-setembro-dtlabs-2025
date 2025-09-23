package dto

import "time"

// @Description Heartbeat message containing device telemetry data
type HeartbeatMessage struct {
    DeviceID     string    `json:"device_id" example:"f8aac3d7-b2ac-43e8-94d5-bbee59dd4ac9"`
    CPU          float64   `json:"cpu" example:"45.67"`                    // CPU usage in %
    RAM          float64   `json:"ram" example:"67.89"`                    // RAM usage in %
    DiskFree     float64   `json:"disk_free" example:"23.45"`              // Free disk space in %
    Temperature  float64   `json:"temperature" example:"35.67"`            // Temperature in Celsius
    Latency      int       `json:"latency" example:"150"`                  // Latency to DNS 8.8.8.8 in milliseconds
    Connectivity int       `json:"connectivity" example:"1"`               // 0 (no connection) or 1 (has connection)
    BootTime     time.Time `json:"boot_time" example:"2023-01-01T00:00:00Z"` // Boot timestamp with UTC+00
}

// @Description Heartbeat response with telemetry data and metadata
type HeartbeatResponse struct {
    ID           string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
    DeviceID     string    `json:"device_id" example:"f8aac3d7-b2ac-43e8-94d5-bbee59dd4ac9"`
    CPU          float64   `json:"cpu" example:"45.67"`
    RAM          float64   `json:"ram" example:"67.89"`
    DiskFree     float64   `json:"disk_free" example:"23.45"`
    Temperature  float64   `json:"temperature" example:"35.67"`
    Latency      int       `json:"latency" example:"150"`
    Connectivity int       `json:"connectivity" example:"1"`
    BootTime     time.Time `json:"boot_time" example:"2023-01-01T00:00:00Z"`
    CreatedAt    time.Time `json:"created_at" example:"2023-01-01T12:00:00Z"`
}

// @Description Request parameters for retrieving device heartbeats
type DeviceHeartbeatsRequest struct {
    StartTime time.Time `form:"start" example:"2023-01-01T00:00:00Z"`  // Start time for filtering
    EndTime   time.Time `form:"end" example:"2023-01-02T00:00:00Z"`    // End time for filtering
}