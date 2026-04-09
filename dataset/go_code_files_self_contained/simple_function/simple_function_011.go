package main

import (
	"os"
)

func ExplicitSecretRingFile() (string, bool) {
	if !secretRingFlagAdded {
		panic("proper use of ExplicitSecretRingFile requires exposing flagSecretRing with AddSecretRingFlag")
	}
	if flagSecretRing != "" {
		return flagSecretRing, true
	}
	if e := os.Getenv("CAMLI_SECRET_RING"); e != "" {
		return e, true
	}
	return "", false
}
