package main

func peersPaths() []string {
	paths := make([]string, 0, 4)

	if env := envPeers(); env != "" {
		paths = append(paths, env)
	}

	if cwd := cwdPeers(); cwd != "" {
		paths = append(paths, cwd)
	}

	if user := userPeers(); user != "" {
		paths = append(paths, user)
	}

	if system := systemPeers(); system != "" {
		paths = append(paths, system)
	}

	return paths
}
