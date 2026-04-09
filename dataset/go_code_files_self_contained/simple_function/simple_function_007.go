package main

import (
	"errors"
)

func GetError() error {
	if err := C.SDL_GetError(); err != nil {
		gostr := C.GoString(err)
		// SDL_GetError returns "an empty string if there hasn't been an error message"
		if len(gostr) > 0 {
			return errors.New(gostr)
		}
	}
	return nil
}
