package web

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/evaluator"
	"go-ut-bench/internal/obs"
	"go-ut-bench/internal/orchestrator"
	"go-ut-bench/internal/reporter"

	"gopkg.in/yaml.v3"
)

//go:embed static
var staticFiles embed.FS

// Server is the HTTP management server.
type Server struct {
	mgr        *RunManager
	bld        *BuildManager
	configPath string
	outputRoot string
	dockerCfg  DockerConfig
	mux        *http.ServeMux
}

// NewServer wires up a Server with the given RunManager and build manager.
// The DockerConfig is used by GET /api/env to report status and by the
// BuildManager to locate the Dockerfile when POST /api/env/build-image fires.
func NewServer(mgr *RunManager, bld *BuildManager, configPath, outputRoot string, cfg DockerConfig) *Server {
	s := &Server{
		mgr:        mgr,
		bld:        bld,
		configPath: configPath,
		outputRoot: outputRoot,
		dockerCfg:  cfg,
	}
	s.mux = http.NewServeMux()
	s.registerRoutes()
	return s
}

// Start begins listening on addr (e.g. ":8080").
func (s *Server) Start(addr string) error {
	fmt.Printf("UTBench Web UI  →  http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	s.mux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	// API routes
	s.mux.HandleFunc("/api/config", s.handleConfig)
	s.mux.HandleFunc("/api/env", s.handleEnv)
	s.mux.HandleFunc("/api/env/build-image", s.handleBuildImage)
	s.mux.HandleFunc("/api/env/build-image/", s.handleBuildImageSub)
	s.mux.HandleFunc("/api/environment/check", s.handleEnvironmentCheck)
	s.mux.HandleFunc("/api/environment/check/", s.handleEnvironmentCheckOne)
	s.mux.HandleFunc("/api/environment/install", s.handleEnvironmentInstall)
	s.mux.HandleFunc("/api/runs", s.handleRuns)
	s.mux.HandleFunc("/api/runs/", s.handleRunSub)
	s.mux.HandleFunc("/api/models", s.handleModels)
	s.mux.HandleFunc("/api/models/test-all", s.handleTestAllModels)
	s.mux.HandleFunc("/api/models/", s.handleModelsSub)

	// Static SPA
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	s.mux.Handle("/", http.FileServer(http.FS(sub)))
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func errJSON(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

// ─── GET /api/config ─────────────────────────────────────────────────────────

type modelInfo struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	ModelID  string `json:"model_id"`
	Enabled  bool   `json:"enabled"`
}

type configResponse struct {
	Models      []modelInfo `json:"models"`
	Languages   []string    `json:"languages"`
	Scenarios   []string    `json:"scenarios"`
	Classes     []string    `json:"classes"`
	DatasetRoot string      `json:"dataset_root"`
	ConfigPath  string      `json:"config_path"`
}

type modelsYAML struct {
	Models map[string]struct {
		Enabled  bool   `yaml:"enabled"`
		Provider string `yaml:"provider"`
		Config   struct {
			Model string `yaml:"model"`
		} `yaml:"config"`
	} `yaml:"models"`
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, "cannot read models.yaml: "+err.Error())
		return
	}
	var mf modelsYAML
	_ = yaml.Unmarshal(raw, &mf)

	models := make([]modelInfo, 0, len(mf.Models))
	for name, m := range mf.Models {
		models = append(models, modelInfo{
			Name:     name,
			Provider: m.Provider,
			ModelID:  m.Config.Model,
			Enabled:  m.Enabled,
		})
	}
	sort.Slice(models, func(i, j int) bool { return models[i].Name < models[j].Name })

	writeJSON(w, http.StatusOK, configResponse{
		Models:      models,
		Languages:   []string{"python", "go", "java", "cpp"},
		Scenarios:   []string{"boundary", "simple_function", "complex_dependency", "interface_mock"},
		Classes:     []string{"self_contained", "module_level"},
		DatasetRoot: s.mgr.datasetRoot,
		ConfigPath:  s.configPath,
	})
}

