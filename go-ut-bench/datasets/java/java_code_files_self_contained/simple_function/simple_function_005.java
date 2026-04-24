import java.util.*;
import java.lang.*;

class Solution {
    /**
    Input is a space-delimited string of numberals from 'zero' to 'nine'.
    Valid choices are 'zero', 'one', 'two', 'three', 'four', 'five', 'six', 'seven', 'eight' and 'nine'.
    Return the string with numbers sorted from smallest to largest
    >>> sortNumbers("three one five")
    "one three five"
     */
    public String sortNumbers(String numbers) {
        String[] nums = numbers.split(" ");
        List<Integer> num = new ArrayList<>();
        for (String string : nums) {
            switch (string) {
                case "zero" -> num.add(0);
                case "one" -> num.add(1);
                case "two" -> num.add(2);
                case "three" -> num.add(3);
                case "four" -> num.add(4);
                case "five" -> num.add(5);
                case "six" -> num.add(6);
                case "seven" -> num.add(7);
                case "eight" -> num.add(8);
                case "nine" -> num.add(9);
            }
        }
        Collections.sort(num);
        List<String> result = new ArrayList<>();
        for (int m : num) {
            switch (m) {
                case 0 -> result.add("zero");
                case 1 -> result.add("one");
                case 2 -> result.add("two");
                case 3 -> result.add("three");
                case 4 -> result.add("four");
                case 5 -> result.add("five");
                case 6 -> result.add("six");
                case 7 -> result.add("seven");
                case 8 -> result.add("eight");
                case 9 -> result.add("nine");
            }
        }
        return String.join(" ", result);
    }
}
