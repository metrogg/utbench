// Converted Java method
import java.util.HashMap;
import java.util.Map;

class StringProcessor {
    
    /**
     * Processes a string to perform multiple transformations and analyses.
     * 
     * @param input The string to be processed
     * @return A map containing various processed versions of the string and statistics
     * @throws IllegalArgumentException if input is null
     */
    public static Map<String, Object> processString(String input) {
        if (input == null) {
            throw new IllegalArgumentException("Input string cannot be null");
        }
        
        Map<String, Object> result = new HashMap<>();
        
        // Remove all whitespace
        String noSpaces = input.replaceAll("\\s+", "");
        result.put("noSpaces", noSpaces);
        
        // Remove punctuation
        String noPunctuation = input.replaceAll("[^a-zA-Z0-9]", "");
        result.put("noPunctuation", noPunctuation);
        
        // Convert to lowercase
        String lowercase = input.toLowerCase();
        result.put("lowercase", lowercase);
        
        // Convert to uppercase
        String uppercase = input.toUpperCase();
        result.put("uppercase", uppercase);
        
        // Count vowels
        int vowelCount = input.replaceAll("[^aeiouAEIOU]", "").length();
        result.put("vowelCount", vowelCount);
        
        // Count consonants
        int consonantCount = input.replaceAll("[^a-zA-Z]", "").replaceAll("[aeiouAEIOU]", "").length();
        result.put("consonantCount", consonantCount);
        
        // Count digits
        int digitCount = input.replaceAll("[^0-9]", "").length();
        result.put("digitCount", digitCount);
        
        // Count words
        int wordCount = input.trim().isEmpty() ? 0 : input.trim().split("\\s+").length;
        result.put("wordCount", wordCount);
        
        return result;
    }
}
