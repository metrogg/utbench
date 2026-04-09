// Converted Java method
import java.util.NoSuchElementException;

class EnhancedQueue<E> {
    private Node<E> head;
    private Node<E> tail;
    private int size;
    private int capacity;
    private static final int DEFAULT_CAPACITY = 10;

    public EnhancedQueue() {
        this(DEFAULT_CAPACITY);
    }

    public EnhancedQueue(int capacity) {
        if (capacity <= 0) {
            throw new IllegalArgumentException("Capacity must be positive");
        }
        this.head = null;
        this.tail = null;
        this.size = 0;
        this.capacity = capacity;
    }

    public void enqueue(E newData) {
        if (isFull()) {
            throw new IllegalStateException("Queue is full");
        }
        Node<E> newNode = new Node<>(newData, null);
        if (isEmpty()) {
            this.head = newNode;
        } else {
            this.tail.setNext(newNode);
        }
        this.tail = newNode;
        size++;
    }

    public E dequeue() {
        if (isEmpty()) {
            throw new NoSuchElementException("Queue is empty");
        }
        E data = this.head.getData();
        this.head = this.head.getNext();
        size--;
        if (isEmpty()) {
            this.tail = null;
        }
        return data;
    }

    public E peek() {
        if (isEmpty()) {
            throw new NoSuchElementException("Queue is empty");
        }
        return this.head.getData();
    }

    public boolean isEmpty() {
        return size == 0;
    }

    public boolean isFull() {
        return size == capacity;
    }

    public int size() {
        return size;
    }

    public void clear() {
        this.head = null;
        this.tail = null;
        this.size = 0;
    }

    public boolean contains(E element) {
        Node<E> current = this.head;
        while (current != null) {
            if (current.getData().equals(element)) {
                return true;
            }
            current = current.getNext();
        }
        return false;
    }

    public void printQueue() {
        if (isEmpty()) {
            System.out.println("Queue is empty");
            return;
        }
        Node<E> current = this.head;
        while (current != null) {
            System.out.print(current.getData() + " ");
            current = current.getNext();
        }
        System.out.println();
    }

    private static class Node<E> {
        private E data;
        private Node<E> next;

        public Node(E data, Node<E> next) {
            this.data = data;
            this.next = next;
        }

        public E getData() {
            return data;
        }

        public Node<E> getNext() {
            return next;
        }

        public void setNext(Node<E> next) {
            this.next = next;
        }
    }
}
