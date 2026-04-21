package main

import (
	"os"
	"strings"
)

func childEnv() []string {
	var env []string
	for _, str := range os.Environ() {
		if strings.HasPrefix(str, "STNORESTART=") {
			continue
		}
		if strings.HasPrefix(str, "STMONITORED=") {
			continue
		}
		env = append(env, str)
	}
	env = append(env, "STMONITORED=yes")
	return env
}
