package main

func GetVersion() (int, error) {
	output, err := cmd("", "--version")
	if err != nil {
		return -1, err
	}

	return parseVersion(output)
}
