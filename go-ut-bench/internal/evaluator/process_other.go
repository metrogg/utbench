//go:build !windows

package evaluator

import "os/exec"

func hideCommandWindow(cmd *exec.Cmd) {}
