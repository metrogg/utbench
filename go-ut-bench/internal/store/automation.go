package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type AutomationSchedule struct {
	ScheduleID              string `json:"schedule_id"`
	Name                    string `json:"name"`
	Description             string `json:"description,omitempty"`
	Enabled                 bool   `json:"enabled"`
	TriggerType             string `json:"trigger_type"`
	CronExpr                string `json:"cron_expr,omitempty"`
	IntervalSeconds         int    `json:"interval_seconds,omitempty"`
	Timezone                string `json:"timezone"`
	ConcurrencyPolicy       string `json:"concurrency_policy"`
	UseDocker               bool   `json:"use_docker"`
	RunSpecJSON             string `json:"run_spec_json"`
	OrchestratorOptionsJSON string `json:"orchestrator_options_json"`
	NotifyPolicyJSON        string `json:"notify_policy_json"`
	LastFireAtUTC           string `json:"last_fire_at_utc,omitempty"`
	NextFireAtUTC           string `json:"next_fire_at_utc,omitempty"`
	CreatedAtUTC            string `json:"created_at_utc"`
	UpdatedAtUTC            string `json:"updated_at_utc"`
	DeletedAtUTC            string `json:"deleted_at_utc,omitempty"`
}

type AutomationRun struct {
	AutomationRunID string `json:"automation_run_id"`
	ScheduleID      string `json:"schedule_id"`
	RunID           string `json:"run_id,omitempty"`
	TriggerSource   string `json:"trigger_source"`
	Status          string `json:"status"`
	StartedAtUTC    string `json:"started_at_utc,omitempty"`
	EndedAtUTC      string `json:"ended_at_utc,omitempty"`
	Error           string `json:"error,omitempty"`
	SummaryJSON     string `json:"summary_json,omitempty"`
	CreatedAtUTC    string `json:"created_at_utc"`
}

