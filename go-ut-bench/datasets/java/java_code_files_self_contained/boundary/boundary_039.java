// Converted Java method
import java.util.HashMap;
import java.util.Map;

class EnhancedUssmsCardErrorHandler {
    private static final Map<Integer, UssmsCardErrType> codeToErrorMap = new HashMap<>();
    
    static {
        for (UssmsCardErrType error : UssmsCardErrType.values()) {
            codeToErrorMap.put(error.getCode(), error);
        }
    }

    /**
     * Enhanced error type enum with additional functionality
     */
    public enum UssmsCardErrType {
        SUCCESS(0, "Activation successful"),
        FORMAT_ERR(1, "Command format error, please re-enter."),
        CARDINFO_ERR(2, "Card number or password incorrect, please verify and try again."),
        CHARGE_ERR(3, "Network error, please try again later."),
        EXPIRED_CARD(4, "This card has expired."),
        ALREADY_USED(5, "This card has already been used."),
        SYSTEM_ERROR(6, "System maintenance in progress, please try again later.");

        private final int code;
        private final String message;

        UssmsCardErrType(int code, String message) {
            this.code = code;
            this.message = message;
        }

        public int getCode() {
            return code;
        }

        public String getMessage() {
            return message;
        }

        /**
         * Returns the error type corresponding to the given code
         * @param code The error code to look up
         * @return The corresponding UssmsCardErrType
         * @throws IllegalArgumentException if the code is not found
         */
        public static UssmsCardErrType fromCode(int code) {
            UssmsCardErrType error = codeToErrorMap.get(code);
            if (error == null) {
                throw new IllegalArgumentException("Invalid error code: " + code);
            }
            return error;
        }

        /**
         * Checks if this error type represents a success state
         * @return true if this is SUCCESS, false otherwise
         */
        public boolean isSuccess() {
            return this == SUCCESS;
        }

        /**
         * Checks if this error type represents a temporary failure that might succeed on retry
         * @return true if the error is temporary (CHARGE_ERR or SYSTEM_ERROR), false otherwise
         */
        public boolean isTemporaryError() {
            return this == CHARGE_ERR || this == SYSTEM_ERROR;
        }
    }

    /**
     * Processes an error code and returns the appropriate response message
     * @param errorCode The error code to process
     * @return The formatted response message
     */
    public static String processError(int errorCode) {
        try {
            UssmsCardErrType error = UssmsCardErrType.fromCode(errorCode);
            return String.format("Error %d: %s", error.getCode(), error.getMessage());
        } catch (IllegalArgumentException e) {
            return "Unknown error occurred. Please contact customer support.";
        }
    }
}
