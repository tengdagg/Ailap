package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"ailap-backend/internal/database"
	"ailap-backend/internal/model"
	"ailap-backend/internal/utils"

	"go.uber.org/zap"
)

type AIHandler struct{}

func NewAIHandler() *AIHandler { return &AIHandler{} }

type analyzeLogsReq struct {
	Prompt string        `json:"prompt"`
	Logs   []interface{} `json:"logs"` // free-form rows
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnalyzeLogs uses the default enabled model to analyze provided logs with a user prompt
func (h *AIHandler) AnalyzeLogs(c *gin.Context) {
	var req analyzeLogsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "bad request"})
		return
	}

	// fetch default model
	var cfg model.MLModel
	if err := database.GetDB().Where("is_default = ? AND enabled = ?", true, true).First(&cfg).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "未找到已启用的默认模型"})
		return
	}
	if strings.TrimSpace(cfg.APIBase) == "" || strings.TrimSpace(cfg.APIKey) == "" || strings.TrimSpace(cfg.Model) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "默认模型配置不完整"})
		return
	}

	// prepare logs snippet (limit length)
	var buf bytes.Buffer
	limit := 8000 // characters limit to avoid hitting provider limits
	for i, row := range req.Logs {
		// compact json for each row
		b, _ := json.Marshal(row)
		if buf.Len()+len(b)+1 > limit {
			fmt.Fprintf(&buf, "\n... (%d more) ...", len(req.Logs)-i)
			break
		}
		if buf.Len() > 0 {
			buf.WriteByte('\n')
		}
		buf.Write(b)
	}

	sysPrompt := "你是资深的日志分析助手。根据提供的日志片段，结合用户问题，用中文给出要点式分析：1) 现象与范围，2) 可能原因，3) 进一步的验证建议，4) 缓解或修复步骤。若上下文不足，请指出需要的关键信息。"
	userPrompt := strings.TrimSpace(req.Prompt)
	if userPrompt == "" {
		userPrompt = "请基于下列日志片段定位可能的问题并给出建议。"
	}
	userContent := fmt.Sprintf("%s\n\n日志片段(截断):\n%s", userPrompt, buf.String())

	endpoint := strings.TrimRight(cfg.APIBase, "/") + "/chat/completions"
	payload := map[string]interface{}{
		"model":       cfg.Model,
		"messages":    []openAIMessage{{Role: "system", Content: sysPrompt}, {Role: "user", Content: userContent}},
		"max_tokens":  chooseInt(cfg.MaxTokens, 512),
		"temperature": chooseFloat(cfg.Temperature, 0.3),
		"stream":      false,
	}
	body, _ := json.Marshal(payload)

	reqHttp, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		utils.GetLogger().Error("ai analyze new request", zap.Error(err), zap.String("endpoint", endpoint))
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": err.Error()})
		return
	}
	reqHttp.Header.Set("Content-Type", "application/json")
	reqHttp.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(reqHttp)
	if err != nil {
		utils.GetLogger().Error("ai analyze provider error", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := string(respBytes)
		if len(msg) > 1024 {
			msg = msg[:1024]
		}
		utils.GetLogger().Error("ai analyze non-2xx", zap.Int("status", resp.StatusCode), zap.String("msg", msg))
		// 规范返回 200，避免前端拦截器将 401/403 作为退出
		c.JSON(http.StatusOK, gin.H{"code": resp.StatusCode, "message": fmt.Sprintf("provider error: %s", msg)})
		return
	}

	// parse openai-compatible response
	var obj struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBytes, &obj); err == nil && len(obj.Choices) > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"reply": obj.Choices[0].Message.Content}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"reply": string(respBytes)}})
}

func chooseInt(v int, def int) int {
	if v > 0 {
		return v
	}
	return def
}

func chooseFloat(v float64, def float64) float64 {
	if v > 0 {
		return v
	}
	return def
}
