package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
)

func NewMAC() (string, error) {
	mac := make([]byte, 6)
	n, err := io.ReadFull(rand.Reader, mac)
	if n != len(mac) || err != nil {
		return "", err
	}
	e := hex.EncodeToString(mac)

	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", e[0:2], e[2:4], e[4:6], e[6:8], e[8:10], e[10:12]), nil
}
