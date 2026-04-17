package main

func machineID() (string, error) {
	id, err := readFile(dbusPath)
	if err != nil {
		// try fallback path
		id, err = readFile(dbusPathEtc)
	}
	if err != nil {
		return "", err
	}
	return trim(string(id)), nil
}
