/*
  Given a positive integer n, return a sorted list that has the odd numbers in collatz sequence.

  The Collatz conjecture is a conjecture in mathematics that concerns a sequence defined
  as follows: start with any positive integer n. Then each term is obtained from the 
  previous term as follows: if the previous term is even, the next term is one half of 
  the previous term. If the previous term is odd, the next term is 3 times the previous
  term plus 1. The conjecture is that no matter what value of n, the sequence will always reach 1.

  Note: 
      1. Collatz(1) is [1].
      2. returned list sorted in increasing order.

  For example:
  getOddCollatz(5) returns [1, 5] # The collatz sequence for 5 is [5, 16, 8, 4, 2, 1], so the odd numbers are only 1, and 5.
  */
const getOddCollatz = (n) => {
  let p = []
  let t = n
  while (1) {
    let u = 0
    for (let i = 0; i < p.length; i++) {
      if (t == p[i]) {
        u = 1
        break;
      }
    }
    if (u == 1) { break }
    if (t % 2 == 1) { p.push(t); t = 3 * t + 1 }
    else { t = t / 2 }
  }
  for (let j = 0; j < p.length; j++) {
    let ind = j
    for (let k = j + 1; k < p.length; k++) {
      if (p[k] < p[ind]) {
        ind = k
      }
    }
    if (ind > j) {
      let tmp = p[j]
      p[j] = p[ind]
      p[ind] = tmp
    }
  }
  return p
}
