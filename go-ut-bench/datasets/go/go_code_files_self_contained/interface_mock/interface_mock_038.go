package main

import (
	"io"
	"os"
)

func ErrorWriter() io.Writer {
	return &ColorWriter{w: os.Stderr, fd: os.Stderr.Fd(), mutex: &stdErrMutex, lastFgColor: fgWhite}
}
