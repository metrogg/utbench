import java.util.ArrayList;
import java.util.List;

class AdvancedBST {
    private Node root;

    private class Node {
        int value;
        Node left, right;

        Node(int value) {
            this.value = value;
            left = right = null;
        }
    }

    /**
     * Inserts a new value into the BST
     * @param value The value to insert
     */
    public void insert(int value) {
        root = insertRec(root, value);
    }

    private Node insertRec(Node root, int value) {
        if (root == null) {
            return new Node(value);
        }

        if (value < root.value) {
            root.left = insertRec(root.left, value);
        } else if (value > root.value) {
            root.right = insertRec(root.right, value);
        }

        return root;
    }

    /**
     * Removes a value from the BST
     * @param value The value to remove
     */
    public void remove(int value) {
        root = removeRec(root, value);
    }

    private Node removeRec(Node root, int value) {
        if (root == null) {
            return null;
        }

        if (value < root.value) {
            root.left = removeRec(root.left, value);
        } else if (value > root.value) {
            root.right = removeRec(root.right, value);
        } else {
            if (root.left == null) {
                return root.right;
            } else if (root.right == null) {
                return root.left;
            }

            root.value = minValue(root.right);
            root.right = removeRec(root.right, root.value);
        }

        return root;
    }

    private int minValue(Node root) {
        int min = root.value;
        while (root.left != null) {
            min = root.left.value;
            root = root.left;
        }
        return min;
    }

    /**
     * Returns inorder traversal of the BST
     * @return List of values in inorder
     */
    public List<Integer> inorder() {
        List<Integer> result = new ArrayList<>();
        inorderRec(root, result);
        return result;
    }

    private void inorderRec(Node root, List<Integer> result) {
        if (root != null) {
            inorderRec(root.left, result);
            result.add(root.value);
            inorderRec(root.right, result);
        }
    }

    /**
     * Returns preorder traversal of the BST
     * @return List of values in preorder
     */
    public List<Integer> preorder() {
        List<Integer> result = new ArrayList<>();
        preorderRec(root, result);
        return result;
    }

    private void preorderRec(Node root, List<Integer> result) {
        if (root != null) {
            result.add(root.value);
            preorderRec(root.left, result);
            preorderRec(root.right, result);
        }
    }

    /**
     * Returns postorder traversal of the BST
     * @return List of values in postorder
     */
    public List<Integer> postorder() {
        List<Integer> result = new ArrayList<>();
        postorderRec(root, result);
        return result;
    }

    private void postorderRec(Node root, List<Integer> result) {
        if (root != null) {
            postorderRec(root.left, result);
            postorderRec(root.right, result);
            result.add(root.value);
        }
    }

    /**
     * Finds the maximum value in the BST
     * @return Maximum value
     */
    public int findMax() {
        if (root == null) {
            throw new IllegalStateException("Tree is empty");
        }
        Node current = root;
        while (current.right != null) {
            current = current.right;
        }
        return current.value;
    }

    /**
     * Finds the minimum value in the BST
     * @return Minimum value
     */
    public int findMin() {
        if (root == null) {
            throw new IllegalStateException("Tree is empty");
        }
        Node current = root;
        while (current.left != null) {
            current = current.left;
        }
        return current.value;
    }

    /**
     * Checks if the BST contains a value
     * @param value The value to search for
     * @return true if value exists, false otherwise
     */
    public boolean contains(int value) {
        return containsRec(root, value);
    }

    private boolean containsRec(Node root, int value) {
        if (root == null) {
            return false;
        }
        if (value == root.value) {
            return true;
        }
        return value < root.value 
            ? containsRec(root.left, value) 
            : containsRec(root.right, value);
    }

    /**
     * Calculates the height of the BST
     * @return Height of the tree
     */
    public int height() {
        return heightRec(root);
    }

    private int heightRec(Node root) {
        if (root == null) {
            return -1;
        }
        return 1 + Math.max(heightRec(root.left), heightRec(root.right));
    }
}
