package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"ailap-backend/internal/config"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		secret := []byte(config.Get().JWTSecret)
		_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) { return secret, nil })
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid token"})
			return
		}
		c.Next()
	}
}













