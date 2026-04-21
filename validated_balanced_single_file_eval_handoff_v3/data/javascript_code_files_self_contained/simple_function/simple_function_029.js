/* Return list of all prefixes from shortest to longest of the input string
  >>> allPrefixes('abc')
  ['a', 'ab', 'abc']
  */
const allPrefixes = (string) => {
  var result = [];
  for (let i = 0; i < string.length; i++) {
    result.push(string.slice(0, i+1));
  }
  return result;
}
