/*Write a function that takes a string and returns true if the string
  length is a prime number or false otherwise
  Examples
  primeLength('Hello') == true
  primeLength('abcdcba') == true
  primeLength('kittens') == true
  primeLength('orange') == false
  */
const primeLength = (string) => {
  let len = string.length
  if (len == 1 || len == 0) { return false }
  for (let i = 2; i * i <= len; i++) {
    if (len % i == 0) { return false }
  }
  return true
}
