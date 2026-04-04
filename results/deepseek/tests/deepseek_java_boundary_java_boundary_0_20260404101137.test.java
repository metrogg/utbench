import org.junit.Test;
import org.junit.Before;
import org.junit.runner.RunWith;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
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
        MockitoAnnotations.initMocks(this);
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
        long[] deadlockedThreadIds = {123L};
        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        StackTraceElement[] stackTrace = {
            new StackTraceElement("com.example.Class", "method1", "Class.java", 10),
            new StackTraceElement("com.example.Class", "method2", "Class.java", 20)
        };
        
        when(mockThreads.findDeadlockedThreads()).thenReturn(deadlockedThreadIds);
        when(mockThreads.getThreadInfo(eq(deadlockedThreadIds), anyInt())).thenReturn(new ThreadInfo[]{mockThreadInfo});
        when(mockThreadInfo.getThreadName()).thenReturn("Thread-1");
        when(mockThreadInfo.getLockName()).thenReturn("Lock-A");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("Thread-2");
        when(mockThreadInfo.getStackTrace()).thenReturn(stackTrace);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        assertEquals(1, result.size());
        String deadlockInfo = result.iterator().next();
        assertTrue(deadlockInfo.contains("Thread-1"));
        assertTrue(deadlockInfo.contains("Lock-A"));
        assertTrue(deadlockInfo.contains("Thread-2"));
        assertTrue(deadlockInfo.contains("com.example.Class.method1"));
        assertTrue(deadlockInfo.contains("com.example.Class.method2"));
        
        verify(mockThreads).findDeadlockedThreads();
        verify(mockThreads).getThreadInfo(eq(deadlockedThreadIds), anyInt());
    }

    @Test
    public void returnsDeadlockInformationForMultipleThreads() {
        long[] deadlockedThreadIds = {123L, 456L, 789L};
        ThreadInfo mockThreadInfo1 = mock(ThreadInfo.class);
        ThreadInfo mockThreadInfo2 = mock(ThreadInfo.class);
        ThreadInfo mockThreadInfo3 = mock(ThreadInfo.class);
        
        StackTraceElement[] stackTrace1 = {new StackTraceElement("Class1", "method1", "Class1.java", 10)};
        StackTraceElement[] stackTrace2 = {new StackTraceElement("Class2", "method2", "Class2.java", 20)};
        StackTraceElement[] stackTrace3 = {new StackTraceElement("Class3", "method3", "Class3.java", 30)};
        
        when(mockThreads.findDeadlockedThreads()).thenReturn(deadlockedThreadIds);
        when(mockThreads.getThreadInfo(eq(deadlockedThreadIds), anyInt()))
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
        verify(mockThreads).getThreadInfo(eq(deadlockedThreadIds), anyInt());
    }

    @Test
    public void handlesThreadWithNullLockNameAndOwner() {
        long[] deadlockedThreadIds = {999L};
        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        StackTraceElement[] stackTrace = {new StackTraceElement("TestClass", "testMethod", "TestClass.java", 5)};
        
        when(mockThreads.findDeadlockedThreads()).thenReturn(deadlockedThreadIds);
        when(mockThreads.getThreadInfo(eq(deadlockedThreadIds), anyInt())).thenReturn(new ThreadInfo[]{mockThreadInfo});
        when(mockThreadInfo.getThreadName()).thenReturn("TestThread");
        when(mockThreadInfo.getLockName()).thenReturn(null);
        when(mockThreadInfo.getLockOwnerName()).thenReturn(null);
        when(mockThreadInfo.getStackTrace()).thenReturn(stackTrace);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        assertEquals(1, result.size());
        String deadlockInfo = result.iterator().next();
        assertTrue(deadlockInfo.contains("TestThread"));
        assertTrue(deadlockInfo.contains("null"));
        assertTrue(deadlockInfo.contains("TestClass.testMethod"));
        
        verify(mockThreads).findDeadlockedThreads();
        verify(mockThreads).getThreadInfo(eq(deadlockedThreadIds), anyInt());
    }

    @Test
    public void handlesEmptyStackTrace() {
        long[] deadlockedThreadIds = {111L};
        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        
        when(mockThreads.findDeadlockedThreads()).thenReturn(deadlockedThreadIds);
        when(mockThreads.getThreadInfo(eq(deadlockedThreadIds), anyInt())).thenReturn(new ThreadInfo[]{mockThreadInfo});
        when(mockThreadInfo.getThreadName()).thenReturn("EmptyStackTraceThread");
        when(mockThreadInfo.getLockName()).thenReturn("SomeLock");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("SomeOwner");
        when(mockThreadInfo.getStackTrace()).thenReturn(new StackTraceElement[0]);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        assertEquals(1, result.size());
        String deadlockInfo = result.iterator().next();
        assertTrue(deadlockInfo.contains("EmptyStackTraceThread"));
        assertTrue(deadlockInfo.contains("SomeLock"));
        assertTrue(deadlockInfo.contains("SomeOwner"));
        
        verify(mockThreads).findDeadlockedThreads();
        verify(mockThreads).getThreadInfo(eq(deadlockedThreadIds), anyInt());
    }

    @Test
    public void returnsUnmodifiableSet() {
        long[] deadlockedThreadIds = {123L};
        ThreadInfo mockThreadInfo = mock(ThreadInfo.class);
        StackTraceElement[] stackTrace = {new StackTraceElement("Test", "test", "Test.java", 1)};
        
        when(mockThreads.findDeadlockedThreads()).thenReturn(deadlockedThreadIds);
        when(mockThreads.getThreadInfo(eq(deadlockedThreadIds), anyInt())).thenReturn(new ThreadInfo[]{mockThreadInfo});
        when(mockThreadInfo.getThreadName()).thenReturn("TestThread");
        when(mockThreadInfo.getLockName()).thenReturn("TestLock");
        when(mockThreadInfo.getLockOwnerName()).thenReturn("TestOwner");
        when(mockThreadInfo.getStackTrace()).thenReturn(stackTrace);
        
        Set<String> result = detector.getDeadlockedThreads();
        
        assertNotNull(result);
        try {
            result.add("Should throw exception");
            fail("Expected UnsupportedOperationException when modifying unmodifiable set");
        } catch (UnsupportedOperationException e) {
            // Expected
        }
        
        verify(mockThreads).findDeadlockedThreads();
        verify(mockThreads).getThreadInfo(eq(deadlockedThreadIds), anyInt());
    }

    @Test
    public void testDefaultConstructor() {
        ThreadDeadlockDetector defaultDetector = new ThreadDeadlockDetector();
        assertNotNull(defaultDetector);
    }

    @Test
    public void testParameterizedConstructor() {
        ThreadDeadlockDetector paramDetector = new ThreadDeadlockDetector(mockThreads);
        assertNotNull(paramDetector);
    }
}