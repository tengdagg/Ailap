package model

import "time"

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// LogQueryHistory stores executed query history
// Minimal model for demo purposes
type LogQueryHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Engine    string    `json:"engine"`
	Mode      string    `json:"mode"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"createdAt"`
}
