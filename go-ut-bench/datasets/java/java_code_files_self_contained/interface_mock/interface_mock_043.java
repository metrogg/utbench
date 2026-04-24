// Converted Java method
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Base64;

class UserAuthenticationService {
    private Map<String, User> userDatabase;
    private MessageDigest digest;

    public UserAuthenticationService() throws NoSuchAlgorithmException {
        this.userDatabase = new HashMap<>();
        this.digest = MessageDigest.getInstance("SHA-256");
        // Initialize with some test users
        initializeTestUsers();
    }

    private void initializeTestUsers() {
        // Admin user
        String adminHash = hashPassword("admin123");
        userDatabase.put("admin@example.com", new User("admin@example.com", adminHash, "ADMIN"));
        
        // Regular user
        String userHash = hashPassword("user123");
        userDatabase.put("user@example.com", new User("user@example.com", userHash, "USER"));
    }

    private String hashPassword(String password) {
        byte[] hashBytes = digest.digest(password.getBytes());
        return Base64.getEncoder().encodeToString(hashBytes);
    }

    /**
     * Authenticates a user with email and password
     * @param email User's email
     * @param password User's password
     * @return Optional containing User if authentication succeeds, empty otherwise
     */
    public Optional<User> authenticate(String email, String password) {
        User user = userDatabase.get(email);
        if (user != null && user.getPasswordHash().equals(hashPassword(password))) {
            return Optional.of(user);
        }
        return Optional.empty();
    }

    /**
     * Checks if a user has a specific role
     * @param email User's email
     * @param role Required role
     * @return true if user exists and has the required role, false otherwise
     */
    public boolean hasRole(String email, String role) {
        User user = userDatabase.get(email);
        return user != null && user.getRole().equals(role);
    }

    /**
     * Finds a user by email (simulating the original findByEmail method)
     * @param email User's email
     * @return Optional containing User if found, empty otherwise
     */
    public Optional<User> findByEmail(String email) {
        return Optional.ofNullable(userDatabase.get(email));
    }

    static class User {
        private String email;
        private String passwordHash;
        private String role;

        public User(String email, String passwordHash, String role) {
            this.email = email;
            this.passwordHash = passwordHash;
            this.role = role;
        }

        public String getEmail() { return email; }
        public String getPasswordHash() { return passwordHash; }
        public String getRole() { return role; }
    }
}
