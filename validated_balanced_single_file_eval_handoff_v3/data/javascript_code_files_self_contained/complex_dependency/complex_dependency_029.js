/* Input to this function is a string represented multiple groups for nested parentheses separated by spaces.
  For each of the group, output the deepest level of nesting of parentheses.
  E.g. (()()) has maximum two levels of nesting while ((())) has three.

  >>> parseNestedParens('(()()) ((())) () ((())()())')
  [2, 3, 1, 3]
  */
const parseNestedParens = (paren_string) => {
  var parseParenGroup = function (s) {
    let depth = 0, max_depth = 0;
    for (const c of s) {
      if (c == '(') {
        depth += 1;
        max_depth = Math.max(max_depth, depth);
      } else {
        depth -= 1;
      }
    }
    return max_depth;
  }
  return paren_string.split(' ')
          .filter(x => x != '')
          .map(x => parseParenGroup(x));
}
