package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ailap-backend/internal/service"
)

type AIHandler struct {
	aiService *service.AIService
}

func NewAIHandler() *AIHandler {
	return &AIHandler{
		aiService: service.NewAIService(),
	}
}

type analyzeLogsReq struct {
	Prompt string        `json:"prompt"`
	Logs   []interface{} `json:"logs"`
}

// AnalyzeLogs uses the default enabled model to analyze provided logs with a user prompt
func (h *AIHandler) AnalyzeLogs(c *gin.Context) {
	var req analyzeLogsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "bad request"})
		return
	}

	reply, err := h.aiService.Analyze(req.Prompt, req.Logs)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"reply": reply}})
}
