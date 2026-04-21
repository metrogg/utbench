/*
  Write a function that accepts two lists of strings and returns the list that has
  total number of chars in the all strings of the list less than the other list.

  if the two lists have the same number of chars, return the first list.

  Examples
  totalMatch([], []) ➞ []
  totalMatch(['hi', 'admin'], ['hI', 'Hi']) ➞ ['hI', 'Hi']
  totalMatch(['hi', 'admin'], ['hi', 'hi', 'admin', 'project']) ➞ ['hi', 'admin']
  totalMatch(['hi', 'admin'], ['hI', 'hi', 'hi']) ➞ ['hI', 'hi', 'hi']
  totalMatch(['4'], ['1', '2', '3', '4', '5']) ➞ ['4']
  */
const totalMatch = (lst1, lst2) => {
  var l1 = lst1.reduce(((prev, item) => prev + item.length), 0);
  var l2 = lst2.reduce(((prev, item) => prev + item.length), 0);
  if (l1 <= l2)
    return lst1;
  else
    return lst2;
}
