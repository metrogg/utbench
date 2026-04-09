import java.util.Arrays;
import java.util.Collection;
import java.util.List;
import java.util.ArrayList;
import java.util.Collections;

class EnhancedTrie {
    private static class TrieNode implements Comparable<Object> {
        private char character;
        private boolean terminal;
        private TrieNode[] children = new TrieNode[0];
        private int frequency;  // Added to track word frequency

        public TrieNode(char character) {
            this.character = character;
            this.frequency = 0;
        }

        public boolean isTerminal() {
            return terminal;
        }

        public void setTerminal(boolean terminal) {
            this.terminal = terminal;
            if (terminal) this.frequency++;
        }

        public char getCharacter() {
            return character;
        }

        public int getFrequency() {
            return frequency;
        }

        public void incrementFrequency() {
            this.frequency++;
        }

        public Collection<TrieNode> getChildren() {
            return Arrays.asList(children);
        }

        public TrieNode getChild(char character) {
            int index = Arrays.binarySearch(children, character);
            return index >= 0 ? children[index] : null;
        }

        public TrieNode getChildIfNotExistThenCreate(char character) {
            TrieNode child = getChild(character);
            if (child == null) {
                child = new TrieNode(character);
                addChild(child);
            }
            return child;
        }

        public void addChild(TrieNode child) {
            children = insert(children, child);
        }

        private TrieNode[] insert(TrieNode[] array, TrieNode element) {
            TrieNode[] newArray = new TrieNode[array.length + 1];
            int i = 0;
            while (i < array.length && array[i].getCharacter() < element.getCharacter()) {
                newArray[i] = array[i];
                i++;
            }
            newArray[i] = element;
            System.arraycopy(array, i, newArray, i + 1, array.length - i);
            return newArray;
        }

        @Override
        public int compareTo(Object o) {
            return this.getCharacter() - (char) o;
        }
    }

    private final TrieNode ROOT_NODE = new TrieNode('/');

    public boolean contains(String item) {
        item = item.trim();
        if (item.isEmpty()) return false;

        TrieNode node = ROOT_NODE;
        for (int i = 0; i < item.length(); i++) {
            char character = item.charAt(i);
            TrieNode child = node.getChild(character);
            if (child == null) return false;
            node = child;
        }
        return node.isTerminal();
    }

    public void addAll(List<String> items) {
        for (String item : items) {
            add(item);
        }
    }

    public void add(String item) {
        item = item.trim();
        if (item.isEmpty()) return;

        TrieNode node = ROOT_NODE;
        for (int i = 0; i < item.length(); i++) {
            char character = item.charAt(i);
            node = node.getChildIfNotExistThenCreate(character);
        }
        node.setTerminal(true);
    }

    public List<String> getAllWords() {
        List<String> words = new ArrayList<>();
        collectWords(ROOT_NODE, new StringBuilder(), words);
        return words;
    }

    private void collectWords(TrieNode node, StringBuilder current, List<String> words) {
        if (node.isTerminal()) {
            words.add(current.toString());
        }
        for (TrieNode child : node.getChildren()) {
            current.append(child.getCharacter());
            collectWords(child, current, words);
            current.deleteCharAt(current.length() - 1);
        }
    }

    public List<String> getWordsWithPrefix(String prefix) {
        List<String> words = new ArrayList<>();
        prefix = prefix.trim();
        if (prefix.isEmpty()) return words;

        TrieNode node = ROOT_NODE;
        for (int i = 0; i < prefix.length(); i++) {
            char character = prefix.charAt(i);
            node = node.getChild(character);
            if (node == null) return words;
        }

        collectWords(node, new StringBuilder(prefix), words);
        return words;
    }

    public int getWordFrequency(String word) {
        word = word.trim();
        if (word.isEmpty()) return 0;

        TrieNode node = ROOT_NODE;
        for (int i = 0; i < word.length(); i++) {
            char character = word.charAt(i);
            node = node.getChild(character);
            if (node == null) return 0;
        }
        return node.isTerminal() ? node.getFrequency() : 0;
    }
}