// ─── /api/runs ───────────────────────────────────────────────────────────────

type runSummaryItem struct {
	RunID     string            `json:"run_id"`
	Status    RunStatus         `json:"status"`
	Paused    bool              `json:"paused,omitempty"`
	StartedAt time.Time         `json:"started_at"`
	EndedAt   *time.Time        `json:"ended_at,omitempty"`
	Error     string            `json:"error,omitempty"`
	Spec      contracts.RunSpec `json:"spec"`
}

func (s *Server) handleRuns(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listRuns(w, r)
	case http.MethodPost:
		s.createRun(w, r)
	default:
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) listRuns(w http.ResponseWriter, _ *http.Request) {
	// Merge in-memory runs + completed runs scanned from artifacts/
	byID := make(map[string]runSummaryItem)

	// Scan artifacts/runs/*/run_summary.json
	pattern := filepath.Join(s.outputRoot, "runs", "*", "run_summary.json")
	matches, _ := filepath.Glob(pattern)
	for _, path := range matches {
		var raw map[string]any
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			continue
		}
		runID, _ := raw["run_id"].(string)
		if runID == "" {
			continue
		}
		if _, exists := byID[runID]; !exists {
			createdStr, _ := raw["created_at_utc"].(string)
			t, _ := time.Parse(time.RFC3339Nano, createdStr)
			byID[runID] = runSummaryItem{
				RunID:     runID,
				Status:    StatusCompleted,
				StartedAt: t,
			}
		}
	}

	// Override / add active in-memory runs
	for _, entry := range s.mgr.List() {
		entry.mu.RLock()
		item := runSummaryItem{
			RunID:     entry.RunID,
			Status:    entry.Status,
			Paused:    entry.Paused,
			StartedAt: entry.StartedAt,
			EndedAt:   entry.EndedAt,
			Error:     entry.Error,
			Spec:      entry.Spec,
		}
		entry.mu.RUnlock()
		byID[entry.RunID] = item
	}

	out := make([]runSummaryItem, 0, len(byID))
	for _, v := range byID {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartedAt.After(out[j].StartedAt)
	})
	writeJSON(w, http.StatusOK, out)
}

type createRunRequest struct {
	RunID           string   `json:"run_id"`
	Models          []string `json:"models"`
	Languages       []string `json:"languages"`
	Class           string   `json:"class"`
	Scenario        string   `json:"scenario"`
	Level           string   `json:"level"`
	MaxSamples      int      `json:"max_samples"`
	Workers         int      `json:"workers"`
	Mode            string   `json:"mode"`
	DryRun          bool     `json:"dry_run"`
	MutationEnabled bool     `json:"mutation_enabled"`
	MutationTimeout int      `json:"mutation_timeout"`
	MutationPolicy  string   `json:"mutation_policy"`
	Ingest          bool     `json:"ingest"`
	// UseDocker selects the Docker execution backend. When true the server
	// shells out to `docker run utbench:latest run ...` instead of running
	// the orchestrator in-process.
	UseDocker bool `json:"use_docker"`
}

