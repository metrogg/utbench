package main


// Implement the Function F that takes n as a parameter,
// and returns a list oF size n, such that the value oF the element at index i is the Factorial oF i iF i is even
// or the sum oF numbers From 1 to i otherwise.
// i starts From 1.
// the Factorial oF i is the multiplication oF the numbers From 1 to i (1 * 2 * ... * i).
// Example:
// F(5) == [1, 2, 6, 24, 15]
func F(n int) []int {
    ret := make([]int, 0, 5)
    for i:=1;i<n+1;i++{
        if i%2 == 0 {
            x := 1
            for j:=1;j<i+1;j++{
                x*=j
            }
            ret = append(ret, x)
        }else {
            x := 0
            for j:=1;j<i+1;j++{
                x+=j
            }
            ret = append(ret, x)
        }
    }
    return ret
}
