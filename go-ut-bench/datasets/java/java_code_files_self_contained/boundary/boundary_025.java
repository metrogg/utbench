import java.util.HashMap;
import java.util.Map;

class AdvancedPalindromeFinder {

    /**
     * Finds the longest palindromic substring and provides additional palindrome statistics.
     * 
     * @param s Input string to analyze
     * @return Map containing:
     *         - "longest": longest palindromic substring
     *         - "length": length of the longest palindrome
     *         - "count": total number of palindromic substrings found
     *         - "positions": map of all palindromes with their start/end positions
     * @throws IllegalArgumentException if input string is null
     */
    public Map<String, Object> analyzePalindromes(String s) {
        if (s == null) {
            throw new IllegalArgumentException("Input string cannot be null");
        }

        Map<String, Object> result = new HashMap<>();
        String longest = "";
        int count = 0;
        Map<String, int[]> positions = new HashMap<>();

        if (s.isEmpty()) {
            result.put("longest", "");
            result.put("length", 0);
            result.put("count", 0);
            result.put("positions", positions);
            return result;
        }

        // Find all palindromes (odd and even length)
        for (int i = 0; i < s.length(); i++) {
            // Odd length palindromes
            String oddPalindrome = expandAroundCenter(s, i, i);
            if (oddPalindrome.length() > longest.length()) {
                longest = oddPalindrome;
            }
            if (oddPalindrome.length() > 1) {
                positions.put(oddPalindrome, new int[]{i - oddPalindrome.length()/2, i + oddPalindrome.length()/2});
                count++;
            }

            // Even length palindromes
            String evenPalindrome = expandAroundCenter(s, i, i + 1);
            if (evenPalindrome.length() > longest.length()) {
                longest = evenPalindrome;
            }
            if (evenPalindrome.length() > 1) {
                positions.put(evenPalindrome, new int[]{i - evenPalindrome.length()/2 + 1, i + evenPalindrome.length()/2});
                count++;
            }
        }

        // Include single characters as palindromes
        count += s.length();

        result.put("longest", longest);
        result.put("length", longest.length());
        result.put("count", count);
        result.put("positions", positions);
        return result;
    }

    private String expandAroundCenter(String s, int left, int right) {
        while (left >= 0 && right < s.length() && s.charAt(left) == s.charAt(right)) {
            left--;
            right++;
        }
        return s.substring(left + 1, right);
    }
}
