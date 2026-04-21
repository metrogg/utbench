/*
  You are given an array arr of integers and you need to return
  sum of magnitudes of integers multiplied by product of all signs
  of each number in the array, represented by 1, -1 or 0.
  Note: return null for empty arr.

  Example:
  >>> prodSigns([1, 2, 2, -4]) == -9
  >>> prodSigns([0, 1]) == 0
  >>> prodSigns([]) == null
  */
const prodSigns = (arr) => {
  if (arr.length == 0) { return null }
  let n = 1
  let s = 0
  for (let i = 0; i < arr.length; i++) {
    s += arr[i]
    if (arr[i] == 0) { return 0 }
    if (arr[i] < 0) { n = -n; s -= 2 * arr[i] }
  }
  return s * n
}
