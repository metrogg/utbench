/*
  Create a function that takes integers, floats, or strings representing
  real numbers, and returns the larger variable in its given variable type.
  Return null if the values are equal.
  Note: If a real number is represented as a string, the floating point might be . or ,

  compareOne(1, 2.5) ➞ 2.5
  compareOne(1, "2,3") ➞ "2,3"
  compareOne("5,1", "6") ➞ "6"
  compareOne("1", 1) ➞ null
  */
const compareOne = (a, b) => {
  let aa = Number(a)
  if (typeof a == 'string') { aa = Number(a.replace(',', '.')) }
  let bb = Number(b)
  if (typeof b == 'string') { bb = Number(b.replace(',', '.')) }
  if (aa > bb) { return a }
  if (aa < bb) { return b }
  return null
}
