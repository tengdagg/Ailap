package model

import (
	"time"
)

// NotificationChannel defines where to send alerts
type NotificationChannel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`   // webhook, email
	Config    string    `json:"config"` // JSON string: {url, token} or {smtp...}
}

// LogMonitor defines a scheduled task to check logs
type LogMonitor struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	Name         string     `json:"name"`
	DatasourceID string     `json:"datasourceId"` // e.g., "1" or "vl_1"
	Engine       string     `json:"engine"`       // loki, elasticsearch, victorialogs
	Cron         string     `json:"cron"`         // e.g., "@every 1h" or "0 * * * *"
	Query        string     `json:"query"`        // base query
	Keywords     string     `json:"keywords"`     // comma separated keywords to filter content
	ChannelID    uint       `json:"channelId"`
	Status       string     `json:"status"` // active, paused
	LastRunAt    *time.Time `json:"lastRunAt"`
	ProjectID    uint       `json:"projectId"` // optional, for multi-tenancy if needed
}
