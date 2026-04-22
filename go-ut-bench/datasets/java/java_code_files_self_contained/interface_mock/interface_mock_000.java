// Converted Java method - Enhanced Authentication Service
import java.util.HashMap;
import java.util.Map;

class AuthenticationService {
    private final Map<String, String> userCredentials;
    private final PasswordValidator passwordValidator;

    public AuthenticationService() {
        this.userCredentials = new HashMap<>();
        this.passwordValidator = new PasswordValidator();
        initializeDefaultUsers();
    }

    /**
     * Initializes the system with default user credentials
     */
    private void initializeDefaultUsers() {
        userCredentials.put("Om", "imagine");
        userCredentials.put("Chinmay", "IMAGINE");
    }

    /**
     * Authenticates a user with the given credentials
     * @param username The username to authenticate
     * @param password The password to verify
     * @return AuthenticationResult containing success status and any messages
     */
    public AuthenticationResult authenticate(String username, String password) {
        if (username == null || username.trim().isEmpty()) {
            return new AuthenticationResult(false, "Username cannot be empty");
        }

        if (password == null || password.trim().isEmpty()) {
            return new AuthenticationResult(false, "Password cannot be empty");
        }

        if (!userCredentials.containsKey(username)) {
            return new AuthenticationResult(false, "Invalid username");
        }

        String storedPassword = userCredentials.get(username);
        if (!storedPassword.equals(password)) {
            return new AuthenticationResult(false, "Invalid password");
        }

        return new AuthenticationResult(true, "Authentication successful");
    }

    /**
     * Adds a new user to the system with password validation
     * @param username The username to add
     * @param password The password for the new user
     * @return RegistrationResult containing success status and any messages
     */
    public RegistrationResult registerUser(String username, String password) {
        if (username == null || username.trim().isEmpty()) {
            return new RegistrationResult(false, "Username cannot be empty");
        }

        if (userCredentials.containsKey(username)) {
            return new RegistrationResult(false, "Username already exists");
        }

        PasswordValidationResult validation = passwordValidator.validate(password);
        if (!validation.isValid()) {
            return new RegistrationResult(false, validation.getMessage());
        }

        userCredentials.put(username, password);
        return new RegistrationResult(true, "User registered successfully");
    }

    /**
     * Nested class for password validation
     */
    private static class PasswordValidator {
        private static final int MIN_LENGTH = 8;
        private static final int MAX_LENGTH = 20;

        public PasswordValidationResult validate(String password) {
            if (password.length() < MIN_LENGTH) {
                return new PasswordValidationResult(false, 
                    "Password must be at least " + MIN_LENGTH + " characters long");
            }

            if (password.length() > MAX_LENGTH) {
                return new PasswordValidationResult(false,
                    "Password cannot exceed " + MAX_LENGTH + " characters");
            }

            if (!password.matches(".*[A-Z].*")) {
                return new PasswordValidationResult(false,
                    "Password must contain at least one uppercase letter");
            }

            if (!password.matches(".*[a-z].*")) {
                return new PasswordValidationResult(false,
                    "Password must contain at least one lowercase letter");
            }

            if (!password.matches(".*\\d.*")) {
                return new PasswordValidationResult(false,
                    "Password must contain at least one digit");
            }

            return new PasswordValidationResult(true, "Password is valid");
        }
    }

    // Result classes for better type safety and information passing
    public static class AuthenticationResult {
        private final boolean success;
        private final String message;

        public AuthenticationResult(boolean success, String message) {
            this.success = success;
            this.message = message;
        }

        public boolean isSuccess() { return success; }
        public String getMessage() { return message; }
    }

    public static class RegistrationResult {
        private final boolean success;
        private final String message;

        public RegistrationResult(boolean success, String message) {
            this.success = success;
            this.message = message;
        }

        public boolean isSuccess() { return success; }
        public String getMessage() { return message; }
    }

    public static class PasswordValidationResult {
        private final boolean valid;
        private final String message;

        public PasswordValidationResult(boolean valid, String message) {
            this.valid = valid;
            this.message = message;
        }

        public boolean isValid() { return valid; }
        public String getMessage() { return message; }
    }
}
