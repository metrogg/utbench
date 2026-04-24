import java.util.*;
import java.nio.file.Path;
import java.nio.file.Paths;

class PathAnalyzer {
    /**
     * Analyzes a file path and extracts various components with validation.
     * 
     * @param filePath The full file path to analyze
     * @return Map containing:
     *         - "isValid": boolean indicating if path is valid
     *         - "filename": extracted filename
     *         - "extension": file extension (empty if none)
     *         - "parentDir": parent directory path
     *         - "depth": depth of path (number of directories)
     * @throws IllegalArgumentException if input is null or empty
     */
    public static Map<String, Object> analyzePath(String filePath) {
        if (filePath == null || filePath.trim().isEmpty()) {
            throw new IllegalArgumentException("Path cannot be null or empty");
        }

        Map<String, Object> result = new HashMap<>();
        Path path = Paths.get(filePath).normalize();
        
        // Basic validation
        boolean isValid = !filePath.contains(" ") && filePath.length() <= 260;
        result.put("isValid", isValid);
        
        // Extract filename
        String filename = path.getFileName() != null ? path.getFileName().toString() : "";
        result.put("filename", filename);
        
        // Extract extension
        String extension = "";
        int dotIndex = filename.lastIndexOf('.');
        if (dotIndex > 0) {
            extension = filename.substring(dotIndex + 1);
        }
        result.put("extension", extension);
        
        // Parent directory
        Path parent = path.getParent();
        result.put("parentDir", parent != null ? parent.toString() : "");
        
        // Calculate depth
        int depth = 0;
        for (Path p : path) {
            depth++;
        }
        result.put("depth", depth - 1); // Subtract 1 for the filename
        
        return result;
    }
}
