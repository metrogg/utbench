import java.util.Arrays;
import java.util.List;

class GameOfLifePatterns {
    private static final int GRID_SIZE = 20;

    /**
     * Applies a predefined pattern to a Game of Life grid
     * @param grid The 2D boolean array representing the game grid
     * @param patternName Name of the pattern to apply
     * @param xOffset Horizontal offset for pattern placement
     * @param yOffset Vertical offset for pattern placement
     * @return The modified grid with the pattern applied
     * @throws IllegalArgumentException if pattern name is invalid or offsets are out of bounds
     */
    public static boolean[][] applyPattern(boolean[][] grid, String patternName, int xOffset, int yOffset) {
        if (grid == null || grid.length == 0 || grid[0].length == 0) {
            throw new IllegalArgumentException("Grid must be initialized");
        }
        
        if (xOffset < 0 || yOffset < 0 || xOffset >= grid.length || yOffset >= grid[0].length) {
            throw new IllegalArgumentException("Invalid pattern offset");
        }

        List<Coordinate> patternCells;
        switch (patternName.toLowerCase()) {
            case "glider":
                patternCells = Arrays.asList(
                    new Coordinate(0, 0), new Coordinate(1, 0), new Coordinate(2, 0),
                    new Coordinate(2, 1), new Coordinate(1, 2)
                );
                break;
            case "small exploder":
                patternCells = Arrays.asList(
                    new Coordinate(0, 1), new Coordinate(1, 0), new Coordinate(1, 1),
                    new Coordinate(1, 2), new Coordinate(2, 0), new Coordinate(2, 2),
                    new Coordinate(3, 1)
                );
                break;
            case "exploder":
                patternCells = Arrays.asList(
                    new Coordinate(0, 0), new Coordinate(1, 0), new Coordinate(2, 0),
                    new Coordinate(3, 0), new Coordinate(4, 0), new Coordinate(0, 2),
                    new Coordinate(4, 2), new Coordinate(0, 4), new Coordinate(1, 4),
                    new Coordinate(2, 4), new Coordinate(3, 4), new Coordinate(4, 4)
                );
                break;
            case "ten cell row":
                patternCells = Arrays.asList(
                    new Coordinate(0, 0), new Coordinate(0, 1), new Coordinate(0, 2),
                    new Coordinate(0, 3), new Coordinate(0, 4), new Coordinate(0, 5),
                    new Coordinate(0, 6), new Coordinate(0, 7), new Coordinate(0, 8),
                    new Coordinate(0, 9)
                );
                break;
            case "lightweight spaceship":
                patternCells = Arrays.asList(
                    new Coordinate(0, 0), new Coordinate(2, 0), new Coordinate(3, 1),
                    new Coordinate(3, 2), new Coordinate(3, 3), new Coordinate(3, 4),
                    new Coordinate(2, 4), new Coordinate(1, 4), new Coordinate(0, 3)
                );
                break;
            case "tumbler":
                patternCells = Arrays.asList(
                    new Coordinate(0, 0), new Coordinate(1, 0), new Coordinate(2, 0),
                    new Coordinate(3, 1), new Coordinate(4, 1), new Coordinate(0, 2),
                    new Coordinate(1, 2), new Coordinate(2, 2), new Coordinate(3, 2),
                    new Coordinate(4, 2), new Coordinate(0, 4), new Coordinate(1, 4),
                    new Coordinate(2, 4), new Coordinate(3, 3), new Coordinate(4, 3),
                    new Coordinate(0, 5), new Coordinate(1, 5), new Coordinate(4, 5),
                    new Coordinate(2, 6), new Coordinate(3, 6), new Coordinate(4, 6)
                );
                break;
            default:
                throw new IllegalArgumentException("Unknown pattern: " + patternName);
        }

        for (Coordinate cell : patternCells) {
            int x = xOffset + cell.x;
            int y = yOffset + cell.y;
            if (x >= 0 && x < grid.length && y >= 0 && y < grid[0].length) {
                grid[x][y] = true;
            }
        }

        return grid;
    }

    private static class Coordinate {
        int x;
        int y;
        Coordinate(int x, int y) {
            this.x = x;
            this.y = y;
        }
    }

    /**
     * Creates a new empty grid of the specified size
     * @param size The size of the grid (size x size)
     * @return A new empty grid
     */
    public static boolean[][] createGrid(int size) {
        return new boolean[size][size];
    }

    /**
     * Counts the number of live cells in a grid
     * @param grid The grid to count cells in
     * @return The count of live cells
     */
    public static int countLiveCells(boolean[][] grid) {
        int count = 0;
        for (boolean[] row : grid) {
            for (boolean cell : row) {
                if (cell) count++;
            }
        }
        return count;
    }
}
