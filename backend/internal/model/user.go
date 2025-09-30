package model

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"uniqueIndex;size:64" json:"username"`
	Password string `json:"-"`
}














