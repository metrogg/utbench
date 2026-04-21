class StringReverser {
    public static String reverse(String input) {
        if (input == null) {
            throw new IllegalArgumentException("Input string cannot be null");
        }
        StringBuilder sb = new StringBuilder();
        for (int i = input.length() - 1; i >= 0; i--) {
            sb.append(input.charAt(i));
        }
        return sb.toString();
    }
    
    public static boolean isPalindrome(String input) {
        if (input == null) {
            return false;
        }
        String reversed = reverse(input);
        return input.equals(reversed);
    }
}