/* Implement the function f that takes n as a parameter,
  and returns a list of size n, such that the value of the element at index i is the factorial of i if i is even
  or the sum of numbers from 1 to i otherwise.
  i starts from 1.
  the factorial of i is the multiplication of the numbers from 1 to i (1 * 2 * ... * i).
  Example:
  f(5) == [1, 2, 6, 24, 15]
  */
const f = (n) => {
  let f = 1
  let p = 0
  let k = []
  for (let i = 1; i <= n; i++) {
    p += i;
    f *= i;
    if (i % 2 == 0) { k.push(f) }
    else { k.push(p) }
  }
  return k
}
