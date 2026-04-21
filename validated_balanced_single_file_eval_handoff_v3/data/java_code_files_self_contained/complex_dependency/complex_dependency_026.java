// Converted Java method - GameEngine core logic
import java.util.HashMap;
import java.util.Map;

class GameEngine {
    private Player player;
    private Map<String, Boolean> gameState;
    private int score;
    
    /**
     * Initializes a new game engine with default player and game state
     */
    public GameEngine() {
        this.player = new Player();
        this.gameState = new HashMap<>();
        this.score = 0;
        initializeGameState();
    }
    
    /**
     * Initializes default game state values
     */
    private void initializeGameState() {
        gameState.put("gameStarted", false);
        gameState.put("gamePaused", false);
        gameState.put("gameOver", false);
    }
    
    /**
     * Starts a new game with initial player position
     * @param startX The starting X coordinate
     * @param startY The starting Y coordinate
     */
    public void startNewGame(int startX, int startY) {
        player.setPosition(startX, startY);
        score = 0;
        gameState.put("gameStarted", true);
        gameState.put("gamePaused", false);
        gameState.put("gameOver", false);
    }
    
    /**
     * Handles player movement based on input direction
     * @param direction The direction to move (UP, DOWN, LEFT, RIGHT)
     * @return New player position as int array [x, y]
     */
    public int[] handlePlayerMovement(String direction) {
        if (!gameState.get("gameStarted") || gameState.get("gamePaused")) {
            return player.getPosition();
        }
        
        switch (direction.toUpperCase()) {
            case "UP":
                player.move(0, -1);
                break;
            case "DOWN":
                player.move(0, 1);
                break;
            case "LEFT":
                player.move(-1, 0);
                break;
            case "RIGHT":
                player.move(1, 0);
                break;
        }
        
        return player.getPosition();
    }
    
    /**
     * Handles mouse click events in game coordinates
     * @param x X coordinate of click
     * @param y Y coordinate of click
     * @return True if click resulted in valid game action
     */
    public boolean handleMouseClick(int x, int y) {
        if (!gameState.get("gameStarted")) {
            return false;
        }
        
        // Simple interaction logic - increase score if click is near player
        int[] playerPos = player.getPosition();
        if (Math.abs(x - playerPos[0]) < 50 && Math.abs(y - playerPos[1]) < 50) {
            score += 10;
            return true;
        }
        return false;
    }
    
    /**
     * Gets current game score
     * @return The current score
     */
    public int getScore() {
        return score;
    }
    
    /**
     * Gets current game state
     * @return Map of game state flags
     */
    public Map<String, Boolean> getGameState() {
        return new HashMap<>(gameState);
    }
}

class Player {
    private int x;
    private int y;
    
    public Player() {
        this.x = 0;
        this.y = 0;
    }
    
    public void setPosition(int x, int y) {
        this.x = x;
        this.y = y;
    }
    
    public void move(int deltaX, int deltaY) {
        this.x += deltaX;
        this.y += deltaY;
    }
    
    public int[] getPosition() {
        return new int[]{x, y};
    }
}
