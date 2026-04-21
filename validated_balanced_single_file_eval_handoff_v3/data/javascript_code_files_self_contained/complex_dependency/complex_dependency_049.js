/*
  removeVowels is a function that takes string and returns string without vowels.
  >>> removeVowels('')
  ''
  >>> removeVowels("abcdef\nghijklm")
  'bcdf\nghjklm'
  >>> removeVowels('abcdef')
  'bcdf'
  >>> removeVowels('aaaaa')
  ''
  >>> removeVowels('aaBAA')
  'B'
  >>> removeVowels('zbcd')
  'zbcd'
  */
const removeVowels = (text) => {
  return text.split("")
             .filter(s => !["a", "e", "i", "o", "u"]
                      .includes(s.toLowerCase())
                    )
             .join("")
}
