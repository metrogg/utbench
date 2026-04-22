import java.util.*;
import java.lang.*;

class Solution {
    /**
    triplesSumToZero takes a list of integers as an input.
    it returns True if there are three distinct elements in the list that
    sum to zero, and False otherwise.

    >>> triplesSumToZero(Arrays.asList(1, 3, 5, 0))
    false
    >>> triplesSumToZero(Arrays.asList(1, 3, -2, 1))
    true
    >>> triplesSumToZero(Arrays.asList(1, 2, 3, 7))
    false
    >>> triplesSumToZero(Arrays.asList(2, 4, -5, 3, 9, 7))
    true
    >>> triplesSumToZero(Arrays.asList(1))
    false
     */
    public boolean triplesSumToZero(List<Integer> l) {
        for (int i = 0; i < l.size(); i++) {
            for (int j = i + 1; j < l.size(); j++) {
                for (int k = j + 1; k < l.size(); k++) {
                    if (l.get(i) + l.get(j) + l.get(k) == 0) {
                        return true;
                    }
                }
            }
        }
        return false;
    }
}
