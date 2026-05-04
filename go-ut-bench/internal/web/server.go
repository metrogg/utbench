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
	"sync"
	"time"

	"go-ut-bench/internal/agentconfig"
	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/obs"
	"go-ut-bench/internal/orchestrator"
	"go-ut-bench/internal/reporter"
	"go-ut-bench/internal/runner"
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

	// 缓存层：避免重复读磁盘/解析YAML
	cacheMu       sync.RWMutex
	catalogCache  *catalogCacheEntry
	runsCache     *runsCacheEntry
	envCheckCache *envCheckCacheEntry
}

type catalogCacheEntry struct {
	data      webCatalog
	loadedAt  time.Time
	configMod time.Time // models.yaml 的 mtime，用于失效判断
}

type runsCacheEntry struct {
	data     []runSummaryItem
	loadedAt time.Time
}

type envCheckCacheEntry struct {
	data     environmentCheckResponse
	loadedAt time.Time
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

// Close cancels all running tasks, cleans up Docker containers, and closes the database connection.
func (s *Server) Close() error {
	// 取消所有活跃任务
	for _, entry := range s.mgr.List() {
		entry.mu.RLock()
		status := entry.Status
		runID := entry.RunID
		useDocker := entry.UseDocker
		entry.mu.RUnlock()
		if status == StatusRunning || status == StatusPending {
			_ = s.mgr.Cancel(runID)
		}
		// 清理 sandbox 子容器
		if useDocker {
			killSandboxContainers(runID)
		}
	}
	// 清理所有可能残留的 utbench sandbox 容器（兜底）
	killAllUtbenchSandboxes()
	// 停止 RunManager 后台清理 goroutine
	s.mgr.Close()
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
	s.mux.HandleFunc("/api/settings/api-keys", s.handleAPIKeys)
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
	s.mux.HandleFunc("/api/db/asset-subjects", s.handleDBAssetSubjects)
	s.mux.HandleFunc("/api/db/subject-versions", s.handleDBSubjectVersions)
	s.mux.HandleFunc("/api/db/asset-generations", s.handleDBAssetGenerations)
	s.mux.HandleFunc("/api/db/asset-evaluations", s.handleDBAssetEvaluations)
	s.mux.HandleFunc("/api/db/asset-explain-reuse", s.handleDBAssetExplainReuse)
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
	// DedupMode 去重模式：
	// - "merge" (默认): 合并所有结果，同名样本可能有多条记录
	// - "overwrite": 按 (model, language, sample_id) 去重，保留最新 run_id 的结果
	DedupMode string `json:"dedup_mode,omitempty"`
}

// dbReportSummary 是写入 _db_report 目录下 run_summary.json 的结构，
// 使合并报告能被 /api/runs 发现并显示在 Web UI 中。
type dbReportSummary struct {
	RunID           string            `json:"run_id"`
	CreatedAtUTC    string            `json:"created_at_utc"`
	SchemaVersion   string            `json:"schema_version"`
	Phase           string            `json:"phase"`
	SourceRunIDs    []string          `json:"source_run_ids"`
	ResultCount     int               `json:"result_count"`
	ReportJSONPath  string            `json:"report_json_path"`
	ReportHTMLPath  string            `json:"report_html_path"`
	IsMergedReport  bool              `json:"is_merged_report"`
	Spec            contracts.RunSpec `json:"spec"`
}

// dedupResultsByLatestRun 按 (model, language, sample_id) 去重，保留最新 run_id 的结果。
// 用于 "overwrite" 模式，确保同一模型+语言+样本只保留最新运行的结果。
func dedupResultsByLatestRun(results []contracts.EvaluationResult) []contracts.EvaluationResult {
	type dedupKey struct {
		Model    string
		Language string
		SampleID string
	}
	// 按 run_id 排序（run_id 是时间戳格式，字典序即时间序），后面的会覆盖前面的
	seen := make(map[dedupKey]int) // key -> index in result slice
	for i, r := range results {
		key := dedupKey{Model: r.Model, Language: r.Language, SampleID: r.SampleID}
		if prevIdx, exists := seen[key]; exists {
			// 比较 run_id，保留最新的
			if r.RunID > results[prevIdx].RunID {
				seen[key] = i
			}
		} else {
			seen[key] = i
		}
	}
	deduped := make([]contracts.EvaluationResult, 0, len(seen))
	for _, idx := range seen {
		deduped = append(deduped, results[idx])
	}
	return deduped
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

	// 增量合并：如果指定了 run_id 且该目录下已有 run_summary.json，
	// 读取其 source_run_ids 与本次新选的合并（去重）。
	allSourceRunIDs := append([]string{}, req.SourceRunIDs...)
	summaryPath := filepath.Join(s.outputRoot, "runs", outRunID, "run_summary.json")
	if existing, err := os.ReadFile(summaryPath); err == nil {
		var prev dbReportSummary
		if json.Unmarshal(existing, &prev) == nil && len(prev.SourceRunIDs) > 0 {
			seen := make(map[string]bool, len(allSourceRunIDs))
			for _, id := range allSourceRunIDs {
				seen[id] = true
			}
			for _, id := range prev.SourceRunIDs {
				if !seen[id] {
					allSourceRunIDs = append(allSourceRunIDs, id)
					seen[id] = true
				}
			}
		}
	}

	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	set, err := db.SelectEvaluationResultSet(r.Context(), outRunID, store.DBReportFilter{
		RunIDs:            allSourceRunIDs,
		EvaluationRunIDs:  req.EvaluationRunIDs,
		Models:            req.Models,
		Languages:         req.Languages,
		ScoreEligibleOnly: req.ScoreEligibleOnly,
	})
	if err != nil {
		errJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	// 按去重模式处理结果
	dedupMode := strings.TrimSpace(req.DedupMode)
	if dedupMode == "" {
		dedupMode = "merge"
	}
	if dedupMode == "overwrite" {
		set.Results = dedupResultsByLatestRun(set.Results)
	}

	logger := obs.NewLogger(false, filepath.Join(s.outputRoot, "runs", outRunID, "logs"))
	out, err := reporter.NewService(logger, &runner.DefaultPromptMetaProvider{}).GenerateFromResultSet(contracts.RunSpec{
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

	// 写入 run_summary.json，使合并报告出现在 /api/runs 列表中。
	// 收集涉及的模型和语言。
	modelSet := make(map[string]bool)
	langSet := make(map[string]bool)
	for _, res := range set.Results {
		modelSet[res.Model] = true
		langSet[res.Language] = true
	}
	models := make([]string, 0, len(modelSet))
	for m := range modelSet {
		models = append(models, m)
	}
	languages := make([]string, 0, len(langSet))
	for l := range langSet {
		languages = append(languages, l)
	}
	sort.Strings(models)
	sort.Strings(languages)

	summary := dbReportSummary{
		RunID:          outRunID,
		CreatedAtUTC:   time.Now().UTC().Format(time.RFC3339Nano),
		SchemaVersion:  "v0.1.0",
		Phase:          "report",
		SourceRunIDs:   allSourceRunIDs,
		ResultCount:    len(set.Results),
		ReportJSONPath: out.ReportJSONPath,
		ReportHTMLPath: out.ReportHTMLPath,
		IsMergedReport: true,
		Spec: contracts.RunSpec{
			RunID:        outRunID,
			Models:       models,
			Languages:    languages,
			OutputRoot:   s.outputRoot,
			ConfigPath:   s.configPath,
			CreatedAtUTC: time.Now().UTC(),
		},
	}
	if data, err := json.MarshalIndent(summary, "", "  "); err == nil {
		_ = os.MkdirAll(filepath.Dir(summaryPath), 0o755)
		_ = os.WriteFile(summaryPath, data, 0o644)
	}

	// 使 runs 缓存失效，下次 loadRuns 能立即看到新合并报告。
	s.cacheMu.Lock()
	s.runsCache = nil
	s.cacheMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]any{
		"run_id":           outRunID,
		"result_count":     len(set.Results),
		"report_json_path": out.ReportJSONPath,
		"report_html_path": out.ReportHTMLPath,
		"source_run_ids":   allSourceRunIDs,
	})
}

// ─── GET /api/config ─────────────────────────────────────────────────────────

type modelInfo struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	ModelID  string `json:"model_id"`
	Enabled  bool   `json:"enabled"`
}

type frameworkInfo struct {
	Name                string   `json:"name"`
	Kind                string   `json:"kind"`
	SandboxMode         string   `json:"sandbox_mode,omitempty"`
	CompatibleModels    []string `json:"compatible_models,omitempty"`
	CompatibleLanguages []string `json:"compatible_languages,omitempty"`
}

type skillInfo struct {
	Name                 string   `json:"name"`
	Version              string   `json:"version,omitempty"`
	InjectMode           string   `json:"inject_mode,omitempty"`
	CompatibleFrameworks []string `json:"compatible_frameworks,omitempty"`
	CompatibleLanguages  []string `json:"compatible_languages,omitempty"`
}

type subjectInfo struct {
	ID          string   `json:"id"`
	Kind        string   `json:"kind"`
	Framework   string   `json:"framework"`
	Model       string   `json:"model"`
	Skill       string   `json:"skill"`
	SandboxMode string   `json:"sandbox_mode,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type configResponse struct {
	Models             []modelInfo     `json:"models"`
	Frameworks         []frameworkInfo `json:"frameworks,omitempty"`
	Skills             []skillInfo     `json:"skills,omitempty"`
	Subjects           []subjectInfo   `json:"subjects,omitempty"`
	Languages          []string        `json:"languages"`
	Scenarios          []string        `json:"scenarios"`
	Classes            []string        `json:"classes"`
	DatasetRoot        string          `json:"dataset_root"`
	ConfigPath         string          `json:"config_path"`
	AgentsConfigPath   string          `json:"agents_config_path,omitempty"`
	AgentsConfigError  string          `json:"agents_config_error,omitempty"`
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

type webCatalog struct {
	models            []modelInfo
	frameworks        []frameworkInfo
	skills            []skillInfo
	subjects          []subjectInfo
	agentsConfigError string
}

// loadWebCatalogCached 带缓存的 catalog 加载，2秒内复用。
// 当 models.yaml 文件 mtime 变化时自动失效。
func (s *Server) loadWebCatalogCached() (webCatalog, error) {
	const cacheTTL = 2 * time.Second

	// 获取 config 文件 mtime
	info, err := os.Stat(s.configPath)
	if err != nil {
		return s.loadWebCatalog()
	}
	modTime := info.ModTime()

	s.cacheMu.RLock()
	if s.catalogCache != nil && time.Since(s.catalogCache.loadedAt) < cacheTTL && s.catalogCache.configMod == modTime {
		cat := s.catalogCache.data
		s.cacheMu.RUnlock()
		return cat, nil
	}
	s.cacheMu.RUnlock()

	cat, err := s.loadWebCatalog()
	if err != nil {
		return cat, err
	}

	s.cacheMu.Lock()
	s.catalogCache = &catalogCacheEntry{data: cat, loadedAt: time.Now(), configMod: modTime}
	s.cacheMu.Unlock()
	return cat, nil
}

func (s *Server) loadWebCatalog() (webCatalog, error) {
	raw, err := os.ReadFile(s.configPath)
	if err != nil {
		return webCatalog{}, fmt.Errorf("cannot read models.yaml: %w", err)
	}
	var mf modelsYAML
	_ = yaml.Unmarshal(raw, &mf)

	models := make([]modelInfo, 0, len(mf.Models))
	modelNames := make([]string, 0, len(mf.Models))
	for name, m := range mf.Models {
		models = append(models, modelInfo{
			Name:     name,
			Provider: m.Provider,
			ModelID:  m.Config.Model,
			Enabled:  m.Enabled,
		})
		if m.Enabled {
			modelNames = append(modelNames, name)
		}
	}
	sort.Slice(models, func(i, j int) bool { return models[i].Name < models[j].Name })
	sort.Strings(modelNames)

	catalog := webCatalog{models: models}
	agentsConfigPath := strings.TrimSpace(s.mgr.agentsConfigPath)
	if agentsConfigPath == "" {
		catalog.agentsConfigError = "未配置 agents config 路径（--agents-config）"
		return catalog, nil
	}
	resolved, err := agentconfig.Load(agentsConfigPath, modelNames, nil)
	if err != nil {
		catalog.agentsConfigError = fmt.Sprintf("加载 agents config 失败: %v", err)
		fmt.Printf("[web] agents config error: %v\n", err)
		return catalog, nil
	}
	fmt.Printf("[web] agents config loaded: %d subjects, models=%v\n", len(resolved), modelNames)
	frameworkSeen := map[string]frameworkInfo{}
	skillSeen := map[string]skillInfo{}
	subjects := make([]subjectInfo, 0, len(resolved))
	for _, item := range resolved {
		if item.Framework.Name != "" {
			frameworkSeen[item.Framework.Name] = frameworkInfo{
				Name:                item.Framework.Name,
				Kind:                item.Framework.Kind,
				SandboxMode:         item.Framework.SandboxMode,
				CompatibleModels:    append([]string{}, item.Framework.CompatibleModels...),
				CompatibleLanguages: append([]string{}, item.Framework.CompatibleLangs...),
			}
		}
		if item.Skill.Name != "" {
			skillSeen[item.Skill.Name] = skillInfo{
				Name:                 item.Skill.Name,
				Version:              item.Skill.Version,
				InjectMode:           item.Skill.InjectMode,
				CompatibleFrameworks: append([]string{}, item.Skill.CompatibleFrameworks...),
				CompatibleLanguages:  append([]string{}, item.Skill.CompatibleLanguages...),
			}
		}
		subjects = append(subjects, subjectInfo{
			ID:          item.Spec.ID,
			Kind:        item.Spec.Kind,
			Framework:   item.Spec.Framework,
			Model:       item.Spec.Model,
			Skill:       item.Spec.Skill,
			SandboxMode: item.Framework.SandboxMode,
			Labels:      append([]string{}, item.Spec.Labels...),
			Tags:        append([]string{}, item.Spec.Tags...),
		})
	}
	catalog.frameworks = make([]frameworkInfo, 0, len(frameworkSeen))
	for _, v := range frameworkSeen {
		catalog.frameworks = append(catalog.frameworks, v)
	}
	sort.Slice(catalog.frameworks, func(i, j int) bool { return catalog.frameworks[i].Name < catalog.frameworks[j].Name })
	catalog.skills = make([]skillInfo, 0, len(skillSeen))
	for _, v := range skillSeen {
		catalog.skills = append(catalog.skills, v)
	}
	sort.Slice(catalog.skills, func(i, j int) bool { return catalog.skills[i].Name < catalog.skills[j].Name })
	sort.Slice(subjects, func(i, j int) bool { return subjects[i].ID < subjects[j].ID })
	catalog.subjects = subjects
	return catalog, nil
}

func deriveModelsFromSubjects(subjectIDs []string, subjects []subjectInfo) []string {
	if len(subjectIDs) == 0 {
		return nil
	}
	selected := map[string]struct{}{}
	for _, id := range subjectIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			selected[id] = struct{}{}
		}
	}
	modelSet := map[string]struct{}{}
	for _, subject := range subjects {
		if _, ok := selected[subject.ID]; ok && strings.TrimSpace(subject.Model) != "" {
			modelSet[subject.Model] = struct{}{}
		}
	}
	out := make([]string, 0, len(modelSet))
	for model := range modelSet {
		out = append(out, model)
	}
	sort.Strings(out)
	return out
}

func subjectRequiresDockerSandbox(subjectIDs []string, subjects []subjectInfo) bool {
	if len(subjectIDs) == 0 {
		return false
	}
	selected := map[string]struct{}{}
	for _, id := range subjectIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			selected[id] = struct{}{}
		}
	}
	for _, subject := range subjects {
		if _, ok := selected[subject.ID]; !ok {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(subject.SandboxMode), "docker") {
			return true
		}
	}
	return false
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	catalog, err := s.loadWebCatalogCached()
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, configResponse{
		Models:            catalog.models,
		Frameworks:        catalog.frameworks,
		Skills:            catalog.skills,
		Subjects:          catalog.subjects,
		Languages:         contracts.SupportedLanguages,
		Scenarios:         contracts.SupportedScenarios,
		Classes:           []string{"self_contained", "repo_level"},
		DatasetRoot:       s.mgr.datasetRoot,
		ConfigPath:        s.configPath,
		AgentsConfigPath:  s.mgr.agentsConfigPath,
		AgentsConfigError: catalog.agentsConfigError,
	})
}

// ─── /api/runs ───────────────────────────────────────────────────────────────

type runSummaryItem struct {
	RunID          string            `json:"run_id"`
	Status         RunStatus         `json:"status"`
	Paused         bool              `json:"paused,omitempty"`
	StartedAt      time.Time         `json:"started_at"`
	EndedAt        *time.Time        `json:"ended_at,omitempty"`
	Error          string            `json:"error,omitempty"`
	Spec           contracts.RunSpec `json:"spec"`
	UseDocker      bool              `json:"use_docker,omitempty"`
	IsMergedReport bool              `json:"is_merged_report,omitempty"`
	SourceRunIDs   []string          `json:"source_run_ids,omitempty"`
	ResultCount    int               `json:"result_count,omitempty"`
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
	// 快速路径：如果全部在内存中（活跃任务），直接返回，不扫描磁盘。
	// 只有当有已完成任务需要从 artifacts/ 恢复时才做文件系统扫描（带缓存）。
	const cacheTTL = 10 * time.Second

	activeRuns := s.mgr.List()
	hasActive := len(activeRuns) > 0

	// 检查缓存
	s.cacheMu.RLock()
	cached := s.runsCache
	s.cacheMu.RUnlock()

	if cached != nil && time.Since(cached.loadedAt) < cacheTTL {
		// 合并缓存的已完成任务 + 实时的活跃任务
		merged := make(map[string]runSummaryItem)
		for _, item := range cached.data {
			merged[item.RunID] = item
		}
		for _, entry := range activeRuns {
			entry.mu.RLock()
			merged[entry.RunID] = runSummaryItem{
				RunID:     entry.RunID,
				Status:    entry.Status,
				Paused:    entry.Paused,
				StartedAt: entry.StartedAt,
				EndedAt:   entry.EndedAt,
				Error:     entry.Error,
				Spec:      entry.Spec,
			}
			entry.mu.RUnlock()
		}
		out := make([]runSummaryItem, 0, len(merged))
		for _, v := range merged {
			out = append(out, v)
		}
		sort.Slice(out, func(i, j int) bool {
			return out[i].StartedAt.After(out[j].StartedAt)
		})
		writeJSON(w, http.StatusOK, out)
		return
	}

	// 缓存失效，重新扫描
	s.listRunsFromDisk(w, activeRuns, hasActive)
}

func (s *Server) listRunsFromDisk(w http.ResponseWriter, activeRuns []*RunEntry, updateCache bool) {
	byID := make(map[string]runSummaryItem)

	// Scan artifacts/runs/*/run_summary.json
	pattern := filepath.Join(s.outputRoot, "runs", "*", "run_summary.json")
	matches, _ := filepath.Glob(pattern)
	for _, path := range matches {
		var raw struct {
			RunID          string            `json:"run_id"`
			CreatedAtUTC   string            `json:"created_at_utc"`
			Spec           contracts.RunSpec `json:"spec"`
			Backend        string            `json:"backend"`
			IsMergedReport bool              `json:"is_merged_report"`
			SourceRunIDs   []string          `json:"source_run_ids"`
			ResultCount    int               `json:"result_count"`
		}
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(data, &raw); err != nil {
			continue
		}
		runID := strings.TrimSpace(raw.RunID)
		if runID == "" {
			continue
		}
		if _, exists := byID[runID]; !exists {
			t, _ := time.Parse(time.RFC3339Nano, raw.CreatedAtUTC)
			byID[runID] = runSummaryItem{
				RunID:          runID,
				Status:         StatusCompleted,
				StartedAt:      t,
				Spec:           raw.Spec,
				UseDocker:      strings.EqualFold(strings.TrimSpace(raw.Backend), "docker"),
				IsMergedReport: raw.IsMergedReport,
				SourceRunIDs:   raw.SourceRunIDs,
				ResultCount:    raw.ResultCount,
			}
		}
	}

	// Override / add active in-memory runs
	for _, entry := range activeRuns {
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

	// 缓存已完成任务（不含活跃任务的实时状态）
	diskItems := make([]runSummaryItem, 0, len(byID))
	for _, v := range byID {
		diskItems = append(diskItems, v)
	}
	s.cacheMu.Lock()
	s.runsCache = &runsCacheEntry{data: diskItems, loadedAt: time.Now()}
	s.cacheMu.Unlock()

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
	Subjects        []string `json:"subjects,omitempty"`
	Languages       []string `json:"languages"`
	Class           string   `json:"class"`
	Scenario        string   `json:"scenario"`
	Level           string   `json:"level"`
	MaxSamples      int      `json:"max_samples"`
	Workers         int      `json:"workers"`
	Mode            string   `json:"mode"`
	DryRun          bool     `json:"dry_run"`
	ReuseGenerated  bool     `json:"reuse_generated"`
	ReuseEvaluation bool     `json:"reuse_evaluation"`
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
	catalog, err := s.loadWebCatalogCached()
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	models := append([]string{}, req.Models...)
	if len(models) == 0 && len(req.Subjects) > 0 {
		models = deriveModelsFromSubjects(req.Subjects, catalog.subjects)
	}
	if len(req.Subjects) > 0 && strings.TrimSpace(s.mgr.agentsConfigPath) == "" {
		errJSON(w, http.StatusBadRequest, "subjects requires agents config")
		return
	}
	if (phase == "full" || phase == "generate") && len(models) == 0 {
		errJSON(w, http.StatusBadRequest, "models or subjects is required for generate/full phase")
		return
	}
	if !req.UseDocker && subjectRequiresDockerSandbox(req.Subjects, catalog.subjects) {
		outputRootSlash := filepath.ToSlash(s.outputRoot)
		if strings.HasPrefix(outputRootSlash, "/app/") {
			if strings.TrimSpace(os.Getenv("UTBENCH_SANDBOX_HOST_OUTPUT_ROOT")) == "" {
				errJSON(w, http.StatusBadRequest, "当前 Web 运行在容器内，且选择了 sandbox_mode=docker 的 subject，但未设置 UTBENCH_SANDBOX_HOST_OUTPUT_ROOT")
				return
			}
			if _, err := os.Stat("/var/run/docker.sock"); err != nil {
				errJSON(w, http.StatusBadRequest, "当前 Web 运行在容器内，且选择了 sandbox_mode=docker 的 subject，但未挂载 /var/run/docker.sock")
				return
			}
		}
	}
	agentsConfigPath := ""
	if len(req.Subjects) > 0 {
		agentsConfigPath = s.mgr.agentsConfigPath
	}

	spec := contracts.RunSpec{
		RunID:            runID,
		Models:           models,
		Subjects:         splitTrim(strings.Join(req.Subjects, ",")),
		AgentsConfigPath: agentsConfigPath,
		Languages:        req.Languages,
		DatasetClasses:   classes,
		DatasetScenario:  req.Scenario,
		DatasetLevel:     req.Level,
		DatasetRoot:      s.mgr.datasetRoot,
		ConfigPath:       s.configPath,
		Mode:             contracts.RunMode(mode),
		DryRun:           req.DryRun,
		ReuseGenerated:   req.ReuseGenerated,
		ReuseEvaluation:  req.ReuseEvaluation,
		DBPath:           s.mgr.dbPath,
		MutationEnabled:  req.MutationEnabled,
		MutationTimeout:  mutTimeout,
		MutationPolicy:   mutPolicy,
		MaxSamples:       req.MaxSamples,
		Workers:          req.Workers,
		OutputRoot:       s.outputRoot,
		CreatedAtUTC:     time.Now().UTC(),
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
	// 新任务提交后清除 runs 缓存
	s.cacheMu.Lock()
	s.runsCache = nil
	s.cacheMu.Unlock()
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
	s.cacheMu.Lock()
	s.runsCache = nil
	s.cacheMu.Unlock()
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
	// 异步执行评测，避免阻塞 HTTP 响应。使用 background context 确保客户端断开不会取消任务。
	go func() {
		out, err := runEvaluateInDocker(context.Background(), runID, spec, s.dockerCfg)
		if err != nil {
			fmt.Printf("[reevaluate] run=%s failed: %v\n%s\n", runID, err, tailString(string(out), 500))
			return
		}
		fmt.Printf("[reevaluate] run=%s completed\n", runID)
	}()
	writeJSON(w, http.StatusAccepted, map[string]any{
		"run_id":        runID,
		"manifest_path": manifestPath,
		"status":        "reevaluation_started",
		"message":       "Re-evaluation is running in the background. Check run status for completion.",
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
	out, err := reporter.NewService(logger, &runner.DefaultPromptMetaProvider{}).Generate(context.Background(), spec, evaluationPath)
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

func (s *Server) handleDBAssetSubjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListAssetSubjects(r.Context(), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBSubjectVersions(w http.ResponseWriter, r *http.Request) {
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
	rows, err := db.ListSubjectVersions(r.Context(), q.Get("subject"), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBAssetGenerations(w http.ResponseWriter, r *http.Request) {
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
	rows, err := db.ListAssetGenerations(r.Context(), q.Get("subject"), q.Get("language"), q.Get("sample"), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBAssetEvaluations(w http.ResponseWriter, r *http.Request) {
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
	rows, err := db.ListAssetEvaluations(r.Context(), q.Get("subject"), q.Get("language"), q.Get("sample"), parseLimit(r, 100))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) handleDBAssetExplainReuse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errJSON(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	q := r.URL.Query()
	subjectID := strings.TrimSpace(q.Get("subject"))
	lang := strings.TrimSpace(q.Get("language"))
	sampleID := strings.TrimSpace(q.Get("sample"))
	if subjectID == "" || lang == "" || sampleID == "" {
		errJSON(w, http.StatusBadRequest, "subject, language, sample are required")
		return
	}
	db, err := s.openStore(r.Context())
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	rows, err := db.ListAssetGenerations(r.Context(), subjectID, lang, sampleID, parseLimit(r, 20))
	if err != nil {
		errJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := map[string]any{
		"subject_id":  subjectID,
		"language":    lang,
		"sample_id":   sampleID,
		"matched":     false,
		"miss_reason": "no_successful_generation_asset",
		"candidates":  rows,
	}
	for _, row := range rows {
		if row.Success && row.GeneratedTestPath != "" && row.GenerationKey != "" {
			response["matched"] = true
			response["miss_reason"] = ""
			response["generation_key"] = row.GenerationKey
			response["latest_reusable"] = row
			response["comparisons"] = map[string]string{
				"stored_subject_version_id":  row.SubjectVersionID,
				"stored_sandbox_fingerprint": row.SandboxFingerprint,
				"stored_sample_uid":          row.SampleUID,
				"dependency_fingerprint":     row.DependencyFingerprint,
				"generation_env_fingerprint": row.GenerationEnvFingerprint,
			}
			break
		}
	}
	writeJSON(w, http.StatusOK, response)
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
