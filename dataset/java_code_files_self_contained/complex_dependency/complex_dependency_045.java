// Converted Java method
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

/**
 * Enhanced AlertingService that tracks alerts and provides statistics.
 * Supports multiple alert levels and callback notifications.
 */
class AlertingService {
    private final Map<String, Consumer<String>> alertListeners = new ConcurrentHashMap<>();
    private final Map<String, AtomicInteger> alertCounts = new HashMap<>();
    private final Map<String, String> lastAlertMessages = new HashMap<>();

    /**
     * Adds an alert listener for a specific alert level.
     * @param listenerId Unique identifier for the listener
     * @param level Alert level (e.g., "SEVERE", "WARNING")
     * @param callback Callback function to handle alerts
     */
    public void addAlertListener(String listenerId, String level, Consumer<String> callback) {
        alertListeners.put(listenerId + ":" + level, callback);
        alertCounts.putIfAbsent(level, new AtomicInteger(0));
    }

    /**
     * Removes an alert listener.
     * @param listenerId Unique identifier for the listener
     * @param level Alert level
     */
    public void removeAlertListener(String listenerId, String level) {
        alertListeners.remove(listenerId + ":" + level);
    }

    /**
     * Triggers an alert at the specified level.
     * @param level Alert level
     * @param message Alert message
     */
    public void triggerAlert(String level, String message) {
        lastAlertMessages.put(level, message);
        alertCounts.computeIfAbsent(level, k -> new AtomicInteger(0)).incrementAndGet();
        
        alertListeners.entrySet().stream()
            .filter(entry -> entry.getKey().endsWith(":" + level))
            .forEach(entry -> entry.getValue().accept(message));
    }

    /**
     * Gets the count of alerts for a specific level.
     * @param level Alert level
     * @return Count of alerts
     */
    public int getAlertCount(String level) {
        return alertCounts.getOrDefault(level, new AtomicInteger(0)).get();
    }

    /**
     * Gets the last alert message for a specific level.
     * @param level Alert level
     * @return Last alert message or null if none
     */
    public String getLastAlertMessage(String level) {
        return lastAlertMessages.get(level);
    }

    /**
     * Resets all alert statistics.
     */
    public void resetStatistics() {
        alertCounts.replaceAll((k, v) -> new AtomicInteger(0));
        lastAlertMessages.clear();
    }
}
