import java.util.Arrays;
import java.util.List;

class GridCycleAnalyzer {
    private char[][] grid;
    private boolean[][] visited;
    private int[][] distance;
    private int rows;
    private int cols;
    private final int[] dx = {0, 0, -1, 1, -1, -1, 1, 1}; // Added diagonal directions
    private final int[] dy = {1, -1, 0, 0, -1, 1, -1, 1}; // Added diagonal directions

    /**
     * Analyzes a grid to find cycles of the same character with optional diagonal movement.
     * 
     * @param grid The 2D character grid to analyze
     * @param allowDiagonals Whether to consider diagonal movements in cycle detection
     * @return A CycleResult object containing cycle information or null if no cycle found
     */
    public CycleResult findCycle(char[][] grid, boolean allowDiagonals) {
        if (grid == null || grid.length == 0 || grid[0].length == 0) {
            return null;
        }

        this.grid = grid;
        this.rows = grid.length;
        this.cols = grid[0].length;
        this.visited = new boolean[rows][cols];
        this.distance = new int[rows][cols];

        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                if (visited[i][j]) continue;
                
                distance = new int[rows][cols];
                CycleResult result = dfsFindCycle(i, j, 1, grid[i][j], allowDiagonals);
                if (result != null) {
                    return result;
                }
            }
        }
        return null;
    }

    private CycleResult dfsFindCycle(int x, int y, int cnt, char color, boolean allowDiagonals) {
        if (visited[x][y]) {
            if (cnt - distance[x][y] >= 4) {
                return new CycleResult(true, color, cnt - distance[x][y], x, y);
            } else {
                return null;
            }
        }

        visited[x][y] = true;
        distance[x][y] = cnt;

        int directions = allowDiagonals ? 8 : 4;
        for (int k = 0; k < directions; k++) {
            int nx = x + dx[k];
            int ny = y + dy[k];
            
            if (0 <= nx && nx < rows && 0 <= ny && ny < cols) {
                if (grid[nx][ny] == color) {
                    CycleResult result = dfsFindCycle(nx, ny, cnt + 1, color, allowDiagonals);
                    if (result != null) {
                        return result;
                    }
                }
            }
        }
        return null;
    }

    /**
     * Inner class to store cycle detection results
     */
    public static class CycleResult {
        public final boolean hasCycle;
        public final char cycleColor;
        public final int cycleLength;
        public final int startX;
        public final int startY;

        public CycleResult(boolean hasCycle, char cycleColor, int cycleLength, int startX, int startY) {
            this.hasCycle = hasCycle;
            this.cycleColor = cycleColor;
            this.cycleLength = cycleLength;
            this.startX = startX;
            this.startY = startY;
        }

        @Override
        public String toString() {
            return hasCycle ? 
                String.format("Cycle found: color='%c', length=%d, start at (%d,%d)", 
                    cycleColor, cycleLength, startX, startY) :
                "No cycle found";
        }
    }
}
