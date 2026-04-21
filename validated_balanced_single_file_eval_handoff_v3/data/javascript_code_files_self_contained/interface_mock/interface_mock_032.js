/*Write a function that accepts a list of strings as a parameter,
  deletes the strings that have odd lengths from it,
  and returns the resulted list with a sorted order,
  The list is always a list of strings and never an array of numbers,
  and it may contain duplicates.
  The order of the list should be ascending by length of each word, and you
  should return the list sorted by that rule.
  If two words have the same length, sort the list alphabetically.
  The function should return a list of strings in sorted order.
  You may assume that all words will have the same length.
  For example:
  assert list_sort(["aa", "a", "aaa"]) => ["aa"]
  assert list_sort(["ab", "a", "aaa", "cd"]) => ["ab", "cd"]
  */
const sortedListSum = (lst) => {
  let p = []
  for (let i = 0; i < lst.length; i++) {
    if (lst[i].length % 2 == 0) {
      p.push(lst[i])
    }
  }
  for (let j = p.length - 2; j >= 0; j--) {
    for (let k = 0; k <= j; k++) {
      let f = 0
      if (p[k].length > p[k + 1].length) { f = 1 }
      if (p[k].length == p[k + 1].length) {
        let r = p[k].length
        for (let l = 0; l < r; l++) {
          if (p[k][l].charCodeAt() > p[k + 1][l].charCodeAt()) {
            f = 1;
            break;
          }
          if (p[k][l].charCodeAt() < p[k + 1][l].charCodeAt()) {
            break;
          }
        }
      }
      if (f == 1) {
        let tmp = p[k]
        p[k] = p[k + 1]
        p[k + 1] = tmp
      }
    }
  }
  return p
}
