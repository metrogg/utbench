package main

import (
	"os"
	"strings"
)

func Environ() []string {
	e := os.Environ()
	for i, n := 0, 0; i < len(e); i++ {
		if !strings.HasPrefix(e[i], EnvPrefix) {
			e[n] = e[i]
			n++
		}
	}
	return e
}
