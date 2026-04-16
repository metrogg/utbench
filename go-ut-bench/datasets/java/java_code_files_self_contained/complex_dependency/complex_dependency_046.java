import java.util.HashMap;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

/**
 * Enhanced asynchronous response handler with multiple callback options
 * and error handling capabilities.
 */
class AsyncResponseHandler {
    private final HashMap<String, Consumer<Object>> callbacks = new HashMap<>();
    private Consumer<Exception> errorHandler;
    
    /**
     * Registers a callback for a specific response type
     * @param responseType The type of response to handle
     * @param callback The callback function to execute
     */
    public void registerCallback(String responseType, Consumer<Object> callback) {
        callbacks.put(responseType, callback);
    }
    
    /**
     * Registers a global error handler
     * @param handler The error handling function
     */
    public void registerErrorHandler(Consumer<Exception> handler) {
        this.errorHandler = handler;
    }
    
    /**
     * Processes an asynchronous response with the appropriate callback
     * @param responseType The type of response received
     * @param response The response data
     */
    public void processResponse(String responseType, Object response) {
        if (callbacks.containsKey(responseType)) {
            callbacks.get(responseType).accept(response);
        } else if (errorHandler != null) {
            errorHandler.accept(new IllegalArgumentException("No handler registered for response type: " + responseType));
        }
    }
    
    /**
     * Executes an asynchronous task and handles the response
     * @param task The task to execute asynchronously
     * @param responseType The expected response type
     * @return CompletableFuture for chaining additional operations
     */
    public CompletableFuture<Void> executeAsync(Runnable task, String responseType) {
        return CompletableFuture.runAsync(() -> {
            try {
                task.run();
                processResponse(responseType, "Task completed successfully");
            } catch (Exception e) {
                if (errorHandler != null) {
                    errorHandler.accept(e);
                }
            }
        });
    }
}
