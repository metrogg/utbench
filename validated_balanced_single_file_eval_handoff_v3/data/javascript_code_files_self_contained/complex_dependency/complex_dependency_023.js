/*You are given two intervals,
  where each interval is a pair of integers. For example, interval = (start, end) = (1, 2).
  The given intervals are closed which means that the interval (start, end)
  includes both start and end.
  For each given interval, it is assumed that its start is less or equal its end.
  Your task is to determine whether the length of intersection of these two 
  intervals is a prime number.
  Example, the intersection of the intervals (1, 3), (2, 4) is (2, 3)
  which its length is 1, which not a prime number.
  If the length of the intersection is a prime number, return "YES",
  otherwise, return "NO".
  If the two intervals don't intersect, return "NO".


  [input/output] samples:
  intersection((1, 2), (2, 3)) ==> "NO"
  intersection((-1, 1), (0, 4)) ==> "NO"
  intersection((-3, -1), (-5, 5)) ==> "YES"
  */
const intersection = (interval1, interval2) => {
  let lo = interval1[0]
  if (interval2[0] > lo) { lo = interval2[0] }
  let hi = interval1[1]
  if (interval2[1] < hi) { hi = interval2[1] }
  let len = 0
  if (hi > lo) { len = hi - lo }
  if (len == 1 || len == 0) { return 'NO' }
  for (let i = 2; i * i <= len; i++) {
    if (len % i == 0) { return 'NO' }
  }
  return 'YES'
}
