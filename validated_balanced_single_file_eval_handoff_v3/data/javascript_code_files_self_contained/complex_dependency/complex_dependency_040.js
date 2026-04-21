/*
  You are given a positive integer n. You have to create an integer array a of length n.
      For each i (1 ≤ i ≤ n), the value of a[i] = i * i - i + 1.
      Return the number of triples (a[i], a[j], a[k]) of a where i < j < k, 
  and a[i] + a[j] + a[k] is a multiple of 3.

  Example :
      Input: n = 5
      Output: 1
      Explanation: 
      a = [1, 3, 7, 13, 21]
      The only valid triple is (1, 7, 13).
  */
const getMaxTriples = (n) => {
  let y = []
  for (let i = 1; i <= n; i++) {
    y.push(i * i - i + 1)
  }
  let u = 0
  for (let i = 0; i < n - 2; i++) {
    for (let j = i + 1; j < n - 1; j++) {
      for (let k = j + 1; k < n; k++) {
        if ((y[i] + y[j] + y[k]) % 3 == 0) { u++ }
      }
    }
  }
  return u
}
