package main

// Check if two words have the same characters.
// >>> SameChars('eabcdzzzz', 'dddzzzzzzzddeddabc')
// true
// >>> SameChars('abcd', 'dddddddabc')
// true
// >>> SameChars('dddddddabc', 'abcd')
// true
// >>> SameChars('eabcd', 'dddddddabc')
// false
// >>> SameChars('abcd', 'dddddddabce')
// false
// >>> SameChars('eabcdzzzz', 'dddzzzzzzzddddabc')
// false
func SameChars(s0 string, s1 string) bool {
    set0 := make(map[int32]interface{})
	set1 := make(map[int32]interface{})
	for _, i := range s0 {
		set0[i] = nil
	}
	for _, i := range s1 {
		set1[i] = nil
	}
	for i, _ := range set0 {
		if _,ok:=set1[i];!ok{
			return false
		}
	}
	for i, _ := range set1 {
		if _,ok:=set0[i];!ok{
			return false
		}
	}
	return true
}
