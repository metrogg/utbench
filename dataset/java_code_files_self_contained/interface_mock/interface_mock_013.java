// Converted Java method
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

class SystemShutdownManager {
    private final Map<String, Service> services;
    private final AtomicBoolean isShuttingDown = new AtomicBoolean(false);
    private final long timeout;
    private final TimeUnit timeoutUnit;

    /**
     * Enhanced system shutdown manager that handles graceful shutdown of multiple services.
     *
     * @param services    Map of service names to Service objects
     * @param timeout     Maximum time to wait for shutdown
     * @param timeoutUnit Time unit for the timeout
     */
    public SystemShutdownManager(Map<String, Service> services, long timeout, TimeUnit timeoutUnit) {
        this.services = new ConcurrentHashMap<>(services);
        this.timeout = timeout;
        this.timeoutUnit = timeoutUnit;
    }

    /**
     * Initiates a graceful shutdown of all registered services.
     * Attempts to stop all services within the configured timeout period.
     *
     * @return true if all services were stopped gracefully, false otherwise
     */
    public boolean shutdownAll() {
        if (isShuttingDown.getAndSet(true)) {
            return false; // Already shutting down
        }

        long startTime = System.currentTimeMillis();
        long timeoutMillis = timeoutUnit.toMillis(timeout);
        boolean allStopped = true;

        for (Map.Entry<String, Service> entry : services.entrySet()) {
            String serviceName = entry.getKey();
            Service service = entry.getValue();

            System.out.println("Initiating shutdown for service: " + serviceName);
            long serviceStartTime = System.currentTimeMillis();

            try {
                service.stop();
                long duration = System.currentTimeMillis() - serviceStartTime;
                System.out.println("Successfully stopped service " + serviceName + " in " + duration + "ms");
            } catch (Exception e) {
                System.out.println("Failed to stop service " + serviceName + ": " + e.getMessage());
                allStopped = false;
            }

            // Check if we've exceeded total timeout
            if (System.currentTimeMillis() - startTime > timeoutMillis) {
                System.out.println("Shutdown timeout exceeded for service: " + serviceName);
                allStopped = false;
                break;
            }
        }

        return allStopped;
    }

    /**
     * Interface for services that can be stopped
     */
    public interface Service {
        void stop() throws Exception;
    }
}
