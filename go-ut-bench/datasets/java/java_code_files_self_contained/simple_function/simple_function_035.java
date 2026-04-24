import java.util.*;
import java.lang.*;

class Solution {
    /**
    Create a function that takes a string as input which contains only square brackets.
    The function should return true if and only if there is a valid subsequence of brackets
    where at least one bracket in the subsequence is nested.

    isNested("[[]]") -> true
    isNested("[]]]]]]][[[[[]") -> false
    isNested("[][]") -> false
    isNested("[]") -> false
    isNested("[[][]]") -> true
    isNested("[[]][[") -> true
     */
    public boolean isNested(String string) {
        List<Integer> opening_bracket_index = new ArrayList<>(), closing_bracket_index = new ArrayList<>();
        for (int i = 0; i < string.length(); i++) {
            if (string.charAt(i) == '[') {
                opening_bracket_index.add(i);
            } else {
                closing_bracket_index.add(i);
            }
        }
        Collections.reverse(closing_bracket_index);
        int i = 0, l = closing_bracket_index.size();
        for (int idx : opening_bracket_index) {
            if (i < l && idx < closing_bracket_index.get(i)) {
                i += 1;
            }
        }
        return i >= 2;
    }
}
