package main

import (
	"fmt"
	"os"
)

func generateConsumerInstanceID() (string, error) {
	uuid, err := generateUUID()
	if err != nil {
		return "", err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", hostname, uuid), nil
}
