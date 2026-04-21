/*Write a function that accepts a list of strings.
  The list contains different words. Return the word with maximum number
  of unique characters. If multiple strings have maximum number of unique
  characters, return the one which comes first in lexicographical order.

  findMax(["name", "of", "string"]) === "string"
  findMax(["name", "enam", "game"]) === "enam"
  findMax(["aaaaaaa", "bb" ,"cc"]) === ""aaaaaaa"
  */
const findMax = (words) => {
  let s = -1
  let u = -1
  if (words.length == 0) { return '' }
  for (let i = 0; i < words.length; i++) {
    let p = 0
    for (let j = 0; j < words[i].length; j++) {
      let y = 1
      for (let k = 0; k < j; k++) {
        if (words[i][j] == words[i][k]) { y = 0 }
      }
      if (y == 1) { p++ }
    }
    if (p > s || (p == s && words[i] < words[u])) {
      u = i;
      s = p;
    }
  }
  return words[u]
}
