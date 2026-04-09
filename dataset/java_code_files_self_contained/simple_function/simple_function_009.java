import java.util.*;
import java.lang.*;

class Solution {
    /**
    In this Kata, you have to sort an array of non-negative integers according to
    number of ones in their binary representation in ascending order.
    For similar number of ones, sort based on decimal value.
    <p>
    It must be implemented like this:
    >>> sortArray(Arrays.asList(1, 5, 2, 3, 4)).equals(Arrays.asList(1, 2, 3, 4, 5))
    >>> sortArray(Arrays.asList(-2, -3, -4, -5, -6)).equals(Arrays.asList(-6, -5, -4, -3, -2))
    >>> sortArray(Arrays.asList(1, 0, 2, 3, 4)).equals(Arrays.asList(0, 1, 2, 3, 4))
     */
    public List<Integer> sortArray(List<Integer> arr) {
                List < Integer > sorted_arr = new ArrayList<>(arr);
        sorted_arr.sort(new Comparator<Integer>() {
            @Override
            public int compare(Integer o1, Integer o2) {
                int cnt1 = (int) Integer.toBinaryString(Math.abs(o1)).chars().filter(ch -> ch == '1').count();
                int cnt2 = (int) Integer.toBinaryString(Math.abs(o2)).chars().filter(ch -> ch == '1').count();
                if (cnt1 > cnt2) {
                    return 1;
                } else if (cnt1 < cnt2) {
                    return -1;
                } else {
                    return o1.compareTo(o2);
                }
            }
        });
        return sorted_arr;
    }
}
