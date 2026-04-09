import java.util.*;

class TreePathAnalyzer {
    static class Node {
        int data;
        List<Node> children;

        public Node(int item) {
            data = item;
            children = new ArrayList<>();
        }

        public void addChild(Node child) {
            children.add(child);
        }

        @Override
        public String toString() {
            return String.valueOf(data);
        }
    }

    private int maxLength = 0;
    private List<Node> longestPath = new ArrayList<>();

    /**
     * Finds the longest path between any two nodes in an N-ary tree.
     * The path can start and end at any node and must follow parent-child relationships.
     *
     * @param root The root node of the N-ary tree
     * @return List of nodes representing the longest path
     */
    public List<Node> findLongestPath(Node root) {
        if (root == null) {
            return new ArrayList<>();
        }

        maxLength = 0;
        longestPath.clear();
        dfs(root);
        return longestPath;
    }

    /**
     * Performs depth-first search to find the longest path.
     * @param node Current node being processed
     * @return The longest path starting from this node
     */
    private List<Node> dfs(Node node) {
        if (node == null) {
            return new ArrayList<>();
        }

        List<Node> currentLongest = new ArrayList<>();
        List<Node> secondLongest = new ArrayList<>();

        for (Node child : node.children) {
            List<Node> childPath = dfs(child);
            if (childPath.size() > currentLongest.size()) {
                secondLongest = new ArrayList<>(currentLongest);
                currentLongest = childPath;
            } else if (childPath.size() > secondLongest.size()) {
                secondLongest = childPath;
            }
        }

        // Check if current path through this node is the longest found so far
        int totalLength = 1 + currentLongest.size() + secondLongest.size();
        if (totalLength > maxLength) {
            maxLength = totalLength;
            longestPath = new ArrayList<>();
            longestPath.addAll(currentLongest);
            longestPath.add(node);
            longestPath.addAll(secondLongest);
        }

        // Return the longest path starting from this node
        List<Node> result = new ArrayList<>(currentLongest);
        result.add(node);
        return result;
    }
}
