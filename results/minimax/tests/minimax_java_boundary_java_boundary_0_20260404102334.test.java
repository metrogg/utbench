import org.junit.Before;
import org.junit.Test;
import org.mockito.Mockito;

import javax.management.ThreadMXBean;
import java.lang.management.ThreadInfo;
import java.lang.management.ManagementFactory;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;

import static org.junit.Assert.*;
import static org.mockito.Mockito.*;

public class ThreadDeadlockDetectorTest {

    private ThreadMXBean mockThreadMXBean;
    private ThreadDeadlockDetector detector;

    @Before
    public void setUp() {
        mockThreadMXBean = mock(ThreadMXBean.class);
        detector = new ThreadDeadlockDetector(mockThreadMXBean);
    }

    @Test
    public void returnsAnEmptySetIfNoThreadsAreDeadlocked() {
        when(mockThreadMXBean.findDeadlockedThreads()).thenReturn(null);

        Set<String> result = detector.getDeadlockedThreads();

        assertNotNull(result);
        assertTrue(result.isEmpty());
        assertEquals(Collections.emptySet(), result);
    }

    @Test
    public void returnsEmptySetWhenNoDeadlockedThreadsIdsReturned() {
        when(mockThreadMXBean.findDeadlockedThreads()).thenReturn(new long[0]);

        Set<String> result = detector.getDeadlockedThreads();

        assertNotNull(result);
        assertTrue(result.isEmpty());
    }

    @Test
    public void returnsDeadlockInfoForSingleDeadlockedThread() {
        long[] deadlockedIds = {1L};
        when(mockThreadMXBean.findDeadlockedThreads()).thenReturn(deadlockedIds);

        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        when(mockThreadInfo.getThreadName()).thenReturn("DeadlockThread");
        when(mockThreadInfo.getLockName()).thenReturn("java.lang.Object@1234");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("OwnerThread");
        when(mockThreadInfo.getStackTrace()).thenReturn(new StackTraceElement[0]);

        when(mockThreadMXBean.getThreadInfo(eq(deadlockedIds), anyInt())).thenReturn(new ThreadInfo[]{mockThreadInfo});

        Set<String> result = detector.getDeadlockedThreads();

        assertNotNull(result);
        assertEquals(1, result.size());
        String deadlockInfo = result.iterator().next();
        assertTrue(deadlockInfo.contains("DeadlockThread"));
        assertTrue(deadlockInfo.contains("java.lang.Object@1234"));
        assertTrue(deadlockInfo.contains("OwnerThread"));
    }

    @Test
    public void returnsDeadlockInfoWithStackTraceElements() {
        long[] deadlockedIds = {1L, 2L};
        when(mockThreadMXBean.findDeadlockedThreads()).thenReturn(deadlockedIds);

        StackTraceElement element1 = new StackTraceElement("ClassA", "method1", "FileA.java", 10);
        StackTraceElement element2 = new StackTraceElement("ClassB", "method2", "FileB.java", 20);

        ThreadInfo mockThreadInfo1 = mock(ThreadInfo.class);
        when(mockThreadInfo1.getThreadName()).thenReturn("Thread-1");
        when(mockThreadInfo1.getLockName()).thenReturn("lock1");
        when(mockThreadInfo1.getLockOwnerName()).thenReturn("Owner-1");
        when(mockThreadInfo1.getStackTrace()).thenReturn(new StackTraceElement[]{element1, element2});

        ThreadInfo mockThreadInfo2 = mock(ThreadInfo.class);
        when(mockThreadInfo2.getThreadName()).thenReturn("Thread-2");
        when(mockThreadInfo2.getLockName()).thenReturn("lock2");
        when(mockThreadInfo2.getLockOwnerName()).thenReturn("Owner-2");
        when(mockThreadInfo2.getStackTrace()).thenReturn(new StackTraceElement[]{element2});

        when(mockThreadMXBean.getThreadInfo(eq(deadlockedIds), anyInt())).thenReturn(new ThreadInfo[]{mockThreadInfo1, mockThreadInfo2});

        Set<String> result = detector.getDeadlockedThreads();

        assertNotNull(result);
        assertEquals(2, result.size());

        boolean foundThread1 = false;
        boolean foundThread2 = false;
        for (String info : result) {
            if (info.contains("Thread-1") && info.contains("ClassA.method1")) {
                foundThread1 = true;
            }
            if (info.contains("Thread-2") && info.contains("ClassB.method2")) {
                foundThread2 = true;
            }
        }
        assertTrue("Thread-1 info not found or incomplete", foundThread1);
        assertTrue("Thread-2 info not found or incomplete", foundThread2);
    }

