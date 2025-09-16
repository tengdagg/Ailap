package handler

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"ailap-backend/internal/database"
	"ailap-backend/internal/model"
)

type LogsHandler struct{}

func NewLogsHandler() *LogsHandler { return &LogsHandler{} }

// Query proxies logs for different engines (loki, elasticsearch)
func (h *LogsHandler) Query(c *gin.Context) {
	engine := c.Query("engine")
	if engine != "loki" && engine != "elasticsearch" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
		return
	}

	mode := c.Query("mode")
	query := c.Query("query")
	lineLimit := c.DefaultQuery("lineLimit", "1000")
	typ := strings.ToLower(c.DefaultQuery("type", "range"))
	qStart := c.Query("start")
	qEnd := c.Query("end")
	qStep := c.DefaultQuery("step", "60s")
	qDirection := c.DefaultQuery("direction", "BACKWARD")

	if engine == "loki" {
		// Build query from builder when needed
		if mode == "builder" && query == "" {
			// Manual parsing of nested URL parameters
			filters := make([]struct {
				Label, Op string
				Values    []string
			}, 0)

			// Parse builder[labelFilters][0][label], builder[labelFilters][0][op], etc.
			for i := 0; i < 10; i++ { // support up to 10 filters
				labelKey := fmt.Sprintf("builder[labelFilters][%d][label]", i)
				opKey := fmt.Sprintf("builder[labelFilters][%d][op]", i)
				label := c.Query(labelKey)
				op := c.Query(opKey)
				if label == "" && op == "" {
					break
				}

				var values []string
				for j := 0; j < 10; j++ { // support up to 10 values per filter
					valueKey := fmt.Sprintf("builder[labelFilters][%d][values][%d]", i, j)
					value := c.Query(valueKey)
					if value == "" {
						break
					}
					values = append(values, value)
				}

				if label != "" && len(values) > 0 {
					filters = append(filters, struct {
						Label, Op string
						Values    []string
					}{Label: label, Op: op, Values: values})
				}
			}

			contains := c.Query("builder[contains]")
			query = buildLokiQuery(filters, contains)
		}

		_, cfg, endpoint, ok := resolveLokiDatasource(c.Query("datasourceId"))
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no loki datasource", "data": gin.H{"items": []interface{}{}}})
			return
		}

		params := url.Values{}
		params.Set("query", query)
		if typ == "range" {
			if qStart == "" || qEnd == "" {
				now := time.Now()
				start := now.Add(-1 * time.Hour).UnixNano()
				end := now.UnixNano()
				params.Set("start", fmt.Sprintf("%d", start))
				params.Set("end", fmt.Sprintf("%d", end))
			} else {
				params.Set("start", qStart)
				params.Set("end", qEnd)
			}
			params.Set("limit", lineLimit)
			params.Set("direction", qDirection)
			params.Set("step", qStep)
		} else {
			params.Set("limit", lineLimit)
		}

		reqURL := endpoint
		if parsed, err := url.Parse(reqURL); err == nil && (parsed.Path == "" || parsed.Path == "/") {
			if typ == "range" {
				reqURL = reqURL + "/loki/api/v1/query_range?" + params.Encode()
			} else {
				reqURL = reqURL + "/loki/api/v1/query?" + params.Encode()
			}
		} else {
			if !strings.Contains(reqURL, "?") {
				reqURL = reqURL + "?" + params.Encode()
			}
		}

		req, _ := http.NewRequest(http.MethodGet, reqURL, nil)
		applyAuthHeaders(req, cfg)
		client := createHTTPClient(cfg, 15*time.Second)
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": err.Error(), "data": gin.H{"items": []interface{}{}}})
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		// Save query history regardless of result
		_ = database.GetDB().Create(&model.LogQueryHistory{Engine: "loki", Mode: mode, Query: query}).Error

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			c.JSON(http.StatusOK, gin.H{"code": 1, "message": string(body), "data": gin.H{"items": []interface{}{}}})
			return
		}

		items := flattenLokiToRows(body)

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
		return
	}

	// Elasticsearch branch
	_, cfg, endpoint, ok := resolveElasticsearchDatasource(c.Query("datasourceId"))
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no elasticsearch datasource", "data": gin.H{"items": []interface{}{}}})
		return
	}
	// resolve fields
	timeField := "@timestamp"
	if es, ok := cfg["es"].(map[string]interface{}); ok {
		if tf, ok := es["timeField"].(string); ok && tf != "" {
			timeField = tf
		}
	}
	messageField := "_source"
	if logsCfg, ok := cfg["logs"].(map[string]interface{}); ok {
		if mf, ok := logsCfg["messageField"].(string); ok && mf != "" {
			messageField = mf
		}
	}
	levelField := ""
	if logsCfg, ok := cfg["logs"].(map[string]interface{}); ok {
		if lf, ok := logsCfg["levelField"].(string); ok && lf != "" {
			levelField = lf
		}
	}
	indexPath := ""
	if es, ok := cfg["es"].(map[string]interface{}); ok {
		if idx, ok := es["index"].(string); ok && idx != "" {
			indexPath = "/" + idx
		}
	}

	// time ns -> ms
	nsToMs := func(s string, def int64) int64 {
		if s == "" {
			return def
		}
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return v / 1e6
		}
		return def
	}
	nowMs := time.Now().UnixMilli()
	startMs := nsToMs(qStart, nowMs-3600*1000)
	endMs := nsToMs(qEnd, nowMs)

	if query == "" {
		query = "*"
	}
	size := 1000
	if v, err := strconv.Atoi(lineLimit); err == nil {
		size = v
	}

	// Build must conditions dynamically
	mustConditions := []interface{}{}
	if query == "" || query == "*" {
		mustConditions = append(mustConditions, map[string]interface{}{"match_all": map[string]interface{}{}})
	} else {
		mustConditions = append(mustConditions, map[string]interface{}{"query_string": map[string]interface{}{"query": query}})
	}
	mustConditions = append(mustConditions, map[string]interface{}{"range": map[string]interface{}{timeField: map[string]interface{}{"gte": startMs, "lte": endMs, "format": "epoch_millis"}}})

	bodyJSON := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
		"sort": []interface{}{map[string]interface{}{timeField: map[string]interface{}{"order": "desc"}}},
		"size": size,
	}

	// Add maxShardRequests if configured
	// Removed: max_concurrent_shard_requests is not a valid field in _search body for Elasticsearch
	// Keep bodyJSON as-is
	payload, _ := json.Marshal(bodyJSON)
	searchURL := endpoint + indexPath + "/_search"
	req, _ := http.NewRequest(http.MethodPost, searchURL, strings.NewReader(string(payload)))
	req.Header.Set("Content-Type", "application/json")

	// Add X-Pack specific headers if enabled
	if es, ok := cfg["es"].(map[string]interface{}); ok {
		if xpack, ok := es["xpack"].(bool); ok && xpack {
			// X-Pack is enabled, add specific headers if needed
			req.Header.Set("X-Elastic-Product", "Elasticsearch")
		}
	}

	applyAuthHeaders(req, cfg)
	client := createHTTPClient(cfg, 15*time.Second)
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": err.Error(), "data": gin.H{"items": []interface{}{}}})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	// Save query history regardless of result
	_ = database.GetDB().Create(&model.LogQueryHistory{Engine: "elasticsearch", Mode: mode, Query: query}).Error

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": string(body), "data": gin.H{"items": []interface{}{}}})
		return
	}

	items := flattenElasticsearchToRows(body, timeField, messageField, levelField)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

