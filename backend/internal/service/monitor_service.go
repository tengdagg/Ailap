package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"ailap-backend/internal/database"
	"ailap-backend/internal/model"
	"ailap-backend/internal/utils"
)

type MonitorService struct {
	cron          *cron.Cron
	logService    *LogService
	aiService     *AIService
	notifyService *NotificationService
	jobMap        map[uint]cron.EntryID
	mu            sync.Mutex
}

func NewMonitorService() *MonitorService {
	// Standard cron parser with support for seconds
	c := cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	c.Start()

	ms := &MonitorService{
		cron:          c,
		logService:    NewLogService(),
		aiService:     NewAIService(),
		notifyService: NewNotificationService(),
		jobMap:        make(map[uint]cron.EntryID),
	}

	// Load active monitors on startup
	go ms.loadJobs()

	return ms
}

func (s *MonitorService) loadJobs() {
	var monitors []model.LogMonitor
	if err := database.GetDB().Where("status = ?", "active").Find(&monitors).Error; err != nil {
		utils.GetLogger().Error("failed to load monitors", zap.Error(err))
		return
	}
	for _, m := range monitors {
		if err := s.AddJob(&m); err != nil {
			utils.GetLogger().Error("failed to start monitor job", zap.String("name", m.Name), zap.Error(err))
		}
	}
}

// AddJob adds or updates a cron job for the monitor
func (s *MonitorService) AddJob(m *model.LogMonitor) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing if any
	if eid, ok := s.jobMap[m.ID]; ok {
		s.cron.Remove(eid)
		delete(s.jobMap, m.ID)
	}

	if m.Status != "active" {
		return nil
	}

	// Wrapper for execution
	job := func() {
		s.ExecuteMonitor(m.ID)
	}

	eid, err := s.cron.AddFunc(m.Cron, job)
	if err != nil {
		return err
	}
	s.jobMap[m.ID] = eid
	return nil
}

// RemoveJob stops the cron job
func (s *MonitorService) RemoveJob(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if eid, ok := s.jobMap[id]; ok {
		s.cron.Remove(eid)
		delete(s.jobMap, id)
	}
}

// ExecuteMonitor is the core logic: Query -> Filter -> AI -> Notify
func (s *MonitorService) ExecuteMonitor(monitorID uint) {
	utils.GetLogger().Info("monitor job started", zap.Uint("monitor_id", monitorID))

	var m model.LogMonitor
	if err := database.GetDB().First(&m, monitorID).Error; err != nil {
		utils.GetLogger().Error("monitor not found during execution", zap.Uint("id", monitorID))
		return
	}

	// 1. Determine Time Range (based on last run or cron interval?)
	// Simplification: query last 1 hour, or parse cron interval?
	// For "Smart Monitoring", usually we look back 'interval' amount of time.
	// But parsing cron to interval is hard.
	// Let's assume we look back 1 hour by default, or maybe making it configurable?
	// Let's try to infer from cron if it starts with "@every"? No.
	// Hardcode lookback to 1h for now, or use a new field `Lookback` in model?
	// Using 1h as reasonable default for periodic checks.
	now := time.Now()
	start := now.Add(-1 * time.Hour)
	end := now

	startNs := fmt.Sprintf("%d", start.UnixNano())
	endNs := fmt.Sprintf("%d", end.UnixNano())

	// 2. Query Logs
	// Need to check keywords?
	// If Engine is ES/Loki, keywords can be part of query.
	// If Engine is VL, same.
	// We append keywords to base query.

	baseQuery := m.Query
	if baseQuery == "" {
		baseQuery = "*" // or empty
	}

	// Construct effective query with keywords
	effectiveQuery := baseQuery
	keywords := strings.Split(m.Keywords, ",")
	validKeywords := []string{}
	for _, k := range keywords {
		k = strings.TrimSpace(k)
		if k != "" {
			validKeywords = append(validKeywords, k)
		}
	}

	// Appending keywords depends on engine syntax
	// Simple approach: perform query "*" then filter in memory?
	// OR: construct query.
	// Loki: {app="foo"} |= "error"
	// ES: app:foo AND (error OR warning)
	// VL: app="foo" (error OR warning)

	// Let's assume user puts base filters in 'Query' and we append keywords as OR condition text search
	if len(validKeywords) > 0 {
		orClause := strings.Join(validKeywords, "|") // Regex style for Loki?
		if m.Engine == "loki" {
			// Loki: ... |~ "k1|k2"
			effectiveQuery = fmt.Sprintf("%s |~ \"%s\"", baseQuery, orClause)
		} else if m.Engine == "elasticsearch" {
			// ES: ... AND (k1 OR k2)
			orClauseES := strings.Join(validKeywords, " OR ")
			if baseQuery == "*" {
				effectiveQuery = fmt.Sprintf("(%s)", orClauseES)
			} else {
				effectiveQuery = fmt.Sprintf("(%s) AND (%s)", baseQuery, orClauseES)
			}
		} else if m.Engine == "victorialogs" {
			// VL: ... (k1 or k2)
			orClauseVL := strings.Join(validKeywords, " or ")
			if baseQuery == "" || baseQuery == "*" {
				effectiveQuery = fmt.Sprintf("(%s)", orClauseVL)
			} else {
				effectiveQuery = fmt.Sprintf("%s (%s)", baseQuery, orClauseVL)
			}
		}
	}

	result, err := s.logService.ExecuteQuery(context.Background(), m.Engine, m.DatasourceID, effectiveQuery, startNs, endNs, 100) // limit 100 for analysis
	if err != nil {
		utils.GetLogger().Error("monitor query failed", zap.Uint("id", m.ID), zap.Error(err))
		return
	}

	if len(result.Items) == 0 {
		// No logs found matches keywords
		utils.GetLogger().Info("monitor found no logs", zap.Uint("id", m.ID))
		return
	}

	// 3. AI Analysis
	// Convert result items to interface slice
	logsInterface := make([]interface{}, len(result.Items))
	for i, v := range result.Items {
		logsInterface[i] = v
	}

	analysis, err := s.aiService.Analyze(fmt.Sprintf("Monitoring Alert: Found %d abnormal logs containing keywords [%s]. Please analyze.", len(result.Items), m.Keywords), logsInterface)
	if err != nil {
		utils.GetLogger().Error("monitor ai analysis failed", zap.Uint("id", m.ID), zap.Error(err))
		analysis = "AI Analysis Failed: " + err.Error()
	}

	// 4. Notify
	var channel model.NotificationChannel
	if err := database.GetDB().First(&channel, m.ChannelID).Error; err != nil {
		utils.GetLogger().Error("monitor channel not found", zap.Uint("id", m.ID), zap.Uint("channel_id", m.ChannelID))
		return
	}

	title := fmt.Sprintf("Smart Alert: %s", m.Name)
	content := fmt.Sprintf("Monitor: %s\nTime: %s\nMatches: %d\nKeywords: %s\n\nAI Analysis:\n%s", m.Name, time.Now().Format(time.RFC3339), len(result.Items), m.Keywords, analysis)

	if err := s.notifyService.SendAlert(&channel, title, content); err != nil {
		utils.GetLogger().Error("monitor notification failed", zap.Error(err))
	} else {
		utils.GetLogger().Info("monitor alert sent", zap.Uint("id", m.ID))
	}

	// Update LastRun
	t := time.Now()
	m.LastRunAt = &t
	database.GetDB().Save(&m)
}
