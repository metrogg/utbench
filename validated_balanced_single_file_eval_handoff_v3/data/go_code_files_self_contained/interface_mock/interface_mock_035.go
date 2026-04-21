package main

import (
	"fmt"
	"math/rand"
)

func NewID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
