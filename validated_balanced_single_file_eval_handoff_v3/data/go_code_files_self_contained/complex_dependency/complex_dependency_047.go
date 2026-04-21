package main

import (
	"fmt"
	"os"
	"strings"
)

func GOPATH() []string {
	paths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	if len(paths) == 0 {
		fmt.Println("GOPATH doesn't exist")
	}
	return paths
}
