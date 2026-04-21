// Converted Java method
import java.util.Arrays;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.Callable;

class AdvancedDotProductCalculator {
    private static final int DEFAULT_THREAD_COUNT = 4;
    
    /**
     * Calculates the dot product of two arrays using multiple threads.
     * The arrays are divided into chunks, and each chunk is processed by a separate thread.
     * 
     * @param array1 First array of integers
     * @param array2 Second array of integers
     * @param threadCount Number of threads to use (must be > 0)
     * @return The dot product of the two arrays
     * @throws IllegalArgumentException if arrays are null, different lengths, or threadCount <= 0
     */
    public static int calculateDotProduct(int[] array1, int[] array2, int threadCount) {
        if (array1 == null || array2 == null) {
            throw new IllegalArgumentException("Input arrays cannot be null");
        }
        if (array1.length != array2.length) {
            throw new IllegalArgumentException("Arrays must be of equal length");
        }
        if (threadCount <= 0) {
            throw new IllegalArgumentException("Thread count must be positive");
        }
        
        final int n = array1.length;
        if (n == 0) {
            return 0;
        }
        
        threadCount = Math.min(threadCount, n);
        final int chunkSize = (n + threadCount - 1) / threadCount;
        
        ExecutorService executor = Executors.newFixedThreadPool(threadCount);
        Future<Integer>[] partialResults = new Future[threadCount];
        
        for (int i = 0; i < threadCount; i++) {
            final int start = i * chunkSize;
            final int end = Math.min(start + chunkSize, n);
            
            partialResults[i] = executor.submit(new Callable<Integer>() {
                @Override
                public Integer call() {
                    int partialSum = 0;
                    for (int j = start; j < end; j++) {
                        partialSum += array1[j] * array2[j];
                    }
                    return partialSum;
                }
            });
        }
        
        int total = 0;
        try {
            for (Future<Integer> future : partialResults) {
                total += future.get();
            }
        } catch (Exception e) {
            throw new RuntimeException("Error during dot product calculation", e);
        } finally {
            executor.shutdown();
        }
        
        return total;
    }
    
    /**
     * Overloaded method that uses default thread count
     */
    public static int calculateDotProduct(int[] array1, int[] array2) {
        return calculateDotProduct(array1, array2, DEFAULT_THREAD_COUNT);
    }
}
