// Converted Java method
import java.util.HashMap;
import java.util.Map;
import java.util.regex.Pattern;

class AuthenticationValidator {
    private static final Pattern USERNAME_PATTERN = Pattern.compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,6}$");
    private static final Pattern PASSWORD_PATTERN = Pattern.compile("^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[@#$%^&+=])(?=\\S+$).{8,}$");
    private static final Map<String, String> VALID_USERS = new HashMap<>();

    static {
        // Simulating a database of valid users
        VALID_USERS.put("user@example.com", "P@ssw0rd123");
        VALID_USERS.put("admin@test.com", "Adm!nP@ss456");
        VALID_USERS.put("test.user@domain.com", "TestP@ss789");
    }

    /**
     * Validates user credentials with multiple checks
     * @param username The username/email to validate
     * @param password The password to validate
     * @return Map containing validation results with keys:
     *         "authenticated" (boolean), 
     *         "usernameValid" (boolean),
     *         "passwordValid" (boolean),
     *         "message" (String)
     */
    public Map<String, Object> validateCredentials(String username, String password) {
        Map<String, Object> result = new HashMap<>();
        
        // Validate username format
        boolean isUsernameValid = USERNAME_PATTERN.matcher(username).matches();
        result.put("usernameValid", isUsernameValid);
        
        // Validate password strength
        boolean isPasswordValid = PASSWORD_PATTERN.matcher(password).matches();
        result.put("passwordValid", isPasswordValid);
        
        // Check authentication
        boolean isAuthenticated = false;
        if (isUsernameValid && isPasswordValid) {
            String storedPassword = VALID_USERS.get(username);
            isAuthenticated = storedPassword != null && storedPassword.equals(password);
        }
        result.put("authenticated", isAuthenticated);
        
        // Set appropriate message
        if (!isUsernameValid) {
            result.put("message", "Invalid username format");
        } else if (!isPasswordValid) {
            result.put("message", "Password doesn't meet complexity requirements");
        } else if (!isAuthenticated) {
            result.put("message", "Invalid credentials");
        } else {
            result.put("message", "Authentication successful");
        }
        
        return result;
    }
}
