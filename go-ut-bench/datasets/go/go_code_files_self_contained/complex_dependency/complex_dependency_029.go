package main

import (
	"os"
	"strconv"
)

func IsWorker() bool {
	masterPid := os.Getenv("EINHORN_MASTER_PID")
	if masterPid == "" {
		return false
	}

	pid, err := strconv.Atoi(masterPid)
	if err != nil {
		return false
	}

	return pid == os.Getppid()
}
