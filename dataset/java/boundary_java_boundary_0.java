class ThreadDeadlockDetector {

    public Set<String> getDeadlockedThreads() {
        final long[] ids = threads.findDeadlockedThreads();
        if (ids != null) {
            final Set<String> deadlocks = new HashSet<>();
            for (ThreadInfo info : threads.getThreadInfo(ids, MAX_STACK_TRACE_DEPTH)) {
                final StringBuilder stackTrace = new StringBuilder();
                for (StackTraceElement element : info.getStackTrace()) {
                    stackTrace.append("\t at ")
                            .append(element.toString())
                            .append(String.format("%n"));
                }

                deadlocks.add(
                        String.format("%s locked on %s (owned by %s):%n%s",
                                info.getThreadName(),
                                info.getLockName(),
                                info.getLockOwnerName(),
                                stackTrace.toString()
                        )
                );
            }
            return Collections.unmodifiableSet(deadlocks);
        }
        return Collections.emptySet();
    }

    public  ThreadDeadlockDetector();
    public  ThreadDeadlockDetector(ThreadMXBean threads);

}

class ThreadDeadlockDetectorTest {

    private final ThreadMXBean threads;
    private final ThreadDeadlockDetector detector;

    @Test
    public void returnsAnEmptySetIfNoThreadsAreDeadlocked() {
