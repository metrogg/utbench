// Converted Java method
import java.util.HashMap;
import java.util.Map;
import java.util.regex.PatternSyntaxException;

class RegexHelper {

    /**
     * Generates regex patterns based on the requested pattern type and optional parameters.
     * Supports various regex constructs including character classes, boundaries, and lookarounds.
     *
     * @param patternType The type of regex pattern to generate (e.g., "characterClass", "boundary")
     * @param parameters Optional parameters for the pattern (e.g., character ranges for character classes)
     * @return A map containing the generated pattern and additional metadata
     * @throws IllegalArgumentException if invalid pattern type or parameters are provided
     * @throws PatternSyntaxException if generated pattern would be invalid
     */
    public Map<String, String> generateRegexPattern(String patternType, Map<String, String> parameters) {
        Map<String, String> result = new HashMap<>();
        String pattern;
        
        switch (patternType.toLowerCase()) {
            case "characterclass":
                pattern = generateCharacterClass(parameters);
                break;
                
            case "boundary":
                pattern = generateBoundary(parameters);
                break;
                
            case "lookaround":
                pattern = generateLookAround(parameters);
                break;
                
            case "specialchar":
                pattern = generateSpecialChar(parameters);
                break;
                
            case "literal":
                pattern = generateLiteral(parameters);
                break;
                
            default:
                throw new IllegalArgumentException("Unknown pattern type: " + patternType);
        }
        
        // Validate the generated pattern
        try {
            java.util.regex.Pattern.compile(pattern);
        } catch (PatternSyntaxException e) {
            throw new PatternSyntaxException("Invalid generated pattern", pattern, -1);
        }
        
        result.put("pattern", pattern);
        result.put("type", patternType);
        result.put("description", getPatternDescription(patternType));
        
        return result;
    }
    
    private String generateCharacterClass(Map<String, String> params) {
        String charClass = params.getOrDefault("class", "any");
        
        switch (charClass.toLowerCase()) {
            case "lowercase": return "[a-z]";
            case "uppercase": return "[A-Z]";
            case "digit": return "\\d";
            case "non-digit": return "\\D";
            case "whitespace": return "\\s";
            case "non-whitespace": return "\\S";
            case "word": return "\\w";
            case "non-word": return "\\W";
            case "custom":
                String chars = params.get("chars");
                if (chars == null || chars.isEmpty()) {
                    throw new IllegalArgumentException("Custom character class requires 'chars' parameter");
                }
                return "[" + escapeRegexChars(chars) + "]";
            case "any": 
            default: return ".";
        }
    }
    
    private String generateBoundary(Map<String, String> params) {
        String boundary = params.getOrDefault("type", "word");
        
        switch (boundary.toLowerCase()) {
            case "start-line": return "^";
            case "end-line": return "$";
            case "word": return "\\b";
            case "non-word": return "\\B";
            case "start-input": return "\\A";
            case "end-input": return "\\z";
            case "end-prev-match": return "\\G";
            case "end-input-final-term": return "\\Z";
            default: 
                throw new IllegalArgumentException("Unknown boundary type: " + boundary);
        }
    }
    
    private String generateLookAround(Map<String, String> params) {
        String type = params.getOrDefault("type", "ahead-positive");
        String pattern = params.get("pattern");
        
        if (pattern == null || pattern.isEmpty()) {
            throw new IllegalArgumentException("Lookaround requires 'pattern' parameter");
        }
        
        switch (type.toLowerCase()) {
            case "ahead-positive": return "(?=" + pattern + ")";
            case "ahead-negative": return "(?!" + pattern + ")";
            case "behind-positive": return "(?<=" + pattern + ")";
            case "behind-negative": return "(?<!" + pattern + ")";
            default: 
                throw new IllegalArgumentException("Unknown lookaround type: " + type);
        }
    }
    
    private String generateSpecialChar(Map<String, String> params) {
        String charType = params.getOrDefault("type", "tab");
        
        switch (charType.toLowerCase()) {
            case "tab": return "\\t";
            case "newline": return "\\n";
            case "return": return "\\r";
            case "windows-newline": return "\\r\\n";
            case "bell": return "\\a";
            case "escape": return "\\e";
            case "form-feed": return "\\f";
            case "vertical-tab": return "\\v";
            default: 
                throw new IllegalArgumentException("Unknown special character type: " + charType);
        }
    }
    
    private String generateLiteral(Map<String, String> params) {
        String text = params.get("text");
        if (text == null || text.isEmpty()) {
            throw new IllegalArgumentException("Literal requires 'text' parameter");
        }
        return escapeRegexChars(text);
    }
    
    private String escapeRegexChars(String input) {
        StringBuilder sb = new StringBuilder();
        for (char c : input.toCharArray()) {
            if ("[]\\(){}.*+?^$|".indexOf(c) != -1) {
                sb.append('\\');
            }
            sb.append(c);
        }
        return sb.toString();
    }
    
    private String getPatternDescription(String patternType) {
        switch (patternType.toLowerCase()) {
            case "characterclass": return "Character class pattern";
            case "boundary": return "Boundary matcher";
            case "lookaround": return "Lookaround assertion";
            case "specialchar": return "Special character";
            case "literal": return "Literal text";
            default: return "Unknown pattern type";
        }
    }
}
