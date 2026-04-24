// Converted Java method
import java.util.HashMap;
import java.util.Map;

class EnhancedApiError {
    private final ErrorResponse response;
    private final ErrorCode errorCode;
    private final Map<String, String> additionalDetails;

    /**
     * Enhanced API error handler with additional features
     * @param response The original error response
     * @param errorCode The error code enum
     * @param additionalDetails Map of additional error details
     */
    public EnhancedApiError(ErrorResponse response, ErrorCode errorCode, Map<String, String> additionalDetails) {
        this.response = response;
        this.errorCode = errorCode;
        this.additionalDetails = new HashMap<>(additionalDetails);
    }

    /**
     * Gets the error code with severity level
     * @return ErrorCode enum value
     */
    public ErrorCode getErrorCode() {
        return errorCode;
    }

    /**
     * Gets the formatted error message with additional details
     * @return Complete error message string
     */
    public String getEnhancedMessage() {
        StringBuilder message = new StringBuilder(getBaseMessage());
        
        if (!additionalDetails.isEmpty()) {
            message.append("\nAdditional Details:");
            for (Map.Entry<String, String> entry : additionalDetails.entrySet()) {
                message.append("\n- ").append(entry.getKey()).append(": ").append(entry.getValue());
            }
        }
        
        return message.toString();
    }

    /**
     * Gets the base error message without additional details
     * @return Base error message string
     */
    public String getBaseMessage() {
        return ErrorMessageResolver.resolve(errorCode);
    }

    /**
     * Gets the original response object
     * @return Original ErrorResponse object
     */
    public ErrorResponse getOriginalResponse() {
        return response;
    }

    /**
     * Checks if this error is recoverable
     * @return true if error is recoverable, false otherwise
     */
    public boolean isRecoverable() {
        return errorCode.getSeverity() != ErrorSeverity.CRITICAL;
    }

    /**
     * Gets all additional details as a map
     * @return Map of additional error details
     */
    public Map<String, String> getAdditionalDetails() {
        return new HashMap<>(additionalDetails);
    }
}

// Supporting enums and classes
enum ErrorCode {
    NON_KEY_ACCOUNT_BALANCE_ERROR(ErrorSeverity.HIGH, "Account balance access denied"),
    INVALID_CREDENTIALS(ErrorSeverity.CRITICAL, "Invalid credentials"),
    SESSION_EXPIRED(ErrorSeverity.MEDIUM, "Session expired"),
    RATE_LIMIT_EXCEEDED(ErrorSeverity.MEDIUM, "Rate limit exceeded");

    private final ErrorSeverity severity;
    private final String defaultMessage;

    ErrorCode(ErrorSeverity severity, String defaultMessage) {
        this.severity = severity;
        this.defaultMessage = defaultMessage;
    }

    public ErrorSeverity getSeverity() {
        return severity;
    }

    public String getDefaultMessage() {
        return defaultMessage;
    }
}

enum ErrorSeverity {
    LOW, MEDIUM, HIGH, CRITICAL
}

class ErrorMessageResolver {
    public static String resolve(ErrorCode errorCode) {
        return errorCode.getDefaultMessage();
    }
}

class ErrorResponse {
    // Simplified error response class
    private final int statusCode;
    private final String rawResponse;

    public ErrorResponse(int statusCode, String rawResponse) {
        this.statusCode = statusCode;
        this.rawResponse = rawResponse;
    }

    public int getStatusCode() {
        return statusCode;
    }

    public String getRawResponse() {
        return rawResponse;
    }
}
