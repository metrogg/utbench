package main

import (
	"errors"
)

func zmqerr() error {
	eno := C.my_errno()
	switch eno {
	case C.ETERM:
		return ErrTerminated
	case C.EAGAIN:
		return ErrTimeout
	case C.EINTR:
		return ErrInterrupted
	}
	str := C.GoString(C.zmq_strerror(eno))
	return errors.New(str)
}