func (s *Server) createRun(w http.ResponseWriter, r *http.Request) {
	var req createRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errJSON(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Models) == 0 {
		errJSON(w, http.StatusBadRequest, "models is required")
		return
	}
	if len(req.Languages) == 0 {
		errJSON(w, http.StatusBadRequest, "languages is required")
		return
	}

	runID := req.RunID
	if runID == "" {
		runID = contracts.NewRunID()
	}
	mode := req.Mode
	if mode == "" {
		mode = "full"
	}
	mutTimeout := req.MutationTimeout
	if mutTimeout == 0 {
		mutTimeout = 1800
	}
	mutPolicy := req.MutationPolicy
	if mutPolicy == "" {
		mutPolicy = "warn"
	}

	classes := splitTrim(req.Class)

	spec := contracts.RunSpec{
		RunID:           runID,
		Models:          req.Models,
		Languages:       req.Languages,
		DatasetClasses:  classes,
		DatasetScenario: req.Scenario,
		DatasetLevel:    req.Level,
		DatasetRoot:     s.mgr.datasetRoot,
		ConfigPath:      s.configPath,
		Mode:            contracts.RunMode(mode),
		DryRun:          req.DryRun,
		MutationEnabled: req.MutationEnabled,
		MutationTimeout: mutTimeout,
		MutationPolicy:  mutPolicy,
		MaxSamples:      req.MaxSamples,
		Workers:         req.Workers,
		OutputRoot:      s.outputRoot,
		CreatedAtUTC:    time.Now().UTC(),
	}
	opts := orchestrator.Options{
		Ingest: req.Ingest,
		DBPath: s.mgr.dbPath,
	}

	entry := s.mgr.Submit(spec, opts, req.UseDocker)
	writeJSON(w, http.StatusCreated, map[string]string{
		"run_id": entry.RunID,
		"status": string(entry.Status),
	})
}

// ─── /api/runs/{id}[/events|/report] ─────────────────────────────────────────

func (s *Server) handleRunSub(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/runs/")
	parts := strings.SplitN(path, "/", 2)
	runID := parts[0]
	sub := ""
	if len(parts) > 1 {
		sub = parts[1]
	}

	switch sub {
	case "events":
		s.handleRunEvents(w, r, runID)
	case "report":
		s.handleRunReport(w, r, runID)
	case "report-html":
		s.handleRunReportHTML(w, r, runID)
	case "rerun":
		s.handleRunRerun(w, r, runID)
	case "reevaluate":
		s.handleRunReevaluate(w, r, runID)
	case "regenerate-report":
		s.handleRunRegenerateReport(w, r, runID)
	case "pause", "resume", "cancel":
		s.handleRunControl(w, r, runID, sub)
	default:
		s.handleRunGet(w, r, runID)
	}
}

// handleRunControl 处理 /api/runs/{id}/{pause|resume|cancel}。
func (s *Server) handleRunControl(w http.ResponseWriter, r *http.Request, runID, action string) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var err error
	switch action {
	case "pause":
		err = s.mgr.Pause(runID)
	case "resume":
		err = s.mgr.Resume(runID)
	case "cancel":
		err = s.mgr.Cancel(runID)
	}
	if err != nil {
		errJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	entry, _ := s.mgr.Get(runID)
	status := ""
	paused := false
	if entry != nil {
		entry.mu.RLock()
		status = string(entry.Status)
		paused = entry.Paused
		entry.mu.RUnlock()
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"run_id": runID,
		"action": action,
		"status": status,
		"paused": paused,
	})
}

// handleRunRerun 以已有 run 的 spec 为模板，生成新 run_id 并提交。
// 数据来源顺序：in-memory RunEntry.Spec → run_summary.json["spec"]。
func (s *Server) handleRunRerun(w http.ResponseWriter, r *http.Request, runID string) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var spec contracts.RunSpec
	found := false
	if entry, ok := s.mgr.Get(runID); ok {
		entry.mu.RLock()
		spec = entry.Spec
		entry.mu.RUnlock()
		found = spec.RunID != ""
	}
	if !found {
		// fallback: 从 artifacts/runs/<id>/run_summary.json 恢复
		path := filepath.Join(s.outputRoot, "runs", runID, "run_summary.json")
		data, err := os.ReadFile(path)
		if err != nil {
			errJSON(w, http.StatusNotFound, "run not found or summary missing: "+runID)
			return
		}
		var raw struct {
			Spec contracts.RunSpec `json:"spec"`
		}
		if err := json.Unmarshal(data, &raw); err != nil || raw.Spec.RunID == "" {
			errJSON(w, http.StatusInternalServerError, "run_summary.json missing spec")
			return
		}
		spec = raw.Spec
	}

	// 重置生成字段：新 run_id、新时间、路径按当前 server 配置
	spec.RunID = contracts.NewRunID()
	spec.CreatedAtUTC = time.Now().UTC()
	spec.OutputRoot = s.outputRoot
	spec.ConfigPath = s.configPath
	spec.DatasetRoot = s.mgr.datasetRoot

	opts := orchestrator.Options{DBPath: s.mgr.dbPath}
	entry := s.mgr.Submit(spec, opts, true)
	writeJSON(w, http.StatusCreated, map[string]any{
		"run_id":        entry.RunID,
		"source_run_id": runID,
		"status":        string(entry.Status),
		"use_docker":    true,
	})
}

