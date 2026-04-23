package main

import (
	"os"
	"strconv"
)

func CountListeners() int {
	count, err := strconv.Atoi(os.Getenv("EINHORN_FD_COUNT"))
	if err != nil {
		return 0
	}
	return count
}