// Suggestions returns label names for Loki when engine=loki
func (h *LogsHandler) Suggestions(c *gin.Context) {
	engine := c.Query("engine")
	if engine != "loki" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
		return
	}
	_, cfg, endpoint, ok := resolveLokiDatasource(c.Query("datasourceId"))
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no loki datasource", "data": gin.H{"items": []interface{}{}}})
		return
	}
	u := endpoint
	if parsed, err := url.Parse(u); err == nil && (parsed.Path == "" || parsed.Path == "/") {
		u = u + "/loki/api/v1/labels"
	}
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	applyAuthHeaders(req, cfg)
	client := createHTTPClient(cfg, 5*time.Second)
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
		return
	}
	var obj map[string]interface{}
	_ = json.Unmarshal(body, &obj)
	var items []interface{}
	if d, ok := obj["data"].([]interface{}); ok {
		items = d
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

// LabelValues returns values for a specific Loki label
// GET /api/logs/label-values?engine=loki&label=service_name[&datasourceId=1]
func (h *LogsHandler) LabelValues(c *gin.Context) {
	if c.Query("engine") != "loki" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
		return
	}
	label := c.Query("label")
	if label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "label is required"})
		return
	}
	_, cfg, endpoint, ok := resolveLokiDatasource(c.Query("datasourceId"))
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no loki datasource", "data": gin.H{"items": []interface{}{}}})
		return
	}
	u := endpoint
	if parsed, err := url.Parse(u); err == nil && (parsed.Path == "" || parsed.Path == "/") {
		u = u + "/loki/api/v1/label/" + url.PathEscape(label) + "/values"
	}
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	applyAuthHeaders(req, cfg)
	client := createHTTPClient(cfg, 5*time.Second)
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
		return
	}
	var obj map[string]interface{}
	_ = json.Unmarshal(body, &obj)
	var items []interface{}
	if d, ok := obj["data"].([]interface{}); ok {
		items = d
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

// History returns query history with auto cleanup
func (h *LogsHandler) History(c *gin.Context) {
	// Auto cleanup: remove non-favorite queries older than 14 days
	cutoff := time.Now().AddDate(0, 0, -14)
	database.GetDB().Where("is_favorite = ? AND created_at < ?", false, cutoff).Delete(&model.LogQueryHistory{})

	queryType := c.DefaultQuery("type", "recent") // recent or favorite
	var items []model.LogQueryHistory

	if queryType == "favorite" {
		database.GetDB().Where("is_favorite = ?", true).Order("updated_at desc").Find(&items)
	} else {
		database.GetDB().Order("id desc").Limit(50).Find(&items)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": items}})
}

// ToggleFavorite toggles the favorite status of a query
func (h *LogsHandler) ToggleFavorite(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "missing id"})
		return
	}

	var item model.LogQueryHistory
	if err := database.GetDB().First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "query not found"})
		return
	}

	item.IsFavorite = !item.IsFavorite
	item.UpdatedAt = time.Now()
	database.GetDB().Save(&item)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"item": item}})
}

