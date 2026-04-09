package main

import (
	"time"
)

func MasterToken() (string, error) {
	masterpublic, err := GetMasterPublicKey()
	if err != nil {
		return "", err
	}
	keypem, err := PEMFromRSAPublicKey(masterpublic, nil)
	if err != nil {
		return "", err
	}

	signed, _, err := CreateJWTIdentity("", "", true, true, keypem, time.Hour)
	if err != nil {
		return "", err
	}

	return signed, nil
}
