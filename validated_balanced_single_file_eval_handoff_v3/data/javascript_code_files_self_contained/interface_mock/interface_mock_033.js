/*Create a function which takes a string representing a file's name, and returns
  'Yes' if the the file's name is valid, and returns 'No' otherwise.
  A file's name is considered to be valid if and only if all the following conditions 
  are met:
  - There should not be more than three digits ('0'-'9') in the file's name.
  - The file's name contains exactly one dot '.'
  - The substring before the dot should not be empty, and it starts with a letter from 
  the latin alphapet ('a'-'z' and 'A'-'Z').
  - The substring after the dot should be one of these: ['txt', 'exe', 'dll']
  Examples:
  fileNameCheck("example.txt") # => 'Yes'
  fileNameCheck("1example.dll") # => 'No' (the name should start with a latin alphapet letter)
  */
const fileNameCheck = (file_name) => {
  let t = file_name.split(/\./)
  if (t.length != 2) { return 'No' }
  if (t[1] != 'txt' && t[1] != 'dll' && t[1] != 'exe') { return 'No' }
  if (t[0] == '') { return 'No' }
  let a = t[0][0].charCodeAt()
  if (!((a >= 65 && a <= 90) || (a >= 97 && a <= 122))) { return 'No' }
  let y = 0
  for (let i = 1; i < t[0].length; i++) {
    if (t[0][i].charCodeAt() >= 48 && t[0][i].charCodeAt() <= 57) { y++ }
    if (y > 3) { return 'No' }
  }
  return 'Yes'
}
