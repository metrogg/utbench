// web 包提供 HTTP 管理服务
// 提供 Web UI 和 API 接口，用于启动评测、查看进度、生成报告
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
	"strconv"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/obs"
	"go-ut-bench/internal/orchestrator"
	"go-ut-bench/internal/reporter"
	"go-ut-bench/internal/store"

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
	db         *store.SQLiteStore // 持久化的数据库连接，避免每次请求重新打开
}

// NewServer wires up a Server with the given RunManager and build manager.
// The DockerConfig is used by GET /api/env to report status and by the
// BuildManager to locate the Dockerfile when POST /api/env/build-image fires.
// dbPath is the SQLite database path; the connection is opened and initialized
// at startup and reused for all requests.
func NewServer(mgr *RunManager, bld *BuildManager, configPath, outputRoot, dbPath string, cfg DockerConfig) (*Server, error) {
	// 启动时打开并初始化数据库连接，避免每次请求重新打开
	db, err := store.OpenSQLite(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := db.Init(context.Background()); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("init sqlite: %w", err)
	}

	s := &Server{
		mgr:        mgr,
		bld:        bld,
		configPath: configPath,
		outputRoot: outputRoot,
		dockerCfg:  cfg,
		db:         db,
	}
	s.mux = http.NewServeMux()
	s.registerRoutes()
	return s, nil
}

// Start begins listening on addr (e.g. ":8080").
func (s *Server) Start(addr string) error {
	fmt.Printf("UTBench Web UI  →  http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s)
}

// Close closes the database connection.
func (s *Server) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
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
	// 数据库管理API
	s.mux.HandleFunc("/api/db/overview", s.handleDBOverview)
	s.mux.HandleFunc("/api/db/runs", s.handleDBRuns)
	s.mux.HandleFunc("/api/db/results", s.handleDBResults)
	s.mux.HandleFunc("/api/db/artifacts", s.handleDBArtifacts)
	s.mux.HandleFunc("/api/db/facets", s.handleDBFacets)
	s.mux.HandleFunc("/api/db/ingest-run", s.handleDBIngestRun)
	s.mux.HandleFunc("/api/db/report", s.handleDBReport)
	// 新增：数据库完整管理API
	s.mux.HandleFunc("/api/db/generation-runs", s.handleDBGenerationRuns)
	s.mux.HandleFunc("/api/db/generated-cases", s.handleDBGeneratedCases)
	s.mux.HandleFunc("/api/db/prompt-renderings", s.handleDBPromptRenderings)
	s.mux.HandleFunc("/api/db/evaluation-runs", s.handleDBEvaluationRuns)
	s.mux.HandleFunc("/api/db/evaluation-stages", s.handleDBEvaluationStages)
	s.mux.HandleFunc("/api/db/dataset-samples", s.handleDBDatasetSamples)
	s.mux.HandleFunc("/api/db/dataset-snapshots", s.handleDBDatasetSnapshots)
	s.mux.HandleFunc("/api/db/model-configs", s.handleDBModelConfigs)
	s.mux.HandleFunc("/api/db/prompt-profiles", s.handleDBPromptProfiles)
	s.mux.HandleFunc("/api/db/evaluation-envs", s.handleDBEvaluationEnvs)
	s.mux.HandleFunc("/api/db/score-policies", s.handleDBScorePolicies)
	s.mux.HandleFunc("/api/db/reports", s.handleDBReports)
	s.mux.HandleFunc("/api/db/run-artifacts", s.handleDBRunArtifacts)
	s.mux.HandleFunc("/api/db/experiments", s.handleDBExperiments)

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

func parseLimit(r *http.Request, def int) int {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return def
	}
	return v
}

// openStore 返回Server启动时初始化的数据库连接。
// 不再每次请求重新打开，避免SQLite锁竞争和性能问题。
func (s *Server) openStore(ctx context.Context) (*store.SQLiteStore, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return s.db, nil
}

// ─── /api/db/* ──────────────────────────────────────────────────────────────

