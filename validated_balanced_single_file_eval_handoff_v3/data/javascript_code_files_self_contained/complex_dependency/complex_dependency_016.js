/*
  Given a positive integer n, return a tuple that has the number of even and odd
  integer palindromes that fall within the range(1, n), inclusive.

  Example 1:

      Input: 3
      Output: (1, 2)
      Explanation:
      Integer palindrome are 1, 2, 3. one of them is even, and two of them are odd.

  Example 2:

      Input: 12
      Output: (4, 6)
      Explanation:
      Integer palindrome are 1, 2, 3, 4, 5, 6, 7, 8, 9, 11. four of them are even, and 6 of them are odd.

  Note:
      1. 1 <= n <= 10^3
      2. returned tuple has the number of even and odd integer palindromes respectively.
  */
const evenOddPalindrome = (n) => {
  let e = 0
  let o = 0
  for (let i = 1; i <= n; i++) {
    let k = i.toString()
    let p = 1
    for (let j = 0; j < k.length; j++) {
      if (k[j] != k[k.length - j - 1]) {
        p = 0;
        break;
      }
    }
    if (p == 1) {
      if (k % 2 == 0) { e++ }
      else { o++ }
    }
  }
  return (e, o)
}
