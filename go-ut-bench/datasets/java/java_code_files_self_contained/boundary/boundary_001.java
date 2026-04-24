// Converted Java method
import java.util.ArrayList;
import java.util.List;

class BSTValidator {

    /**
     * Validates if a sequence could be the post-order traversal of a BST.
     * Enhanced to also check for pre-order and in-order traversals.
     * 
     * @param sequence The sequence to validate
     * @param traversalType Type of traversal to validate against (POST, PRE, IN)
     * @return true if the sequence matches the specified traversal type of a BST
     * @throws IllegalArgumentException if sequence is null or traversalType is invalid
     */
    public static boolean validateBSTSequence(int[] sequence, String traversalType) {
        if (sequence == null) {
            throw new IllegalArgumentException("Sequence cannot be null");
        }
        
        if (!traversalType.equals("POST") && !traversalType.equals("PRE") && !traversalType.equals("IN")) {
            throw new IllegalArgumentException("Invalid traversal type. Use POST, PRE, or IN");
        }

        if (sequence.length == 0) {
            return false;
        }

        List<Integer> list = new ArrayList<>();
        int rootIndex;
        
        if (traversalType.equals("POST")) {
            // For post-order, last element is root
            rootIndex = sequence.length - 1;
            for (int i = 0; i < sequence.length - 1; i++) {
                list.add(sequence[i]);
            }
            return validatePostOrder(list, sequence[rootIndex]);
        } else if (traversalType.equals("PRE")) {
            // For pre-order, first element is root
            rootIndex = 0;
            for (int i = 1; i < sequence.length; i++) {
                list.add(sequence[i]);
            }
            return validatePreOrder(list, sequence[rootIndex]);
        } else {
            // For in-order, should be sorted
            return validateInOrder(sequence);
        }
    }

    private static boolean validatePostOrder(List<Integer> arr, int root) {
        if (arr.size() < 2) {
            return true;
        }

        int splitIndex = arr.size();
        List<Integer> left = new ArrayList<>();
        List<Integer> right = new ArrayList<>();

        for (int i = 0; i < arr.size(); i++) {
            if (arr.get(i) > root) {
                splitIndex = i;
                break;
            } else {
                left.add(arr.get(i));
            }
        }

        for (int i = splitIndex; i < arr.size(); i++) {
            if (arr.get(i) < root) {
                return false;
            } else {
                right.add(arr.get(i));
            }
        }

        boolean leftValid = true;
        boolean rightValid = true;

        if (!left.isEmpty()) {
            int leftRoot = left.get(left.size() - 1);
            left.remove(left.size() - 1);
            leftValid = validatePostOrder(left, leftRoot);
        }

        if (!right.isEmpty()) {
            int rightRoot = right.get(right.size() - 1);
            right.remove(right.size() - 1);
            rightValid = validatePostOrder(right, rightRoot);
        }

        return leftValid && rightValid;
    }

    private static boolean validatePreOrder(List<Integer> arr, int root) {
        if (arr.size() < 2) {
            return true;
        }

        int splitIndex = arr.size();
        List<Integer> left = new ArrayList<>();
        List<Integer> right = new ArrayList<>();

        // First elements < root are left subtree
        // Then elements > root are right subtree
        for (int i = 0; i < arr.size(); i++) {
            if (arr.get(i) > root) {
                splitIndex = i;
                break;
            } else {
                left.add(arr.get(i));
            }
        }

        for (int i = splitIndex; i < arr.size(); i++) {
            if (arr.get(i) < root) {
                return false;
            } else {
                right.add(arr.get(i));
            }
        }

        boolean leftValid = true;
        boolean rightValid = true;

        if (!left.isEmpty()) {
            int leftRoot = left.get(0);
            left.remove(0);
            leftValid = validatePreOrder(left, leftRoot);
        }

        if (!right.isEmpty()) {
            int rightRoot = right.get(0);
            right.remove(0);
            rightValid = validatePreOrder(right, rightRoot);
        }

        return leftValid && rightValid;
    }

    private static boolean validateInOrder(int[] sequence) {
        for (int i = 1; i < sequence.length; i++) {
            if (sequence[i] <= sequence[i - 1]) {
                return false;
            }
        }
        return true;
    }
}
