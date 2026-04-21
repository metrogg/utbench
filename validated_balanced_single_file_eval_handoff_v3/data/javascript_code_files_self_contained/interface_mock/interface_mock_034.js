/*Given a string s and a natural number n, you have been tasked to implement 
  a function that returns a list of all words from string s that contain exactly 
  n consonants, in order these words appear in the string s.
  If the string s is empty then the function should return an empty list.
  Note: you may assume the input string contains only letters and spaces.
  Examples:
  selectWords("Mary had a little lamb", 4) ==> ["little"]
  selectWords("Mary had a little lamb", 3) ==> ["Mary")
  selectWords("simple white space", 2) ==> []
  selectWords("Hello world", 4) ==> ["world"]
  selectWords("Uncle sam", 3) ==> ["Uncle"]
  */
const selectWords = (s, n) => {
  let t = s.split(/\s/)
  if (s == '') { return [] }
  let k = []
  for (let i = 0; i < t.length; i++) {
    let l = 0
    for (let j = 0; j < t[i].length; j++) {
      if (t[i][j] != 'a' && t[i][j] != 'e' && t[i][j] != 'i' && t[i][j] != 'o' && t[i][j] != 'u' && t[i][j] != 'A' &&
        t[i][j] != 'U' && t[i][j] != 'O' && t[i][j] != 'I' && t[i][j] != 'E') {
        l++
      }
    }
    if (l == n) { k.push(t[i]) }
  }
  return k
}
