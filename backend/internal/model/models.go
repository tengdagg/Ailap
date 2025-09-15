package model

type MLModel struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}




