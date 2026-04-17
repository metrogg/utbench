import java.util.ArrayList;
import java.util.List;

class PathFinder {
    /**
     * Finds all unique paths from top-left to bottom-right in a grid with obstacles,
     * and provides additional path analysis information.
     * 
     * @param obstacleGrid 2D array representing the grid (1 = obstacle, 0 = free space)
     * @return PathAnalysis object containing count of paths and other statistics
     * @throws IllegalArgumentException if grid is empty or null
     */
    public PathAnalysis findUniquePathsWithAnalysis(int[][] obstacleGrid) {
        if (obstacleGrid == null || obstacleGrid.length == 0 || obstacleGrid[0].length == 0) {
            throw new IllegalArgumentException("Grid cannot be empty or null");
        }

        int m = obstacleGrid.length;
        int n = obstacleGrid[0].length;
        
        // Check if start or end is blocked
        if (obstacleGrid[0][0] == 1 || obstacleGrid[m-1][n-1] == 1) {
            return new PathAnalysis(0, 0, 0, new ArrayList<>());
        }

        int[][] dp = new int[m][n];
        dp[0][0] = 1;

        // Initialize first column
        for (int i = 1; i < m; i++) {
            dp[i][0] = (obstacleGrid[i][0] == 0) ? dp[i-1][0] : 0;
        }

        // Initialize first row
        for (int j = 1; j < n; j++) {
            dp[0][j] = (obstacleGrid[0][j] == 0) ? dp[0][j-1] : 0;
        }

        // Fill DP table
        for (int i = 1; i < m; i++) {
            for (int j = 1; j < n; j++) {
                dp[i][j] = (obstacleGrid[i][j] == 0) ? dp[i-1][j] + dp[i][j-1] : 0;
            }
        }

        int totalPaths = dp[m-1][n-1];
        int maxPathLength = m + n - 2;
        int minPathLength = maxPathLength; // In grid without diagonal moves, all paths have same length
        
        // For demonstration, we'll find one sample path (not all paths due to complexity)
        List<String> samplePath = findSamplePath(obstacleGrid, dp);
        
        return new PathAnalysis(totalPaths, maxPathLength, minPathLength, samplePath);
    }

    /**
     * Helper method to find one sample path (for demonstration)
     */
    private List<String> findSamplePath(int[][] grid, int[][] dp) {
        List<String> path = new ArrayList<>();
        int i = grid.length - 1;
        int j = grid[0].length - 1;
        
        if (dp[i][j] == 0) {
            return path;
        }

        path.add("(" + i + "," + j + ")");
        
        while (i > 0 || j > 0) {
            if (i > 0 && dp[i-1][j] > 0) {
                i--;
            } else if (j > 0 && dp[i][j-1] > 0) {
                j--;
            }
            path.add(0, "(" + i + "," + j + ")");
        }
        
        return path;
    }

    /**
     * Inner class to hold path analysis results
     */
    public static class PathAnalysis {
        public final int totalPaths;
        public final int maxPathLength;
        public final int minPathLength;
        public final List<String> samplePath;

        public PathAnalysis(int totalPaths, int maxPathLength, int minPathLength, List<String> samplePath) {
            this.totalPaths = totalPaths;
            this.maxPathLength = maxPathLength;
            this.minPathLength = minPathLength;
            this.samplePath = samplePath;
        }

        @Override
        public String toString() {
            return "Total Paths: " + totalPaths + "\n" +
                   "Path Length: " + maxPathLength + "\n" +
                   "Sample Path: " + samplePath;
        }
    }
}
