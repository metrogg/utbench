package main

import (
	"fmt"
	"os"
	"strings"
)

func commandLine() string {
	return fmt.Sprintf("$ %s %s", os.Args[0], strings.Join(os.Args[1:], " "))
}
