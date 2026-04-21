/*Return sorted unique common elements for two lists.
  >>> common([1, 4, 3, 34, 653, 2, 5], [5, 7, 1, 5, 9, 653, 121])
  [1, 5, 653]
  >>> common([5, 3, 2, 8], [3, 2])
  [2, 3]

  */
const common = (l1, l2) => {
  var ret = new Set();
  for (const e1 of l1)
    for (const e2 of l2)
      if (e1 == e2)
        ret.add(e1);
  return [...ret].sort();
}
