// Converted Java method
import java.util.HashMap;
import java.util.Map;
import java.util.regex.Pattern;

class UserManager {
    private Map<String, User> userDatabase;
    private Map<String, String> activeSessions;
    
    public UserManager() {
        this.userDatabase = new HashMap<>();
        this.activeSessions = new HashMap<>();
    }
    
    /**
     * Validates user credentials and creates a session if successful
     * @param username The username to login with
     * @param password The password to verify
     * @return Map containing success status and session token if successful
     */
    public Map<String, Object> loginUser(String username, String password) {
        Map<String, Object> response = new HashMap<>();
        
        if (!userDatabase.containsKey(username)) {
            response.put("success", false);
            response.put("message", "User not found");
            return response;
        }
        
        User user = userDatabase.get(username);
        if (!user.getPassword().equals(hashPassword(password))) {
            response.put("success", false);
            response.put("message", "Invalid password");
            return response;
        }
        
        String sessionToken = generateSessionToken();
        activeSessions.put(sessionToken, username);
        response.put("success", true);
        response.put("sessionToken", sessionToken);
        response.put("user", user);
        return response;
    }
    
    /**
     * Registers a new user with validation
     * @param username The desired username
     * @param password The desired password
     * @param email The user's email address
     * @return Map containing success status and message
     */
    public Map<String, Object> registerUser(String username, String password, String email) {
        Map<String, Object> response = new HashMap<>();
        
        // Validate inputs
        if (userDatabase.containsKey(username)) {
            response.put("success", false);
            response.put("message", "Username already exists");
            return response;
        }
        
        if (!isValidEmail(email)) {
            response.put("success", false);
            response.put("message", "Invalid email format");
            return response;
        }
        
        if (password.length() < 8) {
            response.put("success", false);
            response.put("message", "Password must be at least 8 characters");
            return response;
        }
        
        // Create and store new user
        User newUser = new User(username, hashPassword(password), email);
        userDatabase.put(username, newUser);
        
        response.put("success", true);
        response.put("message", "User registered successfully");
        return response;
    }
    
    /**
     * Logs out a user by invalidating their session
     * @param sessionToken The session token to invalidate
     * @return Map containing success status
     */
    public Map<String, Object> logoutUser(String sessionToken) {
        Map<String, Object> response = new HashMap<>();
        
        if (activeSessions.containsKey(sessionToken)) {
            activeSessions.remove(sessionToken);
            response.put("success", true);
            response.put("message", "Logged out successfully");
        } else {
            response.put("success", false);
            response.put("message", "Invalid session token");
        }
        
        return response;
    }
    
    // Helper methods
    private String hashPassword(String password) {
        // In a real application, use a proper hashing algorithm like BCrypt
        return Integer.toString(password.hashCode());
    }
    
    private String generateSessionToken() {
        return java.util.UUID.randomUUID().toString();
    }
    
    private boolean isValidEmail(String email) {
        String emailRegex = "^[a-zA-Z0-9_+&*-]+(?:\\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,7}$";
        Pattern pattern = Pattern.compile(emailRegex);
        return pattern.matcher(email).matches();
    }
}

class User {
    private String username;
    private String password;
    private String email;
    
    public User(String username, String password, String email) {
        this.username = username;
        this.password = password;
        this.email = email;
    }
    
    // Getters and setters
    public String getUsername() { return username; }
    public String getPassword() { return password; }
    public String getEmail() { return email; }
}
