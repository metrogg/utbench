// Converted Java method
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.stream.Collectors;

class EnhancedMessageProcessor {
    private final Map<String, Message> messages = new LinkedHashMap<>();
    private final Map<String, Integer> messageCounts = new LinkedHashMap<>();

    /**
     * Stores a message for a specific destination and tracks message counts
     * @param message The message to store
     * @param destination The destination identifier
     * @return Processing result with status and metadata
     */
    public synchronized ProcessingResult postMessage(Message message, String destination) {
        // Store the message
        messages.put(destination, message);
        
        // Update message count
        messageCounts.put(destination, messageCounts.getOrDefault(destination, 0) + 1);
        
        // Return processing result with metadata
        return new ProcessingResult(
            Status.SUCCESS,
            "Message stored successfully",
            Map.of(
                "destination", destination,
                "messageLength", message.getContent().length(),
                "totalMessages", messageCounts.get(destination)
            )
        );
    }

    /**
     * Retrieves the latest message for a destination
     * @param destination The destination identifier
     * @return The latest message or null if not found
     */
    public synchronized Message getMessage(String destination) {
        return messages.get(destination);
    }

    /**
     * Gets all messages for a specific content type
     * @param contentType The content type to filter by
     * @return Map of destinations to messages matching the content type
     */
    public synchronized Map<String, Message> getMessagesByType(String contentType) {
        return messages.entrySet().stream()
            .filter(entry -> entry.getValue().getContentType().equals(contentType))
            .collect(Collectors.toMap(
                Map.Entry::getKey,
                Map.Entry::getValue,
                (existing, replacement) -> existing,
                LinkedHashMap::new
            ));
    }

    /**
     * Clears all messages for a specific destination
     * @param destination The destination to clear
     * @return Processing result with status
     */
    public synchronized ProcessingResult clearMessages(String destination) {
        if (messages.containsKey(destination)) {
            messages.remove(destination);
            messageCounts.remove(destination);
            return new ProcessingResult(Status.SUCCESS, "Messages cleared");
        }
        return new ProcessingResult(Status.NOT_FOUND, "Destination not found");
    }

    // Supporting classes
    public static class Message {
        private final String content;
        private final String contentType;
        
        public Message(String content, String contentType) {
            this.content = content;
            this.contentType = contentType;
        }
        
        public String getContent() { return content; }
        public String getContentType() { return contentType; }
    }

    public enum Status { SUCCESS, NOT_FOUND, ERROR }

    public static class ProcessingResult {
        private final Status status;
        private final String message;
        private final Map<String, Object> metadata;
        
        public ProcessingResult(Status status, String message) {
            this(status, message, Map.of());
        }
        
        public ProcessingResult(Status status, String message, Map<String, Object> metadata) {
            this.status = status;
            this.message = message;
            this.metadata = metadata;
        }
        
        // Getters
        public Status getStatus() { return status; }
        public String getMessage() { return message; }
        public Map<String, Object> getMetadata() { return metadata; }
    }
}
