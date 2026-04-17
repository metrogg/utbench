import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Enhanced User Settings Manager with additional functionality
 * - Thread-safe operations using ConcurrentHashMap
 * - Support for multiple setting types (String, Integer, Boolean)
 * - Setting validation
 * - Bulk operations
 */
class UserSettingsManager {
    private final Map<String, Map<String, Object>> userSettings = new ConcurrentHashMap<>();
    
    /**
     * Saves a user setting with validation
     * @param key Setting key
     * @param value Setting value (String, Integer, or Boolean)
     * @param username User identifier
     * @throws IllegalArgumentException if value type is not supported
     */
    public void saveUserSetting(String key, Object value, String username) {
        if (!(value instanceof String || value instanceof Integer || value instanceof Boolean)) {
            throw new IllegalArgumentException("Unsupported value type");
        }
        
        userSettings.computeIfAbsent(username, k -> new ConcurrentHashMap<>())
                   .put(key, value);
    }
    
    /**
     * Gets a user setting with type safety
     * @param key Setting key
     * @param username User identifier
     * @param type Expected class type of the value
     * @return The setting value or null if not found
     * @throws ClassCastException if value type doesn't match expected type
     */
    public <T> T getUserSetting(String key, String username, Class<T> type) {
        Map<String, Object> settings = userSettings.get(username);
        if (settings == null) {
            return null;
        }
        return type.cast(settings.get(key));
    }
    
    /**
     * Bulk update user settings
     * @param settings Map of key-value pairs to update
     * @param username User identifier
     */
    public void updateUserSettings(Map<String, Object> settings, String username) {
        Map<String, Object> userMap = userSettings.computeIfAbsent(username, k -> new ConcurrentHashMap<>());
        settings.forEach((key, value) -> {
            if (value instanceof String || value instanceof Integer || value instanceof Boolean) {
                userMap.put(key, value);
            }
        });
    }
    
    /**
     * Gets all settings for a user
     * @param username User identifier
     * @return Map of all settings or empty map if none exist
     */
    public Map<String, Object> getAllUserSettings(String username) {
        return new HashMap<>(userSettings.getOrDefault(username, new ConcurrentHashMap<>()));
    }
    
    /**
     * Deletes a specific setting for a user
     * @param key Setting key to delete
     * @param username User identifier
     * @return The deleted value or null if not found
     */
    public Object deleteUserSetting(String key, String username) {
        Map<String, Object> settings = userSettings.get(username);
        return settings != null ? settings.remove(key) : null;
    }
    
    /**
     * Clears all settings for a user
     * @param username User identifier
     */
    public void clearUserSettings(String username) {
        userSettings.remove(username);
    }
}