func (s *Server) handleDBOverview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	overview, err := db.Overview(r.Context(), parseLimit(r, 10))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, overview)
}

func (s *Server) handleDBRuns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListRuns(r.Context(), parseLimit(r, 50))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	q := r.URL.Query()
	rows, err := db.ListResults(r.Context(), q.Get("run_id"), q.Get("model"), q.Get("language"), parseLimit(r, 200))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBArtifacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	q := r.URL.Query()
	rows, err := db.ListArtifacts(r.Context(), q.Get("run_id"), q.Get("kind"), parseLimit(r, 200))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBFacets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	facets, err := db.ReportFacets(r.Context(), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, facets)
}

type dbIngestRunRequest struct {
	RunID  string `json:"run_id"`
	RunDir string `json:"run_dir"`
}

func (s *Server) handleDBIngestRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req dbIngestRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errJSON(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}
	runDir := strings.TrimSpace(req.RunDir)
	if runDir == "" {
		if strings.TrimSpace(req.RunID) == "" {
			errJSON(w, http.StatusBadRequest, "run_id or run_dir is required")
			return
		}
		runDir = filepath.Join(s.outputRoot, "runs", req.RunID)
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	sum, err := db.IngestRun(r.Context(), store.IngestRunOptions{RunDir: runDir})
	if err != nil {
		errJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sum)
}

type dbReportRequest struct {
	RunID             string   `json:"run_id"`
	SourceRunIDs      []string `json:"source_run_ids"`
	EvaluationRunIDs  []string `json:"evaluation_run_ids"`
	Models            []string `json:"models"`
	Languages         []string `json:"languages"`
	ScoreEligibleOnly bool     `json:"score_eligible_only"`
}

