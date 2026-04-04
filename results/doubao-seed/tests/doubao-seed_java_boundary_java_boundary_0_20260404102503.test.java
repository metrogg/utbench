import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;
import java.lang.management.ThreadInfo;
import java.lang.management.ThreadMXBean;
import java.util.Arrays;
import java.util.Set;
import static org.junit.Assert.*;
import static org.mockito.ArgumentMatchers.anyInt;
import static org.mockito.ArgumentMatchers.eq;
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
        Set<String> result = detector.getDeadlockedThreads();
        assertTrue(result.isEmpty());
    }

    @Test
    public void returnsEmptySetWhenDeadlockedThreadArrayIsEmpty() {
        when(threads.findDeadlockedThreads()).thenReturn(new long[0]);
        Set<String> result = detector.getDeadlockedThreads();
        assertTrue(result.isEmpty());
    }

    @Test
    public void returnsFormattedStringForSingleDeadlockedThreadWithStackTrace() {
        long threadId = 123L;
        when(threads.findDeadlockedThreads()).thenReturn(new long[]{threadId});
        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        when(mockThreadInfo.getThreadName()).thenReturn("TestThread-1");
        when(mockThreadInfo.getLockName()).thenReturn("lock-a");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("TestThread-2");
        StackTraceElement stackElement = new StackTraceElement("com.example.TestClass", "testMethod", "TestClass.java", 45);
        when(mockThreadInfo.getStackTrace()).thenReturn(new StackTraceElement[]{stackElement});
        when(threads.getThreadInfo(eq(new long[]{threadId}), anyInt())).thenReturn(new ThreadInfo[]{mockThreadInfo});

        Set<String> result = detector.getDeadlockedThreads();
        assertEquals(1, result.size());
        String expected = String.format("TestThread-1 locked on lock-a (owned by TestThread-2):%n\t at com.example.TestClass.testMethod(TestClass.java:45)%n");
        assertEquals(expected, result.iterator().next());
    }

    @Test
    public void returnsMultipleFormattedEntriesForMultipleDeadlockedThreads() {
        long id1 = 1L, id2 = 2L;
        when(threads.findDeadlockedThreads()).thenReturn(new long[]{id1, id2});

        ThreadInfo info1 = mock(ThreadInfo.class);
        when(info1.getThreadName()).thenReturn("Thread-1");
        when(info1.getLockName()).thenReturn("lock2");
        when(info1.getLockOwnerName()).thenReturn("Thread-2");
        when(info1.getStackTrace()).thenReturn(new StackTraceElement[0]);

        ThreadInfo info2 = mock(ThreadInfo.class);
        when(info2.getThreadName()).thenReturn("Thread-2");
        when(info2.getLockName()).thenReturn("lock1");
        when(info2.getLockOwnerName()).thenReturn("Thread-1");
        when(info2.getStackTrace()).thenReturn(new StackTraceElement[0]);

        when(threads.getThreadInfo(eq(new long[]{id1, id2}), anyInt())).thenReturn(new ThreadInfo[]{info1, info2});

        Set<String> result = detector.getDeadlockedThreads();
        assertEquals(2, result.size());
        String expected1 = String.format("Thread-1 locked on lock2 (owned by Thread-2):%n");
        String expected2 = String.format("Thread-2 locked on lock1 (owned by Thread-1):%n");
        assertTrue(result.containsAll(Arrays.asList(expected1, expected2)));
    }

    @Test
    public void correctlyFormatsMultiLineStackTraceForDeadlockedThread() {
        long id = 10L;
        when(threads.findDeadlockedThreads()).thenReturn(new long[]{id});
        ThreadInfo info = mock(ThreadInfo.class);
        when(info.getThreadName()).thenReturn("WorkerThread");
        when(info.getLockName()).thenReturn("db-lock");
        when(info.getLockOwnerName()).thenReturn("MainThread");
        StackTraceElement elem1 = new StackTraceElement("com.example.DbUtil", "getConnection", "DbUtil.java", 12);
        StackTraceElement elem2 = new StackTraceElement("com.example.UserService", "getUser", "UserService.java", 89);
        when(info.getStackTrace()).thenReturn(new StackTraceElement[]{elem1, elem2});
        when(threads.getThreadInfo(eq(new long[]{id}), anyInt())).thenReturn(new ThreadInfo[]{info});

        Set<String> result = detector.getDeadlockedThreads();
        assertEquals(1, result.size());
        String expected = String.format("WorkerThread locked on db-lock (owned by MainThread):%n\t at com.example.DbUtil.getConnection(DbUtil.java:12)%n\t at com.example.UserService.getUser(UserService.java:89)%n");
        assertEquals(expected, result.iterator().next());
    }

    @Test
    public void returnsUnmodifiableSetWhenDeadlocksExist() {
        long threadId = 1L;
        when(threads.findDeadlockedThreads()).thenReturn(new long[]{threadId});
        ThreadInfo mockInfo = mock(ThreadInfo.class);
        when(mockInfo.getThreadName()).thenReturn("Test");
        when(mockInfo.getLockName()).thenReturn("lock");
        when(mockInfo.getLockOwnerName()).thenReturn("Owner");
        when(mockInfo.getStackTrace()).thenReturn(new StackTraceElement[0]);
        when(threads.getThreadInfo(eq(new long[]{threadId}), anyInt())).thenReturn(new ThreadInfo[]{mockInfo});

        Set<String> result = detector.getDeadlockedThreads();
        assertThrows(UnsupportedOperationException.class, () -> result.add("test"));
    }
}