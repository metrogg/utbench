//go:build windows
// +build windows

package evaluator

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

func runCommandWithProcessGroupKill(ctx context.Context, name string, args []string, workdir string, env []string) ([]byte, error) {
	logMutation("DEBUG-3", "run_command_start", "name", name, "args", args, "workdir", workdir)

	cmd := exec.Command(name, args...)
	cmd.Dir = workdir
	if len(env) > 0 {
		cmd.Env = env
	}

	// Windows 不支持进程组，使用 CommandContext 进行超时控制
	// 注意：这可能无法杀死子进程
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
			cmd.Process.Kill()
			logMutation("INFO", "run_command_killed_windows", "pid", cmd.Process.Pid)
		}
		return nil, ctx.Err()
	case r := <-done:
		logMutation("DEBUG-3", "run_command_result", "out_len", len(r.out), "err", r.err)
		return r.out, r.err
	}
}