func (s *Server) handleDBReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req dbReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errJSON(w, http.StatusBadRequest, "invalid json: "+err.Error())
		return
	}
	outRunID := strings.TrimSpace(req.RunID)
	if outRunID == "" {
		outRunID = contracts.NewRunID() + "_db_report"
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	set, err := db.SelectEvaluationResultSet(r.Context(), outRunID, store.DBReportFilter{
		RunIDs:            req.SourceRunIDs,
		EvaluationRunIDs:  req.EvaluationRunIDs,
		Models:            req.Models,
		Languages:         req.Languages,
		ScoreEligibleOnly: req.ScoreEligibleOnly,
	})
	if err != nil {
		errJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	logger := obs.NewLogger(false, filepath.Join(s.outputRoot, "runs", outRunID, "logs"))
	out, err := reporter.NewService(logger).GenerateFromResultSet(contracts.RunSpec{
		RunID:      outRunID,
		OutputRoot: s.outputRoot,
		ConfigPath: s.configPath,
	}, set, "db://"+s.mgr.dbPath)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err := db.IngestReportFile(r.Context(), out.ReportJSONPath); err != nil {
		errJSON(w, http.StatusInternalServerError, "report generated but ingest failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"run_id":           outRunID,
		"result_count":     len(set.Results),
		"report_json_path": out.ReportJSONPath,
		"report_html_path": out.ReportHTMLPath,
	})
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
	ReuseGenerated  bool     `json:"reuse_generated"`
	MutationEnabled bool     `json:"mutation_enabled"`
	MutationTimeout int      `json:"mutation_timeout"`
	MutationPolicy  string   `json:"mutation_policy"`
	Ingest          bool     `json:"ingest"`
	// UseDocker selects the Docker execution backend. When true the server
	// shells out to `docker run utbench:latest run ...` instead of running
	// the orchestrator in-process.
	UseDocker bool `json:"use_docker"`
	// Phase controls which pipeline stage(s) to execute.
	// Supported values: "full" (default), "generate", "evaluate", "report".
	// When phase is not "full", the run requires existing artifacts from previous stages.
	Phase string `json:"phase"`
	// SourceRunID specifies the run ID to use as data source for evaluate/report phases.
	// If empty, uses the current RunID (which must have existing artifacts).
	SourceRunID string `json:"source_run_id"`
	// ManifestPath overrides the default manifest path for evaluate phase.
	// If empty, uses artifacts/runs/<source_run_id>/generated/generated_manifest.json.
	ManifestPath string `json:"manifest_path"`
	// EvaluationPath overrides the default evaluation path for report phase.
	// If empty, uses artifacts/runs/<source_run_id>/evaluation/evaluation_result.json.
	EvaluationPath string `json:"evaluation_path"`
}

func (s *Server) createRun(w http.ResponseWriter, r *http.Request) {
	var req createRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errJSON(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	phase := req.Phase
	if phase == "" {
		phase = "full"
	}

	// Validate based on phase
	if phase == "full" || phase == "generate" {
		if len(req.Models) == 0 {
			errJSON(w, http.StatusBadRequest, "models is required for generate/full phase")
			return
		}
		if len(req.Languages) == 0 {
			errJSON(w, http.StatusBadRequest, "languages is required for generate/full phase")
			return
		}
	}

	// For evaluate/report phases, we need a data source
	if phase == "evaluate" || phase == "report" {
		// Either source_run_id or explicit path must be provided
		if req.SourceRunID == "" && req.ManifestPath == "" && req.EvaluationPath == "" {
			errJSON(w, http.StatusBadRequest, "source_run_id or explicit path is required for evaluate/report phase")
			return
		}
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
		ReuseGenerated:  req.ReuseGenerated,
		DBPath:          s.mgr.dbPath,
		MutationEnabled: req.MutationEnabled,
		MutationTimeout: mutTimeout,
		MutationPolicy:  mutPolicy,
		MaxSamples:      req.MaxSamples,
		Workers:         req.Workers,
		OutputRoot:      s.outputRoot,
		CreatedAtUTC:    time.Now().UTC(),
	}
	opts := orchestrator.Options{
		Ingest:         req.Ingest,
		DBPath:         s.mgr.dbPath,
		Phase:          phase,
		SourceRunID:    req.SourceRunID,
		ManifestPath:   req.ManifestPath,
		EvaluationPath: req.EvaluationPath,
	}

	entry := s.mgr.Submit(spec, opts, req.UseDocker)
	writeJSON(w, http.StatusCreated, map[string]string{
		"run_id":        entry.RunID,
		"status":        string(entry.Status),
		"phase":         phase,
		"source_run_id": req.SourceRunID,
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
	if !isDockerReady(s.dockerCfg) {
		errJSON(w, http.StatusConflict, "Docker image is not ready; reevaluate requires Docker so evaluator tools are complete")
		return
	}
	out, err := runEvaluateInDocker(r.Context(), runID, spec, s.dockerCfg)
	if err != nil {
		errJSON(w, http.StatusInternalServerError, "reevaluate failed: "+err.Error()+": "+tailString(string(out), 2000))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"run_id":          runID,
		"manifest_path":   manifestPath,
		"evaluation_path": filepath.Join(s.outputRoot, "runs", runID, "evaluation", "evaluation_result.json"),
		"backend":         "docker",
		"output_tail":     tailString(string(out), 2000),
	})
}

func tailString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[len(s)-max:]
}

func (s *Server) prepareManifestForHost(runID, manifestPath string) (string, error) {
	manifest, err := contracts.ReadGeneratedManifest(manifestPath)
	if err != nil {
		return "", err
	}

	changed := false
	convert := func(path string) string {
		next := s.containerPathToHost(path)
		if next != path {
			changed = true
		}
		return next
	}
	manifest.Spec = s.normalizeRunSpec(runID, manifest.Spec)
	manifest.PromptSnapshotDir = convert(manifest.PromptSnapshotDir)
	for i := range manifest.Cases {
		manifest.Cases[i].SamplePath = convert(manifest.Cases[i].SamplePath)
		manifest.Cases[i].GeneratedTestPath = convert(manifest.Cases[i].GeneratedTestPath)
		manifest.Cases[i].ResponsePath = convert(manifest.Cases[i].ResponsePath)
		manifest.Cases[i].MetadataPath = convert(manifest.Cases[i].MetadataPath)
		manifest.Cases[i].PromptPath = convert(manifest.Cases[i].PromptPath)
	}
	if !changed {
		return manifestPath, nil
	}

	hostPath := filepath.Join(s.outputRoot, "runs", runID, "generated", "generated_manifest.host.json")
	if err := contracts.WriteJSON(hostPath, manifest); err != nil {
		return "", err
	}
	return hostPath, nil
}

func (s *Server) containerPathToHost(path string) string {
	if path == "" {
		return ""
	}
	clean := filepath.ToSlash(path)
	prefixes := []struct {
		container string
		host      string
	}{
		{"/app/artifacts", s.outputRoot},
		{"/app/datasets", s.mgr.datasetRoot},
		{"/app/configs", filepath.Dir(s.configPath)},
		{"/app/storage", filepath.Dir(s.mgr.dbPath)},
	}
	for _, p := range prefixes {
		if clean == p.container {
			return p.host
		}
		if strings.HasPrefix(clean, p.container+"/") {
			rel := strings.TrimPrefix(clean, p.container+"/")
			return filepath.Join(p.host, filepath.FromSlash(rel))
		}
	}
	return path
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

	snap, ch := entry.SubscribeWithSnapshot()
	defer entry.Unsubscribe(ch)

	sendEvent := func(typ, payload string) {
		// 使用 json.Marshal 以避免 Go %q 在含非 ASCII/控制字符时产出
		// 非 JSON 兼容的 \xNN 转义。
		pb, _ := json.Marshal(payload)
		fmt.Fprintf(w, "data: {\"type\":%q,\"payload\":%s}\n\n", typ, pb)
		if canFlush {
			flusher.Flush()
		}
	}
	sendSnapshot := func(lines []string) {
		// 一次性把已缓冲的所有日志作为单个事件发送，避免 N 条日志触发
		// N 次浏览器端 JSON.parse / DOM 写入，导致首屏卡住。
		b, _ := json.Marshal(map[string]any{"type": "snapshot", "payload": lines})
		fmt.Fprintf(w, "data: %s\n\n", b)
		if canFlush {
			flusher.Flush()
		}
	}

	sendSnapshot(snap)

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

// ─── 新增数据库管理API handlers ─────────────────────────────────────────────

func (s *Server) handleDBGenerationRuns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListGenerationRuns(r.Context(), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBGeneratedCases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	q := r.URL.Query()
	rows, err := db.ListGeneratedCases(r.Context(), q.Get("run_id"), q.Get("model"), q.Get("language"), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBPromptRenderings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	q := r.URL.Query()
	rows, err := db.ListPromptRenderings(r.Context(), q.Get("run_id"), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBEvaluationRuns(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	q := r.URL.Query()
	rows, err := db.ListEvaluationRuns(r.Context(), q.Get("run_id"), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBEvaluationStages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	q := r.URL.Query()
	rows, err := db.ListEvaluationStages(r.Context(), q.Get("evaluation_run_id"), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBDatasetSamples(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	q := r.URL.Query()
	rows, err := db.ListDatasetSamples(r.Context(), q.Get("language"), q.Get("class"), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBDatasetSnapshots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListDatasetSnapshots(r.Context(), parseLimit(r, 50))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBModelConfigs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListModelConfigs(r.Context(), parseLimit(r, 50))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBPromptProfiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListPromptProfiles(r.Context(), parseLimit(r, 20))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBEvaluationEnvs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListEvaluationEnvs(r.Context(), parseLimit(r, 20))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBScorePolicies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListScorePolicies(r.Context(), parseLimit(r, 10))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	q := r.URL.Query()
	rows, err := db.ListReports(r.Context(), q.Get("run_id"), parseLimit(r, 50))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBRunArtifacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	q := r.URL.Query()
	rows, err := db.ListRunArtifacts(r.Context(), q.Get("run_id"), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBExperiments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListExperiments(r.Context(), parseLimit(r, 20))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}
