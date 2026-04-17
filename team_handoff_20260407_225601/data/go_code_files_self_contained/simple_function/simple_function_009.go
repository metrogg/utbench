package main

func DiscoveryToken() string {
	if len(modes) == 0 {
		return Token()
	}
	if _, ok := modes[0].(None); ok {
		return OmitAuthToken
	}
	return Token()
}
