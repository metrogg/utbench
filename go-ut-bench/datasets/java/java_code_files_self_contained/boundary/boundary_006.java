// Converted Java method
import java.util.List;
import java.util.ArrayList;
import java.util.regex.Pattern;
import java.util.regex.Matcher;

class DocumentationValidator {
    private static final Pattern INVALID_CHAR_PATTERN = Pattern.compile("[\\\\]");
    private static final Pattern PLACEHOLDER_PATTERN = Pattern.compile("\\$\\{(.+?)\\}");
    private static final List<String> ALLOWED_PLACEHOLDERS = List.of("method", "class", "param", "return");
    
    /**
     * Validates test documentation for a single item (method or class).
     * 
     * @param doc The documentation string to validate
     * @param qualifiedName The fully qualified name of the item being validated
     * @throws IllegalArgumentException if validation fails
     */
    public static void validateDocumentation(String doc, String qualifiedName) {
        if (doc == null) {
            return; // Null docs are considered valid
        }
        
        // Check for invalid characters
        if (INVALID_CHAR_PATTERN.matcher(doc).find()) {
            throw new IllegalArgumentException(
                String.format("Documentation contains invalid characters: %s", qualifiedName));
        }
        
        // Check for invalid placeholders
        Matcher matcher = PLACEHOLDER_PATTERN.matcher(doc);
        while (matcher.find()) {
            String placeholder = matcher.group(1);
            if (!ALLOWED_PLACEHOLDERS.contains(placeholder)) {
                throw new IllegalArgumentException(
                    String.format("Invalid placeholder '${%s}' in documentation for: %s", 
                                 placeholder, qualifiedName));
            }
        }
        
        // Additional validation rules
        if (doc.trim().isEmpty()) {
            throw new IllegalArgumentException(
                String.format("Empty documentation string for: %s", qualifiedName));
        }
        
        if (doc.length() > 500) {
            throw new IllegalArgumentException(
                String.format("Documentation exceeds 500 characters for: %s", qualifiedName));
        }
    }
    
    /**
     * Validates all documentation in a test class.
     * 
     * @param className The class name
     * @param classDoc The class documentation
     * @param methodDocs List of method documentation strings
     * @throws IllegalArgumentException if any validation fails
     */
    public static void validateClassDocumentation(String className, String classDoc, 
                                                List<String> methodDocs) {
        validateDocumentation(classDoc, className);
        
        for (int i = 0; i < methodDocs.size(); i++) {
            String methodName = String.format("%s.method%d", className, i+1);
            validateDocumentation(methodDocs.get(i), methodName);
        }
    }
}
