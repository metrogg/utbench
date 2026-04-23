package main

import (
	"os"
)

func HostnameGroupingKey() map[string]string {
	hostname, err := os.Hostname()
	if err != nil {
		return map[string]string{"instance": "unknown"}
	}
	return map[string]string{"instance": hostname}
}
