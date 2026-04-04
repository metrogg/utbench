import java.lang.management.ThreadInfo;
import java.lang.management.ThreadMXBean;
import java.util.Set;

import org.junit.Before;
import org.junit.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import static org.junit.Assert.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyInt;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

public class ThreadDeadlockDetectorTest {

    @Mock
    private ThreadMXBean threads;

    private ThreadDeadlockDetector detector;

    @Before
    public void setUp() {
        MockitoAnnotations.initMocks(this);
        detector = new ThreadDeadlockDetector(threads);
    }

    @Test
    public void returnsAnEmptySetIfNoThreadsAreDeadlocked() {
        when(threads.findDeadlockedThreads()).thenReturn(null);

        Set<String> result = detector.getDeadlockedThreads();

        assertTrue("Result should be an empty set", result.isEmpty());
    }

    @Test
    public void returnsAnEmptySetIfDeadlockedThreadIdsArrayIsEmpty() {
        when(threads.findDeadlockedThreads()).thenReturn(new long[0]);

        Set<String> result = detector.getDeadlockedThreads();

        assertTrue("Result should be an empty set", result.isEmpty());
    }

    @Test
    public void returnsFormattedStringForSingleDeadlockedThread() {
        long threadId = 1L;
        String threadName = "Deadlocked-Thread-1";
        String lockName = "java.lang.Object@12345";
        String lockOwner = "Owner-Thread";
        StackTraceElement stackElement = new StackTraceElement("com.example.Service", "process", "Service.java", 42);

        ThreadInfo mockInfo = mock(ThreadInfo.class);
        when(mockInfo.getThreadName()).thenReturn(threadName);
        when(mockInfo.getLockName()).thenReturn(lockName);
        when(mockInfo.getLockOwnerName()).thenReturn(lockOwner);
        when(mockInfo.getStackTrace()).thenReturn(new StackTraceElement[]{stackElement});

        when(threads.findDeadlockedThreads()).thenReturn(new long[]{threadId});
        when(threads.getThreadInfo(eq(new long[]{threadId}), anyInt())).thenReturn(new ThreadInfo[]{mockInfo});

        Set<String> result = detector.getDeadlockedThreads();

        assertEquals(1, result.size());
        String deadlockInfo = result.iterator().next();
        assertTrue(deadlockInfo.contains(threadName));
        assertTrue(deadlockInfo.contains(lockName));
        assertTrue(deadlockInfo.contains(lockOwner));
        assertTrue(deadlockInfo.contains("com.example.Service.process"));
    }

    @Test
    public void returnsMultipleEntriesForMultipleDeadlockedThreads() {
        long id1 = 1L;
        long id2 = 2L;

        ThreadInfo info1 = mock(ThreadInfo.class);
        when(info1.getThreadName()).thenReturn("Thread-1");
        when(info1.getLockName()).thenReturn("Lock-1");
        when(info1.getLockOwnerName()).thenReturn("Owner-1");
        when(info1.getStackTrace()).thenReturn(new StackTraceElement[0]);

        ThreadInfo info2 = mock(ThreadInfo.class);
        when(info2.getThreadName()).thenReturn("Thread-2");
        when(info2.getLockName()).thenReturn("Lock-2");
        when(info2.getLockOwnerName()).thenReturn("Owner-2");
        when(info2.getStackTrace()).thenReturn(new StackTraceElement[0]);

        when(threads.findDeadlockedThreads()).thenReturn(new long[]{id1, id2});
        when(threads.getThreadInfo(any(long[].class), anyInt())).thenReturn(new ThreadInfo[]{info1, info2});

        Set<String> result = detector.getDeadlockedThreads();

        assertEquals(2, result.size());
    }

    @Test(expected = NullPointerException.class)
    public void throwsNullPointerExceptionIfThreadInfoArrayIsNull() {
        long id = 1L;
        when(threads.findDeadlockedThreads()).thenReturn(new long[]{id});
        when(threads.getThreadInfo(any(long[].class), anyInt())).thenReturn(null);

        detector.getDeadlockedThreads();
    }

    @Test(expected = NullPointerException.class)
    public void throwsNullPointerExceptionIfThreadInfoArrayContainsNullElements() {
        long id = 1L;
        when(threads.findDeadlockedThreads()).thenReturn(new long[]{id});
        when(threads.getThreadInfo(any(long[].class), anyInt())).thenReturn(new ThreadInfo[]{null});

        detector.getDeadlockedThreads();
    }
}