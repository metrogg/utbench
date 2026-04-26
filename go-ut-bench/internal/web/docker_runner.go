package web

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
	cmd := []string{"go", "run", "./cmd/utbench", "run"}
	cmd = append(cmd,
		"--run-id", spec.RunID,
		"--models", strings.Join(spec.Models, ","),
		"--langs", strings.Join(spec.Languages, ","),
		"--dataset-root", "/app/datasets",
		"--output-root", "/app/artifacts",
		"--config", "/app/configs/models.yaml",
	)
	if len(spec.DatasetClasses) > 0 {
		cmd = append(cmd, "--class", strings.Join(spec.DatasetClasses, ","))
	}
	if spec.DatasetScenario != "" {
		cmd = append(cmd, "--scenario", spec.DatasetScenario)
	}
	if spec.DatasetLevel != "" {
		cmd = append(cmd, "--level", spec.DatasetLevel)
	}
	if spec.MaxSamples > 0 {
		cmd = append(cmd, "--max-samples", fmt.Sprintf("%d", spec.MaxSamples))
	}
	if spec.Workers > 0 {
		cmd = append(cmd, "--workers", fmt.Sprintf("%d", spec.Workers))
	}
	if spec.Mode != "" {
		cmd = append(cmd, "--mode", string(spec.Mode))
	}
	if spec.DryRun {
		cmd = append(cmd, "--dry-run")
	}
	if spec.MutationEnabled {
		cmd = append(cmd, "--mutation-enabled")
	}
	if spec.MutationTimeout > 0 {
		cmd = append(cmd, "--mutation-timeout", fmt.Sprintf("%d", spec.MutationTimeout))
	}
	if spec.MutationPolicy != "" {
		cmd = append(cmd, "--mutation-policy", spec.MutationPolicy)
	}
	if opts.Ingest {
		cmd = append(cmd, "--ingest", "--db-path", "/app/storage/utbench.db")
	}
	return buildDockerSourceArgs(cfg, cmd)
}

func runEvaluateInDocker(ctx context.Context, runID string, spec contracts.RunSpec, cfg DockerConfig) ([]byte, error) {
	cmdArgs := []string{
		"go", "run", "./cmd/utbench", "evaluate",
		"--run-id", runID,
		"--manifest", "/app/artifacts/runs/" + runID + "/generated/generated_manifest.json",
		"--output-root", "/app/artifacts",
		"--mutation-timeout", strconv.Itoa(defaultInt(spec.MutationTimeout, 600)),
		"--mutation-policy", defaultString(spec.MutationPolicy, "warn"),
	}
	if spec.MutationEnabled {
		cmdArgs = append(cmdArgs, "--mutation-enabled")
	}
	args := buildDockerSourceArgs(cfg, cmdArgs)
	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("docker evaluate failed: %w", err)
	}
	return out, nil
}

func buildDockerSourceArgs(cfg DockerConfig, cmd []string) []string {
	a := buildDockerBaseArgs(cfg)
	root := strings.TrimRight(cfg.ProjectRoot, `/\`)
	a = append(a,
		"-v", root+":/workspace",
		"-w", "/workspace",
		"--entrypoint", "/bin/sh",
		cfg.ImageName,
		"-lc", shellJoin(cmd),
	)
	return a
}

func buildDockerBaseArgs(cfg DockerConfig) []string {
	a := []string{"run", "--rm"}
	if cfg.EnvFile != "" && fileExists(cfg.EnvFile) {
		a = append(a, "--env-file", cfg.EnvFile)
	}
	root := strings.TrimRight(cfg.ProjectRoot, `/\`)
	return append(a,
		"-v", root+`/datasets:/app/datasets`,
		"-v", root+`/artifacts:/app/artifacts`,
		"-v", root+`/configs:/app/configs`,
		"-v", root+`/storage:/app/storage`,
	)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func defaultInt(v, fallback int) int {
	if v == 0 {
		return fallback
	}
	return v
}

func defaultString(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func shellJoin(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = shellQuote(arg)
	}
	return strings.Join(quoted, " ")
}

func shellQuote(arg string) string {
	if arg == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(arg, "'", `'\''`) + "'"
}

// redactArgs produces a single-line human-readable representation of the argv
// without leaking any secrets. Currently we simply join with spaces since the
// arguments don't contain credentials (API keys come from --env-file).
func redactArgs(args []string) string {
	return strings.Join(args, " ")
}
