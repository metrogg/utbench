import java.util.*;
import java.lang.*;

class Solution {
    /**
    You are given a string representing a sentence,
    the sentence contains some words separated by a space,
    and you have to return a string that contains the words from the original sentence,
    whose lengths are prime numbers,
    the order of the words in the new string should be the same as the original one.

    Example 1:
        Input: sentence = "This is a test"
        Output: "is"

    Example 2:
        Input: sentence = "lets go for swimming"
        Output: "go for"

    Constraints:
        * 1 <= len(sentence) <= 100
        * sentence contains only letters
     */
    public String wordsInSentence(String sentence) {
        List<String> new_lst = new ArrayList<>();
        for (String word : sentence.split(" " )) {
            boolean flg = true;
            if (word.length() == 1) {
                continue;
            }
            for (int i = 2; i < word.length(); i++) {
                if (word.length() % i == 0) {
                    flg = false;
                    break;
                }
            }
            if (flg) {
                new_lst.add(word);
            }
        }
        return String.join(" ", new_lst);
    }
}
