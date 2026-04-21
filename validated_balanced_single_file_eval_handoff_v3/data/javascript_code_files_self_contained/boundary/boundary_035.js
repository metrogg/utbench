/*You are given a list of numbers.
  You need to return the sum of squared numbers in the given list,
  round each element in the list to the upper int(Ceiling) first.
  Examples:
  For lst = [1,2,3] the output should be 14
  For lst = [1,4,9] the output should be 98
  For lst = [1,3,5,7] the output should be 84
  For lst = [1.4,4.2,0] the output should be 29
  For lst = [-2.4,1,1] the output should be 6
  */
const sumSquares = (lst) => {
  let p = 0
  for (let i = 0; i < lst.length; i++) {
    let y = lst[i]
    if (y % 1 != 0) {
      if (y > 0) { y = y - y % 1 + 1 }
      else { y = -y; y = y - y % 1 }
    }
    p += y * y
  }
  return p
}
