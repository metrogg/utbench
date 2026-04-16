package main

import (
	"errors"
)

func randomPassword() (string, error) {
	limit := 100
	for ; limit >= 0; limit-- {
		s, err := generatePassword()
		if err != nil {
			return "", err
		}
		if validPassword(s) {
			return s, nil
		}
	}
	return "", errors.New("failed to generate valid Windows password")
}
