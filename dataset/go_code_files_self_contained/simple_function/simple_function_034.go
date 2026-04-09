package main

import (
	"strings"
)

func GetStackTraces() map[string][]string {
	stack := GetStackTrace(true)

	info := make(map[string][]string, 0)
	goroutine := ""

	for _, line := range strings.Split(stack, "\n") {
		line = strings.TrimSpace(line)
		line = strings.Trim(line, "\u0000")
		if strings.HasPrefix(line, "goroutine ") {
			goroutine = line
		} else if line != "" {
			info[goroutine] = append(info[goroutine], line)
		}
	}

	return info
}