func (s *Server) loadRunSpec(runID string) (contracts.RunSpec, error) {
	if entry, ok := s.mgr.Get(runID); ok {
		entry.mu.RLock()
		spec := entry.Spec
		entry.mu.RUnlock()
		if spec.RunID != "" {
			return spec, nil
		}
	}

	path := filepath.Join(s.outputRoot, "runs", runID, "run_summary.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return contracts.RunSpec{}, fmt.Errorf("run not found or summary missing: %s", runID)
	}
	var raw struct {
		Spec contracts.RunSpec `json:"spec"`
	}
	if err := json.Unmarshal(data, &raw); err != nil || raw.Spec.RunID == "" {
		return contracts.RunSpec{}, fmt.Errorf("run_summary.json missing spec")
	}
	return raw.Spec, nil
}

func (s *Server) normalizeRunSpec(runID string, spec contracts.RunSpec) contracts.RunSpec {
	spec.RunID = runID
	spec.OutputRoot = s.outputRoot
	spec.ConfigPath = s.configPath
	spec.DatasetRoot = s.mgr.datasetRoot
	if spec.MutationPolicy == "" {
		spec.MutationPolicy = "warn"
	}
	if spec.MutationTimeout == 0 {
		spec.MutationTimeout = 1800
	}
	return spec
}

func (s *Server) handleRunReevaluate(w http.ResponseWriter, r *http.Request, runID string) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	spec, err := s.loadRunSpec(runID)
	if err != nil {
		errJSON(w, http.StatusNotFound, err.Error())
		return
	}
	spec = s.normalizeRunSpec(runID, spec)
	manifestPath := filepath.Join(s.outputRoot, "runs", runID, "generated", "generated_manifest.json")
	if _, err := os.Stat(manifestPath); err != nil {
		errJSON(w, http.StatusNotFound, "generated manifest not found: "+manifestPath)
		return
	}

	logDir := filepath.Join(s.outputRoot, "runs", runID, "logs")
	logger := obs.NewLogger(true, logDir)
	out, err := evaluator.NewService(logger).Evaluate(r.Context(), spec, manifestPath)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, "reevaluate failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"run_id":          runID,
		"manifest_path":   manifestPath,
		"evaluation_path": out.ResultPath,
		"result_count":    len(out.Result.Results),
	})
}

type regenerateReportRequest struct {
	EvaluationPath string `json:"evaluation_path"`
}

func (s *Server) handleRunRegenerateReport(w http.ResponseWriter, r *http.Request, runID string) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	spec, err := s.loadRunSpec(runID)
	if err != nil {
		errJSON(w, http.StatusNotFound, err.Error())
		return
	}
	spec = s.normalizeRunSpec(runID, spec)

	var req regenerateReportRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	evaluationPath := strings.TrimSpace(req.EvaluationPath)
	if evaluationPath == "" {
		evaluationPath = filepath.Join(s.outputRoot, "runs", runID, "evaluation", "evaluation_result.json")
	}
	if !filepath.IsAbs(evaluationPath) {
		evaluationPath = filepath.Clean(evaluationPath)
	}
	if _, err := os.Stat(evaluationPath); err != nil {
		errJSON(w, http.StatusNotFound, "evaluation JSON not found: "+evaluationPath)
		return
	}

	logDir := filepath.Join(s.outputRoot, "runs", runID, "logs")
	logger := obs.NewLogger(true, logDir)
	out, err := reporter.NewService(logger).Generate(context.Background(), spec, evaluationPath)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, "regenerate report failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"run_id":          runID,
		"evaluation_path": evaluationPath,
		"report_json":     out.ReportJSONPath,
		"report_html":     out.ReportHTMLPath,
	})
}

