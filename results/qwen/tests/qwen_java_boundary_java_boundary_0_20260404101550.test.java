import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;

import java.lang.management.ThreadInfo;
import java.lang.management.ThreadMXBean;
import java.util.Set;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;
import static org.junit.Assert.fail;
import static org.mockito.ArgumentMatchers.anyInt;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@RunWith(MockitoJUnitRunner.class)
class ThreadDeadlockDetectorTest {

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

        assertTrue("Should return an empty set when findDeadlockedThreads returns null", result.isEmpty());
        assertUnmodifiable(result);
    }

    @Test
    public void returnsUnmodifiableEmptySetWhenDeadlockedIdsArrayIsEmpty() {
        when(threads.findDeadlockedThreads()).thenReturn(new long[0]);
        when(threads.getThreadInfo(new long[0], anyInt())).thenReturn(new ThreadInfo[0]);

        Set<String> result = detector.getDeadlockedThreads();

        assertTrue("Should return an empty set for empty deadlock IDs", result.isEmpty());
        assertUnmodifiable(result);
    }

    @Test
    public void returnsFormattedDeadlockDetailsWhenThreadsAreDeadlocked() {
        long[] deadlockedIds = new long[]{1L, 2L};
        when(threads.findDeadlockedThreads()).thenReturn(deadlockedIds);

        ThreadInfo info1 = mock(ThreadInfo.class);
        when(info1.getThreadName()).thenReturn("Thread-1");
        when(info1.getLockName()).thenReturn("Lock-A");
        when(info1.getLockOwnerName()).thenReturn("Thread-2");
        when(info1.getStackTrace()).thenReturn(new StackTraceElement[]{
                new StackTraceElement("com.example.ClassA", "methodA", "ClassA.java", 10)
        });

        ThreadInfo info2 = mock(ThreadInfo.class);
        when(info2.getThreadName()).thenReturn("Thread-2");
        when(info2.getLockName()).thenReturn("Lock-B");
        when(info2.getLockOwnerName()).thenReturn("Thread-1");
        when(info2.getStackTrace()).thenReturn(new StackTraceElement[]{
                new StackTraceElement("com.example.ClassB", "methodB", "ClassB.java", 20),
                new StackTraceElement("com.example.ClassC", "methodC", "ClassC.java", 30)
        });

        when(threads.getThreadInfo(deadlockedIds, anyInt())).thenReturn(new ThreadInfo[]{info1, info2});

        Set<String> result = detector.getDeadlockedThreads();

        assertEquals("Should contain details for both deadlocked threads", 2, result.size());

        String expected1 = String.format("Thread-1 locked on Lock-A (owned by Thread-2):%n\t at com.example.ClassA.methodA(ClassA.java:10)%n");
        String expected2 = String.format("Thread-2 locked on Lock-B (owned by Thread-1):%n\t at com.example.ClassB.methodB(ClassB.java:20)%n\t at com.example.ClassC.methodC(ClassC.java:30)%n");

        assertTrue("Should contain formatted string for Thread-1", result.contains(expected1));
        assertTrue("Should contain formatted string for Thread-2", result.contains(expected2));
        assertUnmodifiable(result);
    }

    private void assertUnmodifiable(Set<String> set) {
        try {
            set.add("should-fail");
            fail("Expected UnsupportedOperationException for unmodifiable set");
        } catch (UnsupportedOperationException e) {
            // Expected behavior
        }
    }
}