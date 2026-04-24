// Converted Java method - Advanced User Management System
import java.util.regex.Pattern;
import java.util.ArrayList;
import java.util.List;

class UserManager {
    private static final Pattern STEAM_ID_PATTERN = Pattern.compile("^\\d{17}$");
    private static final Pattern URL_PATTERN = Pattern.compile("^(https?://)?(www\\.)?steamcommunity\\.com/.*$");
    private List<User> users = new ArrayList<>();

    /**
     * Validates and adds a new user to the system
     * @param steam32Id 17-digit Steam ID
     * @param username User's display name (2-32 characters)
     * @param profileURL Valid Steam community URL
     * @param avatarMedium Medium avatar URL
     * @param avatarFull Full avatar URL
     * @throws IllegalArgumentException if any validation fails
     */
    public void addUser(long steam32Id, String username, String profileURL, 
                       String avatarMedium, String avatarFull) {
        validateSteamId(steam32Id);
        validateUsername(username);
        validateProfileURL(profileURL);
        
        if (avatarMedium == null || avatarMedium.isBlank()) {
            throw new IllegalArgumentException("Medium avatar URL cannot be empty");
        }
        if (avatarFull == null || avatarFull.isBlank()) {
            throw new IllegalArgumentException("Full avatar URL cannot be empty");
        }
        
        if (users.stream().anyMatch(u -> u.getSteam32Id() == steam32Id)) {
            throw new IllegalArgumentException("User with this Steam ID already exists");
        }
        
        users.add(new User(steam32Id, username, profileURL, avatarMedium, avatarFull));
    }

    /**
     * Finds a user by Steam ID
     * @param steam32Id Steam ID to search for
     * @return User object if found, null otherwise
     */
    public User getUserById(long steam32Id) {
        return users.stream()
                   .filter(u -> u.getSteam32Id() == steam32Id)
                   .findFirst()
                   .orElse(null);
    }

    /**
     * Updates user information
     * @param steam32Id Steam ID of user to update
     * @param newUsername New username (null to keep current)
     * @param newProfileURL New profile URL (null to keep current)
     */
    public void updateUser(long steam32Id, String newUsername, String newProfileURL) {
        User user = getUserById(steam32Id);
        if (user == null) {
            throw new IllegalArgumentException("User not found");
        }
        
        if (newUsername != null) {
            validateUsername(newUsername);
            user.setUsername(newUsername);
        }
        
        if (newProfileURL != null) {
            validateProfileURL(newProfileURL);
            user.setProfileURL(newProfileURL);
        }
    }

    private void validateSteamId(long steam32Id) {
        if (!STEAM_ID_PATTERN.matcher(String.valueOf(steam32Id)).matches()) {
            throw new IllegalArgumentException("Invalid Steam ID format");
        }
    }

    private void validateUsername(String username) {
        if (username == null || username.length() < 2 || username.length() > 32) {
            throw new IllegalArgumentException("Username must be 2-32 characters");
        }
    }

    private void validateProfileURL(String profileURL) {
        if (profileURL == null || !URL_PATTERN.matcher(profileURL).matches()) {
            throw new IllegalArgumentException("Invalid Steam profile URL");
        }
    }
}

class User {
    private long steam32Id;
    private String username;
    private String profileURL;
    private String avatarMedium;
    private String avatarFull;

    public User(long steam32Id, String username, String profileURL, 
                String avatarMedium, String avatarFull) {
        this.steam32Id = steam32Id;
        this.username = username;
        this.profileURL = profileURL;
        this.avatarMedium = avatarMedium;
        this.avatarFull = avatarFull;
    }

    // Getters and setters
    public long getSteam32Id() { return steam32Id; }
    public String getUsername() { return username; }
    public String getProfileURL() { return profileURL; }
    public String getAvatarMedium() { return avatarMedium; }
    public String getAvatarFull() { return avatarFull; }
    public void setUsername(String username) { this.username = username; }
    public void setProfileURL(String profileURL) { this.profileURL = profileURL; }
}
