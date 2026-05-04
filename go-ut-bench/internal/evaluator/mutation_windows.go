//go:build windows
// +build windows

package evaluator

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

// runCommandLocal 在宿主机本地执行命令（Windows 实现）。
// Windows 不支持进程组，使用 taskkill /T /F /PID 杀掉整个进程树。
func runCommandLocal(ctx context.Context, name string, args []string, workdir string, env []string) ([]byte, error) {
	logMutation("DEBUG-3", "run_command_start", "name", name, "args", args, "workdir", workdir)

	cmd := exec.Command(name, args...)
	cmd.Dir = workdir
	if len(env) > 0 {
		cmd.Env = env
	}

	type result struct {
		out []byte
		err error
	}
	done := make(chan result, 1)

	go func() {
		startTime := time.Now()
		out, err := cmd.CombinedOutput()
		elapsed := time.Since(startTime)
		logMutation("DEBUG-3", "run_command_goroutine_done", "elapsed_ms", elapsed.Milliseconds(), "err", err)
		done <- result{out: out, err: err}
	}()

	select {
	case <-ctx.Done():
		logMutation("WARN", "run_command_timeout_windows", "ctx_err", ctx.Err())
		fmt.Printf("        [MUTATION-WARN] Windows 进程超时 (context cancelled)\n")
		if cmd.Process != nil {
			// 使用 taskkill /T 杀掉整个进程树（包括子进程）
			pid := strconv.Itoa(cmd.Process.Pid)
			killCmd := exec.Command("taskkill", "/T", "/F", "/PID", pid)
			_ = killCmd.Run()
			logMutation("INFO", "run_command_killed_windows_tree", "pid", pid)
		}
		return nil, ctx.Err()
	case r := <-done:
		logMutation("DEBUG-3", "run_command_result", "out_len", len(r.out), "err", r.err)
		return r.out, r.err
	}
}
