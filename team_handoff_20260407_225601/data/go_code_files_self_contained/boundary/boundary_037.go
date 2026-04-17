package main

import (
	"os"
)

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}

	name, _ := parseHostname(hostname)
	return sanitizeHostname(name)
}
