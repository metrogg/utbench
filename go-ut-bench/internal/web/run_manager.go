package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/ctrl"
	"go-ut-bench/internal/dataset"
	"go-ut-bench/internal/evaluator"
	"go-ut-bench/internal/obs"
	"go-ut-bench/internal/orchestrator"
	"go-ut-bench/internal/reporter"
	"go-ut-bench/internal/runner"
	"go-ut-bench/internal/store"
)

// RunStatus represents the lifecycle state of a benchmark run.
type RunStatus string

const (
	StatusPending   RunStatus = "pending"
	StatusRunning   RunStatus = "running"
	StatusPaused    RunStatus = "paused"
	StatusCompleted RunStatus = "completed"
	StatusFailed    RunStatus = "failed"
	StatusCanceled  RunStatus = "canceled"
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

	// Paused 反映任务是否被用户请求挂起（in-process 模式下由 Gate 实现，
	// Docker 模式下通过 `docker pause/unpause` 实现）。
	// Status 在暂停期间保持 "running"；前端需叠加 Paused 才能显示"已暂停"。
	// 此处单独字段保持 Status 单一职责，避免 paused→running 逻辑到处散落。
	Paused bool `json:"paused,omitempty"`

	logs []string
	mu   sync.RWMutex
	subs []chan string
	Done chan struct{}

	// 运行时控制（不序列化）。
	gate      *ctrl.ChanGate
	cancel    context.CancelFunc
	container string // docker 模式下的容器名，用于 docker pause/unpause/kill
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
//     (Windows mutmut, mull, go-mutesting etc. all work there).
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
		gate:      ctrl.NewChanGate(),
		container: "utbench-" + spec.RunID, // 用于 docker pause/unpause/kill 的稳定名
	}
	m.mu.Lock()
	m.runs[spec.RunID] = entry
	m.mu.Unlock()

	go m.execute(entry, spec, opts)
	return entry
}

func (m *RunManager) execute(entry *RunEntry, spec contracts.RunSpec, opts orchestrator.Options) {
	defer close(entry.Done)

	ctx, cancel := context.WithCancel(context.Background())
	entry.mu.Lock()
	entry.Status = StatusRunning
	entry.cancel = cancel
	entry.mu.Unlock()
	defer cancel()

	// 将 gate 绑定到 ctx，runner/evaluator 的 worker 将在每个任务前调用 ctrl.Wait。
	ctx = ctrl.WithGate(ctx, entry.gate)

	entry.appendLog(fmt.Sprintf("[%s] run started  id=%s  backend=%s",
		logTS(), spec.RunID, backendLabel(entry.UseDocker)))
	entry.appendLog(fmt.Sprintf("[%s] models=%v  langs=%v  dry_run=%v  max_samples=%d  workers=%d",
		logTS(), spec.Models, spec.Languages, spec.DryRun, spec.MaxSamples, spec.Workers))

	var err error
	if entry.UseDocker {
		err = runInDocker(ctx, entry, spec, opts, m.dockerCfg)
	} else {
		err = m.executeInProcess(ctx, entry, spec, opts)
	}

	now := time.Now()
	entry.mu.Lock()
	entry.EndedAt = &now
	entry.Paused = false
	switch {
	case err == nil:
		entry.Status = StatusCompleted
	case ctx.Err() == context.Canceled:
		entry.Status = StatusCanceled
		if entry.Error == "" {
			entry.Error = "canceled by user"
		}
	default:
		entry.Status = StatusFailed
		entry.Error = err.Error()
	}
	finalStatus := entry.Status
	entry.mu.Unlock()

	switch finalStatus {
	case StatusCompleted:
		entry.appendLog(fmt.Sprintf("[%s] run completed", logTS()))
	case StatusCanceled:
		entry.appendLog(fmt.Sprintf("[%s] CANCELED", logTS()))
	default:
		entry.appendLog(fmt.Sprintf("[%s] FAILED: %v", logTS(), err))
	}

	m.ingestRunArtifacts(entry, spec)
}

