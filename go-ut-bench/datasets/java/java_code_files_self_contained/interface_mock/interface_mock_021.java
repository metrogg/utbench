// Converted Java method
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

class UserAuthService {
    private final Map<String, User> accountCache = new ConcurrentHashMap<>();
    private final Map<String, User> emailCache = new ConcurrentHashMap<>();
    private final Map<Long, User> idCache = new ConcurrentHashMap<>();
    
    /**
     * Authenticates a user using either account, email, or ID with password verification.
     * Implements caching mechanism similar to Spring's @Cacheable.
     * 
     * @param identifier Can be account (String), email (String), or ID (Long)
     * @param password The password to verify
     * @return Map containing authentication status and user data if successful
     * @throws IllegalArgumentException if identifier type is invalid
     */
    public Map<String, Object> authenticateUser(Object identifier, String password) {
        Map<String, Object> result = new HashMap<>();
        User user = null;
        
        // Check caches first (simulating @Cacheable behavior)
        if (identifier instanceof String) {
            String strId = (String) identifier;
            if (strId.contains("@")) {
                user = emailCache.get(strId);
            } else {
                user = accountCache.get(strId);
            }
        } else if (identifier instanceof Long) {
            user = idCache.get((Long) identifier);
        } else {
            throw new IllegalArgumentException("Identifier must be String (account/email) or Long (ID)");
        }
        
        // If not found in cache, simulate database lookup
        if (user == null) {
            user = findUserInDatabase(identifier);
            if (user != null) {
                // Cache the user (simulating @CachePut)
                cacheUser(user);
            }
        }
        
        // Verify password
        if (user != null && user.getPassword().equals(password)) {
            result.put("authenticated", true);
            result.put("user", user);
        } else {
            result.put("authenticated", false);
            result.put("message", "Invalid credentials");
        }
        
        return result;
    }
    
    /**
     * Updates user password and clears relevant cache entries.
     * Simulates @CachePut behavior.
     * 
     * @param identifier Can be account, email, or ID
     * @param newPassword The new password to set
     * @return true if update was successful, false otherwise
     */
    public boolean updatePassword(Object identifier, String newPassword) {
        User user = findUserInDatabase(identifier);
        if (user != null) {
            user.setPassword(newPassword);
            // Update database would happen here in real implementation
            cacheUser(user); // Update cache
            return true;
        }
        return false;
    }
    
    // Simulated database access methods
    private User findUserInDatabase(Object identifier) {
        // In a real implementation, this would query the database
        // For demo purposes, we'll return a mock user
        if (identifier instanceof String) {
            String strId = (String) identifier;
            if (strId.equals("testuser") || strId.equals("test@example.com")) {
                return new User(1L, "testuser", "test@example.com", "oldpassword");
            }
        } else if (identifier instanceof Long && (Long)identifier == 1L) {
            return new User(1L, "testuser", "test@example.com", "oldpassword");
        }
        return null;
    }
    
    private void cacheUser(User user) {
        accountCache.put(user.getAccount(), user);
        emailCache.put(user.getEmail(), user);
        idCache.put(user.getId(), user);
    }
    
    public static class User {
        private Long id;
        private String account;
        private String email;
        private String password;
        
        public User(Long id, String account, String email, String password) {
            this.id = id;
            this.account = account;
            this.email = email;
            this.password = password;
        }
        
        // Getters and setters
        public Long getId() { return id; }
        public String getAccount() { return account; }
        public String getEmail() { return email; }
        public String getPassword() { return password; }
        public void setPassword(String password) { this.password = password; }
    }
}
