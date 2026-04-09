import java.util.Arrays;
import java.util.NoSuchElementException;

/**
 * Enhanced BinaryHeap implementation with additional functionality:
 * - Supports both min-heap and max-heap operations
 * - Tracks vertex associations for graph algorithms
 * - Provides heap validation method
 * - Includes heap sort capability
 */
class EnhancedBinaryHeap {
    private static final int DEFAULT_CAPACITY = 100;
    private int currentSize;
    private Comparable[] array;
    private int[] vertex1;
    private int[] vertex2;
    private final boolean isMinHeap;

    /**
     * Constructor for min-heap by default
     */
    public EnhancedBinaryHeap() {
        this(true);
    }

    /**
     * Constructor with heap type specification
     * @param isMinHeap if true, creates a min-heap; otherwise creates a max-heap
     */
    public EnhancedBinaryHeap(boolean isMinHeap) {
        this.isMinHeap = isMinHeap;
        currentSize = 0;
        array = new Comparable[DEFAULT_CAPACITY + 1];
        vertex1 = new int[DEFAULT_CAPACITY + 1];
        vertex2 = new int[DEFAULT_CAPACITY + 1];
    }

    /**
     * Constructor with initial items
     * @param items the initial items in the binary heap
     * @param isMinHeap if true, creates a min-heap; otherwise creates a max-heap
     */
    public EnhancedBinaryHeap(Comparable[] items, boolean isMinHeap) {
        this.isMinHeap = isMinHeap;
        currentSize = items.length;
        array = new Comparable[items.length + 1];
        vertex1 = new int[items.length + 1];
        vertex2 = new int[items.length + 1];
        System.arraycopy(items, 0, array, 1, items.length);
        buildHeap();
    }

    /**
     * Inserts an element with associated vertices into the heap
     * @param x the item to insert
     * @param v1 first vertex
     * @param v2 second vertex
     */
    public void insert(Comparable x, int v1, int v2) {
        if (currentSize + 1 == array.length) {
            doubleArray();
        }

        int hole = ++currentSize;
        array[0] = x; // Sentinel for percolate up

        // Percolate up
        while (hole > 1 && compare(x, array[hole / 2]) < 0) {
            array[hole] = array[hole / 2];
            vertex1[hole] = vertex1[hole / 2];
            vertex2[hole] = vertex2[hole / 2];
            hole /= 2;
        }

        array[hole] = x;
        vertex1[hole] = v1;
        vertex2[hole] = v2;
    }

    /**
     * Removes and returns the root element (min or max depending on heap type)
     * @return the root element
     * @throws NoSuchElementException if heap is empty
     */
    public Comparable deleteRoot() {
        if (isEmpty()) {
            throw new NoSuchElementException("Heap is empty");
        }

        Comparable root = array[1];
        array[1] = array[currentSize];
        vertex1[1] = vertex1[currentSize];
        vertex2[1] = vertex2[currentSize--];
        percolateDown(1);
        return root;
    }

    /**
     * Returns the root element without removing it
     * @return the root element
     * @throws NoSuchElementException if heap is empty
     */
    public Comparable peekRoot() {
        if (isEmpty()) {
            throw new NoSuchElementException("Heap is empty");
        }
        return array[1];
    }

    /**
     * Returns the vertices associated with the root element
     * @return array containing the two vertices [v1, v2]
     * @throws NoSuchElementException if heap is empty
     */
    public int[] peekRootVertices() {
        if (isEmpty()) {
            throw new NoSuchElementException("Heap is empty");
        }
        return new int[]{vertex1[1], vertex2[1]};
    }

    /**
     * Sorts the elements in the heap and returns them in sorted order
     * Note: This operation destroys the heap property
     * @return array of sorted elements
     */
    public Comparable[] heapSort() {
        Comparable[] sorted = new Comparable[currentSize];
        for (int i = 0; i < sorted.length; i++) {
            sorted[i] = deleteRoot();
        }
        return sorted;
    }

    /**
     * Validates the heap property throughout the entire heap
     * @return true if the heap property is satisfied, false otherwise
     */
    public boolean validateHeap() {
        for (int i = 1; i <= currentSize / 2; i++) {
            int left = 2 * i;
            int right = 2 * i + 1;

            if (left <= currentSize && compare(array[i], array[left]) > 0) {
                return false;
            }
            if (right <= currentSize && compare(array[i], array[right]) > 0) {
                return false;
            }
        }
        return true;
    }

    public boolean isEmpty() {
        return currentSize == 0;
    }

    public int size() {
        return currentSize;
    }

    public void makeEmpty() {
        currentSize = 0;
    }

    private void buildHeap() {
        for (int i = currentSize / 2; i > 0; i--) {
            percolateDown(i);
        }
    }

    private void percolateDown(int hole) {
        int child;
        Comparable tmp = array[hole];
        int tmpV1 = vertex1[hole];
        int tmpV2 = vertex2[hole];

        while (hole * 2 <= currentSize) {
            child = hole * 2;
            if (child != currentSize && compare(array[child + 1], array[child]) < 0) {
                child++;
            }
            if (compare(array[child], tmp) < 0) {
                array[hole] = array[child];
                vertex1[hole] = vertex1[child];
                vertex2[hole] = vertex2[child];
                hole = child;
            } else {
                break;
            }
        }
        array[hole] = tmp;
        vertex1[hole] = tmpV1;
        vertex2[hole] = tmpV2;
    }

    private void doubleArray() {
        int newCapacity = array.length * 2;
        array = Arrays.copyOf(array, newCapacity);
        vertex1 = Arrays.copyOf(vertex1, newCapacity);
        vertex2 = Arrays.copyOf(vertex2, newCapacity);
    }

    private int compare(Comparable a, Comparable b) {
        return isMinHeap ? a.compareTo(b) : b.compareTo(a);
    }
}
