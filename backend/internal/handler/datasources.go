package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"ailap-backend/internal/database"
	"ailap-backend/internal/model"
	"ailap-backend/internal/utils"
)

// hasPath reports whether the provided endpoint already contains a non-root path
func hasPath(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	return u.Path != "" && u.Path != "/"
}

type DataSourcesHandler struct{}

func NewDataSourcesHandler() *DataSourcesHandler { return &DataSourcesHandler{} }

func (h *DataSourcesHandler) List(c *gin.Context) {
	var items []model.DataSource
	if err := database.GetDB().Find(&items).Error; err != nil {
		utils.GetLogger().Error("list datasources", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	utils.GetLogger().Info("list datasources", zap.Int("count", len(items)))
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

func (h *DataSourcesHandler) Create(c *gin.Context) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "bad request"})
		return
	}
	cfgBytes, _ := json.Marshal(raw)
	d := model.DataSource{
		Name:     stringOr(raw["name"]),
		Type:     stringOr(raw["type"]),
		Endpoint: stringOr(raw["endpoint"]),
		Config:   string(cfgBytes),
	}
	if err := database.GetDB().Create(&d).Error; err != nil {
		utils.GetLogger().Error("create datasource", zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	utils.GetLogger().Info("create datasource", zap.Uint("id", d.ID), zap.String("type", d.Type))
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"id": d.ID}})
}

func (h *DataSourcesHandler) Update(c *gin.Context) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "bad request"})
		return
	}
	cfgBytes, _ := json.Marshal(raw)
	updates := map[string]interface{}{
		"name":     stringOr(raw["name"]),
		"type":     stringOr(raw["type"]),
		"endpoint": stringOr(raw["endpoint"]),
		"config":   string(cfgBytes),
	}
	id := c.Param("id")
	if err := database.GetDB().Model(&model.DataSource{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		utils.GetLogger().Error("update datasource", zap.String("id", id), zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	utils.GetLogger().Info("update datasource", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func (h *DataSourcesHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := database.GetDB().Delete(&model.DataSource{}, "id = ?", id).Error; err != nil {
		utils.GetLogger().Error("delete datasource", zap.String("id", id), zap.Error(err))
		c.JSON(500, gin.H{"code": 500, "message": err.Error()})
		return
	}
	utils.GetLogger().Info("delete datasource", zap.String("id", id))
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func (h *DataSourcesHandler) Test(c *gin.Context) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "bad request"})
		return
	}
	typ := stringOr(raw["type"])
	endpoint := stringOr(raw["endpoint"])
	client := &http.Client{Timeout: 5 * time.Second}
	reqURL := endpoint
	if typ == "loki" {
		if !hasPath(endpoint) {
			reqURL = endpoint + "/loki/api/v1/labels?limit=1"
		}
	} else if typ == "elasticsearch" {
		if !hasPath(endpoint) {
			reqURL = endpoint + "/_cluster/health"
		}
	}
	req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
	if token := stringOr(raw["token"]); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if apiKey := stringOr(raw["apiKey"]); apiKey != "" {
		req.Header.Set("Authorization", "ApiKey "+apiKey)
	}
	if username := stringOr(raw["username"]); username != "" {
		if password := stringOr(raw["password"]); password != "" {
			req.SetBasicAuth(username, password)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		utils.GetLogger().Error("test datasource", zap.String("type", typ), zap.String("url", reqURL), zap.Error(err))
		c.JSON(200, gin.H{"code": 1, "message": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	utils.GetLogger().Info("test datasource", zap.String("type", typ), zap.String("url", reqURL), zap.Int("status", resp.StatusCode))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		c.JSON(200, gin.H{"code": 0, "message": "ok", "data": gin.H{"status": resp.Status, "body": string(body)}})
		return
	}
	c.JSON(200, gin.H{"code": 1, "message": resp.Status, "data": gin.H{"body": string(body)}})
}

func stringOr(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
