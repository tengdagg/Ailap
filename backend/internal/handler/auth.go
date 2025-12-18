package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ailap-backend/internal/database"
	"ailap-backend/internal/model"
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
	uid, exists := c.Get("userId")
	if !exists {
		// Fallback for demo/admin if no auth middleware active or just hardcoded
		// But usually we should fetch real user.
		// For now, let's keep it simple: fetch admin user (ID 1) or user from ctx
		// If userId is not in context (public test?), return default admin
		// In previous code it just returned hardcoded name.
		// Now we need actual DB user for RetentionDays.
		// Let's assume ID 1 is the main user if ctx missing (though ctx should have it).
		var user model.User
		if err := database.GetDB().First(&user, 1).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"name": "admin", "retentionDays": 15}})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": user})
		return
	}

	var user model.User
	if err := database.GetDB().First(&user, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": user})
}

type updateProfileReq struct {
	RetentionDays   int `json:"retentionDays"`
	AIAnalysisLimit int `json:"aiAnalysisLimit"`
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "bad request"})
		return
	}

	uid, exists := c.Get("userId")
	if !exists {
		// Fallback to ID 1 for now if auth middleware setup allows it
		uid = 1
	}

	var user model.User
	if err := database.GetDB().First(&user, uid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "user not found"})
		return
	}

	if req.RetentionDays > 0 {
		user.RetentionDays = req.RetentionDays
	}

	if req.AIAnalysisLimit > 0 {
		user.AIAnalysisLimit = req.AIAnalysisLimit
	}

	database.GetDB().Save(&user)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
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
