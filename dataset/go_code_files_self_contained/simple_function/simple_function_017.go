package main

func GetDefaultAddr() (string, error) {
	// get the ip address by local hostname
	localIP, err := GetMyAddr()
	if err == nil && IsAddrLocal(localIP) {
		return localIP, nil
	}

	// Return first available address if we could not find by hostname
	return GetFirstLocalAddr()
}
