import java.util.*;
import java.lang.*;

class Solution {
    /**
    Write a function that accepts a list of strings as a parameter,
    deletes the strings that have odd lengths from it,
    and returns the resulted list with a sorted order,
    The list is always a list of strings and never an array of numbers,
    and it may contain duplicates.
    The order of the list should be ascending by length of each word, and you
    should return the list sorted by that rule.
    If two words have the same length, sort the list alphabetically.
    The function should return a list of strings in sorted order.
    You may assume that all words will have the same length.
    For example:
    assert listSort(Arrays.asList("aa", "a", "aaa")) => ["aa"]
    assert listSort(Arrays.asList("ab", "a", "aaa", "cd")) => ["ab", "cd"]
     */
    public List<String> listSort(List<String> lst) {
        List<String> lst_sorted = new ArrayList<>(lst);
        Collections.sort(lst_sorted);
        List<String> new_lst = new ArrayList<>();
        for (String i : lst_sorted) {
            if (i.length() % 2 == 0) {
                new_lst.add(i);
            }
        }
        new_lst.sort(Comparator.comparingInt(String::length));
        return new_lst;
    }
}
