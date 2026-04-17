package main

func ModifyRootKeyLimit() error {
	value, err := readRootKeyLimit(rootKeyFile)
	if err != nil {
		return err
	}
	if value < rootKeyLimit {
		return setRootKeyLimit(rootKeyLimit)
	}
	return nil
}
