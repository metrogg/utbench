// Converted Java method
import java.util.ArrayList;
import java.util.List;

class UserService {
    /**
     * Simulates a user database with enhanced user management functionality.
     * This is a self-contained version that doesn't require Spring or database connections.
     */
    private List<User> userDatabase;

    public UserService() {
        // Initialize with some test users
        this.userDatabase = new ArrayList<>();
        userDatabase.add(new User("kim1", "kim1@example.com", 25));
        userDatabase.add(new User("john2", "john2@example.com", 30));
        userDatabase.add(new User("sarah3", "sarah3@example.com", 28));
    }

    /**
     * Finds all users in the database.
     * @return List of all users
     */
    public List<User> findAllUserInfo() {
        return new ArrayList<>(userDatabase); // Return a copy to prevent modification
    }

    /**
     * Finds a user by username with enhanced validation.
     * @param username The username to search for
     * @return The found user or null if not found
     * @throws IllegalArgumentException if username is null or empty
     */
    public User findByOneUserName(String username) {
        if (username == null || username.trim().isEmpty()) {
            throw new IllegalArgumentException("Username cannot be null or empty");
        }
        
        return userDatabase.stream()
                .filter(user -> user.getUsername().equals(username))
                .findFirst()
                .orElse(null);
    }

    /**
     * Adds a new user to the database with validation.
     * @param username The username
     * @param email The email
     * @param age The age
     * @throws IllegalArgumentException if any parameter is invalid
     */
    public void addUser(String username, String email, int age) {
        if (username == null || username.trim().isEmpty()) {
            throw new IllegalArgumentException("Username cannot be null or empty");
        }
        if (email == null || !email.contains("@")) {
            throw new IllegalArgumentException("Invalid email format");
        }
        if (age <= 0) {
            throw new IllegalArgumentException("Age must be positive");
        }
        
        userDatabase.add(new User(username, email, age));
    }

    /**
     * Inner class representing a User entity
     */
    static class User {
        private String username;
        private String email;
        private int age;

        public User(String username, String email, int age) {
            this.username = username;
            this.email = email;
            this.age = age;
        }

        // Getters
        public String getUsername() { return username; }
        public String getEmail() { return email; }
        public int getAge() { return age; }

        @Override
        public String toString() {
            return "User{" +
                    "username='" + username + '\'' +
                    ", email='" + email + '\'' +
                    ", age=" + age +
                    '}';
        }
    }
}
