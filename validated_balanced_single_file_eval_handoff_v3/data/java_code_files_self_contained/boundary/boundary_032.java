// Converted Java method
import java.util.regex.Pattern;

class EnhancedStringUtil {

    /**
     * Checks if a string is empty (null or only whitespace)
     * @param str The string to check
     * @return true if the string is null or empty/whitespace, false otherwise
     */
    public static boolean isEmpty(String str) {
        return str == null || str.trim().isEmpty();
    }

    /**
     * Checks if a string is not empty
     * @param str The string to check
     * @return true if the string is not null and contains non-whitespace characters
     */
    public static boolean isNotEmpty(String str) {
        return !isEmpty(str);
    }

    /**
     * Checks if a string is blank (null, empty, or whitespace)
     * @param str The string to check
     * @return true if the string is null, empty, or only whitespace
     */
    public static boolean isBlank(String str) {
        return str == null || str.trim().length() == 0;
    }

    /**
     * Checks if a string contains only alphabetic characters
     * @param str The string to check
     * @return true if the string contains only letters, false otherwise
     */
    public static boolean isAlpha(String str) {
        if (isEmpty(str)) return false;
        return str.matches("[a-zA-Z]+");
    }

    /**
     * Checks if a string contains only alphanumeric characters
     * @param str The string to check
     * @return true if the string contains only letters and digits, false otherwise
     */
    public static boolean isAlphanumeric(String str) {
        if (isEmpty(str)) return false;
        return str.matches("[a-zA-Z0-9]+");
    }

    /**
     * Checks if a string is a valid email address
     * @param str The string to check
     * @return true if the string is a valid email format, false otherwise
     */
    public static boolean isEmail(String str) {
        if (isEmpty(str)) return false;
        String emailRegex = "^[a-zA-Z0-9_+&*-]+(?:\\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,7}$";
        return Pattern.compile(emailRegex).matcher(str).matches();
    }

    /**
     * Reverses a string
     * @param str The string to reverse
     * @return The reversed string, or null if input is null
     */
    public static String reverse(String str) {
        if (str == null) return null;
        return new StringBuilder(str).reverse().toString();
    }

    /**
     * Counts the number of words in a string
     * @param str The string to analyze
     * @return The word count, or 0 if the string is empty
     */
    public static int countWords(String str) {
        if (isEmpty(str)) return 0;
        return str.trim().split("\\s+").length;
    }
}
