//go:build windows
// +build windows

package evaluator

import (
	"context"
	"os/exec"
)

func runCommandWithProcessGroupKill(ctx context.Context, name string, args []string, workdir string, env []string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = workdir
	if len(env) > 0 {
		cmd.Env = env
	}

	// Windows 不支持进程组，直接使用 CommandContext
	return cmd.CombinedOutput()
}
