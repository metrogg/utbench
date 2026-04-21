/* Input to this function is a string containing multiple groups of nested parentheses. Your goal is to
  separate those group into separate strings and return the list of those.
  Separate groups are balanced (each open brace is properly closed) and not nested within each other
  Ignore any spaces in the input string.
  >>> separateParenGroups('( ) (( )) (( )( ))')
  ['()', '(())', '(()())']
  */
const separateParenGroups = (paren_string) => {
  var result = [];
  var current_string = [];
  var current_depth = 0;

  for (const c of paren_string) {
    if (c == '(') {
      current_depth += 1;
      current_string.push(c);
    } else if (c == ')') {
      current_depth -= 1;
      current_string.push(c);
      if (current_depth == 0) {
        result.push(current_string.join(''));
        current_string = [];
      }
    }
  }

  return result;
}
