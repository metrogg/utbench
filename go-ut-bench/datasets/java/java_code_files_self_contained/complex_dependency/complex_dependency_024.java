// Converted Java method
import java.util.HashMap;
import java.util.Map;

/**
 * Represents a time-based machine entry with validation and statistics capabilities.
 * Enhanced from the original EntryImpl to include more complex operations.
 */
class MachineEntry {
    private String clock;
    private String machine;
    private Map<String, Integer> operationCounts = new HashMap<>();

    /**
     * Creates a new MachineEntry with the given clock and machine names.
     * @param clock The clock identifier (cannot be null or empty)
     * @param machine The machine identifier (cannot be null or empty)
     * @throws IllegalArgumentException if clock or machine is invalid
     */
    public MachineEntry(String clock, String machine) {
        if (clock == null || clock.trim().isEmpty()) {
            throw new IllegalArgumentException("Clock cannot be null or empty");
        }
        if (machine == null || machine.trim().isEmpty()) {
            throw new IllegalArgumentException("Machine cannot be null or empty");
        }
        this.clock = clock;
        this.machine = machine;
    }

    /**
     * Records an operation performed by this machine.
     * @param operationName Name of the operation (cannot be null or empty)
     * @throws IllegalArgumentException if operationName is invalid
     */
    public void recordOperation(String operationName) {
        if (operationName == null || operationName.trim().isEmpty()) {
            throw new IllegalArgumentException("Operation name cannot be null or empty");
        }
        operationCounts.merge(operationName, 1, Integer::sum);
    }

    /**
     * Gets the count of a specific operation.
     * @param operationName Name of the operation to query
     * @return Count of the operation, or 0 if never recorded
     */
    public int getOperationCount(String operationName) {
        return operationCounts.getOrDefault(operationName, 0);
    }

    /**
     * Calculates and returns operation statistics.
     * @return Map containing:
     *         - "total": total operations recorded
     *         - "unique": number of unique operations
     *         - "mostFrequent": name of most frequent operation
     *         - "maxCount": count of most frequent operation
     */
    public Map<String, Object> getOperationStatistics() {
        Map<String, Object> stats = new HashMap<>();
        int total = 0;
        int maxCount = 0;
        String mostFrequent = null;

        for (Map.Entry<String, Integer> entry : operationCounts.entrySet()) {
            total += entry.getValue();
            if (entry.getValue() > maxCount) {
                maxCount = entry.getValue();
                mostFrequent = entry.getKey();
            }
        }

        stats.put("total", total);
        stats.put("unique", operationCounts.size());
        stats.put("mostFrequent", mostFrequent);
        stats.put("maxCount", maxCount);
        return stats;
    }

    public String getClock() {
        return clock;
    }

    public String getMachine() {
        return machine;
    }

    @Override
    public String toString() {
        return "MachineEntry [clock=" + clock + ", machine=" + machine + 
               ", operations=" + operationCounts + "]";
    }
}
