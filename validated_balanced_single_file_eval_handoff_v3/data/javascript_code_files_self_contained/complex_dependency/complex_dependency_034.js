/* Given a list of numbers, return whether or not they are sorted
  in ascending order. If list has more than 1 duplicate of the same
  number, return false. Assume no negative numbers and only integers.
  Examples
  isSorted([5]) ➞ true
  isSorted([1, 2, 3, 4, 5]) ➞ true
  isSorted([1, 3, 2, 4, 5]) ➞ false
  isSorted([1, 2, 3, 4, 5, 6]) ➞ true
  isSorted([1, 2, 3, 4, 5, 6, 7]) ➞ true
  isSorted([1, 3, 2, 4, 5, 6, 7]) ➞ false
  isSorted([1, 2, 2, 3, 3, 4]) ➞ true
  isSorted([1, 2, 2, 2, 3, 4]) ➞ false
  */
const isSorted = (lst) => {
  if (lst.length == 0) { return true }
  let dup = 1
  let pre = lst[0]
  for (let i = 1; i < lst.length; i++) {
    if (lst[i] < pre) { return false }
    if (lst[i] == pre) {
      dup += 1;
      if (dup == 3) { return false }
    } else {
      pre = lst[i]
      dup = 1
    }
  }
  return true
}
