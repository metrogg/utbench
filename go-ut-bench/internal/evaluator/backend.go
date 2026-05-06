package evaluator

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// EvalBackend 抽象评测命令的执行后端。
// 当前实现：LocalBackend（直接在宿主机执行 shell 命令）、DockerBackend（在评测容器内执行）。
//
// 用途：让 LanguageEvaluator 的编译/测试/覆盖率/变异命令
// 可以在不同的执行环境中运行，而不改变评测逻辑本身。
type EvalBackend interface {
	// Name 返回后端名称，用于日志和 fingerprint。
	Name() string

	// RunCommand 在指定工作目录执行命令，返回 stdout+stderr 合并输出和错误。
	// 超时由 ctx 控制。
	RunCommand(ctx context.Context, workdir string, command []string, env []string) ([]byte, error)
}

// LocalBackend 在宿主机本地执行评测命令。
// 使用 runCommandWithProcessGroupKill 确保超时后正确清理进程组。
type LocalBackend struct{}

func NewLocalBackend() *LocalBackend {
	return &LocalBackend{}
}

func (b *LocalBackend) Name() string {
	return "local"
}

func (b *LocalBackend) RunCommand(ctx context.Context, workdir string, command []string, env []string) ([]byte, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	return runCommandLocal(ctx, command[0], command[1:], workdir, env)
}

// DockerBackend 在评测容器内执行命令。
// 将工作目录挂载到容器内，执行编译/测试/覆盖率/变异命令，收集输出。
// 适用于宿主机缺少某些语言工具链的场景，或需要隔离评测环境的场景。
type DockerBackend struct {
	Image           string // 评测镜像，如 "utbench:latest"
	NetworkDisabled bool   // 是否禁用网络
	CPU             string // CPU 限制
	Memory          string // 内存限制
}

func NewDockerBackend(image string) *DockerBackend {
	return &DockerBackend{Image: image}
}

func (b *DockerBackend) Name() string {
	return "docker"
}

func (b *DockerBackend) RunCommand(ctx context.Context, workdir string, command []string, env []string) ([]byte, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	args := []string{"run", "--rm"}
	if b.NetworkDisabled {
		args = append(args, "--network", "none")
	}
	if b.CPU != "" {
		args = append(args, "--cpus", b.CPU)
	}
	if b.Memory != "" {
		args = append(args, "--memory", b.Memory)
	}
	// 注入环境变量
	keys := make([]string, 0, len(env))
	for _, e := range env {
		if idx := strings.IndexByte(e, '='); idx > 0 {
			keys = append(keys, e[:idx])
		}
	}
	sort.Strings(keys)
	for _, e := range env {
		args = append(args, "-e", e)
	}
	// 挂载工作目录
	args = append(args, "-v", workdir+":/work", "-w", "/work")
	image := strings.TrimSpace(b.Image)
	if image == "" {
		image = "utbench:latest"
	}
	args = append(args, image)
	// 将 command 拼接为 shell 命令
	shellCmd := shellJoinEval(command)
	args = append(args, "/bin/sh", "-c", shellCmd)

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "docker", args...)
	hideCommandWindow(cmd)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	// 合并 stdout + stderr，与 LocalBackend 行为一致
	output := append(stdout.Bytes(), stderr.Bytes()...)
	if ctx.Err() != nil {
		return output, fmt.Errorf("docker eval command timed out: %w", ctx.Err())
	}
	return output, err
}

func shellJoinEval(args []string) string {
	quoted := make([]string, len(args))
	for i, arg := range args {
		if strings.ContainsAny(arg, " \t\"'\\") {
			quoted[i] = "'" + strings.ReplaceAll(arg, "'", "'\"'\"'") + "'"
		} else {
			quoted[i] = arg
		}
	}
	return strings.Join(quoted, " ")
}

// evalBackend 全局评测后端实例。
// 当前默认为 LocalBackend；可通过 SetEvalBackend 切换为 DockerBackend。
var evalBackend EvalBackend = NewLocalBackend()

// SetEvalBackend 设置全局评测后端。
// 在 main 或初始化阶段调用，运行期间不应变更。
func SetEvalBackend(b EvalBackend) {
	if b != nil {
		evalBackend = b
	}
}

// GetEvalBackend 获取当前评测后端。
func GetEvalBackend() EvalBackend {
	return evalBackend
}

// runCommandWithProcessGroupKill 兼容包装函数。
// 所有现有的评测调用点（go_eval / java_eval / cpp_eval / mutation / service）仍然调用此函数，
// 但它内部委托给 evalBackend.RunCommand，从而允许未来切换到 Docker/Remote 后端。
func runCommandWithProcessGroupKill(ctx context.Context, name string, args []string, workdir string, env []string) ([]byte, error) {
	return evalBackend.RunCommand(ctx, workdir, append([]string{name}, args...), env)
}

// defaultCommandTimeout 默认命令超时。
const defaultCommandTimeout = 300 * time.Second
