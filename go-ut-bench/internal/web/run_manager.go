package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/dataset"
	"go-ut-bench/internal/evaluator"
	"go-ut-bench/internal/obs"
	"go-ut-bench/internal/orchestrator"
	"go-ut-bench/internal/reporter"
	"go-ut-bench/internal/runner"
)

// RunStatus represents the lifecycle state of a benchmark run.
type RunStatus string

const (
	StatusPending   RunStatus = "pending"
	StatusRunning   RunStatus = "running"
	StatusCompleted RunStatus = "completed"
	StatusFailed    RunStatus = "failed"
)

// RunEntry holds in-memory state for a single benchmark run.
type RunEntry struct {
	RunID     string            `json:"run_id"`
	Status    RunStatus         `json:"status"`
	StartedAt time.Time         `json:"started_at"`
	EndedAt   *time.Time        `json:"ended_at,omitempty"`
	Error     string            `json:"error,omitempty"`
	Spec      contracts.RunSpec `json:"spec"`
	// UseDocker marks this run as having been executed by shelling out to
	// `docker run utbench:latest ...` instead of running the orchestrator
	// in-process. Mirrors the createRunRequest flag and is surfaced back to
	// the UI so the detail view can label how the run was executed.
	UseDocker bool `json:"use_docker,omitempty"`

	logs []string
	mu   sync.RWMutex
	subs []chan string
	Done chan struct{}
}

func (r *RunEntry) appendLog(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs = append(r.logs, line)
	for _, ch := range r.subs {
		select {
		case ch <- line:
		default:
		}
	}
}

// GetLogs returns a snapshot of all captured log lines.
func (r *RunEntry) GetLogs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cp := make([]string, len(r.logs))
	copy(cp, r.logs)
	return cp
}

// Subscribe returns a channel that receives new log lines.
// Existing buffered lines are replayed first.
func (r *RunEntry) Subscribe() chan string {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan string, 512)
	for _, line := range r.logs {
		select {
		case ch <- line:
		default:
		}
	}
	r.subs = append(r.subs, ch)
	return ch
}

// Unsubscribe removes a subscriber channel and closes it.
func (r *RunEntry) Unsubscribe(ch chan string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, s := range r.subs {
		if s == ch {
			r.subs = append(r.subs[:i], r.subs[i+1:]...)
			close(ch)
			return
		}
	}
}

// lineWriter implements io.Writer: buffers bytes and forwards complete lines to the run.
type lineWriter struct {
	run *RunEntry
	buf []byte
	mu  sync.Mutex
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buf = append(w.buf, p...)
	for {
		idx := bytes.IndexByte(w.buf, '\n')
		if idx < 0 {
			break
		}
		line := string(w.buf[:idx])
		w.buf = w.buf[idx+1:]
		if line != "" {
			w.run.appendLog(line)
		}
	}
	return len(p), nil
}

// RunManager manages the lifecycle of all benchmark runs.
// It supports two execution backends:
//   - In-process: orchestrator.Run() called directly (default, fast path).
//   - Docker:    `docker run utbench:latest ...` forked as a child process
//     so the evaluation runs in a fully-provisioned Linux container
//     (Windows mutmut, mull, gremlins etc. all work there).
type RunManager struct {
	mu          sync.RWMutex
	runs        map[string]*RunEntry
	configPath  string
	datasetRoot string
	outputRoot  string
	dbPath      string
	// dockerCfg is used when a run is submitted with UseDocker=true.
	dockerCfg DockerConfig
}

// NewRunManager creates a RunManager with the given default paths.
// imageName/projectRoot/envFile provide defaults for Docker execution; they
// may be empty if Docker mode is not supported in the user's environment.
func NewRunManager(configPath, datasetRoot, outputRoot, dbPath string, cfg DockerConfig) *RunManager {
	return &RunManager{
		runs:        make(map[string]*RunEntry),
		configPath:  configPath,
		datasetRoot: datasetRoot,
		outputRoot:  outputRoot,
		dbPath:      dbPath,
		dockerCfg:   cfg,
	}
}

// Submit enqueues and immediately starts a run in a goroutine.
// useDocker selects the Docker backend; see RunManager doc.
func (m *RunManager) Submit(spec contracts.RunSpec, opts orchestrator.Options, useDocker bool) *RunEntry {
	entry := &RunEntry{
		RunID:     spec.RunID,
		Status:    StatusPending,
		StartedAt: time.Now(),
		Spec:      spec,
		UseDocker: useDocker,
		Done:      make(chan struct{}),
	}
	m.mu.Lock()
	m.runs[spec.RunID] = entry
	m.mu.Unlock()

	go m.execute(entry, spec, opts)
	return entry
}

func (m *RunManager) execute(entry *RunEntry, spec contracts.RunSpec, opts orchestrator.Options) {
	defer close(entry.Done)

	entry.mu.Lock()
	entry.Status = StatusRunning
	entry.mu.Unlock()

	entry.appendLog(fmt.Sprintf("[%s] run started  id=%s  backend=%s",
		logTS(), spec.RunID, backendLabel(entry.UseDocker)))
	entry.appendLog(fmt.Sprintf("[%s] models=%v  langs=%v  dry_run=%v  max_samples=%d  workers=%d",
		logTS(), spec.Models, spec.Languages, spec.DryRun, spec.MaxSamples, spec.Workers))

	ctx := context.Background()
	var err error
	if entry.UseDocker {
		err = runInDocker(ctx, entry, spec, opts, m.dockerCfg)
	} else {
		err = m.executeInProcess(ctx, entry, spec, opts)
	}

	now := time.Now()
	entry.mu.Lock()
	entry.EndedAt = &now
	if err != nil {
		entry.Status = StatusFailed
		entry.Error = err.Error()
	} else {
		entry.Status = StatusCompleted
	}
	entry.mu.Unlock()

	if err != nil {
		entry.appendLog(fmt.Sprintf("[%s] FAILED: %v", logTS(), err))
	} else {
		entry.appendLog(fmt.Sprintf("[%s] run completed", logTS()))
	}
}

// executeInProcess runs the orchestrator in the same Go process.
// Kept separate so the Docker path has no unused imports and stays testable.
func (m *RunManager) executeInProcess(ctx context.Context, entry *RunEntry, spec contracts.RunSpec, opts orchestrator.Options) error {
	lw := &lineWriter{run: entry}
	mw := io.MultiWriter(os.Stderr, lw)
	logger := obs.NewLoggerWithWriter(true, mw)

	ds := dataset.NewService()
	rn := runner.NewService(logger)
	ev := evaluator.NewService(logger)
	rp := reporter.NewService(logger)
	orch := orchestrator.New(ds, rn, ev, rp)

	_, err := orch.Run(ctx, spec, opts)
	return err
}

func backendLabel(useDocker bool) string {
	if useDocker {
		return "docker"
	}
	return "in-process"
}

// Get returns the RunEntry for the given runID.
func (m *RunManager) Get(runID string) (*RunEntry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.runs[runID]
	return r, ok
}

// List returns all known RunEntries (unordered).
func (m *RunManager) List() []*RunEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*RunEntry, 0, len(m.runs))
	for _, r := range m.runs {
		out = append(out, r)
	}
	return out
}

func logTS() string {
	return time.Now().Format("15:04:05.000")
}