type NotificationChannel struct {
	ChannelID    string `json:"channel_id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Enabled      bool   `json:"enabled"`
	ConfigJSON   string `json:"config_json"`
	CreatedAtUTC string `json:"created_at_utc"`
	UpdatedAtUTC string `json:"updated_at_utc"`
	DeletedAtUTC string `json:"deleted_at_utc,omitempty"`
}

type NotificationDelivery struct {
	DeliveryID      string `json:"delivery_id"`
	AutomationRunID string `json:"automation_run_id,omitempty"`
	ChannelID       string `json:"channel_id"`
	Status          string `json:"status"`
	Attempt         int    `json:"attempt"`
	RequestSummary  string `json:"request_summary,omitempty"`
	ResponseSummary string `json:"response_summary,omitempty"`
	Error           string `json:"error,omitempty"`
	CreatedAtUTC    string `json:"created_at_utc"`
	DeliveredAtUTC  string `json:"delivered_at_utc,omitempty"`
}

func NewAutomationID(prefix string) string {
	var b [5]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s_%s_%s", prefix, time.Now().UTC().Format("20060102T150405"), hex.EncodeToString(b[:]))
}

func (s *SQLiteStore) ListAutomationSchedules(ctx context.Context, includeDisabled bool) ([]AutomationSchedule, error) {
	q := `SELECT schedule_id,name,description,enabled,trigger_type,cron_expr,interval_seconds,timezone,concurrency_policy,use_docker,run_spec_json,orchestrator_options_json,notify_policy_json,last_fire_at_utc,next_fire_at_utc,created_at_utc,updated_at_utc,deleted_at_utc
		FROM automation_schedules WHERE (deleted_at_utc IS NULL OR deleted_at_utc = '')`
	if !includeDisabled {
		q += ` AND enabled = 1`
	}
	q += ` ORDER BY enabled DESC, next_fire_at_utc ASC, updated_at_utc DESC`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AutomationSchedule{}
	for rows.Next() {
		item, err := scanAutomationSchedule(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) DueAutomationSchedules(ctx context.Context, nowUTC string, limit int) ([]AutomationSchedule, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `SELECT schedule_id,name,description,enabled,trigger_type,cron_expr,interval_seconds,timezone,concurrency_policy,use_docker,run_spec_json,orchestrator_options_json,notify_policy_json,last_fire_at_utc,next_fire_at_utc,created_at_utc,updated_at_utc,deleted_at_utc
		FROM automation_schedules
		WHERE (deleted_at_utc IS NULL OR deleted_at_utc = '') AND enabled = 1 AND next_fire_at_utc IS NOT NULL AND next_fire_at_utc <= ?
		ORDER BY next_fire_at_utc ASC LIMIT ?`, nowUTC, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AutomationSchedule{}
	for rows.Next() {
		item, err := scanAutomationSchedule(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) GetAutomationSchedule(ctx context.Context, id string) (AutomationSchedule, error) {
	row := s.db.QueryRowContext(ctx, `SELECT schedule_id,name,description,enabled,trigger_type,cron_expr,interval_seconds,timezone,concurrency_policy,use_docker,run_spec_json,orchestrator_options_json,notify_policy_json,last_fire_at_utc,next_fire_at_utc,created_at_utc,updated_at_utc,deleted_at_utc
		FROM automation_schedules WHERE schedule_id = ? AND (deleted_at_utc IS NULL OR deleted_at_utc = '')`, id)
	return scanAutomationSchedule(row)
}

func (s *SQLiteStore) UpsertAutomationSchedule(ctx context.Context, item AutomationSchedule) (AutomationSchedule, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if item.ScheduleID == "" {
		item.ScheduleID = NewAutomationID("sch")
	}
	if item.CreatedAtUTC == "" {
		item.CreatedAtUTC = now
	}
	item.UpdatedAtUTC = now
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
		item.OrchestratorOptionsJSON = "{}"
	}
	if item.NotifyPolicyJSON == "" {
		item.NotifyPolicyJSON = "{}"
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO automation_schedules (
			schedule_id,name,description,enabled,trigger_type,cron_expr,interval_seconds,timezone,concurrency_policy,use_docker,run_spec_json,orchestrator_options_json,notify_policy_json,last_fire_at_utc,next_fire_at_utc,created_at_utc,updated_at_utc,deleted_at_utc
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(schedule_id) DO UPDATE SET
			name=excluded.name, description=excluded.description, enabled=excluded.enabled,
			trigger_type=excluded.trigger_type, cron_expr=excluded.cron_expr, interval_seconds=excluded.interval_seconds,
			timezone=excluded.timezone, concurrency_policy=excluded.concurrency_policy, use_docker=excluded.use_docker,
			run_spec_json=excluded.run_spec_json, orchestrator_options_json=excluded.orchestrator_options_json,
			notify_policy_json=excluded.notify_policy_json, last_fire_at_utc=excluded.last_fire_at_utc,
			next_fire_at_utc=excluded.next_fire_at_utc, updated_at_utc=excluded.updated_at_utc,
			deleted_at_utc=excluded.deleted_at_utc`,
		item.ScheduleID, item.Name, item.Description, boolInt(item.Enabled), item.TriggerType, item.CronExpr,
		item.IntervalSeconds, item.Timezone, item.ConcurrencyPolicy, boolInt(item.UseDocker), item.RunSpecJSON,
		item.OrchestratorOptionsJSON, item.NotifyPolicyJSON, item.LastFireAtUTC, item.NextFireAtUTC,
		item.CreatedAtUTC, item.UpdatedAtUTC, nullableTimeText(item.DeletedAtUTC))
	return item, err
}

func (s *SQLiteStore) DeleteAutomationSchedule(ctx context.Context, id string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := s.db.ExecContext(ctx, `UPDATE automation_schedules SET deleted_at_utc = ?, enabled = 0, updated_at_utc = ? WHERE schedule_id = ?`, now, now, id)
	return err
}

func (s *SQLiteStore) SetAutomationScheduleFireTimes(ctx context.Context, id, lastUTC, nextUTC string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE automation_schedules SET last_fire_at_utc = ?, next_fire_at_utc = ?, updated_at_utc = ? WHERE schedule_id = ?`,
		lastUTC, nextUTC, time.Now().UTC().Format(time.RFC3339Nano), id)
	return err
}

func (s *SQLiteStore) ClaimAutomationSchedule(ctx context.Context, id, expectedNextUTC, lastUTC, nextUTC string) (bool, error) {
	res, err := s.db.ExecContext(ctx, `UPDATE automation_schedules
		SET last_fire_at_utc = ?, next_fire_at_utc = ?, updated_at_utc = ?
		WHERE schedule_id = ? AND enabled = 1 AND deleted_at_utc IS NULL AND next_fire_at_utc = ?`,
		lastUTC, nextUTC, time.Now().UTC().Format(time.RFC3339Nano), id, expectedNextUTC)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n == 1, err
}

func (s *SQLiteStore) CreateAutomationRun(ctx context.Context, item AutomationRun) (AutomationRun, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if item.AutomationRunID == "" {
		item.AutomationRunID = NewAutomationID("arun")
	}
	if item.CreatedAtUTC == "" {
		item.CreatedAtUTC = now
	}
	if item.Status == "" {
		item.Status = "pending"
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO automation_runs (automation_run_id,schedule_id,run_id,trigger_source,status,started_at_utc,ended_at_utc,error,summary_json,created_at_utc)
		VALUES (?,?,?,?,?,?,?,?,?,?)`, item.AutomationRunID, item.ScheduleID, item.RunID, item.TriggerSource, item.Status, item.StartedAtUTC, item.EndedAtUTC, item.Error, item.SummaryJSON, item.CreatedAtUTC)
	return item, err
}