// UpdateNote updates the note of a query
func (h *LogsHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "missing id"})
		return
	}

	var req struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "invalid request"})
		return
	}

	var item model.LogQueryHistory
	if err := database.GetDB().First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 1, "message": "query not found"})
		return
	}

	item.Note = req.Note
	item.UpdatedAt = time.Now()
	database.GetDB().Save(&item)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"item": item}})
}

// DeleteHistory deletes a query history item
func (h *LogsHandler) DeleteHistory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "missing id"})
		return
	}

	if err := database.GetDB().Delete(&model.LogQueryHistory{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "message": "failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

// Inspect: loki -> URL; elasticsearch -> {url, body}
func (h *LogsHandler) Inspect(c *gin.Context) {
	engine := c.Query("engine")
	if engine == "loki" {
		mode := c.Query("mode")
		query := c.Query("query")
		if mode == "builder" && query == "" {
			// Manual parsing of nested URL parameters
			filters := make([]struct {
				Label, Op string
				Values    []string
			}, 0)

			for i := 0; i < 10; i++ {
				label := c.Query(fmt.Sprintf("builder[labelFilters][%d][label]", i))
				op := c.Query(fmt.Sprintf("builder[labelFilters][%d][op]", i))
				if label == "" && op == "" {
					break
				}

				var values []string
				for j := 0; j < 10; j++ {
					value := c.Query(fmt.Sprintf("builder[labelFilters][%d][values][%d]", i, j))
					if value == "" {
						break
					}
					values = append(values, value)
				}

				if label != "" && len(values) > 0 {
					filters = append(filters, struct {
						Label, Op string
						Values    []string
					}{Label: label, Op: op, Values: values})
				}
			}

			contains := c.Query("builder[contains]")
			query = buildLokiQuery(filters, contains)
		}
		_, _, endpoint, ok := resolveLokiDatasource(c.Query("datasourceId"))
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no loki datasource", "data": gin.H{"url": ""}})
			return
		}
		params := url.Values{}
		params.Set("query", query)
		for _, k := range []string{"start", "end", "step", "direction"} {
			if v := c.Query(k); v != "" {
				params.Set(k, v)
			}
		}
		urlStr := endpoint + "/loki/api/v1/query_range?" + params.Encode()
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"url": urlStr}})
		return
	}
	// ES inspect
	_, cfg, endpoint, ok := resolveElasticsearchDatasource(c.Query("datasourceId"))
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no elasticsearch datasource", "data": gin.H{"url": "", "body": ""}})
		return
	}
	indexPath := ""
	if es, ok := cfg["es"].(map[string]interface{}); ok {
		if idx, ok := es["index"].(string); ok && idx != "" {
			indexPath = "/" + idx
		}
	}
	urlStr := endpoint + indexPath + "/_search"
	// body 基于 query 与时间窗
	timeField := "@timestamp"
	if es, ok := cfg["es"].(map[string]interface{}); ok {
		if tf, ok := es["timeField"].(string); ok && tf != "" {
			timeField = tf
		}
	}
	nsToMs := func(s string, def int64) int64 {
		if s == "" {
			return def
		}
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return v / 1e6
		}
		return def
	}
	nowMs := time.Now().UnixMilli()
	startMs := nsToMs(c.Query("start"), nowMs-3600*1000)
	endMs := nsToMs(c.Query("end"), nowMs)
	q := c.Query("query")
	if q == "" {
		q = "*"
	}
	bodyJSON := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": func() []interface{} {
					conds := make([]interface{}, 0, 2)
					if q == "" || q == "*" {
						conds = append(conds, map[string]interface{}{"match_all": map[string]interface{}{}})
					} else {
						conds = append(conds, map[string]interface{}{"query_string": map[string]interface{}{"query": q}})
					}
					conds = append(conds, map[string]interface{}{"range": map[string]interface{}{timeField: map[string]interface{}{"gte": startMs, "lte": endMs, "format": "epoch_millis"}}})
					return conds
				}(),
			},
		},
	}

	// Add maxShardRequests if configured
	// Removed: max_concurrent_shard_requests is not a valid field in _search body
	b, _ := json.MarshalIndent(bodyJSON, "", "  ")
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"url": urlStr, "body": string(b)}})
}

