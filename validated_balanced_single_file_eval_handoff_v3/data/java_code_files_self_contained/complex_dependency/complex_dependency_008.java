// Converted Java method
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class ApiDocumentationGenerator {
    private String title;
    private String description;
    private String version;
    private String basePackage;
    private List<Map<String, String>> securitySchemes;
    private Map<Class<?>, Class<?>> modelSubstitutes;
    private boolean useDefaultResponseMessages;

    public ApiDocumentationGenerator(String title, String description, String version) {
        this.title = title;
        this.description = description;
        this.version = version;
        this.securitySchemes = new ArrayList<>();
        this.modelSubstitutes = new HashMap<>();
        this.useDefaultResponseMessages = true;
    }

    /**
     * Sets the base package for API scanning
     * @param basePackage the base package to scan for API endpoints
     * @return this generator for method chaining
     */
    public ApiDocumentationGenerator withBasePackage(String basePackage) {
        this.basePackage = basePackage;
        return this;
    }

    /**
     * Adds a security scheme to the API documentation
     * @param name the name of the security scheme
     * @param keyName the key name
     * @param location the location (header, query, etc.)
     * @return this generator for method chaining
     */
    public ApiDocumentationGenerator addSecurityScheme(String name, String keyName, String location) {
        Map<String, String> scheme = new HashMap<>();
        scheme.put("name", name);
        scheme.put("keyName", keyName);
        scheme.put("location", location);
        securitySchemes.add(scheme);
        return this;
    }

    /**
     * Adds a model substitute for documentation
     * @param original the original class
     * @param substitute the substitute class
     * @return this generator for method chaining
     */
    public ApiDocumentationGenerator addModelSubstitute(Class<?> original, Class<?> substitute) {
        modelSubstitutes.put(original, substitute);
        return this;
    }

    /**
     * Sets whether to use default response messages
     * @param useDefault whether to use default messages
     * @return this generator for method chaining
     */
    public ApiDocumentationGenerator withDefaultResponseMessages(boolean useDefault) {
        this.useDefaultResponseMessages = useDefault;
        return this;
    }

    /**
     * Generates the API documentation configuration
     * @return a map representing the documentation configuration
     */
    public Map<String, Object> generateConfiguration() {
        Map<String, Object> config = new HashMap<>();
        config.put("title", title);
        config.put("description", description);
        config.put("version", version);
        config.put("basePackage", basePackage);
        config.put("securitySchemes", securitySchemes);
        config.put("modelSubstitutes", modelSubstitutes);
        config.put("useDefaultResponseMessages", useDefaultResponseMessages);
        return config;
    }

    /**
     * Validates the configuration before generation
     * @throws IllegalStateException if required fields are missing
     */
    public void validateConfiguration() {
        if (title == null || title.isEmpty()) {
            throw new IllegalStateException("Title is required");
        }
        if (basePackage == null || basePackage.isEmpty()) {
            throw new IllegalStateException("Base package is required");
        }
    }
}
