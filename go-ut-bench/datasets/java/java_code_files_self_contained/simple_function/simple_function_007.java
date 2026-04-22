import java.util.*;
import java.lang.*;

class Solution {
    /**
    Write a function that accepts a list of strings.
    The list contains different words. Return the word with maximum number
    of unique characters. If multiple strings have maximum number of unique
    characters, return the one which comes first in lexicographical order.

    findMax(["name", "of", "string"]) == "string"
    findMax(["name", "enam", "game"]) == "enam"
    findMax(["aaaaaaa", "bb" ,"cc"]) == ""aaaaaaa"
     */
    public String findMax(List<String> words) {
        List<String> words_sort = new ArrayList<>(words);
        words_sort.sort(new Comparator<String>() {
            @Override
            public int compare(String o1, String o2) {
                Set<Character> s1 = new HashSet<>();
                for (char ch : o1.toCharArray()) {
                    s1.add(ch);
                }
                Set<Character> s2 = new HashSet<>();
                for (char ch : o2.toCharArray()) {
                    s2.add(ch);
                }
                if (s1.size() > s2.size()) {
                    return 1;
                } else if (s1.size() < s2.size()) {
                    return -1;
                } else {
                    return -o1.compareTo(o2);
                }
            }
        });
        return words_sort.get(words_sort.size() - 1);
    }
}
