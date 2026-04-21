/* Out of list of strings, return the longest one. Return the first one in case of multiple
  strings of the same length. Return null in case the input list is empty.
  >>> longest([])

  >>> longest(['a', 'b', 'c'])
  'a'
  >>> longest(['a', 'bb', 'ccc'])
  'ccc'
  */
const longest = (strings) => {
  if (!Array.isArray(strings) || strings.length == 0)
    return null;
  var maxlen = Math.max(...strings.map(x => x.length));
  for (const s of strings) {
    if (s.length == maxlen) {
      return s;
    }
  }
}
