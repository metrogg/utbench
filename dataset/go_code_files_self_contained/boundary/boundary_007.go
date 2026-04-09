package main

import (
	"crypto/x509"
)

func LoadSystemCAs() (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	for _, keychain := range certKeychains() {
		err := addCertsFromKeychain(pool, keychain)
		if err != nil {
			return nil, err
		}
	}

	return pool, nil
}
