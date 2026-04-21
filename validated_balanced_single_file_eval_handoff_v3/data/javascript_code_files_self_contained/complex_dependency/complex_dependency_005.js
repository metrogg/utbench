/*
  You are given a string representing a sentence,
  the sentence contains some words separated by a space,
  and you have to return a string that contains the words from the original sentence,
  whose lengths are prime numbers,
  the order of the words in the new string should be the same as the original one.

  Example 1:
      Input: sentence = "This is a test"
      Output: "is"

  Example 2:
      Input: sentence = "lets go for swimming"
      Output: "go for"

  Constraints:
      * 1 <= len(sentence) <= 100
      * sentence contains only letters
  */
const wordsInSentence = (sentence) => {
  let t = sentence.split(/\s/)
  let p = ''
  for (let j = 0; j < t.length; j++) {
    let len = t[j].length;
    let u = 1
    if (len == 1 || len == 0) { continue }
    for (let i = 2; i * i <= len; i++) {
      if (len % i == 0) { u = 0 }
    }
    if (u == 0) { continue }
    if (p == '') { p += t[j] }
    else { p = p + ' ' + t[j] }
  }
  return p
}
