package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ailap-backend/internal/service"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler() *AuthHandler { return &AuthHandler{svc: service.NewAuthService()} }

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "bad request"})
		return
	}
	token, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"token": token}})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func (h *AuthHandler) Profile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"name": "admin"}})
}

type changePasswordReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// ChangePassword allows authenticated user to change own password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req changePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil || req.OldPassword == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "bad request"})
		return
	}
	// user id from context (set by AuthRequired)
	uid, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "unauthorized"})
		return
	}
	if err := h.svc.ChangePassword(uid, req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
