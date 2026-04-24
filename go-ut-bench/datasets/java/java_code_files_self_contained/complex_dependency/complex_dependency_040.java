// Converted Java method
import java.util.*;
import java.util.stream.Collectors;

class AdvancedFlagProcessor {
    private final Set<String> validFlags;
    private final Set<String> mutuallyExclusiveGroups;

    /**
     * Creates an AdvancedFlagProcessor with specified valid flags and mutually exclusive groups.
     * @param validFlags Collection of all valid flags
     * @param mutuallyExclusiveGroups Collection of flag groups where only one flag per group can be set
     */
    public AdvancedFlagProcessor(Set<String> validFlags, Set<String> mutuallyExclusiveGroups) {
        this.validFlags = new HashSet<>(validFlags);
        this.mutuallyExclusiveGroups = new HashSet<>(mutuallyExclusiveGroups);
    }

    /**
     * Processes a flag string, validates it, and returns all possible valid combinations.
     * @param flagString Input flag string in format "flag1|flag2|flag3"
     * @return Set of all valid flag combinations that can be formed by adding one more flag
     */
    public Set<String> processFlags(String flagString) {
        Set<String> result = new HashSet<>();
        
        if (flagString == null || flagString.trim().isEmpty()) {
            // Return all single flags if input is empty
            return validFlags.stream()
                    .map(flag -> "|" + flag)
                    .collect(Collectors.toSet());
        }

        String[] currentFlags = flagString.split("\\|");
        Set<String> currentFlagSet = new HashSet<>(Arrays.asList(currentFlags));
        
        // Validate current flags
        if (!validFlags.containsAll(currentFlagSet)) {
            throw new IllegalArgumentException("Invalid flag(s) in input: " + flagString);
        }

        // Check for mutually exclusive flags
        checkMutualExclusivity(currentFlagSet);

        // Generate possible next flags
        for (String flag : validFlags) {
            if (!currentFlagSet.contains(flag)) {
                Set<String> newFlagSet = new HashSet<>(currentFlagSet);
                newFlagSet.add(flag);
                
                // Check if new flag would violate mutual exclusivity
                if (!violatesMutualExclusivity(newFlagSet)) {
                    String newFlagString = String.join("|", newFlagSet);
                    result.add(newFlagString);
                }
            }
        }

        return result;
    }

    private void checkMutualExclusivity(Set<String> flags) {
        if (mutuallyExclusiveGroups.isEmpty()) return;
        
        Map<String, Integer> groupCounts = new HashMap<>();
        for (String flag : flags) {
            if (mutuallyExclusiveGroups.contains(flag)) {
                groupCounts.put(flag, groupCounts.getOrDefault(flag, 0) + 1);
            }
        }
        
        if (groupCounts.size() > 1) {
            throw new IllegalArgumentException(
                "Mutually exclusive flags cannot be combined: " + 
                String.join(", ", groupCounts.keySet())
            );
        }
    }

    private boolean violatesMutualExclusivity(Set<String> flags) {
        if (mutuallyExclusiveGroups.isEmpty()) return false;
        
        long exclusiveFlagsPresent = flags.stream()
            .filter(mutuallyExclusiveGroups::contains)
            .count();
            
        return exclusiveFlagsPresent > 1;
    }
}
