import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.stream.Collectors;

class EnhancedPlayersQueue {
    // Thread-safe storage for player information
    private static final Map<String, PlayerInfo> players = new ConcurrentHashMap<>();
    
    // Thread-safe queue for player turn order
    private static final Queue<String> playerQueue = new LinkedList<>();
    
    // Thread-safe list of output streams
    private static final List<ObjectOutputStream> outputStreams = Collections.synchronizedList(new ArrayList<>());
    
    /**
     * Represents player information including score and status
     */
    private static class PlayerInfo {
        private final String playerId;
        private int score;
        private boolean isActive;
        private long joinTime;
        
        public PlayerInfo(String playerId) {
            this.playerId = playerId;
            this.score = 0;
            this.isActive = true;
            this.joinTime = System.currentTimeMillis();
        }
        
        public void incrementScore(int points) {
            this.score += points;
        }
        
        public void setInactive() {
            this.isActive = false;
        }
        
        public String getPlayerId() {
            return playerId;
        }
        
        public int getScore() {
            return score;
        }
        
        public boolean isActive() {
            return isActive;
        }
        
        public long getJoinTime() {
            return joinTime;
        }
    }
    
    /**
     * Registers a new player with the system
     * @param playerId Unique identifier for the player
     * @return true if registration was successful, false if player already exists
     */
    public static synchronized boolean registerPlayer(String playerId) {
        if (players.containsKey(playerId)) {
            return false;
        }
        players.put(playerId, new PlayerInfo(playerId));
        playerQueue.add(playerId);
        return true;
    }
    
    /**
     * Removes a player from the system
     * @param playerId Unique identifier for the player
     * @return true if removal was successful, false if player didn't exist
     */
    public static synchronized boolean removePlayer(String playerId) {
        if (!players.containsKey(playerId)) {
            return false;
        }
        players.get(playerId).setInactive();
        playerQueue.remove(playerId);
        return true;
    }
    
    /**
     * Updates a player's score
     * @param playerId Unique identifier for the player
     * @param points Points to add (can be negative)
     * @return true if update was successful, false if player doesn't exist
     */
    public static synchronized boolean updateScore(String playerId, int points) {
        if (!players.containsKey(playerId)) {
            return false;
        }
        players.get(playerId).incrementScore(points);
        return true;
    }
    
    /**
     * Gets the next player in the queue
     * @return Player ID or null if queue is empty
     */
    public static synchronized String getNextPlayer() {
        if (playerQueue.isEmpty()) {
            return null;
        }
        String nextPlayer = playerQueue.poll();
        playerQueue.add(nextPlayer); // Move to end of queue
        return nextPlayer;
    }
    
    /**
     * Gets the current player rankings sorted by score
     * @return List of player IDs sorted by score (highest first)
     */
    public static synchronized List<String> getRankings() {
        return players.values().stream()
            .filter(PlayerInfo::isActive)
            .sorted(Comparator.comparingInt(PlayerInfo::getScore).reversed())
            .map(PlayerInfo::getPlayerId)
            .collect(Collectors.toList());
    }
    
    /**
     * Adds an output stream for broadcasting messages
     * @param stream The output stream to add
     */
    public static synchronized void addOutputStream(ObjectOutputStream stream) {
        outputStreams.add(stream);
    }
    
    /**
     * Removes an output stream
     * @param stream The output stream to remove
     */
    public static synchronized void removeOutputStream(ObjectOutputStream stream) {
        outputStreams.remove(stream);
    }
    
    /**
     * Gets all active output streams
     * @return List of active output streams
     */
    public static synchronized List<ObjectOutputStream> getOutputStreams() {
        return new ArrayList<>(outputStreams);
    }
}
