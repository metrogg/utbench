// Converted Java method
import java.util.HashMap;
import java.util.Map;

class MessageQueueConfigurator {
    /**
     * Simulates RabbitMQ configuration and message publishing operations.
     * This is a simplified version that doesn't require Spring or RabbitMQ,
     * but maintains the core logical flow of queue/exchange/binding setup.
     *
     * @param exchangeName Name of the exchange to create
     * @param queueName Name of the queue to create
     * @param routingKey Routing key for binding
     * @return Map containing configuration details and simulated message publishing status
     */
    public Map<String, Object> configureMessageQueue(String exchangeName, String queueName, String routingKey) {
        Map<String, Object> config = new HashMap<>();
        
        // Simulate exchange creation
        config.put("exchange", exchangeName);
        config.put("exchangeType", "topic");
        
        // Simulate queue creation
        config.put("queue", queueName);
        config.put("queueDurable", true);
        config.put("queueExclusive", false);
        config.put("queueAutoDelete", false);
        
        // Simulate binding
        config.put("binding", routingKey);
        
        // Simulate message converter setup
        config.put("messageConverter", "JSON");
        
        // Simulate template creation
        Map<String, String> templateConfig = new HashMap<>();
        templateConfig.put("connectionFactory", "simulated");
        templateConfig.put("messageConverter", "JSON");
        config.put("template", templateConfig);
        
        // Simulate message publishing
        config.put("lastPublished", System.currentTimeMillis());
        config.put("status", "READY");
        
        return config;
    }
    
    /**
     * Simulates publishing a message to the configured queue.
     *
     * @param config The configuration map from configureMessageQueue
     * @param message The message to publish
     * @return Updated configuration with publishing status
     */
    public Map<String, Object> publishMessage(Map<String, Object> config, String message) {
        if (!"READY".equals(config.get("status"))) {
            throw new IllegalStateException("Queue not configured properly");
        }
        
        config.put("lastMessage", message);
        config.put("lastPublished", System.currentTimeMillis());
        config.put("messageCount", (int)config.getOrDefault("messageCount", 0) + 1);
        
        return config;
    }
}
