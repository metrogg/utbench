// Converted Java method
import java.util.*;

class GraphAnalyzer {
    /**
     * Finds shortest paths from a source vertex to all other vertices using Dijkstra's algorithm.
     * Also provides additional graph analysis metrics.
     * 
     * @param graph Adjacency matrix representation of the graph
     * @param source Source vertex index
     * @return A Map containing:
     *         - "distances": array of shortest distances from source
     *         - "paths": list of shortest paths from source to each vertex
     *         - "visitedCount": number of vertices visited during algorithm execution
     *         - "maxDistance": maximum shortest distance found
     * @throws IllegalArgumentException if graph is empty or source is invalid
     */
    public static Map<String, Object> analyzeGraph(int[][] graph, int source) {
        if (graph.length == 0 || graph[0].length == 0) {
            throw new IllegalArgumentException("Graph cannot be empty");
        }
        if (source < 0 || source >= graph.length) {
            throw new IllegalArgumentException("Invalid source vertex");
        }

        int n = graph.length;
        int[] dist = new int[n];
        int[] prev = new int[n];
        boolean[] visited = new boolean[n];
        int visitedCount = 0;
        int maxDistance = 0;

        Arrays.fill(dist, Integer.MAX_VALUE);
        Arrays.fill(prev, -1);
        dist[source] = 0;

        PriorityQueue<Vertex> pq = new PriorityQueue<>(n, Comparator.comparingInt(v -> v.distance));
        pq.add(new Vertex(source, 0));

        while (!pq.isEmpty()) {
            Vertex current = pq.poll();
            int u = current.index;

            if (visited[u]) continue;
            visited[u] = true;
            visitedCount++;

            for (int v = 0; v < n; v++) {
                if (graph[u][v] > 0 && !visited[v]) {
                    int alt = dist[u] + graph[u][v];
                    if (alt < dist[v]) {
                        dist[v] = alt;
                        prev[v] = u;
                        pq.add(new Vertex(v, dist[v]));
                    }
                }
            }
        }

        // Calculate max distance and build paths
        List<List<Integer>> paths = new ArrayList<>(n);
        for (int i = 0; i < n; i++) {
            paths.add(buildPath(prev, source, i));
            if (dist[i] != Integer.MAX_VALUE && dist[i] > maxDistance) {
                maxDistance = dist[i];
            }
        }

        Map<String, Object> result = new HashMap<>();
        result.put("distances", dist);
        result.put("paths", paths);
        result.put("visitedCount", visitedCount);
        result.put("maxDistance", maxDistance);
        return result;
    }

    private static List<Integer> buildPath(int[] prev, int source, int target) {
        List<Integer> path = new LinkedList<>();
        if (prev[target] == -1 && target != source) {
            return path; // No path exists
        }

        for (int at = target; at != -1; at = prev[at]) {
            path.add(0, at);
        }
        return path;
    }

    private static class Vertex {
        int index;
        int distance;

        Vertex(int index, int distance) {
            this.index = index;
            this.distance = distance;
        }
    }
}
