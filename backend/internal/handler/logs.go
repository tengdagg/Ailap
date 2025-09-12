package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogsHandler struct{}

func NewLogsHandler() *LogsHandler { return &LogsHandler{} }

func (h *LogsHandler) Query(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
}

func (h *LogsHandler) Suggestions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
}

