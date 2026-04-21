/*The Fib4 number sequence is a sequence similar to the Fibbonacci sequnece that's defined as follows:
  fib4(0) -> 0
  fib4(1) -> 0
  fib4(2) -> 2
  fib4(3) -> 0
  fib4(n) -> fib4(n-1) + fib4(n-2) + fib4(n-3) + fib4(n-4).
  Please write a function to efficiently compute the n-th element of the fib4 number sequence.  Do not use recursion.
  >>> fib4(5)
  4
  >>> fib4(6)
  8
  >>> fib4(7)
  14
  */
const fib4 = (n) => {
  var results = [0, 0, 2, 0];
  if (n < 4)
    return results[n];
  for (let i = 4; i < n + 1; i++) {
    results.push(results.at(-1) + results.at(-2) +
                results.at(-3) + results.at(-4));
    results.shift();
  }
  return results.pop();
}
