// Converted Java method
import java.util.LinkedList;
import java.util.Queue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

/**
 * An enhanced thread-safe blocking queue implementation with additional features:
 * - Time-bound operations
 * - Capacity management
 * - Statistics tracking
 * - Multiple waiting conditions
 */
class EnhancedBlockingQueue<T> {
    private final Queue<T> queue;
    private final int capacity;
    private final Lock lock = new ReentrantLock();
    private final Condition notFull = lock.newCondition();
    private final Condition notEmpty = lock.newCondition();
    
    // Statistics tracking
    private int totalAdded = 0;
    private int totalRemoved = 0;
    private int maxSizeReached = 0;

    public EnhancedBlockingQueue(int capacity) {
        if (capacity <= 0) {
            throw new IllegalArgumentException("Capacity must be positive");
        }
        this.capacity = capacity;
        this.queue = new LinkedList<>();
    }

    /**
     * Adds an element to the queue, waiting if necessary for space to become available.
     */
    public void put(T element) throws InterruptedException {
        if (element == null) {
            throw new NullPointerException("Element cannot be null");
        }
        
        lock.lock();
        try {
            while (queue.size() == capacity) {
                notFull.await();
            }
            queue.add(element);
            totalAdded++;
            updateMaxSize();
            notEmpty.signal();
        } finally {
            lock.unlock();
        }
    }

    /**
     * Adds an element to the queue, waiting up to the specified time for space to become available.
     * @return true if successful, false if timeout elapsed
     */
    public boolean offer(T element, long timeout, TimeUnit unit) throws InterruptedException {
        if (element == null) {
            throw new NullPointerException("Element cannot be null");
        }
        
        long nanos = unit.toNanos(timeout);
        lock.lock();
        try {
            while (queue.size() == capacity) {
                if (nanos <= 0) {
                    return false;
                }
                nanos = notFull.awaitNanos(nanos);
            }
            queue.add(element);
            totalAdded++;
            updateMaxSize();
            notEmpty.signal();
            return true;
        } finally {
            lock.unlock();
        }
    }

    /**
     * Retrieves and removes the head of the queue, waiting if necessary until an element becomes available.
     */
    public T take() throws InterruptedException {
        lock.lock();
        try {
            while (queue.isEmpty()) {
                notEmpty.await();
            }
            T item = queue.remove();
            totalRemoved++;
            notFull.signal();
            return item;
        } finally {
            lock.unlock();
        }
    }

    /**
     * Retrieves and removes the head of the queue, waiting up to the specified time if necessary.
     * @return the head of the queue, or null if timeout elapsed
     */
    public T poll(long timeout, TimeUnit unit) throws InterruptedException {
        long nanos = unit.toNanos(timeout);
        lock.lock();
        try {
            while (queue.isEmpty()) {
                if (nanos <= 0) {
                    return null;
                }
                nanos = notEmpty.awaitNanos(nanos);
            }
            T item = queue.remove();
            totalRemoved++;
            notFull.signal();
            return item;
        } finally {
            lock.unlock();
        }
    }

    /**
     * Retrieves but does not remove the head of the queue.
     * @return the head of the queue, or null if empty
     */
    public T peek() {
        lock.lock();
        try {
            return queue.peek();
        } finally {
            lock.unlock();
        }
    }

    /**
     * Returns the current size of the queue.
     */
    public int size() {
        lock.lock();
        try {
            return queue.size();
        } finally {
            lock.unlock();
        }
    }

    /**
     * Returns the remaining capacity of the queue.
     */
    public int remainingCapacity() {
        lock.lock();
        try {
            return capacity - queue.size();
        } finally {
            lock.unlock();
        }
    }

    /**
     * Clears the queue and returns the number of elements cleared.
     */
    public int clear() {
        lock.lock();
        try {
            int count = queue.size();
            queue.clear();
            notFull.signalAll();
            return count;
        } finally {
            lock.unlock();
        }
    }

    /**
     * Returns queue statistics.
     */
    public String getStats() {
        lock.lock();
        try {
            return String.format("Added: %d, Removed: %d, Current: %d, Max: %d",
                    totalAdded, totalRemoved, queue.size(), maxSizeReached);
        } finally {
            lock.unlock();
        }
    }

    private void updateMaxSize() {
        if (queue.size() > maxSizeReached) {
            maxSizeReached = queue.size();
        }
    }
}
