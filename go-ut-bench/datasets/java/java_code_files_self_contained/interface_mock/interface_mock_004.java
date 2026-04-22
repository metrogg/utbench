// Converted Java method
import java.util.List;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.stream.Collectors;

class UserManager {
    private List<User> users;
    private Map<String, Integer> loginAttempts;

    public UserManager() {
        this.users = new ArrayList<>();
        this.loginAttempts = new HashMap<>();
    }

    /**
     * Registers a new user with validation checks
     * @param user User object to register
     * @return Map containing success status and message
     */
    public Map<String, Object> registerUser(User user) {
        Map<String, Object> response = new HashMap<>();
        
        if (user.getUsername() == null || user.getUsername().trim().isEmpty()) {
            response.put("success", false);
            response.put("message", "Username cannot be empty");
            return response;
        }
        
        if (user.getPassword() == null || user.getPassword().trim().isEmpty()) {
            response.put("success", false);
            response.put("message", "Password cannot be empty");
            return response;
        }
        
        boolean usernameExists = users.stream()
            .anyMatch(u -> u.getUsername().equalsIgnoreCase(user.getUsername()));
            
        if (usernameExists) {
            response.put("success", false);
            response.put("message", "Username already exists");
            return response;
        }
        
        users.add(user);
        response.put("success", true);
        response.put("message", "User registered successfully");
        return response;
    }

    /**
     * Authenticates a user with login attempt tracking
     * @param username User's username
     * @param password User's password
     * @return Map containing login status and user data if successful
     */
    public Map<String, Object> loginUser(String username, String password) {
        Map<String, Object> response = new HashMap<>();
        
        User user = users.stream()
            .filter(u -> u.getUsername().equals(username))
            .findFirst()
            .orElse(null);
            
        if (user == null) {
            response.put("success", false);
            response.put("message", "User not found");
            return response;
        }
        
        int attempts = loginAttempts.getOrDefault(username, 0);
        if (attempts >= 3) {
            response.put("success", false);
            response.put("message", "Account locked due to too many failed attempts");
            return response;
        }
        
        if (!user.getPassword().equals(password)) {
            loginAttempts.put(username, attempts + 1);
            response.put("success", false);
            response.put("message", "Invalid password");
            response.put("attemptsRemaining", 3 - (attempts + 1));
            return response;
        }
        
        loginAttempts.remove(username);
        response.put("success", true);
        response.put("message", "Login successful");
        response.put("user", user);
        return response;
    }

    /**
     * Searches users based on criteria with pagination
     * @param criteria Search criteria map
     * @param page Page number (1-based)
     * @param pageSize Number of items per page
     * @return List of matching users
     */
    public List<User> searchUsers(Map<String, String> criteria, int page, int pageSize) {
        List<User> filtered = users.stream()
            .filter(user -> {
                boolean matches = true;
                if (criteria.containsKey("username") && 
                    !user.getUsername().contains(criteria.get("username"))) {
                    matches = false;
                }
                if (criteria.containsKey("email") && 
                    !user.getEmail().contains(criteria.get("email"))) {
                    matches = false;
                }
                return matches;
            })
            .collect(Collectors.toList());
            
        int start = (page - 1) * pageSize;
        if (start >= filtered.size()) {
            return new ArrayList<>();
        }
        
        int end = Math.min(start + pageSize, filtered.size());
        return filtered.subList(start, end);
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
    public void setUsername(String username) { this.username = username; }
    public void setPassword(String password) { this.password = password; }
    public void setEmail(String email) { this.email = email; }
}
