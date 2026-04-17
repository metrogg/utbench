package main

import (
	"encoding/binary"
	"io"
	"math/rand"
)

func randomNonZero() uint64 {
	buf := make([]byte, 8)
	n, err := io.ReadFull(rand.Reader, buf)
	if err != nil || n != 8 {
		panic("Unable to fully read from rand.Reader")
	}
	u, x := binary.Uvarint(buf)
	if u == 0 || x == 0 || x < 0 {
		return randomNonZero()
	}
	return u
}
