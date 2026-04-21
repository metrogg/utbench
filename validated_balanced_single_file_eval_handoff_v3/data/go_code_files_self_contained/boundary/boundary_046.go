package main

import (
	"math/rand"
)

func NewSeed() ([]byte, error) {
	// Well that easy, just make go read 256 random bytes into a slice
	s := make([]byte, 256)
	_, err := rand.Read(s)
	return s, err
}
