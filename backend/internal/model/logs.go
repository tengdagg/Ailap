package model

import "time"

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// LogQueryHistory stores executed query history
type LogQueryHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Engine     string    `json:"engine"`
	Mode       string    `json:"mode"`
	Query      string    `json:"query"`
	LineLimit  int       `json:"lineLimit"`
	Note       string    `json:"note"`
	IsFavorite bool      `gorm:"default:false" json:"isFavorite"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
