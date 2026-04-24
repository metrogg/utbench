// Converted Java method
import java.util.*;
import java.util.stream.Collectors;

class GraphAnalyzer {
    private Map<String, List<Edge>> adjacencyList;
    private boolean directed;
    private boolean weighted;

    public GraphAnalyzer(boolean directed, boolean weighted) {
        this.adjacencyList = new HashMap<>();
        this.directed = directed;
        this.weighted = weighted;
    }

    public void addNode(String node) {
        adjacencyList.putIfAbsent(node, new ArrayList<>());
    }

    public void addEdge(String from, String to, int weight) {
        addNode(from);
        addNode(to);
        adjacencyList.get(from).add(new Edge(to, weight));
        if (!directed) {
            adjacencyList.get(to).add(new Edge(from, weight));
        }
    }

    /**
     * Finds the shortest path between two nodes using Dijkstra's algorithm
     * @param start Starting node
     * @param end Destination node
     * @return List of nodes representing the shortest path
     */
    public List<String> findShortestPath(String start, String end) {
        if (!adjacencyList.containsKey(start) || !adjacencyList.containsKey(end)) {
            throw new IllegalArgumentException("Start or end node not found in graph");
        }

        Map<String, Integer> distances = new HashMap<>();
        Map<String, String> previous = new HashMap<>();
        PriorityQueue<NodeDistance> queue = new PriorityQueue<>();

        // Initialize distances
        for (String node : adjacencyList.keySet()) {
            distances.put(node, Integer.MAX_VALUE);
        }
        distances.put(start, 0);
        queue.add(new NodeDistance(start, 0));

        while (!queue.isEmpty()) {
            NodeDistance current = queue.poll();
            String currentNode = current.node;

            if (currentNode.equals(end)) {
                break; // Found the shortest path to end
            }

            for (Edge edge : adjacencyList.get(currentNode)) {
                String neighbor = edge.to;
                int newDist = distances.get(currentNode) + edge.weight;
                if (newDist < distances.get(neighbor)) {
                    distances.put(neighbor, newDist);
                    previous.put(neighbor, currentNode);
                    queue.add(new NodeDistance(neighbor, newDist));
                }
            }
        }

        // Reconstruct path
        List<String> path = new ArrayList<>();
        for (String at = end; at != null; at = previous.get(at)) {
            path.add(at);
        }
        Collections.reverse(path);

        return path.size() == 1 && !start.equals(end) ? 
            Collections.emptyList() : path;
    }

    /**
     * Checks if the graph contains a cycle using DFS
     * @return true if cycle exists, false otherwise
     */
    public boolean hasCycle() {
        Set<String> visited = new HashSet<>();
        Set<String> recursionStack = new HashSet<>();

        for (String node : adjacencyList.keySet()) {
            if (hasCycleUtil(node, visited, recursionStack, null)) {
                return true;
            }
        }
        return false;
    }

    private boolean hasCycleUtil(String node, Set<String> visited, 
                               Set<String> recursionStack, String parent) {
        if (recursionStack.contains(node)) return true;
        if (visited.contains(node)) return false;

        visited.add(node);
        recursionStack.add(node);

        for (Edge edge : adjacencyList.get(node)) {
            if (!directed && edge.to.equals(parent)) continue;
            if (hasCycleUtil(edge.to, visited, recursionStack, node)) {
                return true;
            }
        }

        recursionStack.remove(node);
        return false;
    }

    /**
     * Calculates the degree of a node (number of edges)
     * @param node The node to check
     * @return The degree of the node
     */
    public int getNodeDegree(String node) {
        if (!adjacencyList.containsKey(node)) {
            throw new IllegalArgumentException("Node not found in graph");
        }
        return adjacencyList.get(node).size();
    }

    static class Edge {
        String to;
        int weight;

        public Edge(String to, int weight) {
            this.to = to;
            this.weight = weight;
        }
    }

    static class NodeDistance implements Comparable<NodeDistance> {
        String node;
        int distance;

        public NodeDistance(String node, int distance) {
            this.node = node;
            this.distance = distance;
        }

        @Override
        public int compareTo(NodeDistance other) {
            return Integer.compare(this.distance, other.distance);
        }
    }
}
