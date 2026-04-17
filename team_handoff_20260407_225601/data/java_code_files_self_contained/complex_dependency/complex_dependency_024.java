// Converted Java method
import java.util.Observable;
import java.util.List;
import java.util.ArrayList;

/**
 * Enhanced naval battle game adapter that manages game state changes and notifications.
 * Tracks ship placements, attacks, and game outcomes with detailed observation capabilities.
 */
class NavalBattleAdapter extends Observable {
    private List<String> shipPositions;
    private List<String> attackPositions;
    private boolean gameOver;
    private String winner;

    public NavalBattleAdapter() {
        this.shipPositions = new ArrayList<>();
        this.attackPositions = new ArrayList<>();
        this.gameOver = false;
        this.winner = null;
    }

    /**
     * Adds a ship to the game board at specified coordinates
     * @param coordinates List of coordinates where the ship is placed
     */
    public void placeShips(List<String> coordinates) {
        if (coordinates == null || coordinates.isEmpty()) {
            throw new IllegalArgumentException("Coordinates cannot be null or empty");
        }
        
        shipPositions.addAll(coordinates);
        notifyStateChange("SHIPS_PLACED", coordinates);
    }

    /**
     * Processes an attack at the specified coordinate
     * @param coordinate The coordinate being attacked (e.g., "A5")
     * @return HitResult enum indicating HIT, MISS, or SUNK
     */
    public HitResult receiveAttack(String coordinate) {
        if (coordinate == null || coordinate.isEmpty()) {
            throw new IllegalArgumentException("Coordinate cannot be null or empty");
        }

        attackPositions.add(coordinate);
        
        if (shipPositions.contains(coordinate)) {
            shipPositions.remove(coordinate);
            
            if (shipPositions.isEmpty()) {
                gameOver = true;
                winner = "Opponent";
                notifyStateChange("GAME_OVER", winner);
                return HitResult.SUNK;
            }
            notifyStateChange("HIT", coordinate);
            return HitResult.HIT;
        }
        
        notifyStateChange("MISS", coordinate);
        return HitResult.MISS;
    }

    /**
     * Gets the current game state summary
     * @return GameState object containing ships remaining, attacks made, and game status
     */
    public GameState getGameState() {
        return new GameState(
            new ArrayList<>(shipPositions),
            new ArrayList<>(attackPositions),
            gameOver,
            winner
        );
    }

    private void notifyStateChange(String changeType, Object data) {
        setChanged();
        notifyObservers(new StateChange(changeType, data));
    }

    public enum HitResult {
        HIT, MISS, SUNK
    }

    public static class StateChange {
        public final String changeType;
        public final Object data;

        public StateChange(String changeType, Object data) {
            this.changeType = changeType;
            this.data = data;
        }
    }

    public static class GameState {
        public final List<String> remainingShips;
        public final List<String> attacksMade;
        public final boolean gameOver;
        public final String winner;

        public GameState(List<String> remainingShips, List<String> attacksMade, 
                        boolean gameOver, String winner) {
            this.remainingShips = remainingShips;
            this.attacksMade = attacksMade;
            this.gameOver = gameOver;
            this.winner = winner;
        }
    }
}
