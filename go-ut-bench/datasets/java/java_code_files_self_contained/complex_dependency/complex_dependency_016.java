import java.util.HashMap;
import java.util.Map;

class ServerApprovalHandler {
    private String serverName;
    private String serverIp;
    private String realmId;
    private String imageId;
    private String hardwareProfileId;
    private Map<String, String> approvalStatus;

    public ServerApprovalHandler(String serverName, String serverIp, String realmId, 
                                String imageId, String hardwareProfileId) {
        this.serverName = serverName;
        this.serverIp = serverIp;
        this.realmId = realmId;
        this.imageId = imageId;
        this.hardwareProfileId = hardwareProfileId;
        this.approvalStatus = new HashMap<>();
    }

    /**
     * Validates all server configuration parameters before approval
     * @return Map containing validation results for each parameter
     */
    public Map<String, String> validateConfiguration() {
        approvalStatus.clear();
        
        approvalStatus.put("serverName", validateServerName());
        approvalStatus.put("serverIp", validateIpAddress());
        approvalStatus.put("realmId", validateRealmId());
        approvalStatus.put("imageId", validateImageId());
        approvalStatus.put("hardwareProfileId", validateHardwareProfile());
        
        return approvalStatus;
    }

    /**
     * Processes the server approval with comprehensive validation
     * @return Map containing approval status and validation results
     */
    public Map<String, Object> processApproval() {
        Map<String, Object> result = new HashMap<>();
        
        validateConfiguration();
        
        boolean allValid = approvalStatus.values().stream()
                .allMatch("Valid"::equals);
        
        if (allValid) {
            result.put("approvalStatus", "APPROVED");
            result.put("message", "Server configuration approved successfully");
            result.put("validationResults", approvalStatus);
        } else {
            result.put("approvalStatus", "REJECTED");
            result.put("message", "Server configuration validation failed");
            result.put("validationResults", approvalStatus);
        }
        
        return result;
    }

    private String validateServerName() {
        if (serverName == null || serverName.trim().isEmpty()) {
            return "Invalid: Server name cannot be empty";
        }
        if (serverName.length() > 64) {
            return "Invalid: Server name too long (max 64 chars)";
        }
        return "Valid";
    }

    private String validateIpAddress() {
        if (serverIp == null || !serverIp.matches("^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$")) {
            return "Invalid: Not a valid IPv4 address";
        }
        return "Valid";
    }

    private String validateRealmId() {
        if (realmId == null || !realmId.matches("^[a-zA-Z0-9-]{1,36}$")) {
            return "Invalid: Realm ID must be 1-36 alphanumeric chars with optional hyphens";
        }
        return "Valid";
    }

    private String validateImageId() {
        if (imageId == null || !imageId.matches("^img-[a-zA-Z0-9]{8}$")) {
            return "Invalid: Image ID must follow 'img-xxxxxxxx' format (8 alphanumeric chars)";
        }
        return "Valid";
    }

    private String validateHardwareProfile() {
        if (hardwareProfileId == null || !hardwareProfileId.matches("^hw-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}$")) {
            return "Invalid: Hardware profile must follow 'hw-xxxx-xxxx' format";
        }
        return "Valid";
    }
}
