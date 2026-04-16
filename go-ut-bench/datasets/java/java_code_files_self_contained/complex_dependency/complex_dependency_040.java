// Converted Java method
import java.util.Map;
import java.util.HashMap;
import java.util.List;
import java.util.ArrayList;

class HttpRequestSimulator {
    /**
     * Simulates HTTP requests with enhanced features including:
     * - Request validation
     * - Response simulation
     * - Error handling
     * - Request timing
     * 
     * @param url The endpoint URL
     * @param method HTTP method (GET/POST)
     * @param params Request parameters
     * @param headers Request headers
     * @return Map containing response data, status, and timing information
     * @throws IllegalArgumentException for invalid inputs
     */
    public Map<String, Object> simulateHttpRequest(
            String url, 
            String method, 
            Map<String, String> params, 
            Map<String, String> headers) {
        
        // Validate inputs
        if (url == null || url.trim().isEmpty()) {
            throw new IllegalArgumentException("URL cannot be null or empty");
        }
        
        if (!method.equalsIgnoreCase("GET") && !method.equalsIgnoreCase("POST")) {
            throw new IllegalArgumentException("Unsupported HTTP method: " + method);
        }
        
        // Simulate request processing time (50-200ms)
        long startTime = System.currentTimeMillis();
        try {
            Thread.sleep(50 + (long)(Math.random() * 150));
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        
        // Create response map
        Map<String, Object> response = new HashMap<>();
        
        // Simulate different responses based on inputs
        if (url.contains("error")) {
            response.put("status", 500);
            response.put("body", "Simulated server error");
        } else if (method.equalsIgnoreCase("GET")) {
            response.put("status", 200);
            response.put("body", "GET response for " + url + " with params: " + params);
        } else {
            response.put("status", 201);
            response.put("body", "POST response for " + url + " with params: " + params);
        }
        
        // Add timing information
        response.put("requestTime", System.currentTimeMillis() - startTime);
        
        // Add processed headers if any
        if (headers != null && !headers.isEmpty()) {
            Map<String, String> processedHeaders = new HashMap<>();
            headers.forEach((k, v) -> processedHeaders.put(k.toLowerCase(), v));
            response.put("headers", processedHeaders);
        }
        
        return response;
    }
    
    /**
     * Processes multiple requests in sequence and returns consolidated results
     * 
     * @param requests List of request maps (each containing url, method, params, headers)
     * @return List of response maps for each request
     */
    public List<Map<String, Object>> processBatchRequests(List<Map<String, Object>> requests) {
        List<Map<String, Object>> responses = new ArrayList<>();
        
        for (Map<String, Object> request : requests) {
            try {
                String url = (String) request.get("url");
                String method = (String) request.get("method");
                @SuppressWarnings("unchecked")
                Map<String, String> params = (Map<String, String>) request.get("params");
                @SuppressWarnings("unchecked")
                Map<String, String> headers = (Map<String, String>) request.get("headers");
                
                responses.add(simulateHttpRequest(url, method, params, headers));
            } catch (Exception e) {
                Map<String, Object> errorResponse = new HashMap<>();
                errorResponse.put("error", e.getMessage());
                errorResponse.put("status", 400);
                responses.add(errorResponse);
            }
        }
        
        return responses;
    }
}
