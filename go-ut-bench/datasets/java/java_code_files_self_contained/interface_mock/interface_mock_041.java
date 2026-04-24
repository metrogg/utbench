import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.Set;

class RoleManager {
    private Set<String> validRoles;
    private List<String> allowedRoles;
    private List<String> disallowedRoles;
    private boolean strictMode;

    /**
     * Initializes a RoleManager with configuration for role validation
     * @param validRoles Set of all valid roles in the system
     * @param strictMode If true, will throw exceptions for invalid roles
     */
    public RoleManager(Set<String> validRoles, boolean strictMode) {
        this.validRoles = new HashSet<>(validRoles);
        this.allowedRoles = new ArrayList<>();
        this.disallowedRoles = new ArrayList<>();
        this.strictMode = strictMode;
    }

    /**
     * Sets the allowed roles after validation
     * @param roles Comma-separated string of roles
     * @throws IllegalArgumentException if strictMode is true and invalid roles are found
     */
    public void setAllowedRoles(String roles) {
        List<String> roleList = stringToList(roles);
        if (strictMode && !areRolesValid(roleList)) {
            throw new IllegalArgumentException("Invalid roles in allowed list: " + roles);
        }
        this.allowedRoles = roleList;
    }

    /**
     * Sets the disallowed roles after validation
     * @param roles Comma-separated string of roles
     * @throws IllegalArgumentException if strictMode is true and invalid roles are found
     */
    public void setDisallowedRoles(String roles) {
        List<String> roleList = stringToList(roles);
        if (strictMode && !areRolesValid(roleList)) {
            throw new IllegalArgumentException("Invalid roles in disallowed list: " + roles);
        }
        this.disallowedRoles = roleList;
    }

    /**
     * Checks if a user with given roles should be allowed access
     * @param userRoles Set of roles the user has
     * @return true if access should be granted, false otherwise
     */
    public boolean isAccessAllowed(Set<String> userRoles) {
        // Check disallowed roles first (higher priority)
        for (String role : disallowedRoles) {
            if (userRoles.contains(role)) {
                return false;
            }
        }

        // If allowedRoles is empty, allow by default
        if (allowedRoles.isEmpty()) {
            return true;
        }

        // Check if user has any allowed role
        for (String role : allowedRoles) {
            if (userRoles.contains(role)) {
                return true;
            }
        }

        return false;
    }

    /**
     * Converts comma-separated string to list of roles
     * @param val Comma-separated role string
     * @return List of roles
     */
    protected List<String> stringToList(String val) {
        if (val != null && !val.trim().isEmpty()) {
            return Arrays.asList(val.trim().split("\\s*,\\s*"));
        }
        return Collections.emptyList();
    }

    /**
     * Validates if all roles exist in the system
     * @param roles List of roles to validate
     * @return true if all roles are valid, false otherwise
     */
    protected boolean areRolesValid(List<String> roles) {
        for (String role : roles) {
            if (!validRoles.contains(role)) {
                return false;
            }
        }
        return true;
    }
}
