// Converted Java method
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * A P2P address server that manages peer addresses and their status.
 * Provides functionality to register, unregister, and query peers.
 */
class P2PAddressServer {
    private final Map<String, PeerInfo> peerRegistry;
    private final int maxPeers;
    private int currentPeers;

    /**
     * Constructs a P2P address server with default maximum peers (100).
     */
    public P2PAddressServer() {
        this(100);
    }

    /**
     * Constructs a P2P address server with specified maximum peers.
     * @param maxPeers Maximum number of peers the server can handle
     */
    public P2PAddressServer(int maxPeers) {
        this.peerRegistry = new ConcurrentHashMap<>();
        this.maxPeers = maxPeers;
        this.currentPeers = 0;
    }

    /**
     * Registers a new peer with the server.
     * @param peerId Unique identifier for the peer
     * @param address Network address of the peer
     * @param port Port number the peer is listening on
     * @return true if registration was successful, false otherwise
     */
    public boolean registerPeer(String peerId, String address, int port) {
        if (currentPeers >= maxPeers) {
            return false;
        }
        if (peerRegistry.containsKey(peerId)) {
            return false;
        }
        peerRegistry.put(peerId, new PeerInfo(address, port, System.currentTimeMillis()));
        currentPeers++;
        return true;
    }

    /**
     * Unregisters a peer from the server.
     * @param peerId Unique identifier for the peer
     * @return true if unregistration was successful, false otherwise
     */
    public boolean unregisterPeer(String peerId) {
        if (!peerRegistry.containsKey(peerId)) {
            return false;
        }
        peerRegistry.remove(peerId);
        currentPeers--;
        return true;
    }

    /**
     * Gets information about a specific peer.
     * @param peerId Unique identifier for the peer
     * @return PeerInfo object if peer exists, null otherwise
     */
    public PeerInfo getPeerInfo(String peerId) {
        return peerRegistry.get(peerId);
    }

    /**
     * Gets all registered peers.
     * @return Map of all registered peers
     */
    public Map<String, PeerInfo> getAllPeers() {
        return new HashMap<>(peerRegistry);
    }

    /**
     * Gets the current number of registered peers.
     * @return Number of registered peers
     */
    public int getPeerCount() {
        return currentPeers;
    }

    /**
     * Inner class representing peer information.
     */
    public static class PeerInfo {
        private final String address;
        private final int port;
        private final long registrationTime;

        public PeerInfo(String address, int port, long registrationTime) {
            this.address = address;
            this.port = port;
            this.registrationTime = registrationTime;
        }

        public String getAddress() {
            return address;
        }

        public int getPort() {
            return port;
        }

        public long getRegistrationTime() {
            return registrationTime;
        }

        @Override
        public String toString() {
            return address + ":" + port + " (registered at " + registrationTime + ")";
        }
    }
}
