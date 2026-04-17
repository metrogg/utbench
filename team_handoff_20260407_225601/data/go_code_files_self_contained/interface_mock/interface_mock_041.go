package main

import (
	"io"
)

func NewTruecolorLolWriter() io.Writer {
	colorMode := ColorModeTrueColor
	if noColor {
		colorMode = ColorMode0
	}
	return &Writer{
		Output:    stdout,
		ColorMode: colorMode,
	}
}
