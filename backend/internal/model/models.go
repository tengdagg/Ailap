package model

import "time"

// MLModel represents an AI model configuration
// Provider examples: openai, deepseek, qwen
// Roles is a JSON string storing an array of role definitions
type MLModel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Provider    string    `json:"provider"`
	Model       string    `json:"model"`
	APIBase     string    `json:"apiBase"`
	APIKey      string    `json:"apiKey"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"maxTokens"`
	Roles       string    `gorm:"type:text" json:"roles"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	IsDefault   bool      `gorm:"default:false" json:"isDefault"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
