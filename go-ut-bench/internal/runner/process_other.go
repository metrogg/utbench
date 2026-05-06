//go:build !windows

package runner

import "os/exec"

func hideCommandWindow(cmd *exec.Cmd) {}
