package main

import (
	"os"
)

func SecretRingFile() string {
	if flagSecretRing != "" {
		return flagSecretRing
	}
	if e := os.Getenv("CAMLI_SECRET_RING"); e != "" {
		return e
	}
	return identitySecretRing()
}
