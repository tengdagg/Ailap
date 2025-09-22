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
)

type ModelsHandler struct{}

func NewModelsHandler() *ModelsHandler { return &ModelsHandler{} }

func (h *ModelsHandler) List(c *gin.Context) {
	var items []model.MLModel
	database.GetDB().Find(&items)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

func (h *ModelsHandler) Create(c *gin.Context) {
	var m model.MLModel
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "bad request"})
		return
	}
	// ensure single default
	if m.IsDefault {
		database.GetDB().Model(&model.MLModel{}).Where("is_default = ?", true).Update("is_default", false)
	}
	database.GetDB().Create(&m)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"id": m.ID}})
}

func (h *ModelsHandler) Update(c *gin.Context) {
	var m model.MLModel
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "bad request"})
		return
	}
	id := c.Param("id")
	if m.IsDefault {
		database.GetDB().Model(&model.MLModel{}).Where("is_default = ?", true).Update("is_default", false)
	}
	database.GetDB().Model(&model.MLModel{}).Where("id = ?", id).Updates(m)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func (h *ModelsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	database.GetDB().Delete(&model.MLModel{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// ToggleEnabled updates enabled status
func (h *ModelsHandler) ToggleEnabled(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "bad request"})
		return
	}
	database.GetDB().Model(&model.MLModel{}).Where("id = ?", id).Update("enabled", body.Enabled)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// SetDefault sets a model as default (ensures single default)
func (h *ModelsHandler) SetDefault(c *gin.Context) {
	id := c.Param("id")
	// reset others
	database.GetDB().Model(&model.MLModel{}).Where("is_default = ?", true).Update("is_default", false)
	// set this one
	database.GetDB().Model(&model.MLModel{}).Where("id = ?", id).Update("is_default", true)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// Test calls the configured model provider with a simple prompt for connectivity check
func (h *ModelsHandler) Test(c *gin.Context) {
	var cfg model.MLModel
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "bad request"})
		return
	}

	if strings.TrimSpace(cfg.APIBase) == "" || strings.TrimSpace(cfg.APIKey) == "" || strings.TrimSpace(cfg.Model) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请填写 API Base、API Key 和模型"})
		return
	}

	endpoint := strings.TrimRight(cfg.APIBase, "/") + "/chat/completions"
	payload := map[string]interface{}{
		"model":       cfg.Model,
		"messages":    []map[string]string{{"role": "user", "content": "ping"}},
		"max_tokens":  1,
		"temperature": 0,
		"stream":      false,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "message": err.Error()})
		return
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Optionally ensure response looks like OpenAI-compatible format
		var obj map[string]interface{}
		if err := json.Unmarshal(respBytes, &obj); err == nil {
			if _, ok := obj["choices"]; ok {
				c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
		return
	}

	msg := string(respBytes)
	if len(msg) > 1024 {
		msg = msg[:1024]
	}
	// Normalize provider errors to 200 to avoid FE interceptor logging the user out on 401/403
	c.JSON(http.StatusOK, gin.H{"code": resp.StatusCode, "message": fmt.Sprintf("provider error: %s", msg)})
}
