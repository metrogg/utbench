package web

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/orchestrator"
)

// DockerConfig captures the static inputs required to execute a run inside
// the utbench container image on the local docker daemon.
type DockerConfig struct {
	ImageName   string // e.g. "utbench:latest"
	ProjectRoot string // host absolute path of project root (parent of datasets/, artifacts/, configs/)
	EnvFile     string // optional host path to .env; ignored if empty or missing
}

// runInDocker shells out to `docker run ...` and streams the combined output
// line-by-line into entry's log buffer. The host directories datasets/,
// artifacts/, configs/ and storage/ are mounted so the in-container CLI writes
// artifacts back to the host the same way the in-process runner does.
//
// On success, the container has already produced run_summary.json etc. under
// artifacts/runs/<run-id>/ so the existing listRuns/handleRunReport code keeps
// working unchanged.
func runInDocker(ctx context.Context, entry *RunEntry, spec contracts.RunSpec, opts orchestrator.Options, cfg DockerConfig) error {
	args := buildDockerRunArgs(spec, opts, cfg)
	// 注入稳定容器名，使 docker pause/unpause/kill 可以定位到本次运行。
	// `--name` 必须紧跟在 `docker run` 之后、镜像名之前。
	if entry.container != "" {
		args = append([]string{args[0], "--name", entry.container}, args[1:]...)
	}
	entry.appendLog(fmt.Sprintf("[%s] docker exec → docker %s", logTS(), redactArgs(args)))

	cmd := exec.CommandContext(ctx, "docker", args...)
	// Merge stdout + stderr into the same line-sink so users see everything
	// (build messages, evaluator logs, mutation output) in order.
	lw := &lineWriter{run: entry}
	cmd.Stdout = lw
	cmd.Stderr = lw

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("docker run start: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("docker run failed: %w", err)
	}
	return nil
}

// buildDockerRunArgs assembles the argv for `docker run`.
// It deliberately mirrors the flag set that `utbench run` understands, so
// behaviour matches the in-process backend one-to-one.
func buildDockerRunArgs(spec contracts.RunSpec, opts orchestrator.Options, cfg DockerConfig) []string {
	a := []string{"run", "--rm"}

	if cfg.EnvFile != "" {
		a = append(a, "--env-file", cfg.EnvFile)
	}

	// Mounts: datasets (read-only is safer but writable matches current UX),
	// artifacts, configs, storage. Paths on the container side are fixed and
	// mirror those used in STARTUP_GUIDE.md.
	root := strings.TrimRight(cfg.ProjectRoot, `/\`)
	a = append(a,
		"-v", root+`/datasets:/app/datasets`,
		"-v", root+`/artifacts:/app/artifacts`,
		"-v", root+`/configs:/app/configs`,
		"-v", root+`/storage:/app/storage`,
	)

	a = append(a, cfg.ImageName, "run")
	a = append(a,
		"--run-id", spec.RunID,
		"--models", strings.Join(spec.Models, ","),
		"--langs", strings.Join(spec.Languages, ","),
		"--dataset-root", "/app/datasets",
		"--output-root", "/app/artifacts",
		"--config", "/app/configs/models.yaml",
	)
	if len(spec.DatasetClasses) > 0 {
		a = append(a, "--class", strings.Join(spec.DatasetClasses, ","))
	}
	if spec.DatasetScenario != "" {
		a = append(a, "--scenario", spec.DatasetScenario)
	}
	if spec.DatasetLevel != "" {
		a = append(a, "--level", spec.DatasetLevel)
	}
	if spec.MaxSamples > 0 {
		a = append(a, "--max-samples", fmt.Sprintf("%d", spec.MaxSamples))
	}
	if spec.Workers > 0 {
		a = append(a, "--workers", fmt.Sprintf("%d", spec.Workers))
	}
	if spec.Mode != "" {
		a = append(a, "--mode", string(spec.Mode))
	}
	if spec.DryRun {
		a = append(a, "--dry-run")
	}
	if spec.MutationEnabled {
		a = append(a, "--mutation-enabled")
	}
	if spec.MutationTimeout > 0 {
		a = append(a, "--mutation-timeout", fmt.Sprintf("%d", spec.MutationTimeout))
	}
	if spec.MutationPolicy != "" {
		a = append(a, "--mutation-policy", spec.MutationPolicy)
	}
	if opts.Ingest {
		a = append(a, "--ingest", "--db-path", "/app/storage/utbench.db")
	}
	return a
}

// redactArgs produces a single-line human-readable representation of the argv
// without leaking any secrets. Currently we simply join with spaces since the
// arguments don't contain credentials (API keys come from --env-file).
func redactArgs(args []string) string {
	return strings.Join(args, " ")
}
