package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"ailap-backend/internal/config"
	"ailap-backend/internal/database"
	"ailap-backend/internal/model"
)

type AuthService struct{ db *gorm.DB }

func NewAuthService() *AuthService { return &AuthService{db: database.GetDB()} }

func (s *AuthService) Login(username, password string) (string, error) {
	var u model.User
	if err := s.db.Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	cfg := config.Get()
	claims := jwt.MapClaims{
		"sub":  u.ID,
		"name": u.Username,
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(cfg.JWTSecret))
}


