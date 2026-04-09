package main

import (
	"math/rand"
	"time"
)

func generateNonce() string {
	rand.Seed(time.Now().UTC().UnixNano())
	// For convenience, use a set of chars we don't need to url-escape
	var letters = []rune("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
