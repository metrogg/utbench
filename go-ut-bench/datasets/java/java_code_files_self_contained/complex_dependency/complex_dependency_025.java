// Converted Java method
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.HashSet;
import java.util.Collections;

class EnhancedOptions {
    private Map<String, Object> options = new HashMap<>();
    private Set<String> requiredOptions = new HashSet<>();
    private Map<String, Object> defaultOptions = new HashMap<>();
    private Map<String, List<String>> optionDependencies = new HashMap<>();

    public EnhancedOptions() {
        // Initialize with some default options
        defaultOptions.put("header_footer", true);
        defaultOptions.put("compact", false);
        defaultOptions.put("safe", 0);
        
        // Set some required options
        requiredOptions.add("backend");
        requiredOptions.add("doctype");
        
        // Set option dependencies
        optionDependencies.put("template_engine", Arrays.asList("template_dirs"));
    }

    /**
     * Sets an option with validation
     * @param optionName Name of the option to set
     * @param optionValue Value to set
     * @throws IllegalArgumentException if option is invalid or dependencies aren't met
     */
    public void setOption(String optionName, Object optionValue) {
        validateOption(optionName, optionValue);
        options.put(optionName, optionValue);
    }

    /**
     * Validates an option before setting it
     */
    private void validateOption(String optionName, Object optionValue) {
        // Check for null values
        if (optionValue == null) {
            throw new IllegalArgumentException("Option value cannot be null for: " + optionName);
        }

        // Check dependencies
        if (optionDependencies.containsKey(optionName)) {
            for (String dependency : optionDependencies.get(optionName)) {
                if (!options.containsKey(dependency)) {
                    throw new IllegalArgumentException("Option '" + optionName + 
                            "' requires '" + dependency + "' to be set first");
                }
            }
        }

        // Type-specific validation
        switch (optionName) {
            case "safe":
                if (!(optionValue instanceof Integer) || (Integer)optionValue < 0 || (Integer)optionValue > 3) {
                    throw new IllegalArgumentException("Safe level must be between 0 and 3");
                }
                break;
            case "template_dirs":
                if (!(optionValue instanceof List) && !(optionValue instanceof String[])) {
                    throw new IllegalArgumentException("template_dirs must be a List or String array");
                }
                break;
            case "backend":
                if (!(optionValue instanceof String)) {
                    throw new IllegalArgumentException("backend must be a string");
                }
                break;
        }
    }

    /**
     * Validates all required options are set
     * @throws IllegalStateException if any required options are missing
     */
    public void validateAllOptions() {
        for (String required : requiredOptions) {
            if (!options.containsKey(required)) {
                throw new IllegalStateException("Required option '" + required + "' is missing");
            }
        }
    }

    /**
     * Gets an option value with fallback to default
     */
    public Object getOption(String optionName) {
        if (options.containsKey(optionName)) {
            return options.get(optionName);
        }
        return defaultOptions.get(optionName);
    }

    /**
     * Gets all options with defaults filled in
     */
    public Map<String, Object> getAllOptions() {
        Map<String, Object> allOptions = new HashMap<>(defaultOptions);
        allOptions.putAll(options);
        return Collections.unmodifiableMap(allOptions);
    }

    /**
     * Checks if an option is set (either explicitly or by default)
     */
    public boolean isSet(String optionName) {
        return options.containsKey(optionName) || defaultOptions.containsKey(optionName);
    }

    /**
     * Merges another Options object into this one
     */
    public void merge(EnhancedOptions other) {
        this.options.putAll(other.options);
    }
}
