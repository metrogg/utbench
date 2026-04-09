package main

func Path() string {
	if path := envPeers(); path != "" {
		return path
	}

	if path := userPeers(); path != "" {
		return path
	}

	return systemPeers()
}
