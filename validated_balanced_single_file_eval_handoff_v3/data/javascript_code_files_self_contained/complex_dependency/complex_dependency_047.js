/* Test if gniven strig is a palindrome */
const isPalindrome = (string) => {
  return string == string.split('').reverse().join('');
}

/* Find the shortest palindrome that begins with a supplied string.
  Algorithm idea is simple:
  - Find the longest postfix of supplied string that is a palindrome.
  - Append to the end of the string reverse of a string prefix that comes before the palindromic suffix.
  >>> makePalindrome('')
  ''
  >>> makePalindrome('cat')
  'catac'
  >>> makePalindrome('cata')
  'catac'
  */
const makePalindrome = (string) => {
  if (string == '')
    return '';
  var beginning_of_suffix = 0;
  while (!isPalindrome(string.slice(beginning_of_suffix)))
    beginning_of_suffix += 1;
  return string + string.slice(0, beginning_of_suffix).split('').reverse().join('');
}
