// Converted Java method
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

class MessageSystemSimulator {
    private final Map<String, ServiceNode> nodes = new HashMap<>();
    private final Map<String, String> addressMap = new HashMap<>();
    
    /**
     * Registers a new service node in the message system
     * @param serviceType Type of service (e.g., "Frontend", "DB")
     * @param handler The service implementation
     * @return The generated address for the service
     */
    public String registerService(String serviceType, MessageHandler handler) {
        if (serviceType == null || serviceType.isEmpty()) {
            throw new IllegalArgumentException("Service type cannot be null or empty");
        }
        
        String address = generateAddress(serviceType);
        ServiceNode node = new ServiceNode(address, handler);
        nodes.put(address, node);
        addressMap.put(serviceType, address);
        return address;
    }
    
    /**
     * Sends a message to a specific service
     * @param fromAddress Sender's address
     * @param toServiceType Receiver's service type
     * @param messageContent Message content
     * @return true if message was delivered successfully
     */
    public boolean sendMessage(String fromAddress, String toServiceType, String messageContent) {
        String toAddress = addressMap.get(toServiceType);
        if (toAddress == null) {
            return false;
        }
        
        ServiceNode receiver = nodes.get(toAddress);
        if (receiver == null) {
            return false;
        }
        
        Message message = new Message(fromAddress, toAddress, messageContent);
        return receiver.handleMessage(message);
    }
    
    /**
     * Gets the current status of all nodes
     * @return Map of addresses to status strings
     */
    public Map<String, String> getSystemStatus() {
        Map<String, String> status = new HashMap<>();
        nodes.forEach((addr, node) -> status.put(addr, node.getStatus()));
        return status;
    }
    
    private String generateAddress(String prefix) {
        return prefix + "-" + UUID.randomUUID().toString().substring(0, 8);
    }
    
    static class ServiceNode {
        private final String address;
        private final MessageHandler handler;
        private int messageCount = 0;
        
        public ServiceNode(String address, MessageHandler handler) {
            this.address = address;
            this.handler = handler;
        }
        
        public boolean handleMessage(Message message) {
            messageCount++;
            return handler.process(message);
        }
        
        public String getStatus() {
            return "Active (" + messageCount + " messages processed)";
        }
    }
    
    interface MessageHandler {
        boolean process(Message message);
    }
    
    static class Message {
        private final String from;
        private final String to;
        private final String content;
        
        public Message(String from, String to, String content) {
            this.from = from;
            this.to = to;
            this.content = content;
        }
        
        public String getContent() {
            return content;
        }
    }
}