func (m *RunManager) ingestRunArtifacts(entry *RunEntry, spec contracts.RunSpec) {
	runDir := filepath.Join(m.outputRoot, "runs", spec.RunID)
	sqliteStore, err := store.OpenSQLite(m.dbPath)
	if err != nil {
		entry.appendLog(fmt.Sprintf("[%s] db ingest skipped: open sqlite: %v", logTS(), err))
		return
	}
	defer sqliteStore.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := sqliteStore.Init(ctx); err != nil {
		entry.appendLog(fmt.Sprintf("[%s] db ingest skipped: init sqlite: %v", logTS(), err))
		return
	}
	sum, err := sqliteStore.IngestRun(ctx, store.IngestRunOptions{RunDir: runDir})
	if err != nil {
		entry.appendLog(fmt.Sprintf("[%s] db ingest skipped: %v", logTS(), err))
		return
	}
	entry.appendLog(fmt.Sprintf("[%s] db ingest ok: run=%s generated=%d evaluated=%d artifacts=%d",
		logTS(), sum.RunID, sum.GenerationCases, sum.EvaluationResults, sum.ArtifactsIndexed))
}

// Pause 请求挂起任务。
//   - in-process: 闸门切到暂停态，worker 在下一次任务循环顶端阻塞（当前任务不中断）。
//   - docker: 调用 `docker pause <container>`，直接冻结容器。
//
// 幂等：已暂停时返回 nil。
func (m *RunManager) Pause(runID string) error {
	entry, ok := m.Get(runID)
	if !ok {
		return fmt.Errorf("run not found: %s", runID)
	}
	entry.mu.RLock()
	status := entry.Status
	useDocker := entry.UseDocker
	container := entry.container
	alreadyPaused := entry.Paused
	entry.mu.RUnlock()
	if status != StatusRunning {
		return fmt.Errorf("cannot pause run in status %q", status)
	}
	if alreadyPaused {
		return nil
	}
	if useDocker {
		if err := dockerControl("pause", container); err != nil {
			return err
		}
	} else {
		entry.gate.Pause()
	}
	entry.mu.Lock()
	entry.Paused = true
	entry.mu.Unlock()
	entry.appendLog(fmt.Sprintf("[%s] run paused", logTS()))
	return nil
}

// Resume 解除挂起；未暂停时幂等返回 nil。
func (m *RunManager) Resume(runID string) error {
	entry, ok := m.Get(runID)
	if !ok {
		return fmt.Errorf("run not found: %s", runID)
	}
	entry.mu.RLock()
	useDocker := entry.UseDocker
	container := entry.container
	paused := entry.Paused
	entry.mu.RUnlock()
	if !paused {
		return nil
	}
	if useDocker {
		if err := dockerControl("unpause", container); err != nil {
			return err
		}
	} else {
		entry.gate.Resume()
	}
	entry.mu.Lock()
	entry.Paused = false
	entry.mu.Unlock()
	entry.appendLog(fmt.Sprintf("[%s] run resumed", logTS()))
	return nil
}

// Cancel 终止任务。
//   - in-process: cancel context，worker 快速退出；当前正在进行的 API 调用/子进程会随 ctx 结束被中断。
//   - docker: `docker kill <container>`，容器立即终止。
//
// 如果任务还处于 paused 状态，会先 Resume 再取消，避免阻塞在 gate。
func (m *RunManager) Cancel(runID string) error {
	entry, ok := m.Get(runID)
	if !ok {
		return fmt.Errorf("run not found: %s", runID)
	}
	entry.mu.RLock()
	status := entry.Status
	useDocker := entry.UseDocker
	container := entry.container
	cancel := entry.cancel
	paused := entry.Paused
	entry.mu.RUnlock()
	if status != StatusRunning && status != StatusPending {
		return fmt.Errorf("cannot cancel run in status %q", status)
	}
	if paused {
		// 先放行 gate，否则 in-process worker 无法看到 ctx.Done。
		if useDocker {
			_ = dockerControl("unpause", container)
		} else {
			entry.gate.Resume()
		}
	}
	if useDocker {
		_ = dockerControl("kill", container)
	}
	if cancel != nil {
		cancel()
	}
	entry.appendLog(fmt.Sprintf("[%s] cancel requested", logTS()))
	return nil
}

// dockerControl 封装 `docker <action> <container>`。
func dockerControl(action, container string) error {
	out, err := exec.Command("docker", action, container).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker %s %s: %v: %s", action, container, err, string(out))
	}
	return nil
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
