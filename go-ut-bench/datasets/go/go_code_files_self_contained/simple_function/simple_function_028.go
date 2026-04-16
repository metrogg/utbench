package main

import (
	"os"
	"strings"
)

func environ() map[string]string {
	if env == nil {
		env = make(map[string]string)
		envVars := os.Environ()
		for _, envVar := range envVars {
			kv := strings.SplitN(envVar, "=", 2)
			if len(kv) != 2 {
				continue
			}
			env[kv[0]] = kv[1]
		}
	}
	return env
}
