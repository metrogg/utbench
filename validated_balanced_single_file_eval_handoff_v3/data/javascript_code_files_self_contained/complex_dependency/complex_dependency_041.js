/*Given a list of strings, where each string consists of only digits, return a list.
  Each element i of the output should be "the number of odd elements in the
  string i of the input." where all the i's should be replaced by the number
  of odd digits in the i'th string of the input.

  >>> oddCount(['1234567'])
  ["the number of odd elements 4n the str4ng 4 of the 4nput."]
  >>> oddCount(['3',"11111111"])
  ["the number of odd elements 1n the str1ng 1 of the 1nput.",
   "the number of odd elements 8n the str8ng 8 of the 8nput."]
  */
const oddCount = (lst) => {
  let d = []
  for (let i = 0; i < lst.length; i++) {
    let p = 0;
    let h = lst[i].length
    for (let j = 0; j < h; j++) {
      if (lst[i][j].charCodeAt() % 2 == 1) { p++ }
    }
    p = p.toString()
    d.push('the number of odd elements ' + p + 'n the str' + p + 'ng ' + p + ' of the ' + p + 'nput.')
  }
  return d
}
