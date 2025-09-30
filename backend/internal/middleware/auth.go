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
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) { return secret, nil })
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid token"})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// extract subject (user id) and name into context for downstream handlers
			if sub, ok := claims["sub"]; ok {
				// JWT may encode numbers as float64
				switch v := sub.(type) {
				case float64:
					c.Set("userId", uint(v))
				case int:
					c.Set("userId", uint(v))
				case int64:
					c.Set("userId", uint(v))
				case string:
					// keep as string if cannot parse; downstream may handle
					c.Set("userId", v)
				}
			}
			if name, ok := claims["name"]; ok {
				c.Set("userName", name)
			}
		}
		c.Next()
	}
}
