import java.util.*;
import java.util.concurrent.*;
import java.util.function.*;

class WebSocketMessageHandler {
    private final Map<String, Consumer<Map<String, Object>>> messageHandlers;
    private final RateLimiter rateLimiter;
    private final BiConsumer<String, String> errorLogger;

    public WebSocketMessageHandler(BiConsumer<String, String> errorLogger) {
        this.messageHandlers = new ConcurrentHashMap<>();
        this.rateLimiter = new RateLimiter(100); // Allow 100 messages per second
        this.errorLogger = errorLogger;
    }

    /**
     * Adds a new message handler for a specific message type.
     * @param messageType The type of message to handle
     * @param handler The handler function that processes the message
     */
    public void addMessageHandler(String messageType, Consumer<Map<String, Object>> handler) {
        if (messageType == null || handler == null) {
            throw new IllegalArgumentException("Message type and handler cannot be null");
        }
        messageHandlers.put(messageType, handler);
    }

    /**
     * Processes an incoming WebSocket message with validation and rate limiting.
     * @param messageType The type of message received
     * @param messageData The message payload
     * @return true if message was processed successfully, false otherwise
     */
    public boolean processMessage(String messageType, Map<String, Object> messageData) {
        if (!rateLimiter.tryAcquire()) {
            errorLogger.accept("Rate limit exceeded", "Too many messages received");
            return false;
        }

        if (!validateMessage(messageType, messageData)) {
            errorLogger.accept("Invalid message", "Message validation failed for type: " + messageType);
            return false;
        }

        Consumer<Map<String, Object>> handler = messageHandlers.get(messageType);
        if (handler != null) {
            try {
                handler.accept(messageData);
                return true;
            } catch (Exception e) {
                errorLogger.accept("Handler error", "Error processing message: " + e.getMessage());
                return false;
            }
        }
        return false;
    }

    private boolean validateMessage(String messageType, Map<String, Object> messageData) {
        if (messageType == null || messageType.trim().isEmpty()) {
            return false;
        }
        
        // Basic validation - message data must be a non-null map
        return messageData != null;
    }

    // Nested rate limiter class
    private static class RateLimiter {
        private final int maxRequestsPerSecond;
        private final Queue<Long> requestTimes;

        public RateLimiter(int maxRequestsPerSecond) {
            this.maxRequestsPerSecond = maxRequestsPerSecond;
            this.requestTimes = new ConcurrentLinkedQueue<>();
        }

        public synchronized boolean tryAcquire() {
            long now = System.currentTimeMillis();
            
            // Remove old requests
            while (!requestTimes.isEmpty() && now - requestTimes.peek() > 1000) {
                requestTimes.poll();
            }
            
            // Check if we can add a new request
            if (requestTimes.size() < maxRequestsPerSecond) {
                requestTimes.add(now);
                return true;
            }
            return false;
        }
    }
}
