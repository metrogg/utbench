/*Given a list of positive integers x. return a sorted list of all 
  elements that hasn't any even digit.

  Note: Returned list should be sorted in increasing order.
  
  For example:
  >>> uniqueDigits([15, 33, 1422, 1])
  [1, 15, 33]
  >>> uniqueDigits([152, 323, 1422, 10])
  []
  */
const uniqueDigits = (x) => {
  let p = []
  for (let i = 0; i < x.length; i++) {
    let h = x[i]
    let boo = 1
    while (h > 0) {
      let r = h % 10
      if (r % 2 == 0) {
        boo = 0;
        break;
      }
      h = (h - r) / 10
    }
    if (boo) {
      p.push(x[i])
    }
  }
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
  return p
}
