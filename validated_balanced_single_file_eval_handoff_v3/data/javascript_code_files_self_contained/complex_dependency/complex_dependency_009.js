/*
  In this Kata, you have to sort an array of non-negative integers according to
  number of ones in their binary representation in ascending order.
  For similar number of ones, sort based on decimal value.

  It must be implemented like this:
  >>> sortArray([1, 5, 2, 3, 4]) == [1, 2, 3, 4, 5]
  >>> sortArray([-2, -3, -4, -5, -6]) == [-6, -5, -4, -3, -2]
  >>> sortArray([1, 0, 2, 3, 4]) [0, 1, 2, 3, 4]
  */
const sortArray = (arr) => {
  let p = arr
  for (let j = 0; j < p.length; j++) {
    let ind = j
    for (let k = j + 1; k < p.length; k++) {
      let w1 = p[ind].toString(2)
      let f1 = 0
      for (let u = 0; u < w1.length; u++) {
        if (w1[u] == '1') { f1++ }
      }
      let w2 = p[k].toString(2)
      let f2 = 0
      for (let u = 0; u < w2.length; u++) {
        if (w2[u] == '1') { f2++ }
      }
      if (f2 < f1 || (f1 == f2 && p[k] < p[ind])) {
        ind = k
      }
    }
    if (ind > j) {
      let tmp = p[j]
      p[j] = p[ind]
      p[ind] = tmp
    }
  }
  return p
}
