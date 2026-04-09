package main

func IsWorkingCopyDirty() (bool, error) {
	bare, err := IsBare()
	if bare || err != nil {
		return false, err
	}

	out, err := gitSimple("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return len(out) != 0, nil
}
