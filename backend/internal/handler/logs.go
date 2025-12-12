package handler

import (
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
	"ailap-backend/internal/service"
)

type LogsHandler struct {
	logService *service.LogService
}

func NewLogsHandler() *LogsHandler {
	return &LogsHandler{
		logService: service.NewLogService(),
	}
}

// Query proxies logs for different engines (loki, elasticsearch)
func (h *LogsHandler) Query(c *gin.Context) {
	engine := c.Query("engine")
	if engine != "loki" && engine != "elasticsearch" && engine != "victorialogs" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
		return
	}

	mode := c.Query("mode")
	queryParam := c.Query("query")
	lineLimit := c.DefaultQuery("lineLimit", "1000")
	typ := strings.ToLower(c.DefaultQuery("type", "range"))
	qStart := c.Query("start")
	qEnd := c.Query("end")
	// qStep := c.DefaultQuery("step", "60s")
	// qDirection := c.DefaultQuery("direction", "BACKWARD")
	datasourceID := c.Query("datasourceId")

	finalQuery := queryParam

	if engine == "loki" {
		// Build query from builder when needed
		if mode == "builder" && queryParam == "" {
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
			finalQuery = buildLokiQuery(filters, contains)
		}
	} else if engine == "elasticsearch" {
		// ES logic handled in service, optional query construction here if needed
	}

	// Logic for time range standardization
	limit, _ := strconv.Atoi(lineLimit)

	// For loki range query param prep
	start := qStart
	end := qEnd

	if engine == "loki" && typ == "range" {
		if qStart == "" || qEnd == "" {
			now := time.Now()
			start = fmt.Sprintf("%d", now.Add(-1*time.Hour).UnixNano())
			end = fmt.Sprintf("%d", now.UnixNano())
		}
	}

	// Execute via Service
	result, err := h.logService.ExecuteQuery(c.Request.Context(), engine, datasourceID, finalQuery, start, end, limit)

	// Save history
	_ = database.GetDB().Create(&model.LogQueryHistory{
		Engine:    engine,
		Mode:      mode,
		Query:     finalQuery,
		LineLimit: limit,
	}).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": err.Error(), "data": gin.H{"items": []interface{}{}}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": result.Items}})
}

// Suggestions returns label names for Loki when engine=loki
func (h *LogsHandler) Suggestions(c *gin.Context) {
	engine := c.Query("engine")
	if engine != "loki" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"items": []interface{}{}}})
		return
	}
	_, cfg, endpoint, ok := service.ResolveLokiDatasource(c.Query("datasourceId"))
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no loki datasource", "data": gin.H{"items": []interface{}{}}})
		return
	}
	// Re-implemented logic or just keep as is for metadata functions not yet in service properly
	u := endpoint
	if parsed, err := url.Parse(u); err == nil && (parsed.Path == "" || parsed.Path == "/") {
		u = u + "/loki/api/v1/labels"
	}
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	service.ApplyAuthHeaders(req, cfg)
	client := service.CreateHTTPClient(cfg, 5*time.Second)
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
	_, cfg, endpoint, ok := service.ResolveLokiDatasource(c.Query("datasourceId"))
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no loki datasource", "data": gin.H{"items": []interface{}{}}})
		return
	}
	u := endpoint
	if parsed, err := url.Parse(u); err == nil && (parsed.Path == "" || parsed.Path == "/") {
		u = u + "/loki/api/v1/label/" + url.PathEscape(label) + "/values"
	}
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	service.ApplyAuthHeaders(req, cfg)
	client := service.CreateHTTPClient(cfg, 5*time.Second)
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
	// Keep Inspect logic here or move to service if deep inspection reuse needed.
	// For now, simpler to keep handler logic as it constructs potential requests.
	// Or reuse service's resolve functions.

	engine := c.Query("engine")
	if engine == "loki" {
		// ... (existing helper logic)
		mode := c.Query("mode")
		query := c.Query("query")
		if mode == "builder" && query == "" {
			// ... (construct query) ...
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
		_, _, endpoint, ok := service.ResolveLokiDatasource(c.Query("datasourceId"))
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
	if engine == "victorialogs" {
		_, _, endpoint, ok := service.ResolveVictoriaLogsDatasource(c.Query("datasourceId"))
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no victorialogs datasource", "data": gin.H{"url": ""}})
			return
		}
		params := url.Values{}
		params.Set("query", c.Query("query"))
		if v := c.Query("start"); v != "" {
			params.Set("start", v)
		}
		if v := c.Query("end"); v != "" {
			params.Set("end", v)
		}
		urlStr := endpoint + "/select/logsql/query?" + params.Encode()
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"url": urlStr}})
		return
	}
	// ES inspect
	_, cfg, endpoint, ok := service.ResolveElasticsearchDatasource(c.Query("datasourceId"))
	if !ok {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "no elasticsearch datasource", "data": gin.H{"url": "", "body": ""}})
		return
	}
	// ... (ES Inspect Logic kept mostly same, employing Resolve helper) ...
	indexPath := ""
	if es, ok := cfg["es"].(map[string]interface{}); ok {
		if idx, ok := es["index"].(string); ok && idx != "" {
			indexPath = "/" + idx
		}
	}
	urlStr := endpoint + indexPath + "/_search"
	// ... body constr ...
	// Simplified for brevity in this tool call, keeping original structure logic
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
	b, _ := json.MarshalIndent(bodyJSON, "", "  ")
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"url": urlStr, "body": string(b)}})
}

func buildLokiQuery(filters []struct {
	Label, Op string
	Values    []string
}, contains string) string {
	if len(filters) == 0 && contains == "" {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteString("{")
	for i, f := range filters {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(f.Label)
		sb.WriteString(f.Op)
		sb.WriteString("\"")
		sb.WriteString(strings.Join(f.Values, "|"))
		sb.WriteString("\"")
	}
	sb.WriteString("}")
	if contains != "" {
		sb.WriteString(" |= \"")
		sb.WriteString(contains)
		sb.WriteString("\"")
	}
	return sb.String()
}
