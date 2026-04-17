import java.util.NoSuchElementException;

class LinkedList<D> {
    private Node<D> head;
    private Node<D> tail;
    private int size;

    private static class Node<D> {
        private D data;
        private Node<D> next;
        private Node<D> prev;

        public Node(D data) {
            this.data = data;
            this.next = null;
            this.prev = null;
        }

        public D getData() {
            return data;
        }

        public Node<D> getNext() {
            return next;
        }

        public void setNext(Node<D> next) {
            this.next = next;
        }

        public Node<D> getPrev() {
            return prev;
        }

        public void setPrev(Node<D> prev) {
            this.prev = prev;
        }
    }

    public LinkedList() {
        head = null;
        tail = null;
        size = 0;
    }

    // Add to the end of the list
    public void add(D data) {
        Node<D> newNode = new Node<>(data);
        if (isEmpty()) {
            head = newNode;
            tail = newNode;
        } else {
            tail.setNext(newNode);
            newNode.setPrev(tail);
            tail = newNode;
        }
        size++;
    }

    // Add at specific index
    public void add(int index, D data) {
        if (index < 0 || index > size) {
            throw new IndexOutOfBoundsException("Index: " + index + ", Size: " + size);
        }

        if (index == size) {
            add(data);
            return;
        }

        Node<D> newNode = new Node<>(data);
        if (index == 0) {
            newNode.setNext(head);
            head.setPrev(newNode);
            head = newNode;
        } else {
            Node<D> current = getNode(index);
            newNode.setNext(current);
            newNode.setPrev(current.getPrev());
            current.getPrev().setNext(newNode);
            current.setPrev(newNode);
        }
        size++;
    }

    // Get data at index
    public D get(int index) {
        return getNode(index).getData();
    }

    // Remove from end
    public D remove() {
        if (isEmpty()) {
            throw new NoSuchElementException("List is empty");
        }
        D data = tail.getData();
        if (size == 1) {
            head = null;
            tail = null;
        } else {
            tail = tail.getPrev();
            tail.setNext(null);
        }
        size--;
        return data;
    }

    // Remove at specific index
    public D remove(int index) {
        if (index < 0 || index >= size) {
            throw new IndexOutOfBoundsException("Index: " + index + ", Size: " + size);
        }

        if (index == size - 1) {
            return remove();
        }

        Node<D> toRemove = getNode(index);
        D data = toRemove.getData();

        if (toRemove == head) {
            head = head.getNext();
            if (head != null) {
                head.setPrev(null);
            }
        } else {
            toRemove.getPrev().setNext(toRemove.getNext());
            toRemove.getNext().setPrev(toRemove.getPrev());
        }
        size--;
        return data;
    }

    public int size() {
        return size;
    }

    public boolean isEmpty() {
        return size == 0;
    }

    public void clear() {
        head = null;
        tail = null;
        size = 0;
    }

    public boolean contains(D data) {
        Node<D> current = head;
        while (current != null) {
            if (current.getData().equals(data)) {
                return true;
            }
            current = current.getNext();
        }
        return false;
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder("[");
        Node<D> current = head;
        while (current != null) {
            sb.append(current.getData());
            if (current.getNext() != null) {
                sb.append(", ");
            }
            current = current.getNext();
        }
        sb.append("]");
        return sb.toString();
    }

    private Node<D> getNode(int index) {
        if (index < 0 || index >= size) {
            throw new IndexOutOfBoundsException("Index: " + index + ", Size: " + size);
        }

        Node<D> current;
        if (index < size / 2) {
            current = head;
            for (int i = 0; i < index; i++) {
                current = current.getNext();
            }
        } else {
            current = tail;
            for (int i = size - 1; i > index; i--) {
                current = current.getPrev();
            }
        }
        return current;
    }
}
