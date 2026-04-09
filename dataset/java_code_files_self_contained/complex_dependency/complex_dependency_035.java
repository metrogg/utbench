import java.util.*;

class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;
    TreeNode() {}
    TreeNode(int val) { this.val = val; }
    TreeNode(int val, TreeNode left, TreeNode right) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

class TreeAnalyzer {
    /**
     * Performs comprehensive analysis on a binary tree including:
     * - Maximum ancestor difference
     * - Minimum ancestor difference
     * - Average ancestor difference
     * - Node count
     * 
     * @param root Root of the binary tree
     * @return Map containing analysis results
     */
    public Map<String, Integer> analyzeTree(TreeNode root) {
        Map<String, Integer> result = new HashMap<>();
        if (root == null) {
            result.put("maxDiff", 0);
            result.put("minDiff", 0);
            result.put("avgDiff", 0);
            result.put("nodeCount", 0);
            return result;
        }
        
        AnalysisData data = analyzeHelper(root, root.val, root.val, 0, 0);
        
        result.put("maxDiff", data.maxDiff);
        result.put("minDiff", data.minDiff);
        result.put("avgDiff", data.totalDiff / data.nodeCount);
        result.put("nodeCount", data.nodeCount);
        
        return result;
    }
    
    private AnalysisData analyzeHelper(TreeNode node, int curMax, int curMin, int totalDiff, int nodeCount) {
        if (node == null) {
            return new AnalysisData(curMax - curMin, curMax - curMin, 0, 0);
        }
        
        // Update current max and min
        int newMax = Math.max(curMax, node.val);
        int newMin = Math.min(curMin, node.val);
        int currentDiff = newMax - newMin;
        
        // Recursively analyze left and right subtrees
        AnalysisData left = analyzeHelper(node.left, newMax, newMin, totalDiff + currentDiff, nodeCount + 1);
        AnalysisData right = analyzeHelper(node.right, newMax, newMin, totalDiff + currentDiff, nodeCount + 1);
        
        // Combine results
        int maxDiff = Math.max(currentDiff, Math.max(left.maxDiff, right.maxDiff));
        int minDiff = Math.min(currentDiff, Math.min(left.minDiff, right.minDiff));
        int sumDiff = currentDiff + left.totalDiff + right.totalDiff;
        int count = 1 + left.nodeCount + right.nodeCount;
        
        return new AnalysisData(maxDiff, minDiff, sumDiff, count);
    }
    
    private static class AnalysisData {
        int maxDiff;
        int minDiff;
        int totalDiff;
        int nodeCount;
        
        AnalysisData(int maxDiff, int minDiff, int totalDiff, int nodeCount) {
            this.maxDiff = maxDiff;
            this.minDiff = minDiff;
            this.totalDiff = totalDiff;
            this.nodeCount = nodeCount;
        }
    }
}
