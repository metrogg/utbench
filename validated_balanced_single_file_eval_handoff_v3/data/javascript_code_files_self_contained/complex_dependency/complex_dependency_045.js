/*Write a function vowelsCount which takes a string representing
  a word as input and returns the number of vowels in the string.
  Vowels in this case are 'a', 'e', 'i', 'o', 'u'. Here, 'y' is also a
  vowel, but only when it is at the end of the given word.

  Example:
  >>> vowelsCount("abcde")
  2
  >>> vowelsCount("ACEDY")
  3
  */
const vowelsCount = (s) => {
  var vowels = "aeiouAEIOU";
  var n_vowels = s.split('').reduce((prev, item) => {
    return prev + (vowels.includes(item));
  }, 0);
  if (s.at(-1) == 'y' || s.at(-1) == 'Y')
    n_vowels += 1;
  return n_vowels;
}
