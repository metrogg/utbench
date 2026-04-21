package main

import (
    "bytes"
    "strings"
)

// Given a string s and a natural number n, you have been tasked to implement
// a function that returns a list of all words from string s that contain exactly
// n consonants, in order these words appear in the string s.
// If the string s is empty then the function should return an empty list.
// Note: you may assume the input string contains only letters and spaces.
// Examples:
// SelectWords("Mary had a little lamb", 4) ==> ["little"]
// SelectWords("Mary had a little lamb", 3) ==> ["Mary", "lamb"]
// SelectWords("simple white space", 2) ==> []
// SelectWords("Hello world", 4) ==> ["world"]
// SelectWords("Uncle sam", 3) ==> ["Uncle"]
func SelectWords(s string, n int) []string {
    result := make([]string, 0)
    for _, word := range strings.Fields(s) {
        n_consonants := 0
        lower := strings.ToLower(word)
        for i := 0;i < len(word); i++ {
            if !bytes.Contains([]byte("aeiou"), []byte{lower[i]}) {
                n_consonants++
            }
        }
        if n_consonants == n{
            result = append(result, word)
        }
    }
    return result
}
