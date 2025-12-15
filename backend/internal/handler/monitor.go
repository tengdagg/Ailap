package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"ailap-backend/internal/database"
	"ailap-backend/internal/model"
	"ailap-backend/internal/service"
)

type MonitorHandler struct {
	svc *service.MonitorService
}

func NewMonitorHandler(svc *service.MonitorService) *MonitorHandler {
	return &MonitorHandler{svc: svc}
}

// ---- Monitors ----

func (h *MonitorHandler) ListMonitors(c *gin.Context) {
	var items []model.LogMonitor
	if err := database.GetDB().Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

func (h *MonitorHandler) GetMonitor(c *gin.Context) {
	id := c.Param("id")
	var item model.LogMonitor
	if err := database.GetDB().First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"item": item}})
}

func (h *MonitorHandler) CreateMonitor(c *gin.Context) {
	var req model.LogMonitor
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}
	if err := database.GetDB().Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}
	// Start job if active
	if err := h.svc.AddJob(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "created but failed to start job: " + err.Error(), "data": gin.H{"item": req}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"item": req}})
}

func (h *MonitorHandler) UpdateMonitor(c *gin.Context) {
	id := c.Param("id")
	var item model.LogMonitor
	if err := database.GetDB().First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "not found"})
		return
	}

	var req model.LogMonitor
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	item.Name = req.Name
	item.DatasourceID = req.DatasourceID
	item.Engine = req.Engine
	item.Cron = req.Cron
	item.Query = req.Query
	item.Keywords = req.Keywords
	item.ChannelID = req.ChannelID
	item.Status = req.Status

	if err := database.GetDB().Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}

	// Update job
	if err := h.svc.AddJob(&item); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "updated but failed to restart job: " + err.Error(), "data": gin.H{"item": item}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"item": item}})
}

func (h *MonitorHandler) DeleteMonitor(c *gin.Context) {
	id := c.Param("id")
	uid, _ := strconv.Atoi(id)
	h.svc.RemoveJob(uint(uid))

	if err := database.GetDB().Delete(&model.LogMonitor{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func (h *MonitorHandler) ListAlerts(c *gin.Context) {
	var items []model.AlertHistory
	// Get latest 50 alerts, preload Monitor info
	if err := database.GetDB().Preload("Monitor").Order("created_at desc").Limit(50).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

func (h *MonitorHandler) GetAlert(c *gin.Context) {
	id := c.Param("id")
	var item model.AlertHistory
	if err := database.GetDB().Preload("Monitor").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"item": item}})
}

// ---- Channels ----

func (h *MonitorHandler) ListChannels(c *gin.Context) {
	var items []model.NotificationChannel
	if err := database.GetDB().Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

func (h *MonitorHandler) GetChannel(c *gin.Context) {
	id := c.Param("id")
	var item model.NotificationChannel
	if err := database.GetDB().First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"item": item}})
}

func (h *MonitorHandler) CreateChannel(c *gin.Context) {
	var req model.NotificationChannel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}
	if err := database.GetDB().Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"item": req}})
}

func (h *MonitorHandler) UpdateChannel(c *gin.Context) {
	id := c.Param("id")
	var item model.NotificationChannel
	if err := database.GetDB().First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "not found"})
		return
	}

	var req model.NotificationChannel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	item.Name = req.Name
	item.Type = req.Type
	item.Config = req.Config

	if err := database.GetDB().Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"item": item}})
}

func (h *MonitorHandler) DeleteChannel(c *gin.Context) {
	id := c.Param("id")
	// Check if in use? For now just delete, FK constraints might fail if any.
	// SQLite gorm basic setup usually soft delete or no constraints unless enforced.

	if err := database.GetDB().Delete(&model.NotificationChannel{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func (h *MonitorHandler) TestChannel(c *gin.Context) {
	var req model.NotificationChannel
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": err.Error()})
		return
	}

	svc := service.NewNotificationService()
	testTitle := "AILAP Alert Test (测试告警)"
	testContent := "**Status**: OK\n\nThis is a *test notification* to verify your **Markdown** configuration.\n\n- Channel: " + req.Name + "\n- Time: " + time.Now().Format("15:04:05")

	if err := svc.SendAlert(&req, testTitle, testContent); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "测试发送失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
