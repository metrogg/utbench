package web

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/orchestrator"
)

// DockerConfig captures the static inputs required to execute a run inside
// the utbench container image on the local docker daemon.
type DockerConfig struct {
	EvalImageName string // e.g. "utbench:latest"
	ProjectRoot   string // host absolute path of project root (parent of datasets/, artifacts/, configs/)
	EnvFile       string // optional host path to .env; ignored if empty or missing
}

func (c DockerConfig) EffectiveEvalImage() string {
	if strings.TrimSpace(c.EvalImageName) != "" {
		return strings.TrimSpace(c.EvalImageName)
	}
	return "utbench:latest"
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
	// Windows Git Bash 的 MSYS 会把 -e 传的路径值自动转换（如 /c/Users → C:/Program Files/Git/c/Users），
	// 导致容器内环境变量损坏。设置 MSYS_NO_PATHCONV=1 禁止此转换。
	cmd.Env = append(os.Environ(), "MSYS_NO_PATHCONV=1")
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
// When opts.Phase is set to "generate", "evaluate", or "report", the command
// switches to the corresponding CLI subcommand instead of "run".
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
		// 持久化 Maven 本地仓库，避免每次容器运行都重新下载依赖。
		// Docker 镜像已预下载关键依赖，但此 mount 可缓存运行时新增的依赖。
		"-v", root+`/.m2-cache:/root/.m2/repository`,
	)
	// 不挂载 docker.sock — 评测面不需要 Docker daemon。
	// Agent 沙箱容器由宿主机控制面直接启动（见 executeDockerSplit）。
	// 评测面自身只负责 compile/test/coverage/mutation。

	// Determine the CLI subcommand based on phase.
	// Phase "full" (or empty) uses "run" command.
	// Phase "generate", "evaluate", "report" uses the corresponding subcommand.
	phase := opts.Phase
	if phase == "" {
		phase = "full"
	}
	cliCmd := "run"
	switch phase {
	case "generate":
		cliCmd = "generate"
	case "evaluate":
		cliCmd = "evaluate"
	case "report":
		cliCmd = "report"
	}

	a = append(a, cfg.EffectiveEvalImage(), cliCmd)

	// Common flags for all commands
	a = append(a,
		"--run-id", spec.RunID,
		"--output-root", "/app/artifacts",
	)

	// Determine source run ID for artifact paths
	sourceRunID := opts.SourceRunID
	if sourceRunID == "" {
		sourceRunID = spec.RunID
	}

	// Phase-specific flags
	switch cliCmd {
	case "run":
		a = append(a,
			"--models", strings.Join(spec.Models, ","),
			"--langs", strings.Join(spec.Languages, ","),
			"--dataset-root", "/app/datasets",
			"--config", "/app/configs/models.yaml",
		)
		if spec.AgentsConfigPath != "" {
			a = append(a, "--agents-config", path.Join("/app/configs", filepath.Base(spec.AgentsConfigPath)))
		}
		if len(spec.Subjects) > 0 {
			a = append(a, "--subjects", strings.Join(spec.Subjects, ","))
		}
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
		if spec.ReuseGenerated {
			a = append(a, "--reuse-generated", "--db-path", "/app/storage/utbench.db")
		}
		if spec.ReuseEvaluation {
			a = append(a, "--reuse-evaluation", "--db-path", "/app/storage/utbench.db")
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
		if spec.TestTimeout > 0 {
			a = append(a, "--test-timeout", fmt.Sprintf("%d", spec.TestTimeout))
		}
		if opts.Ingest {
			a = append(a, "--ingest", "--db-path", "/app/storage/utbench.db")
		}

	case "generate":
		a = append(a,
			"--models", strings.Join(spec.Models, ","),
			"--langs", strings.Join(spec.Languages, ","),
			"--dataset-root", "/app/datasets",
			"--config", "/app/configs/models.yaml",
		)
		if spec.AgentsConfigPath != "" {
			a = append(a, "--agents-config", path.Join("/app/configs", filepath.Base(spec.AgentsConfigPath)))
		}
		if len(spec.Subjects) > 0 {
			a = append(a, "--subjects", strings.Join(spec.Subjects, ","))
		}
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
		if spec.ReuseGenerated {
			a = append(a, "--reuse-generated", "--db-path", "/app/storage/utbench.db")
		}

	case "evaluate":
		// Use explicit manifest path if provided, otherwise use source run's manifest
		manifestPath := opts.ManifestPath
		if manifestPath == "" {
			manifestPath = path.Join("/app/artifacts", "runs", sourceRunID, "generated", "generated_manifest.json")
		}
		// Convert host path to container path if it's absolute
		if strings.HasPrefix(manifestPath, root) {
			manifestPath = strings.Replace(manifestPath, root, "/app", 1)
		}
		a = append(a, "--manifest", manifestPath)
		if spec.MutationEnabled {
			a = append(a, "--mutation-enabled")
		}
		if spec.MutationTimeout > 0 {
			a = append(a, "--mutation-timeout", fmt.Sprintf("%d", spec.MutationTimeout))
		}
		if spec.MutationPolicy != "" {
			a = append(a, "--mutation-policy", spec.MutationPolicy)
		}
		if spec.TestTimeout > 0 {
			a = append(a, "--test-timeout", fmt.Sprintf("%d", spec.TestTimeout))
		}
		if spec.Workers > 0 {
			a = append(a, "--workers", fmt.Sprintf("%d", spec.Workers))
		}
		if spec.ReuseEvaluation {
			a = append(a, "--reuse-evaluation", "--db-path", "/app/storage/utbench.db")
		}

	case "report":
		// Use explicit evaluation path if provided, otherwise use source run's evaluation
		evaluationPath := opts.EvaluationPath
		if evaluationPath == "" {
			evaluationPath = path.Join("/app/artifacts", "runs", sourceRunID, "evaluation", "evaluation_result.json")
		}
		// Convert host path to container path if it's absolute
		if strings.HasPrefix(evaluationPath, root) {
			evaluationPath = strings.Replace(evaluationPath, root, "/app", 1)
		}
		a = append(a, "--evaluation", evaluationPath)
	}

	return wrapDockerSourceCommand(a, cfg)
}

// runEvaluateInDocker runs only the evaluation step inside the utbench container.
// Unlike runInDocker (which runs the full pipeline), this is a simpler synchronous
// wrapper that captures output as a string. Used by the reevaluate API handler.
func runEvaluateInDocker(ctx context.Context, runID string, spec contracts.RunSpec, cfg DockerConfig) (string, error) {
	opts := orchestrator.Options{
		Phase:       "evaluate",
		SourceRunID: runID,
	}
	args := buildDockerRunArgs(spec, opts, cfg)
	cmd := exec.CommandContext(ctx, "docker", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func wrapDockerSourceCommand(args []string, cfg DockerConfig) []string {
	imageIdx := -1
	for i, arg := range args {
		if arg == cfg.EffectiveEvalImage() {
			imageIdx = i
			break
		}
	}
	if imageIdx < 0 || imageIdx == len(args)-1 {
		return args
	}
	// Dockerfile ENTRYPOINT 已设置为 ["./utbench"]，命令从子命令开始即可
	out := append([]string{}, args[:imageIdx]...)
	out = append(out, cfg.EffectiveEvalImage())
	out = append(out, args[imageIdx+1:]...)
	return out
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
