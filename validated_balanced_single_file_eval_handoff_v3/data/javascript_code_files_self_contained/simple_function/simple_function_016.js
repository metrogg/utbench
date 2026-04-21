/*Return true if all numbers in the list l are below threshold t.
  >>> belowThreshold([1, 2, 4, 10], 100)
  true
  >>> belowThreshold([1, 20, 4, 10], 5)
  false
  */
const belowThreshold = (l, t) => {
  for (const e of l)
    if (e >= t)
      return false;
  return true;
}