func (s *SQLiteStore) UpdateAutomationRun(ctx context.Context, item AutomationRun) error {
	_, err := s.db.ExecContext(ctx, `UPDATE automation_runs SET run_id = ?, status = ?, started_at_utc = ?, ended_at_utc = ?, error = ?, summary_json = ? WHERE automation_run_id = ?`,
		item.RunID, item.Status, item.StartedAtUTC, item.EndedAtUTC, item.Error, item.SummaryJSON, item.AutomationRunID)
	return err
}

func (s *SQLiteStore) ListAutomationRuns(ctx context.Context, scheduleID string, limit int) ([]AutomationRun, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `SELECT automation_run_id,schedule_id,run_id,trigger_source,status,started_at_utc,ended_at_utc,error,summary_json,created_at_utc
		FROM automation_runs WHERE schedule_id = ? ORDER BY created_at_utc DESC LIMIT ?`, scheduleID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AutomationRun{}
	for rows.Next() {
		item, err := scanAutomationRun(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) CountActiveAutomationRuns(ctx context.Context, scheduleID string) (int, error) {
	var n int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM automation_runs WHERE schedule_id = ? AND status IN ('pending','running')`, scheduleID).Scan(&n)
	return n, err
}

func (s *SQLiteStore) MarkInterruptedAutomationRuns(ctx context.Context) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	res, err := s.db.ExecContext(ctx, `UPDATE automation_runs
		SET status = 'interrupted', ended_at_utc = ?, error = 'web server restarted before automation run completed'
		WHERE status IN ('pending','running')`, now)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *SQLiteStore) ListNotificationChannels(ctx context.Context, includeDisabled bool) ([]NotificationChannel, error) {
	q := `SELECT channel_id,name,type,enabled,config_json,created_at_utc,updated_at_utc,deleted_at_utc FROM notification_channels WHERE (deleted_at_utc IS NULL OR deleted_at_utc = '')`
	if !includeDisabled {
		q += ` AND enabled = 1`
	}
	q += ` ORDER BY enabled DESC, updated_at_utc DESC`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []NotificationChannel{}
	for rows.Next() {
		item, err := scanNotificationChannel(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) GetNotificationChannel(ctx context.Context, id string) (NotificationChannel, error) {
	row := s.db.QueryRowContext(ctx, `SELECT channel_id,name,type,enabled,config_json,created_at_utc,updated_at_utc,deleted_at_utc FROM notification_channels WHERE channel_id = ? AND (deleted_at_utc IS NULL OR deleted_at_utc = '')`, id)
	return scanNotificationChannel(row)
}

func (s *SQLiteStore) UpsertNotificationChannel(ctx context.Context, item NotificationChannel) (NotificationChannel, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if item.ChannelID == "" {
		item.ChannelID = NewAutomationID("chn")
	}
	if item.CreatedAtUTC == "" {
		item.CreatedAtUTC = now
	}
	item.UpdatedAtUTC = now
	if item.ConfigJSON == "" {
		item.ConfigJSON = "{}"
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO notification_channels (channel_id,name,type,enabled,config_json,created_at_utc,updated_at_utc,deleted_at_utc)
		VALUES (?,?,?,?,?,?,?,?)
		ON CONFLICT(channel_id) DO UPDATE SET name=excluded.name,type=excluded.type,enabled=excluded.enabled,config_json=excluded.config_json,updated_at_utc=excluded.updated_at_utc,deleted_at_utc=excluded.deleted_at_utc`,
		item.ChannelID, item.Name, item.Type, boolInt(item.Enabled), item.ConfigJSON, item.CreatedAtUTC, item.UpdatedAtUTC, nullableTimeText(item.DeletedAtUTC))
	return item, err
}

func nullableTimeText(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

func (s *SQLiteStore) DeleteNotificationChannel(ctx context.Context, id string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := s.db.ExecContext(ctx, `UPDATE notification_channels SET deleted_at_utc = ?, enabled = 0, updated_at_utc = ? WHERE channel_id = ?`, now, now, id)
	return err
}

func (s *SQLiteStore) CreateNotificationDelivery(ctx context.Context, item NotificationDelivery) (NotificationDelivery, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if item.DeliveryID == "" {
		item.DeliveryID = NewAutomationID("dlv")
	}
	if item.CreatedAtUTC == "" {
		item.CreatedAtUTC = now
	}
	if item.Attempt <= 0 {
		item.Attempt = 1
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO notification_deliveries (delivery_id,automation_run_id,channel_id,status,attempt,request_summary,response_summary,error,created_at_utc,delivered_at_utc)
		VALUES (?,?,?,?,?,?,?,?,?,?)`, item.DeliveryID, item.AutomationRunID, item.ChannelID, item.Status, item.Attempt, item.RequestSummary, item.ResponseSummary, item.Error, item.CreatedAtUTC, item.DeliveredAtUTC)
	return item, err
}

func (s *SQLiteStore) ListNotificationDeliveries(ctx context.Context, automationRunID string, limit int) ([]NotificationDelivery, error) {
	if limit <= 0 {
		limit = 80
	}
	q := `SELECT delivery_id,automation_run_id,channel_id,status,attempt,request_summary,response_summary,error,created_at_utc,delivered_at_utc FROM notification_deliveries`
	args := []any{}
	if automationRunID != "" {
		q += ` WHERE automation_run_id = ?`
		args = append(args, automationRunID)
	}
	q += ` ORDER BY created_at_utc DESC LIMIT ?`
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []NotificationDelivery{}
	for rows.Next() {
		item, err := scanNotificationDelivery(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanAutomationSchedule(row scanner) (AutomationSchedule, error) {
	var item AutomationSchedule
	var enabled, useDocker int
	var description, cronExpr, lastFire, nextFire, deleted sql.NullString
	if err := row.Scan(&item.ScheduleID, &item.Name, &description, &enabled, &item.TriggerType, &cronExpr, &item.IntervalSeconds, &item.Timezone, &item.ConcurrencyPolicy, &useDocker, &item.RunSpecJSON, &item.OrchestratorOptionsJSON, &item.NotifyPolicyJSON, &lastFire, &nextFire, &item.CreatedAtUTC, &item.UpdatedAtUTC, &deleted); err != nil {
		return item, err
	}
	item.Enabled = enabled != 0
	item.UseDocker = useDocker != 0
	item.Description = sqlNullString(description)
	item.CronExpr = sqlNullString(cronExpr)
	item.LastFireAtUTC = sqlNullString(lastFire)
	item.NextFireAtUTC = sqlNullString(nextFire)
	item.DeletedAtUTC = sqlNullString(deleted)
	return item, nil
}

func scanAutomationRun(row scanner) (AutomationRun, error) {
	var item AutomationRun
	var runID, started, ended, errText, summary sql.NullString
	err := row.Scan(&item.AutomationRunID, &item.ScheduleID, &runID, &item.TriggerSource, &item.Status, &started, &ended, &errText, &summary, &item.CreatedAtUTC)
	item.RunID = sqlNullString(runID)
	item.StartedAtUTC = sqlNullString(started)
	item.EndedAtUTC = sqlNullString(ended)
	item.Error = sqlNullString(errText)
	item.SummaryJSON = sqlNullString(summary)
	return item, err
}

func scanNotificationChannel(row scanner) (NotificationChannel, error) {
	var item NotificationChannel
	var enabled int
	var deleted sql.NullString
	if err := row.Scan(&item.ChannelID, &item.Name, &item.Type, &enabled, &item.ConfigJSON, &item.CreatedAtUTC, &item.UpdatedAtUTC, &deleted); err != nil {
		return item, err
	}
	item.Enabled = enabled != 0
	item.DeletedAtUTC = sqlNullString(deleted)
	return item, nil
}

func scanNotificationDelivery(row scanner) (NotificationDelivery, error) {
	var item NotificationDelivery
	var automationRunID, requestSummary, responseSummary, errText, delivered sql.NullString
	err := row.Scan(&item.DeliveryID, &automationRunID, &item.ChannelID, &item.Status, &item.Attempt, &requestSummary, &responseSummary, &errText, &item.CreatedAtUTC, &delivered)
	item.AutomationRunID = sqlNullString(automationRunID)
	item.RequestSummary = sqlNullString(requestSummary)
	item.ResponseSummary = sqlNullString(responseSummary)
	item.Error = sqlNullString(errText)
	item.DeliveredAtUTC = sqlNullString(delivered)
	return item, err
}

func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func sqlNullString(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return v.String
}
