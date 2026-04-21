package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

func NewIdempotencyKey() string {
	now := time.Now().UnixNano()
	buf := make([]byte, 4)
	rand.Read(buf)
	return fmt.Sprintf("%v_%v", now, base64.URLEncoding.EncodeToString(buf)[:6])
}
