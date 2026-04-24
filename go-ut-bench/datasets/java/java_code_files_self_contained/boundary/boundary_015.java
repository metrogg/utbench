// Converted Java method
import java.util.List;
import java.util.ArrayList;
import java.util.Collections;

class BinaryTree<T extends Comparable<T>> {
    private Node<T> root;
    
    private static class Node<T> {
        T data;
        Node<T> left, right;
        
        Node(T data) {
            this.data = data;
        }
    }
    
    public void insert(T data) {
        root = insertRec(root, data);
    }
    
    private Node<T> insertRec(Node<T> node, T data) {
        if (node == null) {
            return new Node<>(data);
        }
        
        if (data.compareTo(node.data) < 0) {
            node.left = insertRec(node.left, data);
        } else if (data.compareTo(node.data) > 0) {
            node.right = insertRec(node.right, data);
        }
        
        return node;
    }
    
    public List<T> inOrderTraversal() {
        List<T> result = new ArrayList<>();
        inOrderRec(root, result);
        return result;
    }
    
    private void inOrderRec(Node<T> node, List<T> result) {
        if (node != null) {
            inOrderRec(node.left, result);
            result.add(node.data);
            inOrderRec(node.right, result);
        }
    }
}

class AdvancedTreeLoader<T extends Comparable<T>> {
    /**
     * Loads data into a binary tree with additional validation and balancing options
     * 
     * @param tree The binary tree to load
     * @param data List of elements to insert
     * @param shouldSort Whether to sort data before insertion (for balanced tree)
     * @param allowDuplicates Whether to allow duplicate values
     * @throws IllegalArgumentException if duplicates are not allowed and data contains them
     */
    public void load(BinaryTree<T> tree, List<T> data, boolean shouldSort, boolean allowDuplicates) {
        if (!allowDuplicates) {
            checkForDuplicates(data);
        }
        
        if (shouldSort) {
            Collections.sort(data);
            loadBalanced(tree, data, 0, data.size() - 1);
        } else {
            for (T item : data) {
                tree.insert(item);
            }
        }
    }
    
    private void checkForDuplicates(List<T> data) {
        for (int i = 0; i < data.size(); i++) {
            for (int j = i + 1; j < data.size(); j++) {
                if (data.get(i).compareTo(data.get(j)) == 0) {
                    throw new IllegalArgumentException("Duplicate elements found: " + data.get(i));
                }
            }
        }
    }
    
    private void loadBalanced(BinaryTree<T> tree, List<T> data, int start, int end) {
        if (start > end) {
            return;
        }
        
        int mid = (start + end) / 2;
        tree.insert(data.get(mid));
        loadBalanced(tree, data, start, mid - 1);
        loadBalanced(tree, data, mid + 1, end);
    }
}
