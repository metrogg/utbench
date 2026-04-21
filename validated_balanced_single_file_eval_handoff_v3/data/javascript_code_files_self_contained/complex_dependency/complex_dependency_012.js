/*
  Given an array arr of integers and a positive integer k, return a sorted list 
  of length k with the maximum k numbers in arr.

  Example 1:

      Input: arr = [-3, -4, 5], k = 3
      Output: [-4, -3, 5]

  Example 2:

      Input: arr = [4, -4, 4], k = 2
      Output: [4, 4]

  Example 3:

      Input: arr = [-3, 2, 1, 2, -1, -2, 1], k = 1
      Output: [2]

  Note:
      1. The length of the array will be in the range of [1, 1000].
      2. The elements in the array will be in the range of [-1000, 1000].
      3. 0 <= k <= len(arr)
  */
const maximum = (arr, k) => {
  let p = arr
  for (let j = 0; j < p.length; j++) {
    let ind = j
    for (let k = j + 1; k < p.length; k++) {
      if (p[k] < p[ind]) {
        ind = k
      }
    }
    if (ind > j) {
      let tmp = p[j]
      p[j] = p[ind]
      p[ind] = tmp
    }
  }
  if (k == 0) { return [] }
  return p.slice(-k)
}
