package main

import (
	"sync"
)

func New() *Writer {
	termWidth, _ = getTermSize()
	if termWidth != 0 {
		overFlowHandled = true
	}

	return &Writer{
		Out:             Out,
		RefreshInterval: RefreshInterval,

		mtx: &sync.Mutex{},
	}
}
