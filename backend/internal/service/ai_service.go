package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"ailap-backend/internal/database"
	"ailap-backend/internal/model"
	"ailap-backend/internal/utils"

	"go.uber.org/zap"
)

type AIService struct{}

func NewAIService() *AIService {
	return &AIService{}
}

// Analyze performs analysis on provided logs
func (s *AIService) Analyze(prompt string, logs []interface{}) (string, error) {
	// fetch default model
	var cfg model.MLModel
	if err := database.GetDB().Where("is_default = ? AND enabled = ?", true, true).First(&cfg).Error; err != nil {
		return "", fmt.Errorf("no enabled default model found")
	}
	if strings.TrimSpace(cfg.APIBase) == "" || strings.TrimSpace(cfg.APIKey) == "" || strings.TrimSpace(cfg.Model) == "" {
		return "", fmt.Errorf("incomplete model config")
	}

	// prepare logs snippet
	var buf bytes.Buffer
	limit := 8000
	for i, row := range logs {
		b, _ := json.Marshal(row)
		if buf.Len()+len(b)+1 > limit {
			fmt.Fprintf(&buf, "\n... (%d more) ...", len(logs)-i)
			break
		}
		if buf.Len() > 0 {
			buf.WriteByte('\n')
		}
		buf.Write(b)
	}

	sysPrompt := "你是资深的日志分析助手。根据提供的日志片段，结合用户问题，用中文给出要点式分析：1) 现象与范围，2) 可能原因，3) 进一步的验证建议，4) 缓解或修复步骤。"
	userPrompt := strings.TrimSpace(prompt)
	if userPrompt == "" {
		userPrompt = "请基于下列日志片段定位可能的问题并给出建议。"
	}
	userContent := fmt.Sprintf("%s\n\n日志片段(截断):\n%s", userPrompt, buf.String())

	endpoint := strings.TrimRight(cfg.APIBase, "/") + "/chat/completions"
	payload := map[string]interface{}{
		"model":       cfg.Model,
		"messages":    []map[string]string{{"role": "system", "content": sysPrompt}, {"role": "user", "content": userContent}},
		"max_tokens":  chooseInt(cfg.MaxTokens, 512),
		"temperature": chooseFloat(cfg.Temperature, 0.3),
		"stream":      false,
	}
	body, _ := json.Marshal(payload)

	reqHttp, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	reqHttp.Header.Set("Content-Type", "application/json")
	reqHttp.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(reqHttp)
	if err != nil {
		utils.GetLogger().Error("ai analyze provider error", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("provider error: %s", string(respBytes))
	}

	var obj struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBytes, &obj); err == nil && len(obj.Choices) > 0 {
		return obj.Choices[0].Message.Content, nil
	}
	return string(respBytes), nil
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
