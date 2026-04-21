package main

// You are given 2 words. You need to return true if the second word or any of its rotations is a substring in the first word
// CycpatternCheck("abcd","abd") => false
// CycpatternCheck("hello","ell") => true
// CycpatternCheck("whassup","psus") => false
// CycpatternCheck("abab","baa") => true
// CycpatternCheck("efef","eeff") => false
// CycpatternCheck("himenss","simen") => true
func CycpatternCheck(a , b string) bool {
    l := len(b)
    pat := b + b
    for i := 0;i < len(a) - l + 1; i++ {
        for j := 0;j<l + 1;j++ {
            if a[i:i+l] == pat[j:j+l] {
                return true
            }
        }
    }
    return false
}
