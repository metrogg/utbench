package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		stderr("usage: %s [main|sidekick|preStart|postStop]", os.Args[0])
		os.Exit(64)
	}
	mode := os.Args[1]
	var res results
	switch strings.ToLower(mode) {
	case "main":
		res = validateMain()
	case "sidekick":
		res = validateSidekick()
	case "prestart":
		res = validatePrestart()
	case "poststop":
		res = validatePoststop()
	default:
		stderr("unrecognized mode: %s", mode)
		os.Exit(64)
	}
	if len(res) == 0 {
		fmt.Printf("%s OK\n", mode)
		os.Exit(0)
	}
	fmt.Printf("%s FAIL\n", mode)
	for _, err := range res {
		fmt.Fprintln(os.Stderr, "==>", err)
	}
	os.Exit(1)
}
