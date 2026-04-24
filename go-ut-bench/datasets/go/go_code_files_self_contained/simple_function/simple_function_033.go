package main

import (
    "strings"
)

// Create a function Encrypt that takes a string as an argument and
// returns a string Encrypted with the alphabet being rotated.
// The alphabet should be rotated in a manner such that the letters
// shift down by two multiplied to two places.
// For example:
// Encrypt('hi') returns 'lm'
// Encrypt('asdfghjkl') returns 'ewhjklnop'
// Encrypt('gf') returns 'kj'
// Encrypt('et') returns 'ix'
func Encrypt(s string) string {
    d := "abcdefghijklmnopqrstuvwxyz"
    out := make([]rune, 0, len(s))
    for _, c := range s {
        pos := strings.IndexRune(d, c)
        if pos != -1 {
            out = append(out, []rune(d)[(pos+2*2)%26])
        } else {
            out = append(out, c)
        }
    }
    return string(out)
}
