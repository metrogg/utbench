/*
  Checks if given string is a palindrome
  >>> isPalindrome('')
  true
  >>> isPalindrome('aba')
  true
  >>> isPalindrome('aaaaa')
  true
  >>> isPalindrome('zbcd')
  false
  */
const isPalindrome = (text) => {
  for (let i = 0; i < text.length; i++)
    if (text[i] != text.at(-i-1))
      return false;
  return true;
}