// helpers
func resolveLokiDatasource(id string) (ds model.DataSource, cfg map[string]interface{}, endpoint string, ok bool) {
	if id != "" {
		if err := database.GetDB().First(&ds, "id = ?", id).Error; err == nil {
			_ = json.Unmarshal([]byte(ds.Config), &cfg)
			endpoint = ds.Endpoint
			if endpoint == "" && cfg != nil {
				if v, o := cfg["endpoint"].(string); o {
					endpoint = v
				}
			}
			return ds, cfg, endpoint, true
		}
	}
	var items []model.DataSource
	database.GetDB().Where("type = ?", "loki").Find(&items)
	if len(items) == 0 {
		return ds, nil, "", false
	}
	picked := -1
	for i := range items {
		var c map[string]interface{}
		_ = json.Unmarshal([]byte(items[i].Config), &c)
		if c != nil {
			if v, ok := c["isDefault"].(bool); ok && v {
				picked = i
				break
			}
		}
	}
	if picked == -1 {
		picked = 0
	}
	ds = items[picked]
	_ = json.Unmarshal([]byte(ds.Config), &cfg)
	endpoint = ds.Endpoint
	if endpoint == "" && cfg != nil {
		if v, o := cfg["endpoint"].(string); o {
			endpoint = v
		}
	}
	return ds, cfg, endpoint, true
}

