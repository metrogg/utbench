import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import java.util.Queue;

class CourseScheduler {
    private boolean hasCycle;
    private List<Integer> topologicalOrder;

    /**
     * Determines if all courses can be finished and returns the order if possible.
     * Uses Kahn's algorithm for topological sort with cycle detection.
     * 
     * @param numCourses Total number of courses
     * @param prerequisites Array of prerequisite pairs where prerequisites[i] = [a, b] means b must be taken before a
     * @return A Result object containing:
     *         - canFinish: boolean indicating if all courses can be finished
     *         - order: List of course order (empty if can't finish)
     */
    public Result scheduleCourses(int numCourses, int[][] prerequisites) {
        hasCycle = false;
        topologicalOrder = new ArrayList<>();
        
        // Create adjacency list and in-degree count
        List<List<Integer>> adjList = new ArrayList<>(numCourses);
        int[] inDegree = new int[numCourses];
        
        for (int i = 0; i < numCourses; i++) {
            adjList.add(new LinkedList<>());
        }
        
        // Build graph and calculate in-degree for each node
        for (int[] edge : prerequisites) {
            adjList.get(edge[1]).add(edge[0]);
            inDegree[edge[0]]++;
        }
        
        // Initialize queue with all nodes having 0 in-degree
        Queue<Integer> queue = new LinkedList<>();
        for (int i = 0; i < numCourses; i++) {
            if (inDegree[i] == 0) {
                queue.offer(i);
            }
        }
        
        int visitedCount = 0;
        
        // Process nodes in topological order
        while (!queue.isEmpty()) {
            int current = queue.poll();
            topologicalOrder.add(current);
            visitedCount++;
            
            for (int neighbor : adjList.get(current)) {
                if (--inDegree[neighbor] == 0) {
                    queue.offer(neighbor);
                }
            }
        }
        
        // If visitedCount doesn't match numCourses, there's a cycle
        hasCycle = visitedCount != numCourses;
        
        return new Result(!hasCycle, hasCycle ? new ArrayList<>() : topologicalOrder);
    }
    
    public static class Result {
        public final boolean canFinish;
        public final List<Integer> order;
        
        public Result(boolean canFinish, List<Integer> order) {
            this.canFinish = canFinish;
            this.order = order;
        }
    }
}
