// Converted Java method
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Collectors;

class UserManager {
    private final List<User> users = new ArrayList<>();
    private final AtomicLong idCounter = new AtomicLong(1);

    /**
     * Represents a user with basic information and credentials
     */
    static class User {
        private final Long id;
        private String username;
        private String email;
        private String password;
        private boolean active;

        public User(Long id, String username, String email, String password) {
            this.id = id;
            this.username = username;
            this.email = email;
            this.password = password;
            this.active = true;
        }

        // Getters and setters
        public Long getId() { return id; }
        public String getUsername() { return username; }
        public void setUsername(String username) { this.username = username; }
        public String getEmail() { return email; }
        public void setEmail(String email) { this.email = email; }
        public String getPassword() { return password; }
        public void setPassword(String password) { this.password = password; }
        public boolean isActive() { return active; }
        public void setActive(boolean active) { this.active = active; }
    }

    /**
     * Creates a new user with the provided information
     * @param username User's username
     * @param email User's email
     * @param password User's password
     * @return The created user object
     * @throws IllegalArgumentException if username or email already exists
     */
    public User createUser(String username, String email, String password) {
        if (usernameExists(username)) {
            throw new IllegalArgumentException("Username already exists");
        }
        if (emailExists(email)) {
            throw new IllegalArgumentException("Email already exists");
        }

        User newUser = new User(idCounter.getAndIncrement(), username, email, password);
        users.add(newUser);
        return newUser;
    }

    /**
     * Updates an existing user's information
     * @param id User ID to update
     * @param username New username (optional)
     * @param email New email (optional)
     * @param password New password (optional)
     * @return The updated user
     * @throws IllegalArgumentException if user not found or validation fails
     */
    public User updateUser(Long id, String username, String email, String password) {
        User user = findUserById(id)
                .orElseThrow(() -> new IllegalArgumentException("User not found"));

        if (username != null && !username.equals(user.getUsername())) {
            if (usernameExists(username)) {
                throw new IllegalArgumentException("Username already exists");
            }
            user.setUsername(username);
        }

        if (email != null && !email.equals(user.getEmail())) {
            if (emailExists(email)) {
                throw new IllegalArgumentException("Email already exists");
            }
            user.setEmail(email);
        }

        if (password != null) {
            user.setPassword(password);
        }

        return user;
    }

    /**
     * Deactivates a user account
     * @param id User ID to deactivate
     * @return true if deactivated successfully, false if user not found
     */
    public boolean deactivateUser(Long id) {
        Optional<User> user = findUserById(id);
        user.ifPresent(u -> u.setActive(false));
        return user.isPresent();
    }

    /**
     * Retrieves all active users
     * @return List of active users
     */
    public List<User> getAllActiveUsers() {
        return users.stream()
                .filter(User::isActive)
                .collect(Collectors.toList());
    }

    /**
     * Finds a user by ID
     * @param id User ID to search for
     * @return Optional containing the user if found
     */
    public Optional<User> findUserById(Long id) {
        return users.stream()
                .filter(u -> u.getId().equals(id))
                .findFirst();
    }

    private boolean usernameExists(String username) {
        return users.stream().anyMatch(u -> u.getUsername().equals(username));
    }

    private boolean emailExists(String email) {
        return users.stream().anyMatch(u -> u.getEmail().equals(email));
    }
}
