package main

import (
	"fmt"
	"os"
)

func StdinTokenProvider() (string, error) {
	var v string
	fmt.Fprintf(os.Stderr, "Assume Role MFA token code: ")
	_, err := fmt.Scanln(&v)

	return v, err
}
