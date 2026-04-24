package main


// Find how many times a given substring can be found in the original string. Count overlaping cases.
// >>> HowManyTimes('', 'a')
// 0
// >>> HowManyTimes('aaa', 'a')
// 3
// >>> HowManyTimes('aaaa', 'aa')
// 3
func HowManyTimes(str string,substring string) int{
    times := 0
	for i := 0; i < (len(str) - len(substring) + 1); i++ {
		if str[i:i+len(substring)] == substring {
			times += 1
		}
	}
	return times
}
