package main

// Create a function that takes 3 numbers.
// Returns true if one of the numbers is equal to the sum of the other two, and all numbers are integers.
// Returns false in any other cases.
// 
// Examples
// AnyInt(5, 2, 7) ➞ true
// 
// AnyInt(3, 2, 2) ➞ false
// 
// AnyInt(3, -2, 1) ➞ true
// 
// AnyInt(3.6, -2.2, 2) ➞ false
func AnyInt(x, y, z interface{}) bool {
    xx, ok := x.(int)
    if !ok {
        return false
    }
    yy, ok := y.(int)
    if !ok {
        return false
    }
    zz, ok := z.(int)
    if !ok {
        return false
    }

    if (xx+yy == zz) || (xx+zz == yy) || (yy+zz == xx) {
        return true
    }
    return false
}
