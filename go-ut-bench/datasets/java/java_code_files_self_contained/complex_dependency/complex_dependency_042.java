// Converted Java method
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;

class ThreadManager {
    private final AtomicBoolean isInterrupted = new AtomicBoolean(false);
    private final AtomicInteger completedTasks = new AtomicInteger(0);
    private final AtomicInteger interruptedTasks = new AtomicInteger(0);

    /**
     * Creates and manages worker threads that perform tasks with possible interruption.
     * @param threadCount Number of worker threads to create
     * @param workDurationMs Duration each thread should work (sleep) in milliseconds
     * @param interruptAfterMs When to interrupt threads (milliseconds after start)
     * @return Statistics about thread execution
     */
    public ThreadStatistics manageThreads(int threadCount, long workDurationMs, long interruptAfterMs) {
        Thread[] workers = new Thread[threadCount];
        
        // Create and start worker threads
        for (int i = 0; i < threadCount; i++) {
            workers[i] = new Thread(() -> {
                try {
                    System.out.println(Thread.currentThread().getName() + " started working");
                    Thread.sleep(workDurationMs);
                    completedTasks.incrementAndGet();
                    System.out.println(Thread.currentThread().getName() + " completed work");
                } catch (InterruptedException e) {
                    interruptedTasks.incrementAndGet();
                    System.out.println(Thread.currentThread().getName() + " was interrupted");
                }
            });
            workers[i].start();
        }

        // Schedule interruption
        new Thread(() -> {
            try {
                Thread.sleep(interruptAfterMs);
                isInterrupted.set(true);
                for (Thread worker : workers) {
                    if (worker.isAlive()) {
                        worker.interrupt();
                    }
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }).start();

        // Wait for all threads to finish
        for (Thread worker : workers) {
            try {
                worker.join();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }

        return new ThreadStatistics(threadCount, completedTasks.get(), interruptedTasks.get());
    }

    public static class ThreadStatistics {
        public final int totalThreads;
        public final int completedTasks;
        public final int interruptedTasks;

        public ThreadStatistics(int totalThreads, int completedTasks, int interruptedTasks) {
            this.totalThreads = totalThreads;
            this.completedTasks = completedTasks;
            this.interruptedTasks = interruptedTasks;
        }

        @Override
        public String toString() {
            return String.format("ThreadStats{total=%d, completed=%d, interrupted=%d}",
                    totalThreads, completedTasks, interruptedTasks);
        }
    }
}
