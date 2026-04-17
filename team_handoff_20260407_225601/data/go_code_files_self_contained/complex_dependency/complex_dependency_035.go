package main

import (
	"errors"
	"os"
)

func CheckPushConf() error {
	if !PushConf.Ios.Enabled && !PushConf.Android.Enabled {
		return errors.New("Please enable iOS or Android config in yml config")
	}

	if PushConf.Ios.Enabled {
		if PushConf.Ios.KeyPath == "" && PushConf.Ios.KeyBase64 == "" {
			return errors.New("Missing iOS certificate key")
		}

		// check certificate file exist
		if PushConf.Ios.KeyPath != "" {
			if _, err := os.Stat(PushConf.Ios.KeyPath); os.IsNotExist(err) {
				return errors.New("certificate file does not exist")
			}
		}
	}

	if PushConf.Android.Enabled {
		if PushConf.Android.APIKey == "" {
			return errors.New("Missing Android API Key")
		}
	}

	return nil
}
