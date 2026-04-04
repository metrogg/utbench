import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;
import java.lang.management.ThreadInfo;
import java.lang.management.ThreadMXBean;
import java.util.Set;
import static org.hamcrest.CoreMatchers.empty;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.junit.Assert.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@RunWith(MockitoJUnitRunner.class)
public class ThreadDeadlockDetectorTest {

    @Mock
    private ThreadMXBean threads;
    private ThreadDeadlockDetector detector;

    @Before
    public void setUp() {
        detector = new ThreadDeadlockDetector(threads);
    }

    @Test
    public void returnsAnEmptySetIfNoThreadsAreDeadlocked() {
        when(threads.findDeadlockedThreads()).thenReturn(null);
        Set<String> deadlockedThreads = detector.getDeadlockedThreads();
        assertThat(deadlockedThreads, is(empty()));
    }

    @Test
    public void returnsEmptySetWhenDeadlockedThreadsArrayIsEmpty() {
        when(threads.findDeadlockedThreads()).thenReturn(new long[0]);
        when(threads.getThreadInfo(any(long[].class), anyInt())).thenReturn(new ThreadInfo[0]);
        Set<String> deadlockedThreads = detector.getDeadlockedThreads();
        assertThat(deadlockedThreads, is(empty()));
    }

    @Test
    public void returnsCorrectDeadlockInfoWhenThreadsAreDeadlocked() {
        long[] deadlockIds = {1L, 2L};
        ThreadInfo thread1 = createMockThreadInfo("Thread-1", "lockA", "Thread-2", "com.example.ClassA.method1(ClassA.java:123)");
        ThreadInfo thread2 = createMockThreadInfo("Thread-2", "lockB", "Thread-1", "com.example.ClassB.method2(ClassB.java:456)");

        when(threads.findDeadlockedThreads()).thenReturn(deadlockIds);
        when(threads.getThreadInfo(eq(deadlockIds), anyInt())).thenReturn(new ThreadInfo[]{thread1, thread2});

        Set<String> result = detector.getDeadlockedThreads();
        assertEquals(2, result.size());
        assertTrue(result.contains("Thread-1 locked on lockA (owned by Thread-2):\n\t at com.example.ClassA.method1(ClassA.java:123)\n"));
        assertTrue(result.contains("Thread-2 locked on lockB (owned by Thread-1):\n\t at com.example.ClassB.method2(ClassB.java:456)\n"));
    }

    @Test
    public void handlesEmptyStackTraceCorrectly() {
        long[] deadlockId = {3L};
        ThreadInfo thread = createMockThreadInfo("Thread-3", "lockC", "Thread-4");

        when(threads.findDeadlockedThreads()).thenReturn(deadlockId);
        when(threads.getThreadInfo(eq(deadlockId), anyInt())).thenReturn(new ThreadInfo[]{thread});

        Set<String> result = detector.getDeadlockedThreads();
        assertEquals(1, result.size());
        String expectedEntry = "Thread-3 locked on lockC (owned by Thread-4):\n";
        assertTrue(result.contains(expectedEntry));
    }

    @Test(expected = UnsupportedOperationException.class)
    public void returnsUnmodifiableSetWhenDeadlocksExist() {
        long[] deadlockId = {1L};
        ThreadInfo info = createMockThreadInfo("t1", "l1", "t2");
        when(threads.findDeadlockedThreads()).thenReturn(deadlockId);
        when(threads.getThreadInfo(any(long[].class), anyInt())).thenReturn(new ThreadInfo[]{info});

        Set<String> result = detector.getDeadlockedThreads();
        result.add("invalid-entry");
    }

    @Test(expected = UnsupportedOperationException.class)
    public void returnsUnmodifiableSetWhenNoDeadlocks() {
        when(threads.findDeadlockedThreads()).thenReturn(null);
        Set<String> result = detector.getDeadlockedThreads();
        result.add("invalid-entry");
    }

    private ThreadInfo createMockThreadInfo(String threadName, String lockName, String lockOwnerName, String... stackTraceStrs) {
        ThreadInfo mockInfo = mock(ThreadInfo.class);
        when(mockInfo.getThreadName()).thenReturn(threadName);
        when(mockInfo.getLockName()).thenReturn(lockName);
        when(mockInfo.getLockOwnerName()).thenReturn(lockOwnerName);

        StackTraceElement[] stackTrace = new StackTraceElement[stackTraceStrs.length];
        for (int i = 0; i < stackTraceStrs.length; i++) {
            StackTraceElement element = mock(StackTraceElement.class);
            when(element.toString()).thenReturn(stackTraceStrs[i]);
            stackTrace[i] = element;
        }
        when(mockInfo.getStackTrace()).thenReturn(stackTrace);
        return mockInfo;
    }
}