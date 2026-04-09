package main

import (
	"path/filepath"
	"unicode/utf8"
)

func listSeparator() string {
	if pathListSep == "" {
		len := utf8.RuneLen(filepath.ListSeparator)
		if len < 0 {
			panic("filepath.ListSeparator is not valid utf8?!")
		}
		buf := make([]byte, len)
		len = utf8.EncodeRune(buf, filepath.ListSeparator)
		pathListSep = string(buf[:len])
	}

	return pathListSep
}
