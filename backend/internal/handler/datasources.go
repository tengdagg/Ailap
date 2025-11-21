package handler

import (
	"crypto/tls"
	"crypto/x509"
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
	name := stringOr(raw["name"])
	endpoint := stringOr(raw["endpoint"])
	if name == "" || endpoint == "" {
		c.JSON(400, gin.H{"code": 400, "message": "name and endpoint are required"})
		return
	}

	cfgBytes, _ := json.Marshal(raw)
	d := model.DataSource{
		Name:     name,
		Type:     stringOr(raw["type"]),
		Endpoint: endpoint,
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
	name := stringOr(raw["name"])
	endpoint := stringOr(raw["endpoint"])
	if name == "" || endpoint == "" {
		c.JSON(400, gin.H{"code": 400, "message": "name and endpoint are required"})
		return
	}

	cfgBytes, _ := json.Marshal(raw)
	updates := map[string]interface{}{
		"name":     name,
		"type":     stringOr(raw["type"]),
		"endpoint": endpoint,
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
	// allow empty body for id-based testing from list
	_ = c.ShouldBindJSON(&raw)
	if raw == nil {
		raw = map[string]interface{}{}
	}
	// when called as /datasources/:id/test and critical fields missing, load from DB
	if id := c.Param("id"); id != "" {
		if stringOr(raw["type"]) == "" || stringOr(raw["endpoint"]) == "" {
			var d model.DataSource
			if err := database.GetDB().First(&d, "id = ?", id).Error; err == nil {
				var cfg map[string]interface{}
				_ = json.Unmarshal([]byte(d.Config), &cfg)
				if cfg == nil {
					cfg = map[string]interface{}{}
				}
				cfg["name"] = d.Name
				cfg["type"] = d.Type
				cfg["endpoint"] = d.Endpoint
				raw = cfg
			}
		}
	}
	typ := stringOr(raw["type"])
	endpoint := stringOr(raw["endpoint"])
	client := createTestHTTPClient(raw, 5*time.Second)
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

// createTestHTTPClient creates an HTTP client with TLS configuration for testing
func createTestHTTPClient(cfg map[string]interface{}, timeout time.Duration) *http.Client {
	transport := &http.Transport{}

	// TLS configuration
	if tlsCfg := getTestTLSConfig(cfg); tlsCfg != nil {
		transport.TLSClientConfig = tlsCfg
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

// getTestTLSConfig extracts TLS configuration from datasource config for testing
func getTestTLSConfig(cfg map[string]interface{}) *tls.Config {
	if cfg == nil {
		return nil
	}

	tlsData, ok := cfg["tls"].(map[string]interface{})
	if !ok {
		return nil
	}

	tlsConfig := &tls.Config{}

	// Skip certificate verification
	if skipVerify, ok := tlsData["skipVerify"].(bool); ok && skipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	// Server name for TLS
	if serverName, ok := tlsData["serverName"].(string); ok && serverName != "" {
		tlsConfig.ServerName = serverName
	}

	// CA certificate for self-signed certificates
	if caCert, ok := tlsData["caCert"].(string); ok && caCert != "" {
		caCertPool := x509.NewCertPool()
		if caCertPool.AppendCertsFromPEM([]byte(caCert)) {
			tlsConfig.RootCAs = caCertPool
		}
	}

	// Client certificate authentication
	if clientCert, ok := tlsData["clientCert"].(string); ok && clientCert != "" {
		if clientKey, ok := tlsData["clientKey"].(string); ok && clientKey != "" {
			cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
			if err == nil {
				tlsConfig.Certificates = []tls.Certificate{cert}
			}
		}
	}

	return tlsConfig
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
