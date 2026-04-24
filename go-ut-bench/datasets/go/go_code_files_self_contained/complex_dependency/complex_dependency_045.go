package main

import (
	"os"
)

func GetCredentials() (string, string) {
	username := os.Getenv("CF_INT_USERNAME")
	if username == "" {
		username = "admin"
	}
	password := os.Getenv("CF_INT_PASSWORD")
	if password == "" {
		password = "admin"
	}
	return username, password
}
