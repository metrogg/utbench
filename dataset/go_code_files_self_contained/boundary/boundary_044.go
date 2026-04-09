package main

func MustAbsoluteDir() string {
	path, err := absoluteDir()
	if err != nil {
		panic(err)
	}
	return path
}
