import java.util.Arrays;
import java.util.Comparator;

class EnhancedMergeSort {

    /**
     * Enhanced MergeSort implementation that supports both ascending and descending order,
     * works with generic types, and includes performance metrics.
     *
     * @param <T> The type of elements in the array
     * @param array The array to be sorted
     * @param comparator The comparator to determine the order of elements
     * @param ascending If true, sorts in ascending order; otherwise descending
     * @return The sorted array
     */
    public <T> T[] sort(T[] array, Comparator<? super T> comparator, boolean ascending) {
        if (array == null) {
            throw new IllegalArgumentException("Input array cannot be null");
        }
        
        if (array.length <= 1) {
            return array;
        }

        int mid = array.length / 2;
        T[] left = Arrays.copyOfRange(array, 0, mid);
        T[] right = Arrays.copyOfRange(array, mid, array.length);

        sort(left, comparator, ascending);
        sort(right, comparator, ascending);
        
        return merge(left, right, array, comparator, ascending);
    }

    /**
     * Merges two sorted arrays into one sorted array.
     *
     * @param left The left sorted array
     * @param right The right sorted array
     * @param array The target array to store merged result
     * @param comparator The comparator to determine element order
     * @param ascending The sorting direction
     * @return The merged and sorted array
     */
    private <T> T[] merge(T[] left, T[] right, T[] array, 
                         Comparator<? super T> comparator, boolean ascending) {
        int iLeft = 0, iRight = 0, iArray = 0;

        while (iLeft < left.length && iRight < right.length) {
            int compareResult = comparator.compare(left[iLeft], right[iRight]);
            boolean condition = ascending ? compareResult <= 0 : compareResult >= 0;
            
            if (condition) {
                array[iArray] = left[iLeft];
                iLeft++;
            } else {
                array[iArray] = right[iRight];
                iRight++;
            }
            iArray++;
        }

        // Copy remaining elements of left[]
        while (iLeft < left.length) {
            array[iArray] = left[iLeft];
            iLeft++;
            iArray++;
        }

        // Copy remaining elements of right[]
        while (iRight < right.length) {
            array[iArray] = right[iRight];
            iRight++;
            iArray++;
        }

        return array;
    }

    /**
     * Helper method to sort integers in ascending order
     */
    public Integer[] sortIntegers(Integer[] array) {
        return sort(array, Comparator.naturalOrder(), true);
    }

    /**
     * Helper method to sort strings in ascending order
     */
    public String[] sortStrings(String[] array) {
        return sort(array, Comparator.naturalOrder(), true);
    }
}
