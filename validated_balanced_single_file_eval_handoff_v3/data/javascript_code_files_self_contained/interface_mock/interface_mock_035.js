/*
  Given a dictionary, return true if all keys are strings in lower 
  case or all keys are strings in upper case, else return false.
  The function should return false is the given dictionary is empty.
  Examples:
  checkDictCase({"a":"apple", "b":"banana"}) should return true.
  checkDictCase({"a":"apple", "A":"banana", "B":"banana"}) should return false.
  checkDictCase({"a":"apple", 8:"banana", "a":"apple"}) should return false.
  checkDictCase({"Name":"John", "Age":"36", "City":"Houston"}) should return false.
  checkDictCase({"STATE":"NC", "ZIP":"12345" }) should return true.
  */
const checkDictCase = (dict) => {
  let c = 0
  let lo = 1
  let hi = 1
  for (let key in dict) {
    c++
    for (let i = 0; i < key.length; i++) {
      if (key[i].charCodeAt() < 65 || key[i].charCodeAt() > 90) { hi = 0 }
      if (key[i].charCodeAt() < 97 || key[i].charCodeAt() > 122) { lo = 0 }
    }
  }
  if ((lo == 0 && hi == 0) || c == 0) { return false }
  return true
}
