package main

import (
	"bytes"
	"encoding/hex"
	"io"
)

func GenerateV1PSK() (io.Reader, error) {
	psk, err := GenerateV1Bytes()
	if err != nil {
		return nil, err
	}

	hexPsk := make([]byte, len(psk)*2)
	hex.Encode(hexPsk, psk[:])

	// just a shortcut to NewReader
	nr := func(b []byte) io.Reader {
		return bytes.NewReader(b)
	}
	return io.MultiReader(nr(pathPSKv1), newLine(), nr([]byte("/base16/")), newLine(), nr(hexPsk)), nil
}
