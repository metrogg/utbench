// Converted Java method
import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.atomic.AtomicInteger;

class PollingQueue<T> {
    private final Queue<T> queue;
    private final AtomicInteger pollingCount;
    private final int maxPollingAttempts;

    /**
     * Creates a new PollingQueue with the specified maximum polling attempts.
     * @param maxPollingAttempts the maximum number of polling attempts allowed
     */
    public PollingQueue(int maxPollingAttempts) {
        this.queue = new LinkedList<>();
        this.pollingCount = new AtomicInteger(0);
        this.maxPollingAttempts = maxPollingAttempts;
    }

    /**
     * Adds an element to the queue.
     * @param element the element to add
     */
    public synchronized void enqueue(T element) {
        queue.add(element);
        notifyAll();
    }

    /**
     * Polls the queue for an element, incrementing the polling count.
     * @return the next element in the queue, or null if queue is empty
     * @throws IllegalStateException if maximum polling attempts exceeded
     */
    public synchronized T poll() {
        if (pollingCount.incrementAndGet() > maxPollingAttempts) {
            throw new IllegalStateException("Maximum polling attempts exceeded");
        }
        return queue.poll();
    }

    /**
     * Waits for an element to be available in the queue.
     * @param timeoutMillis maximum time to wait in milliseconds
     * @return the next element in the queue
     * @throws InterruptedException if the thread is interrupted while waiting
     */
    public synchronized T waitForElement(long timeoutMillis) throws InterruptedException {
        while (queue.isEmpty()) {
            wait(timeoutMillis);
            if (queue.isEmpty()) {
                return null;
            }
        }
        return queue.poll();
    }

    /**
     * Resets the polling attempt counter.
     */
    public void resetPollingCount() {
        pollingCount.set(0);
    }

    /**
     * Gets the current polling attempt count.
     * @return the current polling count
     */
    public int getPollingCount() {
        return pollingCount.get();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        PollingQueue<?> that = (PollingQueue<?>) o;
        return maxPollingAttempts == that.maxPollingAttempts && 
               queue.equals(that.queue) && 
               pollingCount.get() == that.pollingCount.get();
    }

    @Override
    public int hashCode() {
        int result = queue.hashCode();
        result = 31 * result + pollingCount.get();
        result = 31 * result + maxPollingAttempts;
        return result;
    }
}
