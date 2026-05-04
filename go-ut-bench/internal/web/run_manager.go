package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

const maxRunLogs = 10000

func (r *RunEntry) appendLog(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs = append(r.logs, line)
	if len(r.logs) > maxRunLogs {
		r.logs = r.logs[len(r.logs)-maxRunLogs:]
	}
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

// SubscribeWithSnapshot atomically returns the current buffered log lines and
// a channel that will receive only NEW lines appended after the snapshot.
// This lets the SSE handler emit one batched "snapshot" event up-front and
// avoid flooding the browser with N individual events on connect.
func (r *RunEntry) SubscribeWithSnapshot() ([]string, chan string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	snap := make([]string, len(r.logs))
	copy(snap, r.logs)
	ch := make(chan string, 512)
	r.subs = append(r.subs, ch)
	return snap, ch
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
	mu               sync.RWMutex
	runs             map[string]*RunEntry
	configPath       string
	agentsConfigPath string
	datasetRoot      string
	outputRoot       string
	dbPath           string
	// dockerCfg is used when a run is submitted with UseDocker=true.
	dockerCfg DockerConfig
	stopCleaner chan struct{}
}

// runEntryRetention 是已完成 RunEntry 在内存中的保留时间。
// 超过此时间的已完成/失败/取消条目会被自动清除，释放日志缓冲区内存。
const runEntryRetention = 30 * time.Minute

// NewRunManager creates a RunManager with the given default paths.
// imageName/projectRoot/envFile provide defaults for Docker execution; they
// may be empty if Docker mode is not supported in the user's environment.
func NewRunManager(configPath, agentsConfigPath, datasetRoot, outputRoot, dbPath string, cfg DockerConfig) *RunManager {
	m := &RunManager{
		runs:             make(map[string]*RunEntry),
		configPath:       configPath,
		agentsConfigPath: agentsConfigPath,
		datasetRoot:      datasetRoot,
		outputRoot:       outputRoot,
		dbPath:           dbPath,
		dockerCfg:        cfg,
		stopCleaner:      make(chan struct{}),
	}
	go m.runCleaner()
	return m
}

// runCleaner 周期性清理已完成且超过保留时间的 RunEntry，防止内存无限增长。
func (m *RunManager) runCleaner() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.evictStaleEntries()
		case <-m.stopCleaner:
			return
		}
	}
}

// evictStaleEntries 从 runs map 中移除超过保留时间的终态条目。
func (m *RunManager) evictStaleEntries() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for id, entry := range m.runs {
		entry.mu.RLock()
		status := entry.Status
		endedAt := entry.EndedAt
		entry.mu.RUnlock()
		if status == StatusCompleted || status == StatusFailed || status == StatusCanceled {
			if endedAt != nil && now.Sub(*endedAt) > runEntryRetention {
				delete(m.runs, id)
			}
		}
	}
}

