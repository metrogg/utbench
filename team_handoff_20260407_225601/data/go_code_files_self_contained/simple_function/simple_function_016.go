package main

func IsAccessible() bool {
	if pkcs11Lib == "" {
		return false
	}
	ctx, session, err := SetupHSMEnv(pkcs11Lib, defaultLoader)
	if err != nil {
		return false
	}
	defer cleanup(ctx, session)
	return true
}
