package main

func Terminate() error {
	paErr := C.Pa_Terminate()
	if paErr != C.paNoError {
		return newError(paErr)
	}
	initialized--
	if initialized <= 0 {
		initialized = 0
		cached = false
	}
	return nil
}
