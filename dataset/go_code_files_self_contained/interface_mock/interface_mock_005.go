package main

import (
	"encoding/base64"
	"io"
	"math/rand"
)

func shortID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)
}
