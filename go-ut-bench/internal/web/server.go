package web

import (
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
	"go-ut-bench/internal/orchestrator"

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
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
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
	s.mux.HandleFunc("/api/runs", s.handleRuns)
	s.mux.HandleFunc("/api/runs/", s.handleRunSub)

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
	default:
		s.handleRunGet(w, r, runID)
	}
}

type runDetailResponse struct {
	RunID     string            `json:"run_id"`
	Status    RunStatus         `json:"status"`
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
		errJSON(w, http.StatusNotFound, "run not found: "+runID)
		return
	}
	entry.mu.RLock()
	resp := runDetailResponse{
		RunID:     entry.RunID,
		Status:    entry.Status,
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
