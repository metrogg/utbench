package main

import (
	"os"
)

func GetUsernameFromEnv() string {
	user_env := []string{"USER", "USERNAME", "LOGNAME", "LNAME"}

	for _, val := range user_env {
		name := os.Getenv(val)
		if name != "" {
			return name
		}
	}
	return ""
}
