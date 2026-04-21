package main

import (
	"fmt"
	"os"
)

func defaultUserName() string {
	userName := os.Getenv("USER")
	if userName == "" {
		userName = os.Getenv("USERNAME")
	}
	if userName == "" {
		userName = fmt.Sprintf("user%d", os.Getuid())
	}
	return userName
}
