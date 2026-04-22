import java.util.Date;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Simulates a simplified version of Netty server handler for monitoring connections and messages.
 * This version removes Netty dependencies and focuses on the core message handling logic.
 */
class ConnectionMonitor {

    private final ConcurrentHashMap<String, String> ipChannelCache = new ConcurrentHashMap<>();
    private final ConcurrentHashMap<String, String> ipMessageCache = new ConcurrentHashMap<>();
    private final AtomicInteger idleCount = new AtomicInteger(1);

    /**
     * Simulates a new connection being established.
     * @param remoteAddress The remote address of the client
     * @param channelId The unique channel ID
     */
    public void handleConnection(String remoteAddress, String channelId) {
        System.out.println("New connection from: " + remoteAddress);
        ipChannelCache.put(remoteAddress, channelId);
        ipMessageCache.put(remoteAddress, "");
    }

    /**
     * Processes incoming messages from a client.
     * @param remoteAddress The remote address of the client
     * @param message The received message
     */
    public void handleMessage(String remoteAddress, String message) {
        if (message != null && !message.isEmpty()) {
            String currentValue = ipMessageCache.get(remoteAddress);
            if (currentValue == null || currentValue.isEmpty()) {
                ipMessageCache.put(remoteAddress, message);
            } else {
                ipMessageCache.put(remoteAddress, currentValue + "\n" + message);
            }
            System.out.println("Message from " + remoteAddress + ": " + message);
        }
    }

    /**
     * Simulates a heartbeat timeout event.
     * @param remoteAddress The remote address of the client
     * @return true if connection should be closed, false otherwise
     */
    public boolean handleHeartbeatTimeout(String remoteAddress) {
        System.out.println("Heartbeat timeout from: " + remoteAddress);
        if (idleCount.incrementAndGet() > 2) {
            System.out.println("Closing connection due to multiple timeouts: " + remoteAddress);
            ipChannelCache.remove(remoteAddress);
            ipMessageCache.remove(remoteAddress);
            return true;
        }
        return false;
    }

    /**
     * Simulates connection termination and processes collected messages.
     * @param remoteAddress The remote address of the client
     * @return The concatenated messages received from this client
     */
    public String handleDisconnection(String remoteAddress) {
        String channelId = ipChannelCache.get(remoteAddress);
        String allMessages = ipMessageCache.get(remoteAddress);
        
        if (allMessages != null && !allMessages.isEmpty()) {
            System.out.println("Processing messages from " + remoteAddress);
            // In a real implementation, we would analyze and store the messages here
        }
        
        ipChannelCache.remove(remoteAddress);
        ipMessageCache.remove(remoteAddress);
        System.out.println("Connection closed: " + remoteAddress);
        
        return allMessages != null ? allMessages : "";
    }

    /**
     * Gets the current connection count.
     * @return Number of active connections
     */
    public int getConnectionCount() {
        return ipChannelCache.size();
    }
}
