/*
  Given a string text, replace all spaces in it with underscores, 
  and if a string has more than 2 consecutive spaces, 
  then replace all consecutive spaces with - 
  
  fixSpaces("Example") == "Example"
  fixSpaces("Example 1") == "Example_1"
  fixSpaces(" Example 2") == "_Example_2"
  fixSpaces(" Example   3") == "_Example-3"
  */
const fixSpaces = (text) => {
  let t = ''
  let c = 0
  for (let i = 0; i < text.length; i++) {
    if (text[i] == ' ') { c++ }
    else if (c > 0) {
      if (c == 1) { t += '_' }
      if (c == 2) { t += '__' }
      if (c > 2) { t += '-' }
      t += text[i]
      c = 0;
    } else {
      t += text[i]
    }
  }
  if (c == 1) { t += '_' }
  if (c == 2) { t += '__' }
  if (c > 2) { t += '-' }
  return t
}
