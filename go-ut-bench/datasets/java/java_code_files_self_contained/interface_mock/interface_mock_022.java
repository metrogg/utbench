import java.util.HashMap;
import java.util.Map;

class GraphDatabaseManager {
    private Map<Long, Node> nodes;
    private Map<Long, Relationship> relationships;
    private long nextNodeId = 1;
    private long nextRelationshipId = 1;

    public GraphDatabaseManager() {
        nodes = new HashMap<>();
        relationships = new HashMap<>();
    }

    /**
     * Creates a new node with the given properties
     * @param properties Map of node properties
     * @return ID of the created node
     */
    public long createNode(Map<String, Object> properties) {
        long nodeId = nextNodeId++;
        nodes.put(nodeId, new Node(nodeId, properties));
        return nodeId;
    }

    /**
     * Creates a relationship between two nodes
     * @param startNodeId ID of the starting node
     * @param endNodeId ID of the ending node
     * @param type Type of the relationship
     * @param properties Map of relationship properties
     * @return ID of the created relationship
     * @throws IllegalArgumentException if either node doesn't exist
     */
    public long createRelationship(long startNodeId, long endNodeId, 
                                 String type, Map<String, Object> properties) {
        if (!nodes.containsKey(startNodeId) || !nodes.containsKey(endNodeId)) {
            throw new IllegalArgumentException("One or both nodes do not exist");
        }
        
        long relationshipId = nextRelationshipId++;
        Relationship relationship = new Relationship(relationshipId, 
            startNodeId, endNodeId, type, properties);
        relationships.put(relationshipId, relationship);
        return relationshipId;
    }

    /**
     * Gets information about a node and its relationships
     * @param nodeId ID of the node to query
     * @return String containing node information and relationships
     * @throws IllegalArgumentException if node doesn't exist
     */
    public String getNodeInfo(long nodeId) {
        if (!nodes.containsKey(nodeId)) {
            throw new IllegalArgumentException("Node does not exist");
        }
        
        Node node = nodes.get(nodeId);
        StringBuilder info = new StringBuilder();
        info.append("Node ID: ").append(nodeId).append("\n");
        info.append("Properties: ").append(node.getProperties()).append("\n");
        
        info.append("Relationships:\n");
        for (Relationship rel : relationships.values()) {
            if (rel.getStartNodeId() == nodeId || rel.getEndNodeId() == nodeId) {
                info.append("- ").append(rel.toString()).append("\n");
            }
        }
        
        return info.toString();
    }

    /**
     * Deletes a node and all its relationships
     * @param nodeId ID of the node to delete
     * @throws IllegalArgumentException if node doesn't exist
     */
    public void deleteNode(long nodeId) {
        if (!nodes.containsKey(nodeId)) {
            throw new IllegalArgumentException("Node does not exist");
        }
        
        // Remove all relationships involving this node
        relationships.values().removeIf(rel -> 
            rel.getStartNodeId() == nodeId || rel.getEndNodeId() == nodeId);
        
        // Remove the node
        nodes.remove(nodeId);
    }

    // Inner classes representing simplified Node and Relationship
    private static class Node {
        private long id;
        private Map<String, Object> properties;

        public Node(long id, Map<String, Object> properties) {
            this.id = id;
            this.properties = new HashMap<>(properties);
        }

        public Map<String, Object> getProperties() {
            return new HashMap<>(properties);
        }
    }

    private static class Relationship {
        private long id;
        private long startNodeId;
        private long endNodeId;
        private String type;
        private Map<String, Object> properties;

        public Relationship(long id, long startNodeId, long endNodeId, 
                          String type, Map<String, Object> properties) {
            this.id = id;
            this.startNodeId = startNodeId;
            this.endNodeId = endNodeId;
            this.type = type;
            this.properties = new HashMap<>(properties);
        }

        public long getStartNodeId() { return startNodeId; }
        public long getEndNodeId() { return endNodeId; }

        @Override
        public String toString() {
            return "Relationship ID: " + id + 
                   ", Type: " + type + 
                   ", From: " + startNodeId + 
                   ", To: " + endNodeId + 
                   ", Properties: " + properties;
        }
    }
}
