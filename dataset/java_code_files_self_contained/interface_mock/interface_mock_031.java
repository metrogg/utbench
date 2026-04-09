// Converted Java method
import java.util.List;
import java.util.Map;
import java.util.HashMap;

class TableInfoAnalyzer {
    /**
     * Analyzes table information from a database and provides comprehensive metadata.
     * 
     * @param dbUrl Database connection URL
     * @param username Database username
     * @param password Database password
     * @param tableName Table name to analyze
     * @return Map containing table metadata including:
     *         - "columns": List of column names
     *         - "primaryKeys": List of primary key columns
     *         - "foreignKeys": Map of foreign key relationships
     *         - "columnStats": Map of column statistics (type, size, nullable)
     * @throws IllegalArgumentException if connection parameters are invalid
     */
    public Map<String, Object> analyzeTableMetadata(
            String dbUrl, String username, String password, String tableName) {
        
        // Validate inputs
        if (dbUrl == null || dbUrl.isEmpty() || 
            username == null || username.isEmpty() || 
            tableName == null || tableName.isEmpty()) {
            throw new IllegalArgumentException("Invalid connection parameters");
        }

        // Simulate database connection and metadata retrieval
        Map<String, Object> result = new HashMap<>();
        
        // Simulate column information
        List<String> columns = List.of("id", "name", "email", "created_at", "updated_at");
        result.put("columns", columns);
        
        // Simulate primary keys
        List<String> primaryKeys = List.of("id");
        result.put("primaryKeys", primaryKeys);
        
        // Simulate foreign keys
        Map<String, String> foreignKeys = new HashMap<>();
        foreignKeys.put("user_id", "users(id)");
        result.put("foreignKeys", foreignKeys);
        
        // Simulate column statistics
        Map<String, Map<String, String>> columnStats = new HashMap<>();
        columnStats.put("id", Map.of("type", "INT", "size", "11", "nullable", "NO"));
        columnStats.put("name", Map.of("type", "VARCHAR", "size", "255", "nullable", "NO"));
        columnStats.put("email", Map.of("type", "VARCHAR", "size", "255", "nullable", "YES"));
        columnStats.put("created_at", Map.of("type", "TIMESTAMP", "size", "0", "nullable", "NO"));
        columnStats.put("updated_at", Map.of("type", "TIMESTAMP", "size", "0", "nullable", "YES"));
        result.put("columnStats", columnStats);
        
        return result;
    }
}