    @Test
    public void returnsUnmodifiableSet() {
        when(mockThreadMXBean.findDeadlockedThreads()).thenReturn(null);

        Set<String> result = detector.getDeadlockedThreads();

        try {
            result.add("test");
            fail("Expected UnsupportedOperationException");
        } catch (UnsupportedOperationException e) {
            // Expected
        }
    }

    @Test
    public void defaultConstructorUsesSystemThreadMXBean() {
        ThreadDeadlockDetector defaultDetector = new ThreadDeadlockDetector();

        Set<String> result = defaultDetector.getDeadlockedThreads();

        assertNotNull(result);
    }

    @Test
    public void handlesNullInThreadInfoArray() {
        long[] deadlockedIds = {1L, 2L};
        when(mockThreadMXBean.findDeadlockedThreads()).thenReturn(deadlockedIds);

        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        when(mockThreadInfo.getThreadName()).thenReturn("ValidThread");
        when(mockThreadInfo.getLockName()).thenReturn("lock");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("owner");
        when(mockThreadInfo.getStackTrace()).thenReturn(new StackTraceElement[0]);

        when(mockThreadMXBean.getThreadInfo(eq(deadlockedIds), anyInt())).thenReturn(new ThreadInfo[]{mockThreadInfo, null});

        Set<String> result = detector.getDeadlockedThreads();

        assertNotNull(result);
        assertEquals(1, result.size());
        assertTrue(result.iterator().next().contains("ValidThread"));
    }

    @Test
    public void handlesMultipleDeadlockedThreads() {
        long[] deadlockedIds = {1L, 2L, 3L};
        when(mockThreadMXBean.findDeadlockedThreads()).thenReturn(deadlockedIds);

        ThreadInfo[] threadInfos = new ThreadInfo[3];
        for (int i = 0; i < 3; i++) {
            ThreadInfo info = mock(ThreadInfo.class);
            when(info.getThreadName()).thenReturn("Thread-" + i);
            when(info.getLockName()).thenReturn("lock" + i);
            when(info.getLockOwnerName()).thenReturn("owner" + i);
            when(info.getStackTrace()).thenReturn(new StackTraceElement[0]);
            threadInfos[i] = info;
        }

        when(mockThreadMXBean.getThreadInfo(eq(deadlockedIds), anyInt())).thenReturn(threadInfos);

        Set<String> result = detector.getDeadlockedThreads();

        assertNotNull(result);
        assertEquals(3, result.size());
    }

    @Test
    public void formatsStackTraceCorrectly() {
        long[] deadlockedIds = {1L};
        when(mockThreadMXBean.findDeadlockedThreads()).thenReturn(deadlockedIds);

        StackTraceElement element = new StackTraceElement("com.example.MyClass", "myMethod", "MyClass.java", 42);
        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        when(mockThreadInfo.getThreadName()).thenReturn("TestThread");
        when(mockThreadInfo.getLockName()).thenReturn("TestLock");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("TestOwner");
        when(mockThreadInfo.getStackTrace()).thenReturn(new StackTraceElement[]{element});

        when(mockThreadMXBean.getThreadInfo(eq(deadlockedIds), anyInt())).thenReturn(new ThreadInfo[]{mockThreadInfo});

        Set<String> result = detector.getDeadlockedThreads();

        assertEquals(1, result.size());
        String deadlockInfo = result.iterator().next();
        assertTrue(deadlockInfo.contains("\t at "));
        assertTrue(deadlockInfo.contains("com.example.MyClass.myMethod(MyClass.java:42)"));
    }
}