// Converted Java method
import java.util.regex.Pattern;

class EmployeeValidator {

    /**
     * Validates employee information including email, password, and optional FCM token.
     * Throws IllegalArgumentException with detailed message if validation fails.
     *
     * @param email Employee email to validate
     * @param password Employee password to validate
     * @param fcmToken Optional FCM token to validate (can be null)
     * @throws IllegalArgumentException if any validation fails
     */
    public static void validateEmployee(String email, String password, String fcmToken) {
        validateEmail(email);
        validatePassword(password);
        if (fcmToken != null) {
            validateFcmToken(fcmToken);
        }
    }

    private static void validateEmail(String email) {
        if (email == null || email.isEmpty()) {
            throw new IllegalArgumentException("Email cannot be null or empty");
        }
        
        String emailRegex = "^[a-zA-Z0-9_+&*-]+(?:\\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,7}$";
        Pattern pattern = Pattern.compile(emailRegex);
        if (!pattern.matcher(email).matches()) {
            throw new IllegalArgumentException("Invalid email format");
        }
        
        if (email.length() > 30) {
            throw new IllegalArgumentException("Email must be 30 characters or less");
        }
    }

    private static void validatePassword(String password) {
        if (password == null || password.isEmpty()) {
            throw new IllegalArgumentException("Password cannot be null or empty");
        }
        
        if (password.length() < 8) {
            throw new IllegalArgumentException("Password must be at least 8 characters");
        }
        
        if (password.length() > 255) {
            throw new IllegalArgumentException("Password must be 255 characters or less");
        }
        
        if (!password.matches(".*[A-Z].*")) {
            throw new IllegalArgumentException("Password must contain at least one uppercase letter");
        }
        
        if (!password.matches(".*[a-z].*")) {
            throw new IllegalArgumentException("Password must contain at least one lowercase letter");
        }
        
        if (!password.matches(".*\\d.*")) {
            throw new IllegalArgumentException("Password must contain at least one digit");
        }
        
        if (!password.matches(".*[!@#$%^&*()].*")) {
            throw new IllegalArgumentException("Password must contain at least one special character");
        }
    }

    private static void validateFcmToken(String fcmToken) {
        if (fcmToken.isEmpty()) {
            throw new IllegalArgumentException("FCM token cannot be empty if provided");
        }
        
        if (fcmToken.length() > 255) {
            throw new IllegalArgumentException("FCM token must be 255 characters or less");
        }
        
        // Basic pattern check for FCM token format
        if (!fcmToken.matches("^[a-zA-Z0-9-_:]+$")) {
            throw new IllegalArgumentException("Invalid FCM token format");
        }
    }
}
