package main

import (
	"fmt"
	"os"
)

func StorageRootFromEnv() (string, error) {
	storageRoot, ok := os.LookupEnv(PachRootEnvVar)
	if !ok {
		return "", fmt.Errorf("%s not found", PachRootEnvVar)
	}
	storageBackend, ok := os.LookupEnv(StorageBackendEnvVar)
	if !ok {
		return "", fmt.Errorf("%s not found", StorageBackendEnvVar)
	}
	// These storage backends do not like leading slashes
	switch storageBackend {
	case Amazon:
		fallthrough
	case Minio:
		if len(storageRoot) > 0 && storageRoot[0] == '/' {
			storageRoot = storageRoot[1:]
		}
	}
	return storageRoot, nil
}
