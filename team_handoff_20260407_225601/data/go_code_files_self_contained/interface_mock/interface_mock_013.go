package main

import (
	"os"
)

func Delete() error {
	filename := locateSettingsFile()

	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if err := os.Remove(filename); err != nil {
		return err
	}

	return nil
}
