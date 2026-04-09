package main

import (
	"os"
)

func Pty() (*os.File, *os.File, error) {
	ptm, err := open_pty_master()
	if err != nil {
		return nil, nil, err
	}

	sname, err := Ptsname(ptm)
	if err != nil {
		return nil, nil, err
	}

	err = grantpt(ptm)
	if err != nil {
		return nil, nil, err
	}

	err = unlockpt(ptm)
	if err != nil {
		return nil, nil, err
	}

	pts, err := open_device(sname)
	if err != nil {
		return nil, nil, err
	}
	return os.NewFile(uintptr(ptm), "ptm"), os.NewFile(uintptr(pts), sname), nil
}
