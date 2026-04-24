import java.util.HashMap;
import java.util.Map;

class AuthenticationService {
    /**
     * Simulates user authentication with email/password or Google sign-in.
     * This is a simplified version of the Android authentication logic.
     * 
     * @param email User's email (can be null for Google auth)
     * @param password User's password (can be null for Google auth)
     * @param googleToken Google authentication token (can be null for email/password auth)
     * @return Map containing authentication status and user information
     */
    public Map<String, Object> authenticateUser(String email, String password, String googleToken) {
        Map<String, Object> result = new HashMap<>();
        
        // Validate inputs
        if ((email == null || password == null) && googleToken == null) {
            result.put("success", false);
            result.put("message", "Authentication failed: No credentials provided");
            return result;
        }
        
        // Email/password authentication
        if (email != null && password != null) {
            if (email.isEmpty() || password.isEmpty()) {
                result.put("success", false);
                result.put("message", "Email and password must not be empty");
                return result;
            }
            
            // Simulate authentication (in real app this would call Firebase)
            if (isValidCredentials(email, password)) {
                result.put("success", true);
                result.put("userEmail", email);
                result.put("authMethod", "email");
                return result;
            } else {
                result.put("success", false);
                result.put("message", "Invalid email or password");
                return result;
            }
        }
        
        // Google authentication
        if (googleToken != null) {
            // Simulate Google authentication (in real app this would verify with Google)
            String googleEmail = verifyGoogleToken(googleToken);
            if (googleEmail != null) {
                result.put("success", true);
                result.put("userEmail", googleEmail);
                result.put("authMethod", "google");
                return result;
            } else {
                result.put("success", false);
                result.put("message", "Google authentication failed");
                return result;
            }
        }
        
        // Shouldn't reach here
        result.put("success", false);
        result.put("message", "Unknown authentication error");
        return result;
    }
    
    private boolean isValidCredentials(String email, String password) {
        // Simple validation - in real app this would check with Firebase
        return email.contains("@") && password.length() >= 6;
    }
    
    private String verifyGoogleToken(String token) {
        // Simple simulation - in real app this would verify with Google servers
        if (token.startsWith("google_")) {
            return token.substring(7) + "@gmail.com";
        }
        return null;
    }
}
