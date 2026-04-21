/*
  primeFib returns n-th number that is a Fibonacci number and it's also prime.
  >>> primeFib(1)
  2
  >>> primeFib(2)
  3
  >>> primeFib(3)
  5
  >>> primeFib(4)
  13
  >>> primeFib(5)
  89
  */
const primeFib = (n) => {
  var isPrime = function (p) {
    if (p < 2)
      return false;
    for (let k = 2; k < Math.min(Math.floor(Math.sqrt(p)) + 1, p - 1); k++) {
      if (p % k == 0)
        return false;
    }
    return true;
  }

  var f = [0, 1];
  while (true) {
    f.push(f.at(-1) + f.at(-2));
    if (isPrime(f.at(-1)))
      n -= 1;
    if (n == 0)
      return f.at(-1);
  }
}
