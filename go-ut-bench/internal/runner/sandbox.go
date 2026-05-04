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
	Provider            string
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
	// ReadOnlyMounts 挂载为只读的宿主路径列表（Docker 模式下挂载 :ro）。
	// 用于源码文件保护，防止 Agent 意外修改被测源码。
	ReadOnlyMounts []string
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
	// Docker 要求绝对路径作为 volume mount 源
	if !filepath.IsAbs(mountSource) {
		if absPath, err := filepath.Abs(mountSource); err == nil {
			fmt.Fprintf(os.Stderr, "[sandbox] workspace mount: converted relative %q to absolute %q\n", mountSource, absPath)
			mountSource = absPath
		}
	}
	args := []string{"run", "--rm"}
	// 添加 label 以便取消时能定位和清理孤儿容器
	runID := req.Env["UTBENCH_RUN_ID"]
	if runID != "" {
		args = append(args, "--label", "utbench-run="+runID)
	}
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
	// 挂载 workspace（读写，Agent 需要写入测试文件）
	args = append(args, "-v", mountSource+":/workspace", "-w", "/workspace")
	// 只读挂载源码文件，防止 Agent 意外修改被测源码
	for _, roPath := range req.ReadOnlyMounts {
		roPath = strings.TrimSpace(roPath)
		if roPath == "" {
			continue
		}
		// 将相对路径转换为绝对路径，Docker 要求绝对路径
		if !filepath.IsAbs(roPath) {
			if absPath, err := filepath.Abs(roPath); err == nil {
				fmt.Fprintf(os.Stderr, "[sandbox] ReadOnlyMount: converted relative %q to absolute %q\n", roPath, absPath)
				roPath = absPath
			}
		}
		containerRO := "/workspace/readonly_sources/" + filepath.Base(roPath)
		args = append(args, "-v", roPath+":"+containerRO+":ro")
	}
	args = append(args, image, sandboxShell(), "-c", req.Command)
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
	fmt.Fprintf(os.Stderr, "[sandbox] resolveDockerWorkspaceMount: workspace=%q hostOutputRoot=%q containerOutputRoot=%q containerOutputRoot_fallback=%q\n",
		workspace, hostOutputRoot, containerOutputRoot, strings.TrimSpace(req.ContainerOutputRoot))
	if containerOutputRoot == "" {
		containerOutputRoot = strings.TrimSpace(req.ContainerOutputRoot)
	}
	if containerOutputRoot == "" {
		containerOutputRoot = "/app/artifacts"
	}
	if hostOutputRoot != "" {
		rel, err := filepath.Rel(containerOutputRoot, workspace)
		fmt.Fprintf(os.Stderr, "[sandbox] filepath.Rel(%q, %q) = %q, err=%v\n", containerOutputRoot, workspace, rel, err)
		if err == nil {
			if rel == "." {
				result := filepath.ToSlash(hostOutputRoot)
				fmt.Fprintf(os.Stderr, "[sandbox] resolved to hostOutputRoot: %q\n", result)
				return result, nil
			}
			if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				result := filepath.ToSlash(filepath.Join(hostOutputRoot, rel))
				fmt.Fprintf(os.Stderr, "[sandbox] resolved to joined path: %q\n", result)
				return result, nil
			}
			fmt.Fprintf(os.Stderr, "[sandbox] rel path rejected: %q (starts with ..)\n", rel)
		} else {
			fmt.Fprintf(os.Stderr, "[sandbox] filepath.Rel failed: %v\n", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "[sandbox] hostOutputRoot is empty, returning workspace as-is\n")
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
		name = sandboxShell()
		args = []string{"-c", req.Command}
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

// isDockerAvailable 检测当前环境是否有可用的 Docker daemon。
// Linux/macOS 检查 /var/run/docker.sock 是否存在；
// Windows 检查 docker 命令是否可用（Docker Desktop 使用命名管道而非 Unix socket）。
func isDockerAvailable() bool {
	if runtime.GOOS == "windows" {
		_, err := exec.LookPath("docker")
		return err == nil
	}
	_, err := os.Stat("/var/run/docker.sock")
	return err == nil
}

// sandboxShell 返回沙箱内执行命令用的 shell 路径。
// Windows 上 Git Bash 的 MSYS 会把 "/bin/sh" 自动转换成 "C:/Program Files/Git/usr/bin/sh"，
// 导致 Docker daemon 收到无效路径。用 "//bin/sh"（双斜杠）可以绕过 MSYS 路径转换。
func sandboxShell() string {
	if runtime.GOOS == "windows" {
		return "//bin/sh"
	}
	return "/bin/sh"
}

func sandboxFingerprintForRequest(req SandboxRunRequest) string {
	payload := map[string]any{
		"provider":         req.Provider,
		"mode":             req.Mode,
		"docker_image":     req.DockerImage,
		"network_disabled": req.NetworkDisabled,
		"cpu":              req.CPU,
		"memory":           req.Memory,
		"env_keys":         sortedMapKeys(req.Env),
		"env_from_host":    uniqueSortedStrings(req.EnvFromHost),
		"readonly_mounts":  uniqueSortedStrings(req.ReadOnlyMounts),
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
