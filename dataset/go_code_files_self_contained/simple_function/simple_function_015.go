package main

import (
	"fmt"
)

func NewTarTimeoutError() error {
	return Error{
		Message:    fmt.Sprintf("timeout waiting for tar stream"),
		Details:    nil,
		ErrorCode:  TarTimeoutError,
		Suggestion: "check the Source-To-Image scripts if it accepts tar stream for assemble and sends for save-artifacts",
	}
}
