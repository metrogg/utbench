package web

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// EnvStatus describes the host environment relevant to benchmark execution.
// Returned from GET /api/env so the frontend can decide whether to recommend
// Docker execution, show a "Build Image" button, etc.
type EnvStatus struct {
	OS                 string          `json:"os"`
	ProjectRoot        string          `json:"project_root"`
	DockerAvailable    bool            `json:"docker_available"`
	DockerVersion      string          `json:"docker_version,omitempty"`
	DockerError        string          `json:"docker_error,omitempty"`
	EvalImageName      string          `json:"eval_image_name,omitempty"`
	EvalImagePresent   bool            `json:"eval_image_present,omitempty"`
	EvalImageID        string          `json:"eval_image_id,omitempty"`
	AgentImageName     string          `json:"agent_image_name,omitempty"`
	AgentImagePresent  bool            `json:"agent_image_present,omitempty"`
	AgentImageID       string          `json:"agent_image_id,omitempty"`
	EnvFilePresent     bool            `json:"env_file_present"`
	EnvFilePath        string          `json:"env_file_path,omitempty"`
	NativeTools        map[string]bool `json:"native_tools"`
	RunningInContainer bool            `json:"running_in_container"`
	TopologyMode       string          `json:"topology_mode,omitempty"`
	Recommendation     string          `json:"recommendation,omitempty"`
}

// detectTimeout bounds every external command we shell out to.
const detectTimeout = 4 * time.Second

const defaultAgentImageName = "utbench-agent-base:latest"

// DetectEnv probes the host for Docker + the image + native tool chain.
// It never returns an error; any per-check failure is recorded in the fields.
func DetectEnv(cfg DockerConfig) EnvStatus {
	st := EnvStatus{
		OS:                 runtime.GOOS,
		ProjectRoot:        cfg.ProjectRoot,
		EvalImageName:      cfg.EffectiveEvalImage(),
		AgentImageName:     defaultAgentImageName,
		NativeTools:        detectNativeTools(),
		RunningInContainer: detectContainerRuntime(),
	}
	st.DockerAvailable, st.DockerVersion, st.DockerError = detectDocker()
	if st.DockerAvailable {
		st.EvalImagePresent, st.EvalImageID = detectImage(st.EvalImageName)
		st.AgentImagePresent, st.AgentImageID = detectImage(st.AgentImageName)
	}
	if cfg.EnvFile != "" {
		if _, err := os.Stat(cfg.EnvFile); err == nil {
			st.EnvFilePresent = true
			st.EnvFilePath = cfg.EnvFile
		}
	}
	st.TopologyMode = detectTopologyMode(st)
	st.Recommendation = buildRecommendation(st)
	return st
}

// detectDocker runs `docker version --format {{.Server.Version}}`.
func detectDocker() (bool, string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), detectTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "version", "--format", "{{.Server.Version}}")
	hideCommandWindow(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, "", strings.TrimSpace(string(out)) + " " + err.Error()
	}
	return true, strings.TrimSpace(string(out)), ""
}

// detectImage checks whether the named image exists locally.
// `docker image inspect <name> --format {{.Id}}` exits non-zero if missing.
func detectImage(name string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), detectTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "image", "inspect", name, "--format", "{{.Id}}")
	hideCommandWindow(cmd)
	out, err := cmd.Output()
	if err != nil {
		return false, ""
	}
	return true, strings.TrimSpace(string(out))
}

// detectNativeTools probes common evaluator toolchains via exec.LookPath.
// These are advisory only; failure to find a tool just tells the UI to
// suggest Docker execution.
func detectNativeTools() map[string]bool {
	tools := []string{
		"python3", "pytest", "coverage", "mutmut",
		"go", "go-mutesting",
		"mvn", "java", "javac",
		"clang", "clang++", "cmake", "mull-runner-19", "mull-runner",
	}
	out := make(map[string]bool, len(tools))
	for _, t := range tools {
		_, err := exec.LookPath(t)
		out[t] = err == nil
	}
	return out
}

func detectContainerRuntime() bool {
	if fileExistsLocal("/.dockerenv") || fileExistsLocal("/run/.containerenv") {
		return true
	}
	if strings.TrimSpace(os.Getenv("container")) != "" {
		return true
	}
	return false
}

func fileExistsLocal(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func detectTopologyMode(st EnvStatus) string {
	switch {
	case st.RunningInContainer:
		return "container"
	case st.DockerAvailable:
		return "host+docker"
	default:
		return "host"
	}
}

// buildRecommendation crafts a short advisory message for the UI.
func buildRecommendation(st EnvStatus) string {
	switch {
	case !st.DockerAvailable:
		if st.OS == "windows" && !st.NativeTools["mutmut"] {
			return "Docker is not available. Mutation testing (mutmut) is not supported natively on Windows; install Docker Desktop or use WSL."
		}
		return ""
	case !st.EvalImagePresent:
		return "Docker 已就绪，但评测镜像 '" + st.EvalImageName + "' 尚未构建。运行 ./build.sh eval 构建。"
	case st.OS == "windows" && !st.NativeTools["mutmut"]:
		return "Docker image is ready. Recommend enabling 'Execute in Docker' when running mutation tests on Windows."
	default:
		return ""
	}
}
