// Converted Java method
import java.util.*;

class EnhancedBinarySearchTree<T extends Comparable<T>> {
    private Node root;
    private int size;

    private class Node {
        T data;
        Node left, right;

        Node(T data) {
            this.data = data;
        }
    }

    /**
     * Adds an element to the tree if it's not already present
     * @param element The element to add
     * @return true if the element was added, false if it was already present
     */
    public boolean add(T element) {
        if (element == null) {
            throw new NullPointerException("Cannot add null elements");
        }
        
        if (root == null) {
            root = new Node(element);
            size++;
            return true;
        }
        
        return add(root, element);
    }

    private boolean add(Node node, T element) {
        int cmp = element.compareTo(node.data);
        
        if (cmp < 0) {
            if (node.left == null) {
                node.left = new Node(element);
                size++;
                return true;
            } else {
                return add(node.left, element);
            }
        } else if (cmp > 0) {
            if (node.right == null) {
                node.right = new Node(element);
                size++;
                return true;
            } else {
                return add(node.right, element);
            }
        } else {
            return false; // element already exists
        }
    }

    /**
     * Finds the smallest element in the tree
     * @return The smallest element
     * @throws NoSuchElementException if the tree is empty
     */
    public T first() {
        if (root == null) {
            throw new NoSuchElementException("Tree is empty");
        }
        
        Node current = root;
        while (current.left != null) {
            current = current.left;
        }
        return current.data;
    }

    /**
     * Finds the largest element in the tree
     * @return The largest element
     * @throws NoSuchElementException if the tree is empty
     */
    public T last() {
        if (root == null) {
            throw new NoSuchElementException("Tree is empty");
        }
        
        Node current = root;
        while (current.right != null) {
            current = current.right;
        }
        return current.data;
    }

    /**
     * Checks if the tree contains an element
     * @param element The element to search for
     * @return true if the element is found, false otherwise
     */
    public boolean contains(T element) {
        return contains(root, element);
    }

    private boolean contains(Node node, T element) {
        if (node == null) return false;
        
        int cmp = element.compareTo(node.data);
        if (cmp < 0) {
            return contains(node.left, element);
        } else if (cmp > 0) {
            return contains(node.right, element);
        } else {
            return true;
        }
    }

    /**
     * Removes an element from the tree
     * @param element The element to remove
     * @return true if the element was removed, false if it wasn't found
     */
    public boolean remove(T element) {
        if (element == null) {
            throw new NullPointerException("Cannot remove null elements");
        }
        
        int oldSize = size;
        root = remove(root, element);
        return size < oldSize;
    }

    private Node remove(Node node, T element) {
        if (node == null) return null;
        
        int cmp = element.compareTo(node.data);
        if (cmp < 0) {
            node.left = remove(node.left, element);
        } else if (cmp > 0) {
            node.right = remove(node.right, element);
        } else {
            // Found the node to remove
            size--;
            
            // Case 1: No children
            if (node.left == null && node.right == null) {
                return null;
            }
            // Case 2: One child
            else if (node.left == null) {
                return node.right;
            } else if (node.right == null) {
                return node.left;
            }
            // Case 3: Two children
            else {
                // Find the smallest node in the right subtree
                Node successor = findMin(node.right);
                node.data = successor.data;
                node.right = remove(node.right, successor.data);
            }
        }
        return node;
    }

    private Node findMin(Node node) {
        while (node.left != null) {
            node = node.left;
        }
        return node;
    }

    /**
     * Checks if the tree is empty
     * @return true if the tree is empty, false otherwise
     */
    public boolean isEmpty() {
        return size == 0;
    }

    /**
     * Gets the number of elements in the tree
     * @return The size of the tree
     */
    public int size() {
        return size;
    }

    /**
     * Clears all elements from the tree
     */
    public void clear() {
        root = null;
        size = 0;
    }
}
