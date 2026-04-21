// Converted Java method
import java.util.logging.Logger;
import java.io.StringWriter;
import java.io.PrintWriter;
import java.util.HashMap;
import java.util.Map;

class EnhancedExceptionHandler {
    private static final Logger logger = Logger.getLogger("EnhancedExceptionHandler");
    
    /**
     * Handles exceptions with enhanced logging and context information.
     * Creates a detailed error report including stack trace, timestamp,
     * and custom context data.
     * 
     * @param exception The exception to handle
     * @param context Additional context information about the error
     * @return A map containing detailed error information
     */
    public static Map<String, String> handleException(Exception exception, Map<String, String> context) {
        Map<String, String> errorReport = new HashMap<>();
        
        // Capture stack trace
        StringWriter stackTraceWriter = new StringWriter();
        exception.printStackTrace(new PrintWriter(stackTraceWriter));
        String stackTrace = stackTraceWriter.toString();
        
        // Add basic error info
        errorReport.put("exceptionType", exception.getClass().getName());
        errorReport.put("message", exception.getMessage());
        errorReport.put("stackTrace", stackTrace);
        
        // Add timestamp
        errorReport.put("timestamp", String.valueOf(System.currentTimeMillis()));
        
        // Add context if provided
        if (context != null) {
            errorReport.putAll(context);
        }
        
        // Log the error
        StringBuilder logMessage = new StringBuilder();
        logMessage.append("Exception Occurred:\n");
        logMessage.append("Type: ").append(errorReport.get("exceptionType")).append("\n");
        logMessage.append("Message: ").append(errorReport.get("message")).append("\n");
        logMessage.append("Context: ").append(context != null ? context.toString() : "None");
        
        logger.severe(logMessage.toString());
        
        return errorReport;
    }
    
    /**
     * Simulates a business operation that might fail with different exception types
     * @param operationType Type of operation to simulate
     * @throws Exception Different types of exceptions based on input
     */
    public static void performBusinessOperation(String operationType) throws Exception {
        if ("divideByZero".equals(operationType)) {
            int result = 10 / 0;
        } else if ("nullPointer".equals(operationType)) {
            String str = null;
            str.length();
        } else if ("arrayIndex".equals(operationType)) {
            int[] arr = new int[5];
            int val = arr[10];
        } else if ("custom".equals(operationType)) {
            throw new CustomBusinessException("Business rule violation");
        } else {
            throw new IllegalArgumentException("Invalid operation type");
        }
    }
}

class CustomBusinessException extends Exception {
    public CustomBusinessException(String message) {
        super(message);
    }
}
