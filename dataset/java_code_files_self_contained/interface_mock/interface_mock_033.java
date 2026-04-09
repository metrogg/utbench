import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * EnhancedKeyValueStore provides a thread-safe key-value storage system
 * with additional statistics tracking and advanced data operations.
 */
class EnhancedKeyValueStore {
    private final Map<String, Object> store;
    private final Map<String, AtomicInteger> accessCounts;
    private final Map<String, Long> lastAccessTimes;
    private final Map<Class<?>, Integer> typeCounts;

    public EnhancedKeyValueStore() {
        this.store = new ConcurrentHashMap<>();
        this.accessCounts = new ConcurrentHashMap<>();
        this.lastAccessTimes = new ConcurrentHashMap<>();
        this.typeCounts = new ConcurrentHashMap<>();
    }

    /**
     * Stores a value with the given key and tracks its type
     * @param key The key to store the value under
     * @param value The value to store (can be any Object)
     */
    public void put(String key, Object value) {
        if (key == null || value == null) {
            throw new IllegalArgumentException("Key and value cannot be null");
        }
        
        store.put(key, value);
        accessCounts.putIfAbsent(key, new AtomicInteger(0));
        lastAccessTimes.put(key, System.currentTimeMillis());
        typeCounts.merge(value.getClass(), 1, Integer::sum);
    }

    /**
     * Retrieves a value by key and updates access statistics
     * @param key The key to retrieve
     * @return The stored value or null if not found
     */
    public Object get(String key) {
        if (key == null) {
            throw new IllegalArgumentException("Key cannot be null");
        }
        
        Object value = store.get(key);
        if (value != null) {
            accessCounts.get(key).incrementAndGet();
            lastAccessTimes.put(key, System.currentTimeMillis());
        }
        return value;
    }

    /**
     * Gets the access count for a specific key
     * @param key The key to check
     * @return Number of times the key has been accessed
     */
    public int getAccessCount(String key) {
        AtomicInteger count = accessCounts.get(key);
        return count != null ? count.get() : 0;
    }

    /**
     * Gets the last access time for a key
     * @param key The key to check
     * @return Last access time in milliseconds since epoch
     */
    public long getLastAccessTime(String key) {
        Long time = lastAccessTimes.get(key);
        return time != null ? time : 0L;
    }

    /**
     * Gets the count of values stored by type
     * @param type The class type to count
     * @return Number of values of the specified type
     */
    public int getTypeCount(Class<?> type) {
        return typeCounts.getOrDefault(type, 0);
    }

    /**
     * Performs a bulk put operation
     * @param entries Map of key-value pairs to store
     */
    public void putAll(Map<String, Object> entries) {
        if (entries == null) {
            throw new IllegalArgumentException("Entries map cannot be null");
        }
        entries.forEach(this::put);
    }

    /**
     * Computes statistics for all stored values
     * @return Map containing various statistics
     */
    public Map<String, Object> computeStatistics() {
        Map<String, Object> stats = new HashMap<>();
        
        // Basic counts
        stats.put("totalEntries", store.size());
        stats.put("keyWithMaxAccess", 
            accessCounts.entrySet().stream()
                .max(Map.Entry.comparingByValue(Comparator.comparing(AtomicInteger::get)))
                .map(Map.Entry::getKey)
                .orElse("N/A"));
        
        // Type distribution
        stats.put("typeDistribution", new HashMap<>(typeCounts));
        
        // Memory estimate (very rough)
        stats.put("estimatedMemoryBytes", store.values().stream()
            .mapToInt(this::estimateObjectSize)
            .sum());
            
        return stats;
    }

    private int estimateObjectSize(Object obj) {
        if (obj == null) return 0;
        if (obj instanceof String) return ((String) obj).length() * 2;
        if (obj instanceof Number) return 8; // rough estimate for numbers
        if (obj instanceof Boolean) return 1;
        return 16; // default estimate for other objects
    }

    /**
     * Finds keys by value type
     * @param type The class type to search for
     * @return List of keys storing values of the specified type
     */
    public List<String> findKeysByType(Class<?> type) {
        List<String> keys = new ArrayList<>();
        store.forEach((key, value) -> {
            if (value != null && type.isInstance(value)) {
                keys.add(key);
            }
        });
        return keys;
    }

    /**
     * Clears all stored data and statistics
     */
    public void clear() {
        store.clear();
        accessCounts.clear();
        lastAccessTimes.clear();
        typeCounts.clear();
    }
}
