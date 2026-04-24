import java.util.Arrays;

class MatrixProcessor {
    
    /**
     * Processes a matrix by zeroing out rows and columns containing zeros, with additional features:
     * - Tracks the original positions of zeros
     * - Provides statistics about the transformation
     * - Supports both space-efficient and straightforward approaches
     * 
     * @param matrix The input matrix to process
     * @param useSpaceEfficient If true, uses O(1) space algorithm; otherwise uses O(n+m) space
     * @return An object containing the processed matrix and transformation statistics
     * @throws IllegalArgumentException if the input matrix is null or empty
     */
    public MatrixResult processMatrix(int[][] matrix, boolean useSpaceEfficient) {
        if (matrix == null || matrix.length == 0 || matrix[0].length == 0) {
            throw new IllegalArgumentException("Matrix cannot be null or empty");
        }

        int[][] originalMatrix = deepCopy(matrix);
        int zeroCount = 0;
        long startTime = System.nanoTime();

        if (useSpaceEfficient) {
            zeroMatrixSpaceEfficient(matrix);
        } else {
            zeroMatrixWithSets(matrix);
        }

        long endTime = System.nanoTime();
        double processingTimeMs = (endTime - startTime) / 1_000_000.0;

        // Count zeros in the original matrix
        for (int[] row : originalMatrix) {
            for (int val : row) {
                if (val == 0) zeroCount++;
            }
        }

        return new MatrixResult(
            originalMatrix,
            matrix,
            zeroCount,
            processingTimeMs,
            useSpaceEfficient ? "Space-Efficient" : "HashSet-Based"
        );
    }

    private void zeroMatrixWithSets(int[][] matrix) {
        int rows = matrix.length;
        int cols = matrix[0].length;
        boolean[] zeroRows = new boolean[rows];
        boolean[] zeroCols = new boolean[cols];

        // Mark rows and columns to be zeroed
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                if (matrix[i][j] == 0) {
                    zeroRows[i] = true;
                    zeroCols[j] = true;
                }
            }
        }

        // Zero out marked rows
        for (int i = 0; i < rows; i++) {
            if (zeroRows[i]) {
                Arrays.fill(matrix[i], 0);
            }
        }

        // Zero out marked columns
        for (int j = 0; j < cols; j++) {
            if (zeroCols[j]) {
                for (int i = 0; i < rows; i++) {
                    matrix[i][j] = 0;
                }
            }
        }
    }

    private void zeroMatrixSpaceEfficient(int[][] matrix) {
        int rows = matrix.length;
        int cols = matrix[0].length;
        boolean firstRowHasZero = false;
        boolean firstColHasZero = false;

        // Check if first row has zero
        for (int j = 0; j < cols; j++) {
            if (matrix[0][j] == 0) {
                firstRowHasZero = true;
                break;
            }
        }

        // Check if first column has zero
        for (int i = 0; i < rows; i++) {
            if (matrix[i][0] == 0) {
                firstColHasZero = true;
                break;
            }
        }

        // Use first row and column as markers
        for (int i = 1; i < rows; i++) {
            for (int j = 1; j < cols; j++) {
                if (matrix[i][j] == 0) {
                    matrix[i][0] = 0;
                    matrix[0][j] = 0;
                }
            }
        }

        // Zero out cells based on markers
        for (int i = 1; i < rows; i++) {
            for (int j = 1; j < cols; j++) {
                if (matrix[i][0] == 0 || matrix[0][j] == 0) {
                    matrix[i][j] = 0;
                }
            }
        }

        // Zero out first row if needed
        if (firstRowHasZero) {
            Arrays.fill(matrix[0], 0);
        }

        // Zero out first column if needed
        if (firstColHasZero) {
            for (int i = 0; i < rows; i++) {
                matrix[i][0] = 0;
            }
        }
    }

    private int[][] deepCopy(int[][] matrix) {
        return Arrays.stream(matrix).map(int[]::clone).toArray(int[][]::new);
    }

    static class MatrixResult {
        final int[][] originalMatrix;
        final int[][] processedMatrix;
        final int zeroCount;
        final double processingTimeMs;
        final String methodUsed;

        public MatrixResult(int[][] originalMatrix, int[][] processedMatrix, 
                          int zeroCount, double processingTimeMs, String methodUsed) {
            this.originalMatrix = originalMatrix;
            this.processedMatrix = processedMatrix;
            this.zeroCount = zeroCount;
            this.processingTimeMs = processingTimeMs;
            this.methodUsed = methodUsed;
        }
    }
}
