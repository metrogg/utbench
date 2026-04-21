// Converted Java method
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;

class UserManager {
    private static final Pattern EMAIL_PATTERN = Pattern.compile("^[A-Za-z0-9+_.-]+@(.+)$");
    private final Map<String, User> usersByEmail = new HashMap<>();
    private final Map<Integer, User> usersById = new HashMap<>();
    private int nextUserId = 1;

    /**
     * Represents a user with basic information
     */
    public static class User {
        private final int id;
        private final String email;
        private final String name;
        private boolean active;

        public User(int id, String email, String name) {
            this.id = id;
            this.email = email;
            this.name = name;
            this.active = true;
        }

        // Getters and setters
        public int getId() { return id; }
        public String getEmail() { return email; }
        public String getName() { return name; }
        public boolean isActive() { return active; }
        public void setActive(boolean active) { this.active = active; }
    }

    /**
     * Adds a new user to the system after validation
     * @param email User's email address
     * @param name User's full name
     * @return The created User object
     * @throws IllegalArgumentException if email is invalid or already exists
     */
    public User addUser(String email, String name) {
        if (!isValidEmail(email)) {
            throw new IllegalArgumentException("Invalid email format");
        }
        if (usersByEmail.containsKey(email)) {
            throw new IllegalArgumentException("Email already exists");
        }
        if (name == null || name.trim().isEmpty()) {
            throw new IllegalArgumentException("Name cannot be empty");
        }

        User newUser = new User(nextUserId++, email, name);
        usersByEmail.put(email, newUser);
        usersById.put(newUser.getId(), newUser);
        return newUser;
    }

    /**
     * Finds a user by email address
     * @param email Email address to search for
     * @return The User object if found
     * @throws IllegalArgumentException if email is null or user not found
     */
    public User findUserByEmail(String email) {
        if (email == null) {
            throw new IllegalArgumentException("Email cannot be null");
        }
        User user = usersByEmail.get(email);
        if (user == null) {
            throw new IllegalArgumentException("User not found");
        }
        return user;
    }

    /**
     * Gets all active users in the system
     * @return List of active users
     */
    public List<User> getActiveUsers() {
        List<User> activeUsers = new ArrayList<>();
        for (User user : usersByEmail.values()) {
            if (user.isActive()) {
                activeUsers.add(user);
            }
        }
        return activeUsers;
    }

    /**
     * Validates email format using regex pattern
     * @param email Email address to validate
     * @return true if email is valid, false otherwise
     */
    private boolean isValidEmail(String email) {
        return email != null && EMAIL_PATTERN.matcher(email).matches();
    }
}
