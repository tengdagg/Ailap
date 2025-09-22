package utils

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger     *zap.Logger
	loggerOnce sync.Once
)

// GetLogger returns a singleton zap logger
func GetLogger() *zap.Logger {
	loggerOnce.Do(func() {
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		l, err := cfg.Build()
		if err != nil {
			panic(err)
		}
		logger = l
	})
	return logger
}











