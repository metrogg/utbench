package main

import (
	"os"
)

func Sniff() string {
	if bind := os.Getenv("GOJI_BIND"); bind != "" {
		return bind
	} else if usingEinhorn() {
		return "einhorn@0"
	} else if usingSystemd() {
		return "fd@3"
	} else if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	return ""
}
