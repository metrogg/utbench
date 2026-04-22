// Converted Java method
import java.util.HashMap;
import java.util.Map;

class TLSExtensionProcessor {

    /**
     * Simulates processing of TLS extensions with enhanced functionality.
     * This includes validation, security checks, and state management.
     *
     * @param extensionType The type of TLS extension being processed
     * @param extensionData Map containing extension data (keys: "maxEarlyDataSize", "connectionEnd", etc.)
     * @param context Current TLS context state
     * @return Map containing processing results and any warnings/errors
     * @throws IllegalArgumentException if invalid parameters are provided
     */
    public Map<String, Object> processExtension(String extensionType, 
                                              Map<String, Object> extensionData,
                                              Map<String, Object> context) {
        if (extensionType == null || extensionType.isEmpty()) {
            throw new IllegalArgumentException("Extension type cannot be null or empty");
        }

        Map<String, Object> result = new HashMap<>();
        result.put("extensionType", extensionType);
        result.put("originalData", extensionData);

        // Validate extension data
        if (extensionData == null) {
            result.put("status", "error");
            result.put("message", "Extension data cannot be null");
            return result;
        }

        // Process based on extension type
        switch (extensionType) {
            case "EARLY_DATA":
                processEarlyDataExtension(extensionData, context, result);
                break;
            case "SERVER_NAME":
                processServerNameExtension(extensionData, context, result);
                break;
            default:
                result.put("status", "warning");
                result.put("message", "Unsupported extension type");
        }

        // Update context if processing was successful
        if (!result.containsKey("status") || !result.get("status").equals("error")) {
            updateTlsContext(extensionType, extensionData, context, result);
        }

        return result;
    }

    private void processEarlyDataExtension(Map<String, Object> data, 
                                         Map<String, Object> context,
                                         Map<String, Object> result) {
        // Validate max early data size
        if (data.containsKey("maxEarlyDataSize")) {
            try {
                int maxSize = (int) data.get("maxEarlyDataSize");
                if (maxSize < 0) {
                    result.put("status", "error");
                    result.put("message", "Max early data size cannot be negative");
                    return;
                }
                result.put("maxEarlyDataSize", maxSize);
            } catch (ClassCastException e) {
                result.put("status", "error");
                result.put("message", "Invalid maxEarlyDataSize format");
                return;
            }
        }

        // Check connection end type
        String connectionEnd = (String) data.getOrDefault("connectionEnd", "CLIENT");
        if (!connectionEnd.equals("CLIENT") && !connectionEnd.equals("SERVER")) {
            result.put("status", "error");
            result.put("message", "Invalid connection end type");
            return;
        }

        result.put("connectionEnd", connectionEnd);
        result.put("status", "success");
    }

    private void processServerNameExtension(Map<String, Object> data,
                                          Map<String, Object> context,
                                          Map<String, Object> result) {
        // Validate server name
        if (!data.containsKey("serverName") || ((String)data.get("serverName")).isEmpty()) {
            result.put("status", "error");
            result.put("message", "Server name cannot be empty");
            return;
        }

        result.put("serverName", data.get("serverName"));
        result.put("status", "success");
    }

    private void updateTlsContext(String extensionType,
                                Map<String, Object> data,
                                Map<String, Object> context,
                                Map<String, Object> result) {
        // Update context with processed extension data
        context.put(extensionType, data);

        // Special handling for EARLY_DATA
        if (extensionType.equals("EARLY_DATA")) {
            if (data.containsKey("maxEarlyDataSize")) {
                context.put("maxEarlyDataSize", data.get("maxEarlyDataSize"));
            } else if (data.get("connectionEnd").equals("SERVER")) {
                context.put("negotiatedExtensions", extensionType);
            }
        }

        result.put("updatedContext", context);
    }
}
