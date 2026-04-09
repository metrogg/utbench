package main

func isGRPCEnabled() bool {
	if GRPCPort != nil && *GRPCPort != 0 {
		return true
	}

	if SocketFile != nil && *SocketFile != "" {
		return true
	}

	return false
}
