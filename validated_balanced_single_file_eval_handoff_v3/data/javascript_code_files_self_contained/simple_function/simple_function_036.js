/* Filter an input list of strings only for ones that contain given substring
  >>> filterBySubstring([], 'a')
  []
  >>> filterBySubstring(['abc', 'bacd', 'cde', 'array'], 'a')
  ['abc', 'bacd', 'array']
  */
const filterBySubstring = (strings, substring) => {
  return strings.filter(x => x.indexOf(substring) != -1);
}
