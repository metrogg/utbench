package main

import (
    "strings"
)

// Write a function that takes a message, and Encodes in such a
// way that it swaps case of all letters, replaces all vowels in
// the message with the letter that appears 2 places ahead of that
// vowel in the english alphabet.
// Assume only letters.
// 
// Examples:
// >>> Encode('test')
// 'TGST'
// >>> Encode('This is a message')
// 'tHKS KS C MGSSCGG'
func Encode(message string) string {
    vowels := "aeiouAEIOU"
    vowels_replace := make(map[rune]rune)
    for _, c := range vowels {
        vowels_replace[c] = c + 2
    }
    result := make([]rune, 0, len(message))
    for _, c := range message {
        if 'a' <= c && c <= 'z' {
            c += 'A' - 'a'
        } else if 'A' <= c && c <= 'Z' {
            c += 'a' - 'A'
        }
        if strings.ContainsRune(vowels, c) {
            result = append(result, vowels_replace[c])
        } else {
            result = append(result, c)
        }
    }
    return string(result)
}
