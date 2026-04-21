/*Return the number of times the digit 7 appears in integers less than n which are divisible by 11 or 13.
  >>> fizzBuzz(50)
  0
  >>> fizzBuzz(78)
  2
  >>> fizzBuzz(79)
  3
  */
const fizzBuzz = (n) => {
  var ns = [], ans = 0;
  for (let i = 0; i < n; i++)
    if (i % 11 == 0 || i % 13 == 0)
      ns.push(i);
  var s = ns.map(x => x.toString()).join('');
  for (const c of s)
    ans += (c == '7');
  return ans;
}
