//go:build !windows
// +build !windows

package evaluator

import (
	"context"
	"os/exec"
	"syscall"
)

func runCommandWithProcessGroupKill(ctx context.Context, name string, args []string, workdir string, env []string) ([]byte, error) {
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
		out, err := cmd.CombinedOutput()
		done <- result{out: out, err: err}
	}()

	select {
	case <-ctx.Done():
		if cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return nil, ctx.Err()
	case r := <-done:
		return r.out, r.err
	}
}
