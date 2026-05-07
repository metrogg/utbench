package web

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/orchestrator"
	"go-ut-bench/internal/store"
)

type AutomationScheduler struct {
	server *Server
	stop   chan struct{}
	once   sync.Once
}

type automationNotifyPolicy struct {
	ChannelIDs  []string `json:"channel_ids"`
	Channels    []string `json:"channels"`
	OnSuccess   bool     `json:"on_success"`
	OnFailure   bool     `json:"on_failure"`
	OnCanceled  bool     `json:"on_canceled"`
	MaxAttempts int      `json:"max_attempts"`
}

func NewAutomationScheduler(s *Server) *AutomationScheduler {
	return &AutomationScheduler{server: s, stop: make(chan struct{})}
}

func (a *AutomationScheduler) Start() {
	if a == nil {
		return
	}
	a.once.Do(func() {
		a.markInterrupted()
		a.reconcileNextFire()
		go a.loop()
	})
}

func (a *AutomationScheduler) markInterrupted() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if a.server.db == nil {
		return
	}
	n, err := a.server.db.MarkInterruptedAutomationRuns(ctx)
	if err != nil {
		fmt.Printf("[automation] mark interrupted failed: %v\n", err)
		return
	}
	if n > 0 {
		fmt.Printf("[automation] marked %d stale automation run(s) as interrupted\n", n)
	}
}

func (a *AutomationScheduler) reconcileNextFire() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if a.server.db == nil {
		return
	}
	rows, err := a.server.db.ListAutomationSchedules(ctx, false)
	if err != nil {
		fmt.Printf("[automation] reconcile schedules failed: %v\n", err)
		return
	}
	for _, schedule := range rows {
		if strings.TrimSpace(schedule.NextFireAtUTC) != "" {
			continue
		}
		next := computeNextAutomationFire(schedule, time.Now())
		if next == nil {
			continue
		}
		_ = a.server.db.SetAutomationScheduleFireTimes(ctx, schedule.ScheduleID, schedule.LastFireAtUTC, formatTimePtr(next))
	}
}

func (a *AutomationScheduler) Stop() {
	if a == nil {
		return
	}
	select {
	case <-a.stop:
	default:
		close(a.stop)
	}
}

func (a *AutomationScheduler) loop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	a.tick()
	for {
		select {
		case <-ticker.C:
			a.tick()
		case <-a.stop:
			return
		}
	}
}

func (a *AutomationScheduler) tick() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := a.server.db
	if db == nil {
		return
	}
	due, err := db.DueAutomationSchedules(ctx, time.Now().UTC().Format(time.RFC3339Nano), 10)
	if err != nil {
		fmt.Printf("[automation] list due failed: %v\n", err)
		return
	}
	for _, schedule := range due {
		if strings.TrimSpace(schedule.NextFireAtUTC) == "" {
			continue
		}
		now := time.Now()
		nowUTC := now.UTC().Format(time.RFC3339Nano)
		next := computeNextAutomationFire(schedule, now)
		nextUTC := formatTimePtr(next)
		if strings.EqualFold(schedule.ConcurrencyPolicy, "skip") {
			active, err := db.CountActiveAutomationRuns(ctx, schedule.ScheduleID)
			if err == nil && active > 0 {
				_ = db.SetAutomationScheduleFireTimes(ctx, schedule.ScheduleID, schedule.LastFireAtUTC, nextUTC)
				continue
			}
		}
		claimed, err := db.ClaimAutomationSchedule(ctx, schedule.ScheduleID, schedule.NextFireAtUTC, nowUTC, nextUTC)
		if err != nil {
			fmt.Printf("[automation] claim %s failed: %v\n", schedule.ScheduleID, err)
			continue
		}
		if !claimed {
			continue
		}
		schedule.LastFireAtUTC = nowUTC
		schedule.NextFireAtUTC = nextUTC
		if _, err := a.server.triggerAutomationSchedule(ctx, schedule, "schedule"); err != nil {
			fmt.Printf("[automation] trigger %s failed: %v\n", schedule.ScheduleID, err)
			a.server.recordAutomationTriggerFailure(ctx, schedule, "schedule", err)
		}
	}
}

