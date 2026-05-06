//go:build !windows

package dataset

import "os/exec"

func hideCommandWindow(cmd *exec.Cmd) {}
