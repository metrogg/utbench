/*Circular shift the digits of the integer x, shift the digits right by shift
  and return the result as a string.
  If shift > number of digits, return digits reversed.
  >>> circularShift(12, 1)
  "21"
  >>> circularShift(12, 2)
  "12"
  */
const circularShift = (x, shift) => {
  s = x.toString();
  if (shift > s.length)
    return s.split('').reverse().join('');
  else
    return s.slice(-shift) + s.slice(0, -shift);
}
