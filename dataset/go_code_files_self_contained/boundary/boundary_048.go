package main

import (
	"fmt"
)

func gMSARegistryValueName() (string, error) {
	randomSuffix, err := randomString(gMSARegistryValueNameSuffixRandomBytes)

	if err != nil {
		return "", fmt.Errorf("error when generating gMSA registry value name: %v", err)
	}

	return gMSARegistryValueNamePrefix + randomSuffix, nil
}
