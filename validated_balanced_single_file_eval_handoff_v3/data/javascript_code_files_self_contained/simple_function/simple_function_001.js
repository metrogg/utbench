/*Given a positive integer n, return the product of the odd digits.
  Return 0 if all digits are even.
  For example:
  digits(1)  == 1
  digits(4)  == 0
  digits(235) == 15
  */
const digits = (n) => {
  let p = 1
  let k = 1
  while (n > 0) {
    let y = n % 10
    if (y % 2 == 1) {
      p *= y; k = 0;
    }
    n = (n - n % 10) / 10
  }
  if (k == 0) { return p }
  return 0
}
