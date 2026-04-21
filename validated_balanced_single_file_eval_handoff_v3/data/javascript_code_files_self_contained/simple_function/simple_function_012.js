/*A simple program which should return the value of x if n is 
  a prime number and should return the value of y otherwise.

  Examples:
  for xOrY(7, 34, 12) == 34
  for xOrY(15, 8, 5) == 5
  
  */
const xOrY = (n, x, y) => {
  let len = n
  if (len == 1 || len == 0) { return y }
  for (let i = 2; i * i <= len; i++) {
    if (len % i == 0) { return y }
  }
  return x
}
