// Converted Java method
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

class AuthTokenValidator {
    private Map<String, Token> tokenDatabase;
    private Map<String, User> userDatabase;

    public AuthTokenValidator() {
        // Initialize in-memory databases for testing
        this.tokenDatabase = new HashMap<>();
        this.userDatabase = new HashMap<>();
        
        // Add test data
        User testUser = new User("user123", "Test User");
        userDatabase.put("user123", testUser);
        
        Token validToken = new Token("auth_user192.168.1.1", "user123", 
                                   new Date(System.currentTimeMillis() + 3600000), 1);
        tokenDatabase.put("auth_user192.168.1.1", validToken);
        
        Token expiredToken = new Token("expired_user192.168.1.2", "user123", 
                                     new Date(System.currentTimeMillis() - 3600000), 1);
        tokenDatabase.put("expired_user192.168.1.2", expiredToken);
        
        Token invalidToken = new Token("invalid_user192.168.1.3", "user123", 
                                     new Date(System.currentTimeMillis() + 3600000), 0);
        tokenDatabase.put("invalid_user192.168.1.3", invalidToken);
    }

    /**
     * Validates an authentication token against IP address and token status
     * 
     * @param token The authentication token to validate
     * @param clientIp The client's IP address
     * @return Map containing validation result and user information if valid
     */
    public Map<String, Object> validateToken(String token, String clientIp) {
        Map<String, Object> result = new HashMap<>();
        result.put("valid", false);
        
        if (token == null || token.isEmpty()) {
            result.put("message", "Token is empty");
            return result;
        }
        
        if (!token.contains("user")) {
            result.put("message", "Token format invalid");
            return result;
        }
        
        String tokenIp = token.substring(token.indexOf("user") + 4);
        if (!clientIp.equals(tokenIp)) {
            result.put("message", "IP address mismatch");
            return result;
        }
        
        Token tokenRecord = tokenDatabase.get(token);
        if (tokenRecord == null) {
            result.put("message", "Token not found");
            return result;
        }
        
        if (tokenRecord.getExpired().before(new Date())) {
            result.put("message", "Token expired");
            return result;
        }
        
        if (tokenRecord.getStatus() != 1) {
            result.put("message", "Token invalid");
            return result;
        }
        
        User user = userDatabase.get(tokenRecord.getOpenid());
        if (user == null) {
            result.put("message", "User not found");
            return result;
        }
        
        result.put("valid", true);
        result.put("user", user);
        result.put("message", "Validation successful");
        return result;
    }
}

// Supporting classes
class Token {
    private String token;
    private String openid;
    private Date expired;
    private int status;
    
    public Token(String token, String openid, Date expired, int status) {
        this.token = token;
        this.openid = openid;
        this.expired = expired;
        this.status = status;
    }
    
    // Getters
    public String getToken() { return token; }
    public String getOpenid() { return openid; }
    public Date getExpired() { return expired; }
    public int getStatus() { return status; }
}

class User {
    private String id;
    private String name;
    
    public User(String id, String name) {
        this.id = id;
        this.name = name;
    }
    
    // Getters
    public String getId() { return id; }
    public String getName() { return name; }
    
    @Override
    public String toString() {
        return "User{id='" + id + "', name='" + name + "'}";
    }
}
