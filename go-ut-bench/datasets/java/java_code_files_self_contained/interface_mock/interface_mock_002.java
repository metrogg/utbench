// Converted Java method
import java.util.*;
import java.util.stream.Collectors;

class AuthService {
    /**
     * Simulates OAuth2 user registration and authentication process.
     * Handles both new user registration and existing user authentication.
     * Includes enhanced validation and role management.
     * 
     * @param provider The OAuth2 provider name (e.g., "google", "facebook")
     * @param userAttributes Map of user attributes from OAuth2 provider
     * @param existingUsers Map simulating user repository (email -> User)
     * @param availableRoles Set of available roles in the system
     * @return Map containing authentication result and user details
     * @throws IllegalArgumentException for invalid inputs
     */
    public Map<String, Object> processOAuth2User(String provider, 
                                               Map<String, Object> userAttributes,
                                               Map<String, User> existingUsers,
                                               Set<Role> availableRoles) {
        // Validate inputs
        if (provider == null || provider.trim().isEmpty()) {
            throw new IllegalArgumentException("OAuth2 provider cannot be null or empty");
        }
        
        if (userAttributes == null || userAttributes.isEmpty()) {
            throw new IllegalArgumentException("User attributes cannot be null or empty");
        }
        
        // Extract user info
        String email = (String) userAttributes.get("email");
        if (email == null || email.trim().isEmpty()) {
            throw new IllegalArgumentException("Email not found in OAuth2 user attributes");
        }
        
        // Process user
        User user;
        boolean isNewUser = !existingUsers.containsKey(email);
        
        if (isNewUser) {
            // Register new user
            user = registerNewUser(provider, userAttributes, availableRoles);
            existingUsers.put(email, user);
        } else {
            // Get existing user
            user = existingUsers.get(email);
            
            // Update user attributes if needed
            user.setName((String) userAttributes.getOrDefault("given_name", user.getName()));
            user.setLastname((String) userAttributes.getOrDefault("family_name", user.getLastname()));
        }
        
        // Prepare result
        Map<String, Object> result = new HashMap<>();
        result.put("status", isNewUser ? "registered" : "authenticated");
        result.put("user", user);
        result.put("provider", provider);
        result.put("roles", user.getRoles().stream()
                              .map(Role::getName)
                              .collect(Collectors.toSet()));
        
        return result;
    }
    
    private User registerNewUser(String provider, 
                               Map<String, Object> attributes,
                               Set<Role> availableRoles) {
        User user = new User();
        user.setEmail((String) attributes.get("email"));
        user.setName((String) attributes.getOrDefault("given_name", ""));
        user.setLastname((String) attributes.getOrDefault("family_name", ""));
        user.setEnabled(true);
        
        // Generate username from email and provider
        String emailPrefix = user.getEmail().split("@")[0];
        user.setUsername(emailPrefix + "_" + provider);
        
        // Assign default role
        Set<Role> roles = new HashSet<>();
        Role defaultRole = availableRoles.stream()
            .filter(r -> "ROLE_USER".equals(r.getName()))
            .findFirst()
            .orElseThrow(() -> new IllegalStateException("Default role not found"));
        roles.add(defaultRole);
        user.setRoles(roles);
        
        return user;
    }
}

// Supporting classes
class User {
    private String username;
    private String name;
    private String lastname;
    private String email;
    private boolean enabled;
    private Set<Role> roles;
    
    // Getters and setters
    public String getUsername() { return username; }
    public void setUsername(String username) { this.username = username; }
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
    public String getLastname() { return lastname; }
    public void setLastname(String lastname) { this.lastname = lastname; }
    public String getEmail() { return email; }
    public void setEmail(String email) { this.email = email; }
    public boolean isEnabled() { return enabled; }
    public void setEnabled(boolean enabled) { this.enabled = enabled; }
    public Set<Role> getRoles() { return roles; }
    public void setRoles(Set<Role> roles) { this.roles = roles; }
    
    @Override
    public String toString() {
        return "User{" +
               "username='" + username + '\'' +
               ", email='" + email + '\'' +
               ", name='" + name + '\'' +
               ", lastname='" + lastname + '\'' +
               ", enabled=" + enabled +
               ", roles=" + roles.stream().map(Role::getName).collect(Collectors.toSet()) +
               '}';
    }
}

class Role {
    private String name;
    
    public Role(String name) { this.name = name; }
    
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
}
