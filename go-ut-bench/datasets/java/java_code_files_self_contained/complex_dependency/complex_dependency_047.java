// Converted Java method
import java.util.HashMap;
import java.util.Map;

class EnhancedResponseHandler {
    /**
     * Enhanced response handler that can manage multiple types of responses
     * including success, error, and custom responses with additional metadata.
     */
    private Map<String, Object> responseData;
    private int statusCode;
    private String statusMessage;
    
    public EnhancedResponseHandler() {
        this.responseData = new HashMap<>();
        this.statusCode = 200; // Default to success
        this.statusMessage = "Success";
    }
    
    /**
     * Creates a success response with optional data payload
     * @param data The data to include in the response (can be null)
     * @return Map containing the full response structure
     */
    public Map<String, Object> createSuccessResponse(Map<String, Object> data) {
        Map<String, Object> response = new HashMap<>();
        response.put("status", "success");
        response.put("code", 200);
        response.put("message", "Operation completed successfully");
        if (data != null) {
            response.put("data", data);
        }
        return response;
    }
    
    /**
     * Creates an error response with specific error code and message
     * @param errorCode The error code
     * @param errorMessage The error message
     * @return Map containing the error response structure
     */
    public Map<String, Object> createErrorResponse(int errorCode, String errorMessage) {
        Map<String, Object> response = new HashMap<>();
        response.put("status", "error");
        response.put("code", errorCode);
        response.put("message", errorMessage);
        return response;
    }
    
    /**
     * Creates a custom response with additional metadata
     * @param status The custom status (e.g., "warning", "partial_success")
     * @param code The status code
     * @param message The status message
     * @param data The data payload (can be null)
     * @param metadata Additional metadata (can be null)
     * @return Map containing the full custom response structure
     */
    public Map<String, Object> createCustomResponse(String status, int code, 
            String message, Map<String, Object> data, Map<String, Object> metadata) {
        Map<String, Object> response = new HashMap<>();
        response.put("status", status);
        response.put("code", code);
        response.put("message", message);
        if (data != null) {
            response.put("data", data);
        }
        if (metadata != null) {
            response.put("metadata", metadata);
        }
        return response;
    }
    
    /**
     * Adds a red point indicator to an existing response
     * @param response The existing response map
     * @param showRedPoint Whether to show the red point (true/false)
     * @return The modified response with red point indicator
     */
    public Map<String, Object> addRedPointIndicator(Map<String, Object> response, boolean showRedPoint) {
        response.put("redPoint", showRedPoint ? 1 : 0);
        return response;
    }
}
