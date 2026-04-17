package main

func GetIoThreads() (int, error) {
	if initVersionError != nil {
		return 0, initVersionError
	}
	if initContextError != nil {
		return 0, initContextError
	}
	return nr_of_threads, nil
}
