package main

func Summary() string {
	if Version != "" && GitInfo != "" {
		return Version + ", " + GitInfo
	}
	if GitInfo != "" {
		return GitInfo
	}
	if Version != "" {
		return Version
	}
	return "unknown"
}
