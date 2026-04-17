// Converted Java method
import java.util.*;
import java.util.function.Supplier;

/**
 * A generic configuration manager that simulates Spring-like dependency injection
 * and configuration loading without requiring the Spring framework.
 */
class ConfigManager {
    private final Map<String, Object> beans = new HashMap<>();
    private final Properties properties = new Properties();
    private final List<String> configFiles = new ArrayList<>();

    /**
     * Registers a bean in the configuration manager
     * @param name The bean name
     * @param supplier A supplier that provides the bean instance
     */
    public <T> void registerBean(String name, Supplier<T> supplier) {
        beans.put(name, supplier.get());
    }

    /**
     * Loads properties from a file path (simulated)
     * @param path The property file path
     */
    public void loadProperties(String path) {
        configFiles.add(path);
        // Simulate loading properties (in real implementation would read from file)
        if (path.contains("config.properties")) {
            properties.setProperty("url", "http://example.com");
            properties.setProperty("timeout", "5000");
        }
    }

    /**
     * Gets a property value
     * @param key The property key
     * @return The property value or null if not found
     */
    public String getProperty(String key) {
        return properties.getProperty(key);
    }

    /**
     * Gets a bean by its name
     * @param name The bean name
     * @return The bean instance or null if not found
     */
    @SuppressWarnings("unchecked")
    public <T> T getBean(String name) {
        return (T) beans.get(name);
    }

    /**
     * Simulates importing XML configuration (just tracks imported files)
     * @param path The XML file path
     */
    public void importXmlConfig(String path) {
        configFiles.add(path);
    }

    /**
     * Gets all loaded configuration files
     * @return List of configuration file paths
     */
    public List<String> getConfigFiles() {
        return Collections.unmodifiableList(configFiles);
    }
}
