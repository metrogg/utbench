// Converted Java method
import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

class FileUploadUtil {
    
    /**
     * Extracts the filename from a content-disposition header string.
     * Handles various formats of content-disposition headers and multiple quotes.
     * 
     * @param contentDisposition The content-disposition header string
     * @return The extracted filename or empty string if not found
     */
    public static String getFilename(String contentDisposition) {
        if (contentDisposition == null || contentDisposition.isEmpty()) {
            return "";
        }
        
        String[] parts = contentDisposition.split(";");
        for (String part : parts) {
            part = part.trim();
            if (part.startsWith("filename") || part.startsWith("filename*")) {
                String[] keyValue = part.split("=", 2);
                if (keyValue.length == 2) {
                    String filename = keyValue[1].trim();
                    // Remove all types of quotes
                    filename = filename.replaceAll("^[\"']|[\"']$", "");
                    return filename;
                }
            }
        }
        return "";
    }
    
    /**
     * Extracts the file extension from a filename.
     * Handles edge cases like no extension, multiple dots, and null input.
     * 
     * @param filename The filename to process
     * @return The file extension in lowercase, or empty string if none found
     */
    public static String getExtension(String filename) {
        if (filename == null || filename.isEmpty()) {
            return "";
        }
        
        int lastDotIndex = filename.lastIndexOf('.');
        if (lastDotIndex == -1 || lastDotIndex == filename.length() - 1) {
            return "";
        }
        
        return filename.substring(lastDotIndex + 1).toLowerCase();
    }
    
    /**
     * Parses a complete content-disposition header into its components.
     * Extracts type, name, filename, and all parameters.
     * 
     * @param contentDisposition The content-disposition header string
     * @return Map containing all parsed components
     */
    public static Map<String, String> parseContentDisposition(String contentDisposition) {
        Map<String, String> result = new HashMap<>();
        if (contentDisposition == null || contentDisposition.isEmpty()) {
            return result;
        }
        
        String[] parts = contentDisposition.split(";");
        if (parts.length > 0) {
            result.put("type", parts[0].trim());
        }
        
        for (int i = 1; i < parts.length; i++) {
            String[] keyValue = parts[i].trim().split("=", 2);
            if (keyValue.length == 2) {
                String value = keyValue[1].trim().replaceAll("^[\"']|[\"']$", "");
                result.put(keyValue[0].trim().toLowerCase(), value);
            }
        }
        
        return result;
    }
}
