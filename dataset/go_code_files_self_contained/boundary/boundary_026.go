package main

func getIptablesHasCheckCommand() (bool, error) {
	vstring, err := getIptablesVersionString()
	if err != nil {
		return false, err
	}

	v1, v2, v3, err := extractIptablesVersion(vstring)
	if err != nil {
		return false, err
	}

	return iptablesHasCheckCommand(v1, v2, v3), nil
}
