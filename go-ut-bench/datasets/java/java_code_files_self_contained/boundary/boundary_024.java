import java.util.ArrayList;
import java.util.List;

class LinkedListNode<T> {
    T data;
    LinkedListNode<T> next;

    public LinkedListNode(T data) {
        this.data = data;
        this.next = null;
    }

    public T getData() {
        return data;
    }
}

class LinkedListSwapper {
    /**
     * Swaps two nodes in a linked list at given positions i and j.
     * Handles various edge cases including adjacent nodes, head/tail swaps,
     * and invalid positions.
     *
     * @param head Head of the linked list
     * @param i First position (0-based)
     * @param j Second position (0-based)
     * @return Head of the modified linked list after swap
     */
    public static LinkedListNode<Integer> swapNodes(LinkedListNode<Integer> head, int i, int j) {
        if (head == null || i == j) {
            return head;
        }

        // Ensure i is always the smaller index to simplify logic
        if (i > j) {
            int temp = i;
            i = j;
            j = temp;
        }

        LinkedListNode<Integer> prevI = null, currI = head;
        int pos = 0;
        while (currI != null && pos < i) {
            prevI = currI;
            currI = currI.next;
            pos++;
        }

        LinkedListNode<Integer> prevJ = null, currJ = head;
        pos = 0;
        while (currJ != null && pos < j) {
            prevJ = currJ;
            currJ = currJ.next;
            pos++;
        }

        // If either node wasn't found (invalid positions)
        if (currI == null || currJ == null) {
            return head;
        }

        // Special case: nodes are adjacent
        if (currI.next == currJ) {
            currI.next = currJ.next;
            currJ.next = currI;
            if (prevI != null) {
                prevI.next = currJ;
            } else {
                head = currJ;
            }
            return head;
        }

        // General case: nodes are not adjacent
        LinkedListNode<Integer> temp = currI.next;
        currI.next = currJ.next;
        currJ.next = temp;

        if (prevI != null) {
            prevI.next = currJ;
        } else {
            head = currJ;
        }

        if (prevJ != null) {
            prevJ.next = currI;
        } else {
            head = currI;
        }

        return head;
    }

    /**
     * Helper method to create a linked list from a list of integers.
     *
     * @param elements List of integers
     * @return Head of the created linked list
     */
    public static LinkedListNode<Integer> createLinkedList(List<Integer> elements) {
        if (elements.isEmpty()) {
            return null;
        }

        LinkedListNode<Integer> head = new LinkedListNode<>(elements.get(0));
        LinkedListNode<Integer> current = head;

        for (int i = 1; i < elements.size(); i++) {
            current.next = new LinkedListNode<>(elements.get(i));
            current = current.next;
        }

        return head;
    }

    /**
     * Helper method to convert a linked list to a list of integers.
     *
     * @param head Head of the linked list
     * @return List of integers representing the linked list
     */
    public static List<Integer> linkedListToList(LinkedListNode<Integer> head) {
        List<Integer> result = new ArrayList<>();
        LinkedListNode<Integer> current = head;
        while (current != null) {
            result.add(current.getData());
            current = current.next;
        }
        return result;
    }
}
