package evaluator

import (
	"fmt"
	"sync"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/obs"
)

type activeEvalTracker struct {
	mu     sync.Mutex
	tasks  map[string]activeEvalTask
	logger *obs.Logger
}

type activeEvalTask struct {
	Model      string
	Language   string
	SampleID   string
	Phase      string
	StartedAt  time.Time
	PhaseSince time.Time
}

func newActiveEvalTracker(logger *obs.Logger) *activeEvalTracker {
	return &activeEvalTracker{
		tasks:  make(map[string]activeEvalTask),
		logger: logger,
	}
}

func (t *activeEvalTracker) start(item contracts.GeneratedCase) (string, func(string), func()) {
	key := fmt.Sprintf("%s/%s/%s", item.Model, item.Language, item.SampleID)
	now := time.Now()
	t.mu.Lock()
	t.tasks[key] = activeEvalTask{
		Model:      item.Model,
		Language:   item.Language,
		SampleID:   item.SampleID,
		Phase:      "start",
		StartedAt:  now,
		PhaseSince: now,
	}
	t.mu.Unlock()

	update := func(phase string) {
		t.update(key, phase)
	}
	done := func() {
		t.done(key)
	}
	return key, update, done
}

func (t *activeEvalTracker) update(key, phase string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	task, ok := t.tasks[key]
	if !ok {
		return
	}
	task.Phase = phase
	task.PhaseSince = time.Now()
	t.tasks[key] = task
}

func (t *activeEvalTracker) done(key string) {
	t.mu.Lock()
	delete(t.tasks, key)
	t.mu.Unlock()
}

func (t *activeEvalTracker) printStalled(minAge time.Duration) {
	now := time.Now()
	stalled := make([]activeEvalTask, 0)

	t.mu.Lock()
	for _, task := range t.tasks {
		if now.Sub(task.PhaseSince) >= minAge || now.Sub(task.StartedAt) >= minAge {
			stalled = append(stalled, task)
		}
	}
	t.mu.Unlock()

	for _, task := range stalled {
		totalElapsed := now.Sub(task.StartedAt).Round(time.Second)
		phaseElapsed := now.Sub(task.PhaseSince).Round(time.Second)
		fmt.Printf("        [EVAL-WARN] still running | %s | %s | %s | phase=%s | phase_elapsed=%s | total_elapsed=%s\n",
			task.Model, task.Language, task.SampleID, task.Phase, phaseElapsed, totalElapsed)
		if t.logger != nil {
			t.logger.ToFile("evaluator").Trace("eval_task_still_running",
				"model", task.Model,
				"language", task.Language,
				"sample_id", task.SampleID,
				"phase", task.Phase,
				"phase_elapsed_seconds", int(phaseElapsed.Seconds()),
				"total_elapsed_seconds", int(totalElapsed.Seconds()),
			)
		}
	}
}
