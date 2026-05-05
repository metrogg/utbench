//go:build !windows

package web

import "os/exec"

func hideCommandWindow(cmd *exec.Cmd) {}
