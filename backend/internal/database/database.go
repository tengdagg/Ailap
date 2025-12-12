package database

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ailap-backend/internal/config"
	"ailap-backend/internal/model"
)

var db *gorm.DB

func GetDB() *gorm.DB { return db }

func Init() error {
	cfg := config.Get()
	dsn := cfg.DBDSN
	if dsn == "" {
		_ = os.MkdirAll("data", 0o755)
		dsn = "data/ailap.db"
	}
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	db = gdb

	if err := db.AutoMigrate(&model.User{}, &model.MLModel{}, &model.DataSource{}, &model.LogQueryHistory{}, &model.LogMonitor{}, &model.NotificationChannel{}); err != nil {
		return err
	}

	var cnt int64
	db.Model(&model.User{}).Count(&cnt)
	if cnt == 0 {
		adminUser := os.Getenv("AILAP_ADMIN_USER")
		if adminUser == "" {
			adminUser = "admin"
		}
		adminPass := os.Getenv("AILAP_ADMIN_PASS")
		if adminPass == "" {
			adminPass = "admin123"
		}
		pwdHash, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
		if err := db.Create(&model.User{Username: adminUser, Password: string(pwdHash)}).Error; err != nil {
			return fmt.Errorf("seed admin failed: %w", err)
		}
	}
	return nil
}
