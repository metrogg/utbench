// Converted Java method
import java.util.*;
import java.util.concurrent.TimeUnit;

class ConfigurationManager {
    private final Map<String, Object> configMap;
    private final Set<String> requiredKeys;
    private final Map<String, List<Validator>> keyValidators;

    /**
     * Initializes configuration manager with required keys and validators
     * @param requiredKeys Set of configuration keys that must be present
     * @param keyValidators Map of validators for each configuration key
     */
    public ConfigurationManager(Set<String> requiredKeys, 
                              Map<String, List<Validator>> keyValidators) {
        this.configMap = new HashMap<>();
        this.requiredKeys = new HashSet<>(requiredKeys);
        this.keyValidators = new HashMap<>(keyValidators);
    }

    /**
     * Loads configuration from properties map after validation
     * @param properties Map of configuration properties
     * @throws ConfigurationException if validation fails
     */
    public void loadConfiguration(Map<String, Object> properties) throws ConfigurationException {
        validateConfiguration(properties);
        configMap.putAll(properties);
    }

    /**
     * Validates configuration against required keys and validators
     * @param properties Map of configuration properties to validate
     * @throws ConfigurationException if validation fails
     */
    private void validateConfiguration(Map<String, Object> properties) throws ConfigurationException {
        // Check for missing required keys
        Set<String> missingKeys = new HashSet<>(requiredKeys);
        missingKeys.removeAll(properties.keySet());
        
        if (!missingKeys.isEmpty()) {
            throw new ConfigurationException("Missing required configuration keys: " + missingKeys);
        }

        // Validate each property that has validators
        for (Map.Entry<String, Object> entry : properties.entrySet()) {
            String key = entry.getKey();
            Object value = entry.getValue();
            
            if (keyValidators.containsKey(key)) {
                for (Validator validator : keyValidators.get(key)) {
                    if (!validator.validate(value)) {
                        throw new ConfigurationException(
                            String.format("Validation failed for key '%s': %s", 
                                        key, validator.getErrorMessage()));
                    }
                }
            }
        }
    }

    /**
     * Gets configuration value with type conversion
     * @param key Configuration key
     * @param type Expected return type class
     * @return Configuration value converted to requested type
     * @throws ConfigurationException if conversion fails or key doesn't exist
     */
    @SuppressWarnings("unchecked")
    public <T> T getConfigValue(String key, Class<T> type) throws ConfigurationException {
        if (!configMap.containsKey(key)) {
            throw new ConfigurationException("Configuration key not found: " + key);
        }

        Object value = configMap.get(key);
        try {
            if (type == Integer.class) {
                return (T) Integer.valueOf(value.toString());
            } else if (type == Long.class) {
                return (T) Long.valueOf(value.toString());
            } else if (type == Double.class) {
                return (T) Double.valueOf(value.toString());
            } else if (type == Boolean.class) {
                return (T) Boolean.valueOf(value.toString());
            } else if (type == String.class) {
                return (T) value.toString();
            } else if (type == TimeUnit.class) {
                return (T) TimeUnit.valueOf(value.toString().toUpperCase());
            } else {
                return type.cast(value);
            }
        } catch (Exception e) {
            throw new ConfigurationException(
                String.format("Failed to convert value '%s' for key '%s' to type %s",
                            value, key, type.getSimpleName()));
        }
    }

    public interface Validator {
        boolean validate(Object value);
        String getErrorMessage();
    }

    public static class ConfigurationException extends Exception {
        public ConfigurationException(String message) {
            super(message);
        }
    }
}
