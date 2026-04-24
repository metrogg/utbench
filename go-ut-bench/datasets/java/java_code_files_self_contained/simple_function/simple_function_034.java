import java.util.*;
import java.lang.*;

class Solution {
    /**
    Create a function that takes a value (string) representing a number
    and returns the closest integer to it. If the number is equidistant
    from two integers, round it away from zero.

    Examples
    >>> closest_integer("10")
    10
    >>> closest_integer("15.3")
    15

    Note:
    Rounding away from zero means that if the given number is equidistant
    from two integers, the one you should return is the one that is the
    farthest from zero. For example closest_integer("14.5") should
    return 15 and closest_integer("-14.5") should return -15.
     */
    public int countUpper(String value) {
        if (value.contains(".")) {
            while (value.charAt(value.length() - 1) == '0') {
                value = value.substring(0, value.length() - 1);
            }
        }
        double num = Double.parseDouble(value);
        int res = 0;
        if (value.substring(Math.max(value.length() - 2, 0)).equals(".5")) {
            if (num > 0) {
                res = (int) Math.ceil(num);
            } else {
                res = (int) Math.floor(num);
            }
        } else if(value.length() > 0) {
            res = (int) Math.round(num);
        }
        return res;
    }
}
