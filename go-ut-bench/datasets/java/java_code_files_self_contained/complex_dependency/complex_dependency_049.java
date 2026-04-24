// Converted Java method
import java.util.*;
import java.util.concurrent.*;

class ThreadResourceManager {
    private final Map<String, Thread> threadResources;
    private final ScheduledExecutorService monitorExecutor;

    /**
     * Manages daemon thread resources with monitoring capabilities.
     * Tracks created threads and provides monitoring functionality.
     */
    public ThreadResourceManager() {
        this.threadResources = new ConcurrentHashMap<>();
        this.monitorExecutor = Executors.newSingleThreadScheduledExecutor();
        startMonitoring();
    }

    /**
     * Creates a new daemon thread with the given name and runnable task.
     * @param threadName Name for the new thread
     * @param task Runnable task for the thread to execute
     * @return true if thread was created successfully, false otherwise
     */
    public boolean createDaemonThread(String threadName, Runnable task) {
        if (threadName == null || threadName.trim().isEmpty() || task == null) {
            return false;
        }

        if (threadResources.containsKey(threadName)) {
            return false;
        }

        Thread thread = new Thread(task);
        thread.setName(threadName);
        thread.setDaemon(true);
        thread.start();
        threadResources.put(threadName, thread);
        return true;
    }

    /**
     * Stops and removes a daemon thread by name.
     * @param threadName Name of the thread to stop
     * @return true if thread was found and stopped, false otherwise
     */
    public boolean stopDaemonThread(String threadName) {
        Thread thread = threadResources.get(threadName);
        if (thread != null) {
            thread.interrupt();
            threadResources.remove(threadName);
            return true;
        }
        return false;
    }

    /**
     * Gets the count of active daemon threads being managed.
     * @return Number of active threads
     */
    public int getActiveThreadCount() {
        return threadResources.size();
    }

    /**
     * Starts periodic monitoring of thread status.
     */
    private void startMonitoring() {
        monitorExecutor.scheduleAtFixedRate(() -> {
            threadResources.forEach((name, thread) -> {
                if (!thread.isAlive()) {
                    threadResources.remove(name);
                }
            });
        }, 1, 1, TimeUnit.SECONDS);
    }

    /**
     * Shuts down the thread manager and cleans up resources.
     */
    public void shutdown() {
        threadResources.values().forEach(Thread::interrupt);
        threadResources.clear();
        monitorExecutor.shutdown();
    }
}
