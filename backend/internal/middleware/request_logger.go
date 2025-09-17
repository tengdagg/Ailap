package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"ailap-backend/internal/utils"
)

func RequestLogger() gin.HandlerFunc {
	logger := utils.GetLogger().Sugar()
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errs := c.Errors.ByType(gin.ErrorTypeAny)
		if len(errs) > 0 {
			for _, e := range errs {
				logger.Errorw("request error", zap.String("method", method), zap.String("path", path), zap.String("query", query), zap.Int("status", status), zap.Duration("latency", latency), zap.String("ip", clientIP), zap.Error(e))
			}
		} else {
			logger.Infow("request", zap.String("method", method), zap.String("path", path), zap.String("query", query), zap.Int("status", status), zap.Duration("latency", latency), zap.String("ip", clientIP))
		}
	}
}








