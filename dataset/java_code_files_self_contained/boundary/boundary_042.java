import java.util.HashMap;
import java.util.Map;

class AdvancedTileCalculator {

    /**
     * Calculates the tile configuration based on neighboring tiles.
     * This expanded version supports more complex tile configurations and includes
     * diagonal neighbors for more precise tile mapping.
     *
     * @param neighbors A map representing presence of neighbors in 8 directions:
     *                  "UP", "DOWN", "LEFT", "RIGHT", "UP_LEFT", "UP_RIGHT", "DOWN_LEFT", "DOWN_RIGHT"
     * @return A TileConfiguration object containing x and y coordinates for the tile sprite
     */
    public TileConfiguration calculateAdvancedTile(Map<String, Boolean> neighbors) {
        boolean up = neighbors.getOrDefault("UP", false);
        boolean down = neighbors.getOrDefault("DOWN", false);
        boolean left = neighbors.getOrDefault("LEFT", false);
        boolean right = neighbors.getOrDefault("RIGHT", false);
        boolean upLeft = neighbors.getOrDefault("UP_LEFT", false);
        boolean upRight = neighbors.getOrDefault("UP_RIGHT", false);
        boolean downLeft = neighbors.getOrDefault("DOWN_LEFT", false);
        boolean downRight = neighbors.getOrDefault("DOWN_RIGHT", false);

        // Handle all possible combinations with diagonal checks
        if (!up && !down && !left && !right) {
            return new TileConfiguration(0, 0); // Isolated tile
        }

        // Full surrounded tile
        if (up && down && left && right && upLeft && upRight && downLeft && downRight) {
            return new TileConfiguration(4, 2);
        }

        // Vertical passage
        if (up && down && !left && !right) {
            return new TileConfiguration(5, 2);
        }

        // Horizontal passage
        if (left && right && !up && !down) {
            return new TileConfiguration(1, 2);
        }

        // Dead ends
        if (down && !up && !left && !right) return new TileConfiguration(0, 2);
        if (right && !up && !left && !down) return new TileConfiguration(5, 0);
        if (left && !up && !right && !down) return new TileConfiguration(5, 1);
        if (up && !left && !right && !down) return new TileConfiguration(1, 0);

        // Corners with diagonal checks
        if (up && left && !right && !down && upLeft) return new TileConfiguration(3, 0);
        if (up && right && !left && !down && upRight) return new TileConfiguration(2, 0);
        if (down && left && !right && !up && downLeft) return new TileConfiguration(3, 1);
        if (down && right && !left && !up && downRight) return new TileConfiguration(2, 1);

        // T-junctions
        if (up && down && right && !left) return new TileConfiguration(2, 2);
        if (up && down && left && !right) return new TileConfiguration(3, 2);
        if (left && right && up && !down) return new TileConfiguration(4, 0);
        if (left && right && down && !up) return new TileConfiguration(4, 1);

        // Cross
        if (up && down && left && right) {
            // Check if it's a full cross or missing some diagonals
            if (!upLeft || !upRight || !downLeft || !downRight) {
                return new TileConfiguration(6, 0); // Special cross with missing diagonals
            }
            return new TileConfiguration(4, 2);
        }

        // Default case for unhandled configurations
        return new TileConfiguration(1, 1);
    }

    public static class TileConfiguration {
        private final int x;
        private final int y;

        public TileConfiguration(int x, int y) {
            this.x = x;
            this.y = y;
        }

        public int getX() {
            return x;
        }

        public int getY() {
            return y;
        }

        @Override
        public String toString() {
            return "TileConfiguration{x=" + x + ", y=" + y + "}";
        }
    }
}
