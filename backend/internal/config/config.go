package config

import (
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	HTTPPort     int
	JWTSecret    string
	DBDriver     string
	DBDSN        string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	AllowOrigins []string
}

var cfg AppConfig

func init() {
	viper.SetEnvPrefix("AILAP")
	viper.AutomaticEnv()

	viper.SetDefault("HTTP_PORT", 8080)
	viper.SetDefault("JWT_SECRET", "dev-secret-change-me")
	viper.SetDefault("DB_DRIVER", "")
	viper.SetDefault("DB_DSN", "")
	viper.SetDefault("READ_TIMEOUT", 120)
	viper.SetDefault("WRITE_TIMEOUT", 300)
	viper.SetDefault("ALLOW_ORIGINS", []string{"*"})

	cfg = AppConfig{
		HTTPPort:     viper.GetInt("HTTP_PORT"),
		JWTSecret:    viper.GetString("JWT_SECRET"),
		DBDriver:     viper.GetString("DB_DRIVER"),
		DBDSN:        viper.GetString("DB_DSN"),
		ReadTimeout:  time.Duration(viper.GetInt("READ_TIMEOUT")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("WRITE_TIMEOUT")) * time.Second,
		AllowOrigins: viper.GetStringSlice("ALLOW_ORIGINS"),
	}
}

func Get() AppConfig { return cfg }

