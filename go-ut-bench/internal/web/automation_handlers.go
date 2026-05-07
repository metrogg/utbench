package web

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/orchestrator"
	"go-ut-bench/internal/store"
)

func (s *Server) handleAutomations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := s.db.ListAutomationSchedules(r.Context(), true)
		if err != nil {
			errJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, rows)
	case http.MethodPost:
		var item store.AutomationSchedule
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			errJSON(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		s.prepareAutomationSchedule(&item)
		if err := validateAutomationSchedule(item); err != nil {
			errJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		saved, err := s.db.UpsertAutomationSchedule(r.Context(), item)
		if err != nil {
			errJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, saved)
	default:
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleAutomationSub(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/automations/"), "/")
	if path == "" {
		errJSON(w, http.StatusNotFound, "automation schedule id required")
		return
	}
	parts := strings.Split(path, "/")
	id := parts[0]
	sub := ""
	if len(parts) > 1 {
		sub = parts[1]
	}
	if sub == "trigger" {
		s.handleAutomationTrigger(w, r, id)
		return
	}
	if sub == "runs" {
		s.handleAutomationRuns(w, r, id)
		return
	}
	if sub == "deliveries" {
		s.handleAutomationDeliveries(w, r, id)
		return
	}
	switch r.Method {
	case http.MethodGet:
		item, err := s.db.GetAutomationSchedule(r.Context(), id)
		if err != nil {
			errJSON(w, http.StatusNotFound, "automation schedule not found")
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodPut, http.MethodPatch:
		var item store.AutomationSchedule
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			errJSON(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		existing, _ := s.db.GetAutomationSchedule(r.Context(), id)
		item.ScheduleID = id
		item.CreatedAtUTC = existing.CreatedAtUTC
		s.prepareAutomationSchedule(&item)
		if err := validateAutomationSchedule(item); err != nil {
			errJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		saved, err := s.db.UpsertAutomationSchedule(r.Context(), item)
		if err != nil {
			errJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, saved)
	case http.MethodDelete:
		if err := s.db.DeleteAutomationSchedule(r.Context(), id); err != nil {
			errJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	default:
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleAutomationTrigger(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	item, err := s.db.GetAutomationSchedule(r.Context(), id)
	if err != nil {
		errJSON(w, http.StatusNotFound, "automation schedule not found")
		return
	}
	record, err := s.triggerAutomationSchedule(r.Context(), item, "manual")
	if err != nil {
		errJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, record)
}

func (s *Server) handleAutomationRuns(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rows, err := s.db.ListAutomationRuns(r.Context(), id, parseLimit(r, 50))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleAutomationDeliveries(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	runs, err := s.db.ListAutomationRuns(r.Context(), id, 200)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	runIDs := make(map[string]bool)
	for _, run := range runs {
		runIDs[run.AutomationRunID] = true
	}
	rows, err := s.db.ListNotificationDeliveries(r.Context(), "", parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := rows[:0]
	for _, row := range rows {
		if runIDs[row.AutomationRunID] {
			out = append(out, row)
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleNotificationChannels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := s.db.ListNotificationChannels(r.Context(), true)
		if err != nil {
			errJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, rows)
	case http.MethodPost:
		var item store.NotificationChannel
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			errJSON(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		if item.Name == "" || item.Type == "" {
			errJSON(w, http.StatusBadRequest, "name and type are required")
			return
		}
		if !validNotificationChannelType(item.Type) {
			errJSON(w, http.StatusBadRequest, "unsupported notification channel type")
			return
		}
		if item.ConfigJSON == "" {
			item.ConfigJSON = "{}"
		}
		if !json.Valid([]byte(item.ConfigJSON)) {
			errJSON(w, http.StatusBadRequest, "config_json must be valid JSON")
			return
		}
		saved, err := s.db.UpsertNotificationChannel(r.Context(), item)
		if err != nil {
			errJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, saved)
	default:
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) handleNotificationDeliveries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rows, err := s.db.ListNotificationDeliveries(r.Context(), strings.TrimSpace(r.URL.Query().Get("automation_run_id")), parseLimit(r, 80))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleNotificationChannelSub(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/notification-channels/"), "/")
	if path == "" {
		errJSON(w, http.StatusNotFound, "notification channel id required")
		return
	}
	parts := strings.Split(path, "/")
	id := parts[0]
	if len(parts) > 1 && parts[1] == "test" {
		s.handleNotificationChannelTest(w, r, id)
		return
	}
	switch r.Method {
	case http.MethodGet:
		item, err := s.db.GetNotificationChannel(r.Context(), id)
		if err != nil {
			errJSON(w, http.StatusNotFound, "notification channel not found")
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodPut, http.MethodPatch:
		var item store.NotificationChannel
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			errJSON(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}
		existing, _ := s.db.GetNotificationChannel(r.Context(), id)
		item.ChannelID = id
		item.CreatedAtUTC = existing.CreatedAtUTC
		if !validNotificationChannelType(item.Type) {
			errJSON(w, http.StatusBadRequest, "unsupported notification channel type")
			return
		}
		if item.ConfigJSON == "" {
			item.ConfigJSON = "{}"
		}
		if !json.Valid([]byte(item.ConfigJSON)) {
			errJSON(w, http.StatusBadRequest, "config_json must be valid JSON")
			return
		}
		saved, err := s.db.UpsertNotificationChannel(r.Context(), item)
		if err != nil {
			errJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, saved)
	case http.MethodDelete:
		if err := s.db.DeleteNotificationChannel(r.Context(), id); err != nil {
			errJSON(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	default:
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func validNotificationChannelType(t string) bool {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "webhook", "feishu", "lark", "dingtalk", "wecom", "email_smtp":
		return true
	default:
		return false
	}
}

func (s *Server) handleNotificationChannelTest(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	channel, err := s.db.GetNotificationChannel(r.Context(), id)
	if err != nil {
		errJSON(w, http.StatusNotFound, "notification channel not found")
		return
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	reportHTMLPath := s.latestReportHTMLPath()
	if reportHTMLPath == "" {
		reportHTMLPath = filepath.Join(s.outputRoot, "runs", "test", "report", "report.html")
	}
	payload := map[string]any{
		"type":              "utbench.notification_channel.test",
		"channel_id":        channel.ChannelID,
		"channel_name":      channel.Name,
		"schedule_name":     "通知渠道测试",
		"automation_run_id": "test",
		"run_id":            "test",
		"status":            "test",
		"error":             "",
		"sent_at_utc":       now,
		"report_html_path":  reportHTMLPath,
		"message":           "UTBench notification channel test",
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Minute)
	defer cancel()
	responseSummary, err := s.sendNotification(ctx, channel, payload)
	if err != nil {
		_, _ = s.db.CreateNotificationDelivery(r.Context(), store.NotificationDelivery{ChannelID: id, Status: "failed", Attempt: 1, RequestSummary: summarizeNotificationRequest(channel, payload), Error: err.Error(), CreatedAtUTC: now})
		errJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	_, _ = s.db.CreateNotificationDelivery(r.Context(), store.NotificationDelivery{ChannelID: id, Status: "sent", Attempt: 1, RequestSummary: summarizeNotificationRequest(channel, payload), ResponseSummary: responseSummary, CreatedAtUTC: now, DeliveredAtUTC: time.Now().UTC().Format(time.RFC3339Nano)})
	writeJSON(w, http.StatusOK, map[string]string{"status": "sent", "response_summary": responseSummary})
}

func (s *Server) latestReportHTMLPath() string {
	root := filepath.Join(s.outputRoot, "runs")
	var newestPath string
	var newestTime time.Time
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Base(path) != "report.html" {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if newestPath == "" || info.ModTime().After(newestTime) {
			newestPath = path
			newestTime = info.ModTime()
		}
		return nil
	})
	return newestPath
}

func (s *Server) prepareAutomationSchedule(item *store.AutomationSchedule) {
	if item.TriggerType == "" {
		item.TriggerType = "interval"
	}
	if item.Timezone == "" {
		item.Timezone = "Asia/Shanghai"
	}
	if item.ConcurrencyPolicy == "" {
		item.ConcurrencyPolicy = "skip"
	}
	if item.OrchestratorOptionsJSON == "" {
		opts, _ := json.Marshal(orchestrator.Options{Ingest: true, DBPath: s.mgr.dbPath, Phase: "full"})
		item.OrchestratorOptionsJSON = string(opts)
	}
	if item.NotifyPolicyJSON == "" {
		item.NotifyPolicyJSON = `{"on_success":true,"on_failure":true,"on_canceled":true,"max_attempts":3,"channel_ids":[]}`
	}
	if item.NextFireAtUTC == "" && item.Enabled {
		item.NextFireAtUTC = formatTimePtr(computeNextAutomationFire(*item, time.Now()))
	}
	if item.RunSpecJSON == "" {
		spec := contracts.RunSpec{
			Models:          []string{"deepseek-v4-flash"},
			Languages:       []string{"python"},
			DatasetClasses:  []string{"self_contained"},
			DatasetRoot:     s.mgr.datasetRoot,
			ConfigPath:      s.configPath,
			Mode:            contracts.RunMode("full"),
			ReuseGenerated:  true,
			MutationEnabled: true,
			MutationTimeout: 1800,
			MutationPolicy:  "warn",
			MaxSamples:      1,
			Workers:         4,
			OutputRoot:      s.outputRoot,
			DBPath:          s.mgr.dbPath,
			CreatedAtUTC:    time.Now().UTC(),
		}
		raw, _ := json.MarshalIndent(spec, "", "  ")
		item.RunSpecJSON = string(raw)
	}
	item.UseDocker = true
}

func validateAutomationSchedule(item store.AutomationSchedule) error {
	if strings.TrimSpace(item.Name) == "" {
		return errText("name is required")
	}
	if item.TriggerType == "cron" && strings.TrimSpace(item.CronExpr) == "" {
		return errText("cron_expr is required for cron trigger")
	}
	if item.TriggerType != "cron" && item.IntervalSeconds <= 0 {
		return errText("interval_seconds must be greater than 0")
	}
	if !json.Valid([]byte(item.RunSpecJSON)) {
		return errText("run_spec_json must be valid JSON")
	}
	if !json.Valid([]byte(item.OrchestratorOptionsJSON)) {
		return errText("orchestrator_options_json must be valid JSON")
	}
	if !json.Valid([]byte(item.NotifyPolicyJSON)) {
		return errText("notify_policy_json must be valid JSON")
	}
	return nil
}

type errText string

func (e errText) Error() string { return string(e) }
