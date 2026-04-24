import java.util.List;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;

class TrianglePathAnalyzer {
    
    /**
     * Calculates the minimum path sum from top to bottom of a triangle (bottom-up approach).
     * Also provides additional analysis about the path.
     * 
     * @param triangle List of lists representing the triangle
     * @return Analysis result containing min path sum, max path sum, and path indices
     * @throws IllegalArgumentException if triangle is null or empty
     */
    public TriangleAnalysis analyzeTriangle(List<List<Integer>> triangle) {
        if (triangle == null || triangle.isEmpty()) {
            throw new IllegalArgumentException("Triangle cannot be null or empty");
        }
        
        // Create deep copies to avoid modifying original input
        List<List<Integer>> triangleCopy = deepCopy(triangle);
        List<List<Integer>> maxTriangle = deepCopy(triangle);
        
        // Bottom-up approach for minimum path
        for (int i = triangleCopy.size() - 2; i >= 0; i--) {
            for (int j = 0; j <= i; j++) {
                int min = Math.min(triangleCopy.get(i + 1).get(j), triangleCopy.get(i + 1).get(j + 1));
                triangleCopy.get(i).set(j, triangleCopy.get(i).get(j) + min);
            }
        }
        
        // Bottom-up approach for maximum path
        for (int i = maxTriangle.size() - 2; i >= 0; i--) {
            for (int j = 0; j <= i; j++) {
                int max = Math.max(maxTriangle.get(i + 1).get(j), maxTriangle.get(i + 1).get(j + 1));
                maxTriangle.get(i).set(j, maxTriangle.get(i).get(j) + max);
            }
        }
        
        // Reconstruct the minimum path indices
        List<Integer> minPathIndices = new ArrayList<>();
        minPathIndices.add(0);
        int currentIndex = 0;
        
        for (int i = 1; i < triangle.size(); i++) {
            if (triangleCopy.get(i).get(currentIndex) < triangleCopy.get(i).get(currentIndex + 1)) {
                minPathIndices.add(currentIndex);
            } else {
                minPathIndices.add(currentIndex + 1);
                currentIndex++;
            }
        }
        
        return new TriangleAnalysis(
            triangleCopy.get(0).get(0),
            maxTriangle.get(0).get(0),
            minPathIndices
        );
    }
    
    /**
     * Helper method to create a deep copy of the triangle
     */
    private List<List<Integer>> deepCopy(List<List<Integer>> triangle) {
        List<List<Integer>> copy = new ArrayList<>();
        for (List<Integer> row : triangle) {
            copy.add(new ArrayList<>(row));
        }
        return copy;
    }
    
    /**
     * Inner class to hold analysis results
     */
    public static class TriangleAnalysis {
        public final int minPathSum;
        public final int maxPathSum;
        public final List<Integer> minPathIndices;
        
        public TriangleAnalysis(int minPathSum, int maxPathSum, List<Integer> minPathIndices) {
            this.minPathSum = minPathSum;
            this.maxPathSum = maxPathSum;
            this.minPathIndices = Collections.unmodifiableList(minPathIndices);
        }
        
        @Override
        public String toString() {
            return String.format(
                "Min Path Sum: %d, Max Path Sum: %d, Min Path Indices: %s",
                minPathSum, maxPathSum, minPathIndices
            );
        }
    }
}
