// Converted Java method
import java.util.*;
import java.util.function.Function;

class MessageProcessor {
    private final Map<String, Function<String, String>> messageHandlers;
    private final Map<String, String> routingRules;
    private final Queue<String> messageQueue;
    private final Set<String> processedMessageIds;

    /**
     * Initializes a MessageProcessor with routing rules and message handlers.
     */
    public MessageProcessor() {
        this.messageHandlers = new HashMap<>();
        this.routingRules = new HashMap<>();
        this.messageQueue = new LinkedList<>();
        this.processedMessageIds = new HashSet<>();
    }

    /**
     * Adds a new message handler for a specific message type.
     * @param messageType The type of message this handler can process
     * @param handler The function that processes the message
     */
    public void addHandler(String messageType, Function<String, String> handler) {
        messageHandlers.put(messageType, handler);
    }

    /**
     * Adds a routing rule that maps message patterns to handlers.
     * @param pattern The pattern to match against message types
     * @param handlerKey The key of the handler to use
     */
    public void addRoutingRule(String pattern, String handlerKey) {
        routingRules.put(pattern, handlerKey);
    }

    /**
     * Processes a single message using the appropriate handler.
     * @param messageId Unique identifier for the message
     * @param messageType Type of the message
     * @param content Content of the message
     * @return Processing result or null if no handler found
     */
    public String processMessage(String messageId, String messageType, String content) {
        if (processedMessageIds.contains(messageId)) {
            return "Message already processed";
        }

        String handlerKey = findHandlerKey(messageType);
        if (handlerKey == null || !messageHandlers.containsKey(handlerKey)) {
            return "No handler found for message type: " + messageType;
        }

        String result = messageHandlers.get(handlerKey).apply(content);
        processedMessageIds.add(messageId);
        return result;
    }

    /**
     * Processes all messages in the queue.
     * @return Map of message IDs to their processing results
     */
    public Map<String, String> processQueue() {
        Map<String, String> results = new HashMap<>();
        while (!messageQueue.isEmpty()) {
            String message = messageQueue.poll();
            // Simplified parsing - in real implementation would parse message parts
            String[] parts = message.split("\\|", 3);
            if (parts.length == 3) {
                String result = processMessage(parts[0], parts[1], parts[2]);
                results.put(parts[0], result);
            }
        }
        return results;
    }

    /**
     * Adds a message to the processing queue.
     * @param messageId Unique identifier for the message
     * @param messageType Type of the message
     * @param content Content of the message
     */
    public void enqueueMessage(String messageId, String messageType, String content) {
        messageQueue.offer(messageId + "|" + messageType + "|" + content);
    }

    private String findHandlerKey(String messageType) {
        for (Map.Entry<String, String> entry : routingRules.entrySet()) {
            if (messageType.matches(entry.getKey().replace("#", ".*").replace("*", "[^.]*"))) {
                return entry.getValue();
            }
        }
        return null;
    }
}
