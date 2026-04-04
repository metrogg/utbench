import org.junit.Test;
import org.junit.Before;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.junit.MockitoJUnitRunner;
import java.lang.management.ThreadInfo;
import java.lang.management.ThreadMXBean;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;
import static org.junit.Assert.*;
import static org.mockito.Mockito.*;

@RunWith(MockitoJUnitRunner.class)
public class ThreadDeadlockDetectorTest {

    @Mock
    private ThreadMXBean mockThreads;

    private ThreadDeadlockDetector detector;

    @Before
    public void setUp() {
        detector = new ThreadDeadlockDetector(mockThreads);
    }

    @Test
    public void returnsAnEmptySetIfNoThreadsAreDeadlocked() {
        when(mockThreads.findDeadlockedThreads()).thenReturn(null);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        assertTrue(result.isEmpty());
        verify(mockThreads).findDeadlockedThreads();
    }

    @Test
    public void returnsEmptySetWhenDeadlockedThreadIdsIsEmptyArray() {
        when(mockThreads.findDeadlockedThreads()).thenReturn(new long[0]);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        assertTrue(result.isEmpty());
        verify(mockThreads).findDeadlockedThreads();
    }

    @Test
    public void returnsDeadlockInformationForSingleThread() {
        long[] deadlockedIds = {123L};
        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        StackTraceElement[] stackTrace = {
            new StackTraceElement("com.example.Class1", "method1", "Class1.java", 10),
            new StackTraceElement("com.example.Class2", "method2", "Class2.java", 20)
        };
        
        when(mockThreads.findDeadlockedThreads()).thenReturn(deadlockedIds);
        when(mockThreads.getThreadInfo(deadlockedIds, ThreadDeadlockDetector.MAX_STACK_TRACE_DEPTH))
            .thenReturn(new ThreadInfo[]{mockThreadInfo});
        when(mockThreadInfo.getThreadName()).thenReturn("Thread-1");
        when(mockThreadInfo.getLockName()).thenReturn("java.lang.Object@123456");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("Thread-2");
        when(mockThreadInfo.getStackTrace()).thenReturn(stackTrace);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        assertEquals(1, result.size());
        String deadlockInfo = result.iterator().next();
        assertTrue(deadlockInfo.contains("Thread-1"));
        assertTrue(deadlockInfo.contains("java.lang.Object@123456"));
        assertTrue(deadlockInfo.contains("Thread-2"));
        assertTrue(deadlockInfo.contains("com.example.Class1.method1"));
        assertTrue(deadlockInfo.contains("com.example.Class2.method2"));
        verify(mockThreads).findDeadlockedThreads();
        verify(mockThreads).getThreadInfo(deadlockedIds, ThreadDeadlockDetector.MAX_STACK_TRACE_DEPTH);
    }

    @Test
    public void returnsDeadlockInformationForMultipleThreads() {
        long[] deadlockedIds = {123L, 456L, 789L};
        ThreadInfo mockThreadInfo1 = mock(ThreadInfo.class);
        ThreadInfo mockThreadInfo2 = mock(ThreadInfo.class);
        ThreadInfo mockThreadInfo3 = mock(ThreadInfo.class);
        
        StackTraceElement[] stackTrace1 = {
            new StackTraceElement("com.example.ClassA", "methodA", "ClassA.java", 15)
        };
        StackTraceElement[] stackTrace2 = {
            new StackTraceElement("com.example.ClassB", "methodB", "ClassB.java", 25)
        };
        StackTraceElement[] stackTrace3 = {
            new StackTraceElement("com.example.ClassC", "methodC", "ClassC.java", 35)
        };
        
        when(mockThreads.findDeadlockedThreads()).thenReturn(deadlockedIds);
        when(mockThreads.getThreadInfo(deadlockedIds, ThreadDeadlockDetector.MAX_STACK_TRACE_DEPTH))
            .thenReturn(new ThreadInfo[]{mockThreadInfo1, mockThreadInfo2, mockThreadInfo3});
        
        when(mockThreadInfo1.getThreadName()).thenReturn("Thread-A");
        when(mockThreadInfo1.getLockName()).thenReturn("Lock-1");
        when(mockThreadInfo1.getLockOwnerName()).thenReturn("Thread-B");
        when(mockThreadInfo1.getStackTrace()).thenReturn(stackTrace1);
        
        when(mockThreadInfo2.getThreadName()).thenReturn("Thread-B");
        when(mockThreadInfo2.getLockName()).thenReturn("Lock-2");
        when(mockThreadInfo2.getLockOwnerName()).thenReturn("Thread-C");
        when(mockThreadInfo2.getStackTrace()).thenReturn(stackTrace2);
        
        when(mockThreadInfo3.getThreadName()).thenReturn("Thread-C");
        when(mockThreadInfo3.getLockName()).thenReturn("Lock-3");
        when(mockThreadInfo3.getLockOwnerName()).thenReturn("Thread-A");
        when(mockThreadInfo3.getStackTrace()).thenReturn(stackTrace3);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        assertEquals(3, result.size());
        verify(mockThreads).findDeadlockedThreads();
        verify(mockThreads).getThreadInfo(deadlockedIds, ThreadDeadlockDetector.MAX_STACK_TRACE_DEPTH);
    }

