// Converted Java method
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.TimeUnit;

class MessageProducer {
    private final BlockingQueue<String> messageQueue;
    private final List<String> sentMessages;
    private final int maxRetries;
    private final long retryDelayMillis;

    /**
     * Creates a MessageProducer with default configuration.
     */
    public MessageProducer() {
        this(new LinkedBlockingQueue<>(), 3, 100);
    }

    /**
     * Creates a MessageProducer with custom configuration.
     * @param messageQueue The queue to store messages before sending
     * @param maxRetries Maximum number of retry attempts for failed sends
     * @param retryDelayMillis Delay between retry attempts in milliseconds
     */
    public MessageProducer(BlockingQueue<String> messageQueue, int maxRetries, long retryDelayMillis) {
        this.messageQueue = messageQueue;
        this.sentMessages = new ArrayList<>();
        this.maxRetries = maxRetries;
        this.retryDelayMillis = retryDelayMillis;
    }

    /**
     * Sends multiple messages to the queue with retry logic.
     * @param messagePrefix Prefix for each message
     * @param count Number of messages to send
     * @return List of successfully sent messages
     * @throws InterruptedException if interrupted during sending
     */
    public List<String> sendMessages(String messagePrefix, int count) throws InterruptedException {
        sentMessages.clear();
        
        for (int i = 1; i <= count; i++) {
            String message = messagePrefix + " " + i;
            boolean sent = false;
            int attempts = 0;
            
            while (!sent && attempts < maxRetries) {
                try {
                    // Simulate message sending by putting in queue
                    boolean offered = messageQueue.offer(message, 100, TimeUnit.MILLISECONDS);
                    if (offered) {
                        sentMessages.add(message);
                        sent = true;
                    } else {
                        attempts++;
                        Thread.sleep(retryDelayMillis);
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    throw e;
                }
            }
            
            if (!sent) {
                System.err.println("Failed to send message after " + maxRetries + " attempts: " + message);
            }
        }
        
        return new ArrayList<>(sentMessages);
    }

    /**
     * Gets the current message queue for testing purposes.
     * @return The message queue
     */
    public BlockingQueue<String> getMessageQueue() {
        return messageQueue;
    }
}