func resolveElasticsearchDatasource(id string) (ds model.DataSource, cfg map[string]interface{}, endpoint string, ok bool) {
	if id != "" {
		if err := database.GetDB().First(&ds, "id = ?", id).Error; err == nil {
			_ = json.Unmarshal([]byte(ds.Config), &cfg)
			endpoint = ds.Endpoint
			if endpoint == "" && cfg != nil {
				if v, o := cfg["endpoint"].(string); o {
					endpoint = v
				}
			}
			return ds, cfg, endpoint, true
		}
	}
	var items []model.DataSource
	database.GetDB().Where("type = ?", "elasticsearch").Find(&items)
	if len(items) == 0 {
		return ds, nil, "", false
	}
	picked := -1
	for i := range items {
		var c map[string]interface{}
		_ = json.Unmarshal([]byte(items[i].Config), &c)
		if c != nil {
			if v, ok := c["isDefault"].(bool); ok && v {
				picked = i
				break
			}
		}
	}
	if picked == -1 {
		picked = 0
	}
	ds = items[picked]
	_ = json.Unmarshal([]byte(ds.Config), &cfg)
	endpoint = ds.Endpoint
	if endpoint == "" && cfg != nil {
		if v, o := cfg["endpoint"].(string); o {
			endpoint = v
		}
	}
	return ds, cfg, endpoint, true
}

