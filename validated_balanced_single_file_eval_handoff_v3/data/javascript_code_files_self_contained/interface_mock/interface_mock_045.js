/*You are given a string s.
  if s[i] is a letter, reverse its case from lower to upper or vise versa, 
  otherwise keep it as it is.
  If the string contains no letters, reverse the string.
  The function should return the resulted string.
  Examples
  solve("1234") = "4321"
  solve("ab") = "AB"
  solve("#a@C") = "#A@c"
  */
const solve = (s) => {
  let t = 0
  let p = ''
  for (let i = 0; i < s.length; i++) {
    let y = s[i].charCodeAt()
    if (y >= 65 && y <= 90) {
      y += 32;
      t = 1;
    } else if (y >= 97 && y <= 122) {
      y -= 32;
      t = 1;
    }
    p += String.fromCharCode(y)
  }
  if (t == 1) { return p }
  let u = ''
  for (let i = 0; i < p.length; i++) {
    u += p[p.length - i - 1]
  }
  return u
}
