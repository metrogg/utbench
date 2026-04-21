/*
  Write a function that takes a message, and encodes in such a 
  way that it swaps case of all letters, replaces all vowels in 
  the message with the letter that appears 2 places ahead of that 
  vowel in the english alphabet. 
  Assume only letters. 
  
  Examples:
  >>> encode('test')
  'TGST'
  >>> encode('This is a message')
  'tHKS KS C MGSSCGG'
  */
const encode = (message) => {
  let t = ''
  for (let i = 0; i < message.length; i++) {
    let p = message[i].charCodeAt()
    if (p > 96) { p -= 32 }
    else if (p!=32 && p < 96) { p += 32 }
    if (p == 65 || p == 97 || p == 69 || p == 101 || p == 73 || p == 105 || p == 79 || p == 111 || p == 85 || p == 117) { p += 2 }
    t += String.fromCharCode(p)
  }
  return t
}
