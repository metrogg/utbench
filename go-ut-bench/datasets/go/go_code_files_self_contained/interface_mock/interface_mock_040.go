package main

import (
	"io"
)

func NewLolWriter() io.Writer {
	colorMode := ColorMode256
	if noColor {
		colorMode = ColorMode0
	}
	return &Writer{
		Output:    stdout,
		ColorMode: colorMode,
	}
}
