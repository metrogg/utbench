package web

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// BuildStatus mirrors RunStatus but for image-build jobs.
type BuildStatus string

const (
	BuildPending   BuildStatus = "pending"
	BuildRunning   BuildStatus = "running"
	BuildCompleted BuildStatus = "completed"
	BuildFailed    BuildStatus = "failed"
	BuildCanceled  BuildStatus = "canceled"
)

// BuildJob represents a running `docker build` invocation.
// It deliberately mirrors RunEntry's log + subscribe shape so the SSE handler
// can reuse the same pattern.
type BuildJob struct {
	BuildID    string        `json:"build_id"`
	Target     string        `json:"target"`
	Dockerfile string        `json:"dockerfile,omitempty"`
	ImageName  string        `json:"image_name"`
	Status     BuildStatus   `json:"status"`
	StartedAt  time.Time     `json:"started_at"`
	EndedAt    *time.Time    `json:"ended_at,omitempty"`
	Error      string        `json:"error,omitempty"`
	Done       chan struct{} `json:"-"`

	logs   []string
	mu     sync.RWMutex
	subs   []chan string
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

func (b *BuildJob) appendLog(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.logs = append(b.logs, line)
	for _, ch := range b.subs {
		select {
		case ch <- line:
		default:
		}
	}
}

// GetLogs returns a snapshot of all captured log lines.
func (b *BuildJob) GetLogs() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	cp := make([]string, len(b.logs))
	copy(cp, b.logs)
	return cp
}

// Subscribe returns a channel that receives new log lines.
func (b *BuildJob) Subscribe() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 512)
	for _, line := range b.logs {
		select {
		case ch <- line:
		default:
		}
	}
	b.subs = append(b.subs, ch)
	return ch
}

// Unsubscribe removes a subscriber channel and closes it.
func (b *BuildJob) Unsubscribe(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, s := range b.subs {
		if s == ch {
			b.subs = append(b.subs[:i], b.subs[i+1:]...)
			close(ch)
			return
		}
	}
}

func (b *BuildJob) Cancel() bool {
	b.mu.Lock()
	if b.Status != BuildPending && b.Status != BuildRunning {
		b.mu.Unlock()
		return false
	}
	cancel := b.cancel
	b.Status = BuildCanceled
	b.Error = "canceled by user"
	b.mu.Unlock()

	b.appendLog(fmt.Sprintf("[%s] build canceled by user", logTS()))
	if cancel != nil {
		cancel()
	}
	return true
}

// buildLineWriter is an io.Writer that splits bytes on '\n' and forwards each
// complete line to the BuildJob. Identical in spirit to the lineWriter used by
// RunEntry — kept separate to avoid cross-type coupling.
type buildLineWriter struct {
	job *BuildJob
	buf []byte
	mu  sync.Mutex
}

func (w *buildLineWriter) Write(p []byte) (int, error) {
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
			w.job.appendLog(line)
		}
	}
	return len(p), nil
}

// BuildManager holds the set of image-build jobs started through the web UI.
// Only one successful build per process is typical, but we keep the map so the
// SSE log stream can be re-opened on reconnect.
type BuildManager struct {
	mu          sync.RWMutex
	jobs        map[string]*BuildJob
	projectRoot string
}

type BuildProfile struct {
	Target     string
	ImageName  string
	Dockerfile string
}

// NewBuildManager creates a BuildManager rooted at projectRoot (the directory
// containing Dockerfile).
func NewBuildManager(projectRoot string) *BuildManager {
	return &BuildManager{
		jobs:        make(map[string]*BuildJob),
		projectRoot: projectRoot,
	}
}

// Submit starts `docker build -f <dockerfile> -t <imageName> <projectRoot>` in a goroutine and
// returns the BuildJob immediately. The caller should subscribe to Done or to
// the log stream to observe progress.
func (m *BuildManager) Submit(profile BuildProfile) *BuildJob {
	id := "build_" + time.Now().UTC().Format("20060102T150405.000000000Z")
	target := strings.TrimSpace(profile.Target)
	if target == "" {
		target = "eval"
	}
	dockerfile := strings.TrimSpace(profile.Dockerfile)
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}
	job := &BuildJob{
		BuildID:    id,
		Target:     target,
		Dockerfile: dockerfile,
		ImageName:  profile.ImageName,
		Status:     BuildPending,
		StartedAt:  time.Now(),
		Done:       make(chan struct{}),
	}
	m.mu.Lock()
	m.jobs[id] = job
	m.mu.Unlock()

	go m.execute(job)
	return job
}

// Get returns the BuildJob with the given id.
func (m *BuildManager) Get(id string) (*BuildJob, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	j, ok := m.jobs[id]
	return j, ok
}

// Latest returns the most recently started job, or nil if none exist.
func (m *BuildManager) Latest() *BuildJob {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var latest *BuildJob
	for _, j := range m.jobs {
		if latest == nil || j.StartedAt.After(latest.StartedAt) {
			latest = j
		}
	}
	return latest
}

func (m *BuildManager) execute(job *BuildJob) {
	defer close(job.Done)

	job.mu.Lock()
	job.Status = BuildRunning
	job.mu.Unlock()

	job.appendLog(fmt.Sprintf("[%s] docker build -f %s -t %s %s", logTS(), job.Dockerfile, job.ImageName, m.projectRoot))

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	job.mu.Lock()
	job.cancel = cancel
	job.mu.Unlock()

	cmd := exec.CommandContext(ctx, "docker", "build", "-f", job.Dockerfile, "-t", job.ImageName, m.projectRoot)
	lw := &buildLineWriter{job: job}
	cmd.Stdout = lw
	cmd.Stderr = lw
	job.cmd = cmd

	err := cmd.Run()

	now := time.Now()
	job.mu.Lock()
	job.EndedAt = &now
	if job.Status == BuildCanceled {
		job.Error = "canceled by user"
	} else if err != nil {
		job.Status = BuildFailed
		job.Error = err.Error()
	} else {
		job.Status = BuildCompleted
	}
	job.mu.Unlock()

	if job.Status == BuildCanceled {
		job.appendLog(fmt.Sprintf("[%s] BUILD CANCELED", logTS()))
	} else if err != nil {
		job.appendLog(fmt.Sprintf("[%s] BUILD FAILED: %v", logTS(), err))
	} else {
		job.appendLog(fmt.Sprintf("[%s] build completed", logTS()))
	}
}

func defaultBuildProfile(target string, cfg DockerConfig) (BuildProfile, error) {
	switch strings.ToLower(strings.TrimSpace(target)) {
	case "", "eval", "evaluation":
		return BuildProfile{
			Target:     "eval",
			ImageName:  cfg.EffectiveEvalImage(),
			Dockerfile: "Dockerfile",
		}, nil
	default:
		return BuildProfile{}, fmt.Errorf("unknown build target: %s (supported: eval)", target)
	}
}
