// Converted Java method
import java.util.ArrayList;
import java.util.List;

class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;
    
    public TreeNode(int val) {
        this.val = val;
    }
}

class TreePathAnalyzer {
    /**
     * Finds all root-to-leaf paths in a binary tree and performs additional analysis.
     * 
     * @param root The root node of the binary tree
     * @return A map containing:
     *         - "paths": List of all root-to-leaf paths (each path as List<Integer>)
     *         - "maxDepth": Maximum depth of the tree
     *         - "leafCount": Number of leaf nodes
     *         - "hasPathSum": Whether any path sums to a target value (if provided)
     */
    public TreeAnalysisResult analyzeTreePaths(TreeNode root, Integer targetSum) {
        TreeAnalysisResult result = new TreeAnalysisResult();
        if (root == null) {
            return result;
        }
        
        List<Integer> currentPath = new ArrayList<>();
        traverseTree(root, currentPath, result, targetSum);
        
        return result;
    }
    
    private void traverseTree(TreeNode node, List<Integer> currentPath, 
                            TreeAnalysisResult result, Integer targetSum) {
        if (node == null) {
            return;
        }
        
        currentPath.add(node.val);
        
        if (node.left == null && node.right == null) {
            // Found a leaf node
            result.leafCount++;
            result.maxDepth = Math.max(result.maxDepth, currentPath.size());
            
            // Add the current path to results
            result.paths.add(new ArrayList<>(currentPath));
            
            // Check path sum if target is provided
            if (targetSum != null) {
                int pathSum = currentPath.stream().mapToInt(Integer::intValue).sum();
                if (pathSum == targetSum) {
                    result.hasPathSum = true;
                }
            }
        }
        
        traverseTree(node.left, currentPath, result, targetSum);
        traverseTree(node.right, currentPath, result, targetSum);
        
        currentPath.remove(currentPath.size() - 1);
    }
}

class TreeAnalysisResult {
    List<List<Integer>> paths = new ArrayList<>();
    int maxDepth = 0;
    int leafCount = 0;
    boolean hasPathSum = false;
    
    @Override
    public String toString() {
        return "Paths: " + paths + "\n" +
               "Max Depth: " + maxDepth + "\n" +
               "Leaf Count: " + leafCount + "\n" +
               "Has Path Sum: " + hasPathSum;
    }
}
