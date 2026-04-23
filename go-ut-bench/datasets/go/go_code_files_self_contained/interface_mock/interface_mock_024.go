package main

import (
	"os"
	"strings"
)

func GetSysLoggerTag() string {
	procName := os.Args[0]
	if strings.ContainsRune(procName, os.PathSeparator) {
		parts := strings.FieldsFunc(procName, func(c rune) bool {
			return c == os.PathSeparator
		})
		procName = parts[len(parts)-1]
	}
	return procName
}
