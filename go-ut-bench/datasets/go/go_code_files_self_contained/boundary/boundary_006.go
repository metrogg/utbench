package main

import (
	"math/rand"
)

func uuid() []byte {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic("cue/hosted: uuid() failed to read random bytes")
	}

	// The following bit twiddling is outlined in RFC 4122.  In short, it
	// identifies the UUID as a v4 random UUID.
	uuid[6] = (4 << 4) | (0xf & uuid[6])
	uuid[8] = (8 << 4) | (0x3f & uuid[8])
	return uuid
}