// Close 停止后台清理 goroutine。在服务器关闭时调用。
func (m *RunManager) Close() {
	select {
	case <-m.stopCleaner:
	default:
		close(m.stopCleaner)
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
		// Docker 模式下拆分流水线：generation 在宿主机执行（agent 沙箱需要宿主机 Docker），
		// evaluation 在 eval 容器内执行（需要语言工具链）。
		// 单独的 evaluate/report phase 直接进容器。
		err = m.executeDockerSplit(ctx, entry, spec, opts)
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

	// 清理可能残留的 sandbox 子容器（无论成功/失败/取消）
	if entry.UseDocker {
		go killSandboxContainers(spec.RunID)
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
		// 清理内层 sandbox 子容器（通过 Docker socket 创建的兄弟容器）
		go killSandboxContainers(runID)
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

// killSandboxContainers 清理指定 run 的所有 sandbox 子容器。
// sandbox 容器通过 --label utbench-run=<runID> 标记，此函数查找并 kill+rm 它们。
func killSandboxContainers(runID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 查找所有带 utbench-run=<runID> label 的容器（包括已停止的）
	out, err := exec.CommandContext(ctx, "docker", "ps", "-a", "--filter", "label=utbench-run="+runID, "-q").CombinedOutput()
	if err != nil {
		return
	}
	ids := strings.TrimSpace(string(out))
	if ids == "" {
		return
	}
	args := strings.Fields(ids)
	// kill 正在运行的
	exec.CommandContext(ctx, "docker", append([]string{"kill"}, args...)...).CombinedOutput()
	// rm 已停止的
	exec.CommandContext(ctx, "docker", append([]string{"rm", "-f"}, args...)...).CombinedOutput()
}

// killAllUtbenchSandboxes 兜底清理：kill 所有带 utbench-run label 的容器。
// 用于服务器关闭时清理可能残留的孤儿容器。
func killAllUtbenchSandboxes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "docker", "ps", "-a", "--filter", "label=utbench-run", "-q").CombinedOutput()
	if err != nil {
		return
	}
	ids := strings.TrimSpace(string(out))
	if ids == "" {
		return
	}
	args := strings.Fields(ids)
	exec.CommandContext(ctx, "docker", append([]string{"kill"}, args...)...).CombinedOutput()
	exec.CommandContext(ctx, "docker", append([]string{"rm", "-f"}, args...)...).CombinedOutput()
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
	rp := reporter.NewService(logger, &runner.DefaultPromptMetaProvider{})
	orch := orchestrator.New(ds, rn, ev, rp)

	_, err := orch.Run(ctx, spec, opts)
	return err
}

// executeDockerSplit 拆分流水线：generation 在宿主机执行，evaluation 在 eval 容器内执行。
//
// Agent 沙箱（opencode 等 CLI Agent）需要宿主机 Docker 来启动沙箱容器，
// 而评测工具链（compile/test/coverage/mutation）在 eval 容器内。
// 因此完整流水线拆为三步：
//  1. generate — 宿主机 in-process（agent 沙箱可用宿主机 Docker）
//  2. evaluate — eval 容器（语言工具链齐全）
//  3. report   — 宿主机 in-process（纯数据聚合，不需要工具链）
//
// 单独的 evaluate/report phase 直接进容器执行。
func (m *RunManager) executeDockerSplit(ctx context.Context, entry *RunEntry, spec contracts.RunSpec, opts orchestrator.Options) error {
	phase := opts.Phase
	if phase == "" {
		phase = "full"
	}

	// 单独的 report — 直接进 eval 容器。
	if phase == "report" {
		return runInDocker(ctx, entry, spec, opts, m.dockerCfg)
	}

	// 单独的 evaluate — 容器内评测 + 宿主机报告。
	if phase == "evaluate" {
		entry.appendLog(fmt.Sprintf("[%s] phase 1/2: evaluate (docker: %s)", logTS(), m.dockerCfg.EffectiveEvalImage()))
		if err := runInDocker(ctx, entry, spec, opts, m.dockerCfg); err != nil {
			return fmt.Errorf("evaluate phase failed: %w", err)
		}
		entry.appendLog(fmt.Sprintf("[%s] evaluate phase completed", logTS()))

		entry.appendLog(fmt.Sprintf("[%s] phase 2/2: report (in-process)", logTS()))
		reportOpts := orchestrator.Options{
			Phase:       "report",
			SourceRunID: spec.RunID,
		}
		if err := m.executeInProcess(ctx, entry, spec, reportOpts); err != nil {
			return fmt.Errorf("report phase failed: %w", err)
		}
		return nil
	}

	// phase == "full" 或 "generate"：先在宿主机跑 generation。
	// 设置路径映射环境变量，使 sandbox 代码能将相对 workspace 路径转换为 Docker 需要的绝对路径。
	hostOutputRoot := filepath.Join(strings.TrimRight(m.dockerCfg.ProjectRoot, `/\`), "artifacts")
	os.Setenv("UTBENCH_SANDBOX_HOST_OUTPUT_ROOT", hostOutputRoot)
	os.Setenv("UTBENCH_SANDBOX_CONTAINER_OUTPUT_ROOT", spec.OutputRoot)
	defer func() {
		os.Unsetenv("UTBENCH_SANDBOX_HOST_OUTPUT_ROOT")
		os.Unsetenv("UTBENCH_SANDBOX_CONTAINER_OUTPUT_ROOT")
	}()
	entry.appendLog(fmt.Sprintf("[%s] phase 1/2: generate (in-process, agent sandbox uses host Docker)", logTS()))
	genOpts := orchestrator.Options{
		Phase:       "generate",
		SourceRunID: spec.RunID,
	}
	if err := m.executeInProcess(ctx, entry, spec, genOpts); err != nil {
		return fmt.Errorf("generate phase failed: %w", err)
	}
	entry.appendLog(fmt.Sprintf("[%s] generate phase completed", logTS()))

	// 清理 generate 阶段可能残留的 sandbox 子容器
	killSandboxContainers(spec.RunID)

	// generate-only 模式到此结束。
	if phase == "generate" {
		return nil
	}

	// phase == "full"：进 eval 容器跑 evaluate + report。
	entry.appendLog(fmt.Sprintf("[%s] phase 2/2: evaluate+report (docker: %s)", logTS(), m.dockerCfg.EffectiveEvalImage()))
	evalOpts := orchestrator.Options{
		Phase:       "evaluate",
		SourceRunID: spec.RunID,
		Ingest:      opts.Ingest,
		DBPath:      opts.DBPath,
	}
	if err := runInDocker(ctx, entry, spec, evalOpts, m.dockerCfg); err != nil {
		return fmt.Errorf("evaluate phase failed: %w", err)
	}
	entry.appendLog(fmt.Sprintf("[%s] evaluate phase completed", logTS()))

	// 报告生成在宿主机 in-process 执行（纯数据聚合 + 模板渲染，不需要语言工具链）。
	entry.appendLog(fmt.Sprintf("[%s] phase 3/3: report (in-process)", logTS()))
	reportOpts := orchestrator.Options{
		Phase:       "report",
		SourceRunID: spec.RunID,
	}
	if err := m.executeInProcess(ctx, entry, spec, reportOpts); err != nil {
		return fmt.Errorf("report phase failed: %w", err)
	}

	return nil
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