func (s *Server) triggerAutomationSchedule(ctx context.Context, schedule store.AutomationSchedule, source string) (store.AutomationRun, error) {
	var spec contracts.RunSpec
	if err := json.Unmarshal([]byte(schedule.RunSpecJSON), &spec); err != nil {
		return store.AutomationRun{}, fmt.Errorf("invalid run_spec_json: %w", err)
	}
	var opts orchestrator.Options
	if strings.TrimSpace(schedule.OrchestratorOptionsJSON) != "" {
		if err := json.Unmarshal([]byte(schedule.OrchestratorOptionsJSON), &opts); err != nil {
			return store.AutomationRun{}, fmt.Errorf("invalid orchestrator_options_json: %w", err)
		}
	}
	spec.RunID = contracts.NewRunID()
	spec.CreatedAtUTC = time.Now().UTC()
	spec.OutputRoot = s.outputRoot
	spec.ConfigPath = s.configPath
	spec.DatasetRoot = s.mgr.datasetRoot
	spec.DBPath = s.mgr.dbPath
	if spec.MutationTimeout == 0 {
		spec.MutationTimeout = 1800
	}
	if spec.MutationPolicy == "" {
		spec.MutationPolicy = "warn"
	}
	if len(spec.Subjects) > 0 {
		spec.AgentsConfigPath = s.mgr.agentsConfigPath
	}
	if opts.DBPath == "" {
		opts.DBPath = s.mgr.dbPath
	}
	if opts.Phase == "" {
		opts.Phase = "full"
	}
	useDocker := schedule.UseDocker
	if !useDocker {
		useDocker = true
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)
	record := store.AutomationRun{
		ScheduleID:    schedule.ScheduleID,
		RunID:         spec.RunID,
		TriggerSource: source,
		Status:        string(StatusPending),
		StartedAtUTC:  now,
		CreatedAtUTC:  now,
	}
	record, err := s.db.CreateAutomationRun(ctx, record)
	if err != nil {
		return record, err
	}
	if source == "schedule" && schedule.LastFireAtUTC == "" {
		next := computeNextAutomationFire(schedule, time.Now())
		_ = s.db.SetAutomationScheduleFireTimes(ctx, schedule.ScheduleID, now, formatTimePtr(next))
	}

	entry := s.mgr.Submit(spec, opts, useDocker)
	s.cacheMu.Lock()
	s.runsCache = nil
	s.cacheMu.Unlock()

	go s.watchAutomationRun(record, schedule, entry)
	return record, nil
}

func (s *Server) recordAutomationTriggerFailure(ctx context.Context, schedule store.AutomationSchedule, source string, triggerErr error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	record := store.AutomationRun{
		ScheduleID:    schedule.ScheduleID,
		TriggerSource: source,
		Status:        string(StatusFailed),
		StartedAtUTC:  now,
		EndedAtUTC:    now,
		Error:         triggerErr.Error(),
		CreatedAtUTC:  now,
	}
	record, err := s.db.CreateAutomationRun(ctx, record)
	if err != nil {
		return
	}
	s.dispatchAutomationNotifications(ctx, schedule, record)
}

func (s *Server) watchAutomationRun(record store.AutomationRun, schedule store.AutomationSchedule, entry *RunEntry) {
	<-entry.Done
	entry.mu.RLock()
	status := string(entry.Status)
	errText := entry.Error
	ended := entry.EndedAt
	entry.mu.RUnlock()
	endedUTC := time.Now().UTC().Format(time.RFC3339Nano)
	if ended != nil {
		endedUTC = ended.UTC().Format(time.RFC3339Nano)
	}
	record.Status = status
	record.EndedAtUTC = endedUTC
	record.Error = errText
	record.SummaryJSON = s.buildAutomationSummary(entry.RunID, status, errText)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	_ = s.db.UpdateAutomationRun(ctx, record)
	s.dispatchAutomationNotifications(ctx, schedule, record)
}

func computeNextAutomationFire(schedule store.AutomationSchedule, from time.Time) *time.Time {
	loc, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		loc = time.Local
	}
	base := from.In(loc)
	switch strings.ToLower(schedule.TriggerType) {
	case "once":
		return nil
	case "cron":
		if next, ok := nextSimpleCron(schedule.CronExpr, base); ok {
			utc := next.UTC()
			return &utc
		}
	default:
		seconds := schedule.IntervalSeconds
		if seconds <= 0 {
			seconds = 3600
		}
		next := from.Add(time.Duration(seconds) * time.Second).UTC()
		return &next
	}
	return nil
}

func nextSimpleCron(expr string, from time.Time) (time.Time, bool) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return time.Time{}, false
	}
	t := from.Truncate(time.Minute).Add(time.Minute)
	limit := t.Add(366 * 24 * time.Hour)
	for t.Before(limit) {
		if cronFieldMatch(fields[0], t.Minute(), 0, 59) &&
			cronFieldMatch(fields[1], t.Hour(), 0, 23) &&
			cronFieldMatch(fields[2], t.Day(), 1, 31) &&
			cronFieldMatch(fields[3], int(t.Month()), 1, 12) &&
			cronFieldMatch(fields[4], int(t.Weekday()), 0, 6) {
			return t, true
		}
		t = t.Add(time.Minute)
	}
	return time.Time{}, false
}

func cronFieldMatch(field string, value, min, max int) bool {
	field = strings.TrimSpace(field)
	if field == "*" || field == "?" {
		return true
	}
	if strings.HasPrefix(field, "*/") {
		step, err := strconv.Atoi(strings.TrimPrefix(field, "*/"))
		return err == nil && step > 0 && (value-min)%step == 0
	}
	for _, part := range strings.Split(field, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err == nil && n >= min && n <= max && n == value {
			return true
		}
	}
	return false
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}
