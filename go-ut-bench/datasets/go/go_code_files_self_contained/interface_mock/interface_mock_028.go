package main

import (
	"fmt"
	"math/rand"
	"time"
)

func NewTraceID() string {
	b := make([]byte, 12)
	rand.Read(b)
	return fmt.Sprintf("%d-%x-%s", 1, time.Now().Unix(), fmt.Sprintf("%x", b))
}
