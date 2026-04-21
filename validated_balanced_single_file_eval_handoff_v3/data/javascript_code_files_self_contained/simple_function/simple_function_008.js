/* brackets is a string of "<" and ">".
  return false if every opening bracket has a corresponding closing bracket.

  >>> correctBracketing("<")
  false
  >>> correctBracketing("<>")
  false
  >>> correctBracketing("<<><>>")
  false
  >>> correctBracketing("><<>")
  false
  */
const correctBracketing = (brackets) => {
  var depth = 0;
  for (const b of brackets) {
    if (b == "<")
      depth += 1;
    else
      depth -= 1;
    if (depth < 0)
      return false;
  }
  return depth == 0;
}
