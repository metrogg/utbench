package main

import (
	"math/rand"
	"time"
)

func SimpleRand() *rand.Rand {
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	return rand.New(source)
}