    @Test
    public void returnsUnmodifiableSet() {
        when(mockThreads.findDeadlockedThreads()).thenReturn(null);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        try {
            result.add("test");
            fail("Should throw UnsupportedOperationException");
        } catch (UnsupportedOperationException e) {
            // Expected
        }
    }

    @Test
    public void handlesThreadInfoWithNullStackTrace() {
        long[] deadlockedIds = {999L};
        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        
        when(mockThreads.findDeadlockedThreads()).thenReturn(deadlockedIds);
        when(mockThreads.getThreadInfo(deadlockedIds, ThreadDeadlockDetector.MAX_STACK_TRACE_DEPTH))
            .thenReturn(new ThreadInfo[]{mockThreadInfo});
        when(mockThreadInfo.getThreadName()).thenReturn("NullStackTraceThread");
        when(mockThreadInfo.getLockName()).thenReturn("SomeLock");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("OwnerThread");
        when(mockThreadInfo.getStackTrace()).thenReturn(null);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        assertEquals(1, result.size());
        String deadlockInfo = result.iterator().next();
        assertTrue(deadlockInfo.contains("NullStackTraceThread"));
        assertTrue(deadlockInfo.contains("SomeLock"));
        assertTrue(deadlockInfo.contains("OwnerThread"));
        verify(mockThreads).findDeadlockedThreads();
        verify(mockThreads).getThreadInfo(deadlockedIds, ThreadDeadlockDetector.MAX_STACK_TRACE_DEPTH);
    }

    @Test
    public void handlesThreadInfoWithEmptyStackTrace() {
        long[] deadlockedIds = {888L};
        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        
        when(mockThreads.findDeadlockedThreads()).thenReturn(deadlockedIds);
        when(mockThreads.getThreadInfo(deadlockedIds, ThreadDeadlockDetector.MAX_STACK_TRACE_DEPTH))
            .thenReturn(new ThreadInfo[]{mockThreadInfo});
        when(mockThreadInfo.getThreadName()).thenReturn("EmptyStackTraceThread");
        when(mockThreadInfo.getLockName()).thenReturn("EmptyLock");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("EmptyOwner");
        when(mockThreadInfo.getStackTrace()).thenReturn(new StackTraceElement[0]);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        assertEquals(1, result.size());
        String deadlockInfo = result.iterator().next();
        assertTrue(deadlockInfo.contains("EmptyStackTraceThread"));
        assertTrue(deadlockInfo.contains("EmptyLock"));
        assertTrue(deadlockInfo.contains("EmptyOwner"));
        verify(mockThreads).findDeadlockedThreads();
        verify(mockThreads).getThreadInfo(deadlockedIds, ThreadDeadlockDetector.MAX_STACK_TRACE_DEPTH);
    }

    @Test
    public void testDefaultConstructor() {
        ThreadDeadlockDetector defaultDetector = new ThreadDeadlockDetector();
        assertNotNull(defaultDetector);
    }

    @Test
    public void testConstructorWithThreadMXBean() {
        ThreadDeadlockDetector customDetector = new ThreadDeadlockDetector(mockThreads);
        assertNotNull(customDetector);
    }
}