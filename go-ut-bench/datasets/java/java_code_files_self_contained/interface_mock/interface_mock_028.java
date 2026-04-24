import java.util.HashMap;
import java.util.Iterator;
import java.util.Locale;
import java.util.Map;
import java.util.Map.Entry;

/**
 * Enhanced Activity Manager that simulates Android activity navigation with additional features
 */
class EnhancedActivityManager {
    private Map<String, String> activityParams;
    private Locale currentLocale;
    private String lastErrorMessage;

    public EnhancedActivityManager() {
        this.activityParams = new HashMap<>();
        this.currentLocale = Locale.getDefault();
    }

    /**
     * Navigates to a new activity with parameters
     * @param targetActivity The target activity class name
     * @param params Parameters to pass to the new activity
     * @return Navigation success status
     */
    public boolean navigateTo(String targetActivity, Map<String, String> params) {
        if (targetActivity == null || targetActivity.isEmpty()) {
            lastErrorMessage = "Target activity cannot be null or empty";
            return false;
        }

        // Store parameters for the next activity
        if (params != null) {
            Iterator<Entry<String, String>> iter = params.entrySet().iterator();
            while (iter.hasNext()) {
                Map.Entry<String, String> param = iter.next();
                activityParams.put(param.getKey(), param.getValue());
            }
        }

        // Simulate activity navigation success
        return true;
    }

    /**
     * Changes the application locale
     * @param language Language code (e.g., "en")
     * @param country Country code (e.g., "US")
     * @return true if locale was changed successfully
     */
    public boolean changeLocale(String language, String country) {
        if (language == null || country == null) {
            lastErrorMessage = "Language and country cannot be null";
            return false;
        }

        try {
            Locale newLocale = new Locale(language, country);
            Locale.setDefault(newLocale);
            currentLocale = newLocale;
            return true;
        } catch (Exception e) {
            lastErrorMessage = "Failed to change locale: " + e.getMessage();
            return false;
        }
    }

    /**
     * Handles an error condition
     * @param errorMessage The error message
     * @param isBlocking Whether the error should block further execution
     * @return The formatted error response
     */
    public Map<String, Object> handleError(String errorMessage, boolean isBlocking) {
        Map<String, Object> errorResponse = new HashMap<>();
        errorResponse.put("message", errorMessage);
        errorResponse.put("isBlocking", isBlocking);
        errorResponse.put("timestamp", System.currentTimeMillis());
        
        lastErrorMessage = errorMessage;
        
        if (isBlocking) {
            errorResponse.put("action", "terminate");
        } else {
            errorResponse.put("action", "continue");
        }
        
        return errorResponse;
    }

    /**
     * Gets the current locale
     * @return The current locale object
     */
    public Locale getCurrentLocale() {
        return currentLocale;
    }

    /**
     * Gets the last error message
     * @return The last error message or null if none
     */
    public String getLastErrorMessage() {
        return lastErrorMessage;
    }

    /**
     * Clears all activity parameters
     */
    public void clearParameters() {
        activityParams.clear();
    }
}
