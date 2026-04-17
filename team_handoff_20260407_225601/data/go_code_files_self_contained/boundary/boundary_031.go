package main

import (
	"os"
)

func ReadTarget() (string, error) {
	if target := os.Getenv("TSURU_TARGET"); target != "" {
		targets, err := getTargets()
		if err == nil {
			if val, ok := targets[target]; ok {
				return val, nil
			}
		}
		return target, nil
	}
	targetPath := JoinWithUserDir(".tsuru", "target")
	target, err := readTarget(targetPath)
	if err == errUndefinedTarget {
		copyTargetFiles()
		target, err = readTarget(JoinWithUserDir(".tsuru_target"))
	}
	return target, err
}
