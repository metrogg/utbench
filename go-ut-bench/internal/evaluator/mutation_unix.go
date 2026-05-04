//go:build !windows
// +build !windows

package evaluator

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

// runCommandLocal 在宿主机本地执行命令（Unix 实现）。
// 使用 Setpgid 创建进程组，超时后通过 Kill(-pid) 清理整个进程树。
func runCommandLocal(ctx context.Context, name string, args []string, workdir string, env []string) ([]byte, error) {
	logMutation("DEBUG-3", "run_command_start", "name", name, "args", args, "workdir", workdir)

	cmd := exec.Command(name, args...)
	cmd.Dir = workdir
	if len(env) > 0 {
		cmd.Env = env
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
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
		// 超时或取消
		if cmd.Process != nil {
			pid := cmd.Process.Pid
			logMutation("WARN", "run_command_timeout", "pid", pid, "ctx_err", ctx.Err())
			fmt.Printf("        [MUTATION-WARN] 进程超时，正在杀死进程组 (PID=%d)\n", pid)

			// 先尝试 SIGTERM，再 SIGKILL
			if err := syscall.Kill(-pid, syscall.SIGTERM); err != nil {
				logMutation("WARN", "sigterm_failed", "pid", pid, "error", err.Error())
			}

			// 等待一小段时间让进程清理
			time.Sleep(100 * time.Millisecond)

			// 强制杀死
			if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil {
				logMutation("WARN", "sigkill_failed", "pid", pid, "error", err.Error())
				// 尝试只杀死主进程
				if err2 := syscall.Kill(pid, syscall.SIGKILL); err2 != nil {
					logMutation("ERROR", "kill_failed", "pid", pid, "error", err2.Error())
				}
			}

			// 等待进程真正退出
			cmd.Process.Wait()
			logMutation("INFO", "run_command_killed", "pid", pid)
		}
		return nil, ctx.Err()
	case r := <-done:
		logMutation("DEBUG-3", "run_command_result", "out_len", len(r.out), "err", r.err)
		return r.out, r.err
	}
}