type runDetailResponse struct {
	RunID     string            `json:"run_id"`
	Status    RunStatus         `json:"status"`
	Paused    bool              `json:"paused,omitempty"`
	StartedAt time.Time         `json:"started_at"`
	EndedAt   *time.Time        `json:"ended_at,omitempty"`
	Error     string            `json:"error,omitempty"`
	Spec      contracts.RunSpec `json:"spec"`
	Logs      []string          `json:"logs"`
}

func (s *Server) handleRunGet(w http.ResponseWriter, r *http.Request, runID string) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	entry, ok := s.mgr.Get(runID)
	if !ok {
		spec, err := s.loadRunSpec(runID)
		if err != nil {
			errJSON(w, http.StatusNotFound, "run not found: "+runID)
			return
		}
		writeJSON(w, http.StatusOK, runDetailResponse{
			RunID:  runID,
			Status: StatusCompleted,
			Spec:   spec,
			Logs:   []string{},
		})
		return
	}
	entry.mu.RLock()
	resp := runDetailResponse{
		RunID:     entry.RunID,
		Status:    entry.Status,
		Paused:    entry.Paused,
		StartedAt: entry.StartedAt,
		EndedAt:   entry.EndedAt,
		Error:     entry.Error,
		Spec:      entry.Spec,
		Logs:      entry.GetLogs(),
	}
	entry.mu.RUnlock()
	writeJSON(w, http.StatusOK, resp)
}

// handleRunEvents streams log lines via Server-Sent Events.
func (s *Server) handleRunEvents(w http.ResponseWriter, r *http.Request, runID string) {
	entry, ok := s.mgr.Get(runID)
	if !ok {
		errJSON(w, http.StatusNotFound, "run not found: "+runID)
		return
	}

	flusher, canFlush := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	ch := entry.Subscribe()
	defer entry.Unsubscribe(ch)

	sendEvent := func(typ, payload string) {
		fmt.Fprintf(w, "data: {\"type\":%q,\"payload\":%q}\n\n", typ, payload)
		if canFlush {
			flusher.Flush()
		}
	}

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-entry.Done:
			// drain remaining lines
			for {
				select {
				case line, ok := <-ch:
					if !ok {
						goto streamDone
					}
					sendEvent("log", line)
				default:
					goto streamDone
				}
			}
		streamDone:
			entry.mu.RLock()
			status := string(entry.Status)
			entry.mu.RUnlock()
			sendEvent("done", status)
			return
		case line, ok := <-ch:
			if !ok {
				return
			}
			sendEvent("log", line)
		}
	}
}

func (s *Server) handleRunReport(w http.ResponseWriter, r *http.Request, runID string) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	reportPath := filepath.Join(s.outputRoot, "runs", runID, "report", "report_summary.json")
	data, err := os.ReadFile(reportPath)
	if err != nil {
		errJSON(w, http.StatusNotFound, "report not available yet")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(data)
}

func (s *Server) handleRunReportHTML(w http.ResponseWriter, r *http.Request, runID string) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	htmlPath := filepath.Join(s.outputRoot, "runs", runID, "report", "report.html")
	data, err := os.ReadFile(htmlPath)
	if err != nil {
		errJSON(w, http.StatusNotFound, "HTML report not available yet")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(data)
}

func splitTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
