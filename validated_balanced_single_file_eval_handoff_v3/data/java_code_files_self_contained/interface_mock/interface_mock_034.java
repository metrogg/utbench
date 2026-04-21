// Converted Java method
import java.util.HashMap;
import java.util.Map;

class UserAuthenticator {
    private Map<String, String> userDatabase; // Simulating a user database
    private Map<String, String> verificationCodes; // For storing verification codes

    public UserAuthenticator() {
        // Initialize with some test users
        userDatabase = new HashMap<>();
        userDatabase.put("1234567890", "password123");
        userDatabase.put("0987654321", "securepass");
        
        verificationCodes = new HashMap<>();
    }

    /**
     * Authenticates a user with phone number and password
     * @param phoneNumber User's phone number
     * @param password User's password
     * @return AuthenticationResult containing status and user information
     * @throws IllegalArgumentException if inputs are invalid
     */
    public AuthenticationResult authenticate(String phoneNumber, String password) {
        // Input validation
        if (phoneNumber == null || phoneNumber.trim().isEmpty()) {
            throw new IllegalArgumentException("Phone number cannot be empty");
        }
        if (password == null || password.trim().isEmpty()) {
            throw new IllegalArgumentException("Password cannot be empty");
        }
        if (phoneNumber.length() != 10) {
            throw new IllegalArgumentException("Phone number must be 10 digits");
        }

        // Check if user exists
        if (!userDatabase.containsKey(phoneNumber)) {
            return new AuthenticationResult(404, "User not found", null);
        }

        // Verify password
        if (!userDatabase.get(phoneNumber).equals(password)) {
            return new AuthenticationResult(401, "Invalid password", null);
        }

        // Generate verification code for 2FA (simplified)
        String verificationCode = generateVerificationCode();
        verificationCodes.put(phoneNumber, verificationCode);

        // Return success with verification required status
        User user = new User(phoneNumber, "Test User");
        return new AuthenticationResult(250, "Verification required", user);
    }

    /**
     * Verifies a user's verification code
     * @param phoneNumber User's phone number
     * @param code Verification code to check
     * @return AuthenticationResult containing status and user information
     */
    public AuthenticationResult verifyCode(String phoneNumber, String code) {
        if (!verificationCodes.containsKey(phoneNumber)) {
            return new AuthenticationResult(404, "No pending verification for this user", null);
        }

        if (verificationCodes.get(phoneNumber).equals(code)) {
            User user = new User(phoneNumber, "Test User");
            verificationCodes.remove(phoneNumber);
            return new AuthenticationResult(200, "Verification successful", user);
        }

        return new AuthenticationResult(401, "Invalid verification code", null);
    }

    private String generateVerificationCode() {
        // Simplified code generation - in real app would send SMS
        return "123456";
    }

    // Nested classes for structured responses
    class AuthenticationResult {
        private int statusCode;
        private String message;
        private User user;

        public AuthenticationResult(int statusCode, String message, User user) {
            this.statusCode = statusCode;
            this.message = message;
            this.user = user;
        }

        // Getters
        public int getStatusCode() { return statusCode; }
        public String getMessage() { return message; }
        public User getUser() { return user; }
    }

    class User {
        private String phoneNumber;
        private String fullName;

        public User(String phoneNumber, String fullName) {
            this.phoneNumber = phoneNumber;
            this.fullName = fullName;
        }

        // Getters
        public String getPhoneNumber() { return phoneNumber; }
        public String getFullName() { return fullName; }
    }
}
