import java.util.*;
import java.lang.*;

class Solution {
    /**
    Write a function countNums which takes an array of integers and returns
    the number of elements which has a sum of digits > 0.
    If a number is negative, then its first signed digit will be negative:
    e.g. -123 has signed digits -1, 2, and 3.
    >>> countNums(Arrays.asList()) == 0
    >>> countNums(Arrays.asList(-1, 11, -11)) == 1
    >>> countNums(Arrays.asList(1, 1, 2)) == 3
     */
    public int countNums(List<Integer> arr) {
        int count = 0;
        for (int n: arr) {
            int neg = 1;
            if (n < 0) {
                n = -n;
                neg = -1;
            }
            List<Integer> digits = new ArrayList<>();
            for (char digit : String.valueOf(n).toCharArray()) {
                digits.add(digit - '0');
            }
            digits.set(0, digits.get(0) * neg);
            if (digits.stream().reduce(0, Integer::sum) > 0) {
                count += 1;
            }
        }
        return count;
    }
}
