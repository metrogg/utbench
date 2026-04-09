#include <string>
#include <vector>
#include <algorithm>
#include <map>
#include <utility>
using namespace std;

class Node {
    string word;
    int count;
    Node *left, *right;
    int height;  // For AVL tree implementation

public:
    Node(const string& w) : word(w), count(1), left(nullptr), right(nullptr), height(1) {}
    
    string getWord() const { return word; }
    int getCount() const { return count; }
    void increment() { count++; }
    void decrement() { if (count > 0) count--; }
    
    friend class Tree;
    friend ostream& operator<<(ostream& os, const Node& n);
};

class Tree {
    Node* root;
    
    // AVL tree helper functions
    int height(Node* n) { return n ? n->height : 0; }
    int balanceFactor(Node* n) { return n ? height(n->left) - height(n->right) : 0; }
    void updateHeight(Node* n) { n->height = 1 + max(height(n->left), height(n->right)); }
    
    Node* rotateRight(Node* y) {
        Node* x = y->left;
        Node* T2 = x->right;
        
        x->right = y;
        y->left = T2;
        
        updateHeight(y);
        updateHeight(x);
        
        return x;
    }
    
    Node* rotateLeft(Node* x) {
        Node* y = x->right;
        Node* T2 = y->left;
        
        y->left = x;
        x->right = T2;
        
        updateHeight(x);
        updateHeight(y);
        
        return y;
    }
    
    Node* balance(Node* n) {
        if (!n) return n;
        
        updateHeight(n);
        int bf = balanceFactor(n);
        
        if (bf > 1) {
            if (balanceFactor(n->left) < 0)
                n->left = rotateLeft(n->left);
            return rotateRight(n);
        }
        if (bf < -1) {
            if (balanceFactor(n->right) > 0)
                n->right = rotateRight(n->right);
            return rotateLeft(n);
        }
        return n;
    }
    
    Node* insert(Node* node, const string& word) {
        if (!node) return new Node(word);
        
        if (word < node->word)
            node->left = insert(node->left, word);
        else if (word > node->word)
            node->right = insert(node->right, word);
        else {
            node->increment();
            return node;
        }
        
        return balance(node);
    }
    
    Node* minValueNode(Node* node) {
        Node* current = node;
        while (current && current->left)
            current = current->left;
        return current;
    }
    
    Node* remove(Node* root, const string& word) {
        if (!root) return root;
        
        if (word < root->word)
            root->left = remove(root->left, word);
        else if (word > root->word)
            root->right = remove(root->right, word);
        else {
            if (root->count > 1) {
                root->decrement();
                return root;
            }
            
            if (!root->left || !root->right) {
                Node* temp = root->left ? root->left : root->right;
                if (!temp) {
                    temp = root;
                    root = nullptr;
                } else {
                    *root = *temp;
                }
                delete temp;
            } else {
                Node* temp = minValueNode(root->right);
                root->word = temp->word;
                root->count = temp->count;
                root->right = remove(root->right, temp->word);
            }
        }
        
        if (!root) return root;
        return balance(root);
    }
    
    void inOrderTraversal(Node* n, vector<pair<string, int>>& result) const {
        if (!n) return;
        inOrderTraversal(n->left, result);
        result.emplace_back(n->word, n->count);
        inOrderTraversal(n->right, result);
    }
    
    void clear(Node* n) {
        if (!n) return;
        clear(n->left);
        clear(n->right);
        delete n;
    }
    
public:
    Tree() : root(nullptr) {}
    ~Tree() { clear(root); }
    
    void Add(const string& word) { root = insert(root, word); }
    void Remove(const string& word) { root = remove(root, word); }
    bool Find(const string& word) const {
        Node* current = root;
        while (current) {
            if (word == current->word) return true;
            if (word < current->word) current = current->left;
            else current = current->right;
        }
        return false;
    }
    
    int getWordCount(const string& word) const {
        Node* current = root;
        while (current) {
            if (word == current->word) return current->count;
            if (word < current->word) current = current->left;
            else current = current->right;
        }
        return 0;
    }
    
    vector<pair<string, int>> getAllWords() const {
        vector<pair<string, int>> result;
        inOrderTraversal(root, result);
        return result;
    }
    
    int getTotalWordCount() const {
        vector<pair<string, int>> words = getAllWords();
        int total = 0;
        for (const auto& p : words) total += p.second;
        return total;
    }
    
    int getDistinctWordCount() const {
        vector<pair<string, int>> words = getAllWords();
        return words.size();
    }
    
    map<int, vector<string>> getWordsByFrequency() const {
        vector<pair<string, int>> words = getAllWords();
        map<int, vector<string>> freqMap;
        for (const auto& p : words) {
            freqMap[p.second].push_back(p.first);
        }
        return freqMap;
    }
    
    vector<string> getMostFrequentWords(int n) const {
        vector<pair<string, int>> words = getAllWords();
        sort(words.begin(), words.end(), [](const auto& a, const auto& b) {
            return a.second > b.second || (a.second == b.second && a.first < b.first);
        });
        vector<string> result;
        for (int i = 0; i < min(n, (int)words.size()); i++) {
            result.push_back(words[i].first);
        }
        return result;
    }
    
    Node* getRoot() const { return root; }
};

ostream& operator<<(ostream& os, const Tree& tree) {
    auto words = tree.getAllWords();
    os << "Word Count Statistics:\n";
    os << "Total words: " << tree.getTotalWordCount() << "\n";
    os << "Distinct words: " << tree.getDistinctWordCount() << "\n\n";
    os << "Word frequencies:\n";
    for (const auto& p : words) {
        os << p.first << ": " << p.second << "\n";
    }
    return os;
}
