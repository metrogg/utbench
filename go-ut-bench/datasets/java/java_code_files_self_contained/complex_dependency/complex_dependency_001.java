class TreeNode {
    int value;
    TreeNode left;
    TreeNode right;
    
    TreeNode(int value) {
        this.value = value;
        this.left = null;
        this.right = null;
    }
}

class BinarySearchTree {
    private TreeNode root;
    
    public BinarySearchTree() {
        this.root = null;
    }
    
    public void insert(int value) {
        root = insertRecursive(root, value);
    }
    
    private TreeNode insertRecursive(TreeNode node, int value) {
        if (node == null) {
            return new TreeNode(value);
        }
        if (value < node.value) {
            node.left = insertRecursive(node.left, value);
        } else {
            node.right = insertRecursive(node.right, value);
        }
        return node;
    }
    
    public boolean search(int value) {
        return searchRecursive(root, value);
    }
    
    private boolean searchRecursive(TreeNode node, int value) {
        if (node == null) {
            return false;
        }
        if (node.value == value) {
            return true;
        }
        if (value < node.value) {
            return searchRecursive(node.left, value);
        }
        return searchRecursive(node.right, value);
    }
    
    public int[] inorderTraversal() {
        int[] result = new int[100];
        int index = 0;
        index = inorderRecursive(root, result, index);
        int[] finalResult = new int[index];
        for (int i = 0; i < index; i++) {
            finalResult[i] = result[i];
        }
        return finalResult;
    }
    
    private int inorderRecursive(TreeNode node, int[] result, int index) {
        if (node != null) {
            index = inorderRecursive(node.left, result, index);
            result[index++] = node.value;
            index = inorderRecursive(node.right, result, index);
        }
        return index;
    }
}