// createHTTPClient creates an HTTP client with TLS configuration
func createHTTPClient(cfg map[string]interface{}, timeout time.Duration) *http.Client {
	transport := &http.Transport{}

	// TLS configuration
	if tlsCfg := getTLSConfig(cfg); tlsCfg != nil {
		transport.TLSClientConfig = tlsCfg
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

// getTLSConfig extracts TLS configuration from datasource config
func getTLSConfig(cfg map[string]interface{}) *tls.Config {
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

func applyAuthHeaders(req *http.Request, cfg map[string]interface{}) {
	if cfg == nil {
		return
	}
	if token, ok := cfg["token"].(string); ok && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if apiKey, ok := cfg["apiKey"].(string); ok && apiKey != "" {
		req.Header.Set("Authorization", "ApiKey "+apiKey)
	}
	if username, ok := cfg["username"].(string); ok && username != "" {
		if password, ok := cfg["password"].(string); ok && password != "" {
			req.SetBasicAuth(username, password)
		}
	}

	// Apply custom headers and cookies
	if httpCfg, ok := cfg["http"].(map[string]interface{}); ok {
		if cookies, ok := httpCfg["allowedCookies"].([]interface{}); ok {
			for _, cookie := range cookies {
				if cookieStr, ok := cookie.(string); ok && cookieStr != "" {
					req.Header.Add("Cookie", cookieStr)
				}
			}
		}
	}
}

func buildLokiQuery(filters []struct {
	Label, Op string
	Values    []string
}, contains string) string {
	var sb strings.Builder
	sb.WriteString("{")
	for i, f := range filters {
		if f.Label == "" || len(f.Values) == 0 {
			continue
		}
		if i > 0 {
			sb.WriteString(",")
		}
		op := f.Op
		if op == "" {
			op = "="
		}
		v := f.Values[0]
		sb.WriteString(f.Label)
		sb.WriteString(op)
		sb.WriteString("\"")
		sb.WriteString(strings.ReplaceAll(v, "\"", "\\\""))
		sb.WriteString("\"")
	}
	sb.WriteString("}")
	if contains != "" {
		sb.WriteString(" |~ \"")
		sb.WriteString(strings.ReplaceAll(contains, "\"", "\\\""))
		sb.WriteString("\"")
	}
	return sb.String()
}

func flattenLokiToRows(respBody []byte) []map[string]interface{} {
	var obj map[string]interface{}
	_ = json.Unmarshal(respBody, &obj)
	rows := make([]map[string]interface{}, 0)
	data, ok := obj["data"].(map[string]interface{})
	if !ok {
		return rows
	}
	results, ok := data["result"].([]interface{})
	if !ok {
		results, ok = data["streams"].([]interface{})
		if !ok {
			return rows
		}
	}
	for _, r := range results {
		m, _ := r.(map[string]interface{})
		if values, ok := m["values"].([]interface{}); ok {
			for _, vv := range values {
				pair, _ := vv.([]interface{})
				if len(pair) >= 2 {
					rows = append(rows, map[string]interface{}{"timestamp": pair[0], "level": "", "message": pair[1]})
				}
			}
			continue
		}
		if entries, ok := m["entries"].([]interface{}); ok {
			for _, e := range entries {
				em, _ := e.(map[string]interface{})
				rows = append(rows, map[string]interface{}{"timestamp": em["ts"], "level": "", "message": em["line"]})
			}
		}
	}
	return rows
}

// getNestedField extracts nested field values from a map using dot notation
func getNestedField(data map[string]interface{}, fieldPath string) interface{} {
	if fieldPath == "" {
		return nil
	}

	// Split field path by dots
	parts := strings.Split(fieldPath, ".")
	current := data

	for i, part := range parts {
		if val, exists := current[part]; exists {
			// If this is the last part, return the value
			if i == len(parts)-1 {
				return val
			}
			// If not the last part, continue traversing
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				// Can't traverse further
				return nil
			}
		} else {
			// Field doesn't exist
			return nil
		}
	}

	return nil
}

func flattenElasticsearchToRows(respBody []byte, timeField, messageField, levelField string) []map[string]interface{} {
	rows := make([]map[string]interface{}, 0)

	var obj map[string]interface{}
	_ = json.Unmarshal(respBody, &obj)

	hitsWrap, ok := obj["hits"].(map[string]interface{})
	if !ok {
		return rows
	}
	hits, ok := hitsWrap["hits"].([]interface{})
	if !ok {
		return rows
	}

	for _, h := range hits {
		hm, _ := h.(map[string]interface{})
		source, _ := hm["_source"].(map[string]interface{})

		// timestamp
		var timestamp interface{}
		if timeField != "" {
			if val, exists := source[timeField]; exists {
				timestamp = val
			} else {
				timestamp = getNestedField(source, timeField)
			}
		}

		// message
		var message interface{}
		if messageField == "_source" {
			if b, err := json.Marshal(source); err == nil {
				message = string(b)
			} else {
				message = fmt.Sprintf("%v", source)
			}
		} else if messageField != "" {
			if val, exists := source[messageField]; exists {
				message = val
			} else {
				message = getNestedField(source, messageField)
				if message == nil {
					for _, f := range []string{"message", "log.message", "msg", "text"} {
						if val := getNestedField(source, f); val != nil {
							message = val
							break
						}
					}
					if message == nil {
						if b, err := json.Marshal(source); err == nil {
							message = string(b)
						} else {
							message = fmt.Sprintf("%v", source)
						}
					}
				}
			}
		}

		// level
		var level interface{} = ""
		if levelField != "" {
			if val := getNestedField(source, levelField); val != nil {
				level = val
			} else {
				for _, f := range []string{"level", "log.level", "severity", "priority"} {
					if val := getNestedField(source, f); val != nil {
						level = val
						break
					}
				}
			}
		} else {
			for _, f := range []string{"level", "log.level", "severity", "priority"} {
				if val := getNestedField(source, f); val != nil {
					level = val
					break
				}
			}
		}

		// Create raw data with all ES metadata and source fields
		rawData := make(map[string]interface{})

		// Add ES metadata fields
		if id, exists := hm["_id"]; exists {
			rawData["_id"] = id
		}
		if index, exists := hm["_index"]; exists {
			rawData["_index"] = index
		}
		if docType, exists := hm["_type"]; exists {
			rawData["_type"] = docType
		}
		if score, exists := hm["_score"]; exists {
			rawData["_score"] = score
		}
		if sort, exists := hm["sort"]; exists {
			rawData["sort"] = sort
		}
		if highlight, exists := hm["highlight"]; exists {
			rawData["highlight"] = highlight
		}

		// Add all source fields with flattened keys
		flattenMap(source, "", rawData)

		rows = append(rows, map[string]interface{}{
			"timestamp": timestamp,
			"level":     level,
			"message":   message,
			"__raw":     rawData, // Now includes both metadata and flattened source
		})
	}

	return rows
}

// flattenMap recursively flattens nested maps with dot notation keys
func flattenMap(data map[string]interface{}, prefix string, result map[string]interface{}) {
	for key, value := range data {
		var newKey string
		if prefix == "" {
			newKey = key
		} else {
			newKey = prefix + "." + key
		}

		if nestedMap, ok := value.(map[string]interface{}); ok {
			// Recursively flatten nested maps
			flattenMap(nestedMap, newKey, result)
		} else {
			// Add the flattened key-value pair
			result[newKey] = value
		}
	}
}
