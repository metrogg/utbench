// Converted Java method
import java.util.Calendar;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.ReentrantLock;

class ThreadSynchronizer {
    private final ReentrantLock lock = new ReentrantLock();
    private final Condition condition = lock.newCondition();
    private boolean isSignalled = false;

    /**
     * Waits until either the specified deadline is reached or the thread is notified.
     * @param waitSeconds Number of seconds to wait as deadline
     * @return true if notified before deadline, false if deadline reached
     */
    public boolean waitWithDeadline(int waitSeconds) {
        try {
            Calendar calendar = Calendar.getInstance();
            calendar.add(Calendar.SECOND, waitSeconds);
            
            lock.lock();
            System.out.println("Thread " + Thread.currentThread().getName() + 
                             " waiting until " + calendar.getTime());
            
            while (!isSignalled) {
                if (!condition.awaitUntil(calendar.getTime())) {
                    System.out.println("Thread " + Thread.currentThread().getName() + 
                                     " reached deadline");
                    return false;
                }
            }
            System.out.println("Thread " + Thread.currentThread().getName() + 
                             " was notified");
            return true;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            System.out.println("Thread " + Thread.currentThread().getName() + 
                             " was interrupted");
            return false;
        } finally {
            isSignalled = false;
            lock.unlock();
        }
    }

    /**
     * Notifies all waiting threads
     */
    public void notifyAllThreads() {
        try {
            lock.lock();
            isSignalled = true;
            condition.signalAll();
            System.out.println("Notified all waiting threads");
        } finally {
            lock.unlock();
        }
    }
}
