/* Input are two strings a and b consisting only of 1s and 0s.
  Perform binary XOR on these inputs and return result also as a string.
  >>> stringXor('010', '110')
  '100'
  */
const stringXor = (a, b) => {
  var xor = function (i, j) {
    if (i == j)
      return '0';
    else
      return '1';
  }
  return a.split('').map((item, index) => xor(item, b[index])).join('');
}
