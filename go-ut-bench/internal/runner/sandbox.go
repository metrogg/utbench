package runner

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type SandboxRunRequest struct {
	Mode                string
	Workspace           string
	ContainerOutputRoot string
	Command             string
	Env                 map[string]string
	EnvFromHost         []string
	DockerImage         string
	NetworkDisabled     bool
	CPU                 string
	Memory              string
	TimeoutSeconds      int
}

type SandboxRunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type SandboxRunner interface {
	Run(ctx context.Context, req SandboxRunRequest) (SandboxRunResult, error)
}

type defaultSandboxRunner struct{}

func NewSandboxRunner() SandboxRunner {
	return defaultSandboxRunner{}
}

func (defaultSandboxRunner) Run(ctx context.Context, req SandboxRunRequest) (SandboxRunResult, error) {
	timeout := req.TimeoutSeconds
	if timeout <= 0 {
		timeout = 600
	}
	runCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	if strings.EqualFold(req.Mode, "local") {
		return runLocalSandbox(runCtx, req)
	}
	return runDockerSandbox(runCtx, req, timeout)
}

func runDockerSandbox(ctx context.Context, req SandboxRunRequest, timeout int) (SandboxRunResult, error) {
	mountSource, err := resolveDockerWorkspaceMount(req)
	if err != nil {
		return SandboxRunResult{}, err
	}
	args := []string{"run", "--rm"}
	if req.NetworkDisabled {
		args = append(args, "--network", "none")
	}
	if req.CPU != "" {
		args = append(args, "--cpus", req.CPU)
	}
	if req.Memory != "" {
		args = append(args, "--memory", req.Memory)
	}
	for _, key := range uniqueSortedStrings(req.EnvFromHost) {
		if value, ok := os.LookupEnv(key); ok {
			args = append(args, "-e", key+"="+value)
		}
	}
	keys := make([]string, 0, len(req.Env))
	for key := range req.Env {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		args = append(args, "-e", key+"="+req.Env[key])
	}
	image := strings.TrimSpace(req.DockerImage)
	if image == "" {
		image = "utbench-agent:latest"
	}
	args = append(args, "-v", mountSource+":/workspace", "-w", "/workspace", image, "/bin/sh", "-lc", req.Command)
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	result := SandboxRunResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}
	if ctx.Err() != nil {
		return result, fmt.Errorf("agent command timed out after %ds", timeout)
	}
	return result, err
}

func resolveDockerWorkspaceMount(req SandboxRunRequest) (string, error) {
	workspace := strings.TrimSpace(req.Workspace)
	if workspace == "" {
		return "", fmt.Errorf("docker sandbox workspace is empty")
	}
	hostOutputRoot := strings.TrimSpace(os.Getenv("UTBENCH_SANDBOX_HOST_OUTPUT_ROOT"))
	containerOutputRoot := strings.TrimSpace(os.Getenv("UTBENCH_SANDBOX_CONTAINER_OUTPUT_ROOT"))
	if containerOutputRoot == "" {
		containerOutputRoot = strings.TrimSpace(req.ContainerOutputRoot)
	}
	if containerOutputRoot == "" {
		containerOutputRoot = "/app/artifacts"
	}
	if hostOutputRoot != "" {
		rel, err := filepath.Rel(containerOutputRoot, workspace)
		if err == nil {
			if rel == "." {
				return hostOutputRoot, nil
			}
			if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return filepath.Join(hostOutputRoot, rel), nil
			}
		}
	}
	if strings.HasPrefix(filepath.ToSlash(workspace), "/app/") {
		return "", fmt.Errorf("docker sandbox workspace %s looks container-local; set UTBENCH_SANDBOX_HOST_OUTPUT_ROOT to the host artifacts path when using DOOD", workspace)
	}
	return workspace, nil
}

func runLocalSandbox(ctx context.Context, req SandboxRunRequest) (SandboxRunResult, error) {
	var name string
	var args []string
	if runtime.GOOS == "windows" {
		name = "powershell"
		args = []string{"-NoProfile", "-Command", req.Command}
	} else {
		name = "/bin/sh"
		args = []string{"-lc", req.Command}
	}
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = req.Workspace
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if len(req.Env) > 0 || len(req.EnvFromHost) > 0 {
		cmd.Env = os.Environ()
		for _, key := range uniqueSortedStrings(req.EnvFromHost) {
			if value, ok := os.LookupEnv(key); ok {
				cmd.Env = append(cmd.Env, key+"="+value)
			}
		}
		keys := make([]string, 0, len(req.Env))
		for key := range req.Env {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			cmd.Env = append(cmd.Env, key+"="+req.Env[key])
		}
	}
	err := cmd.Run()
	result := SandboxRunResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}
	if ctx.Err() != nil {
		return result, ctx.Err()
	}
	return result, err
}

func sandboxFingerprintForRequest(req SandboxRunRequest) string {
	payload := map[string]any{
		"mode":             req.Mode,
		"docker_image":     req.DockerImage,
		"network_disabled": req.NetworkDisabled,
		"cpu":              req.CPU,
		"memory":           req.Memory,
		"env_keys":         sortedMapKeys(req.Env),
		"env_from_host":    uniqueSortedStrings(req.EnvFromHost),
	}
	raw, _ := json.Marshal(payload)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func sortedMapKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
