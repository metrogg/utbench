package main

import (
	"errors"
	"os"
	"strconv"
)

func setManagerSessionTTL() error {
	ttlStr := os.Getenv("tidb_manager_ttl")
	if len(ttlStr) == 0 {
		return nil
	}
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil {
		return errors.Trace(err)
	}
	ManagerSessionTTL = ttl
	return nil
}
