#include <vector>
#include <string>
#include <map>
#include <stdexcept>
#include <algorithm>

template <typename T>
class TrieMap {
private:
    struct TrieNode {
        std::map<char, TrieNode*> children;
        bool isEnd = false;
        T value;
    };

    TrieNode* root;
    size_t count = 0;

    void clearHelper(TrieNode* node) {
        if (!node) return;
        for (auto& pair : node->children) {
            clearHelper(pair.second);
        }
        delete node;
    }

    TrieNode* cloneHelper(const TrieNode* node) {
        if (!node) return nullptr;
        TrieNode* newNode = new TrieNode();
        newNode->isEnd = node->isEnd;
        newNode->value = node->value;
        for (const auto& pair : node->children) {
            newNode->children[pair.first] = cloneHelper(pair.second);
        }
        return newNode;
    }

public:
    TrieMap() : root(new TrieNode()) {}
    
    TrieMap(const TrieMap& other) : root(cloneHelper(other.root)), count(other.count) {}
    
    ~TrieMap() {
        clearHelper(root);
    }
    
    TrieMap& operator=(const TrieMap& other) {
        if (this != &other) {
            clearHelper(root);
            root = cloneHelper(other.root);
            count = other.count;
        }
        return *this;
    }
    
    bool contains(const std::string& key) const {
        TrieNode* current = root;
        for (char ch : key) {
            if (current->children.find(ch) == current->children.end()) {
                return false;
            }
            current = current->children[ch];
        }
        return current->isEnd;
    }
    
    T& operator[](const std::string& key) {
        TrieNode* current = root;
        for (char ch : key) {
            if (current->children.find(ch) == current->children.end()) {
                current->children[ch] = new TrieNode();
            }
            current = current->children[ch];
        }
        if (!current->isEnd) {
            current->isEnd = true;
            count++;
        }
        return current->value;
    }
    
    void remove(const std::string& key) {
        if (!contains(key)) return;
        
        std::vector<TrieNode*> nodes;
        TrieNode* current = root;
        nodes.push_back(current);
        
        for (char ch : key) {
            current = current->children[ch];
            nodes.push_back(current);
        }
        
        current->isEnd = false;
        count--;
        
        for (int i = nodes.size() - 1; i > 0; --i) {
            if (nodes[i]->children.empty() && !nodes[i]->isEnd) {
                char ch = key[i-1];
                delete nodes[i];
                nodes[i-1]->children.erase(ch);
            } else {
                break;
            }
        }
    }
    
    std::vector<std::string> keysWithPrefix(const std::string& prefix) const {
        std::vector<std::string> result;
        TrieNode* current = root;
        
        for (char ch : prefix) {
            if (current->children.find(ch) == current->children.end()) {
                return result;
            }
            current = current->children[ch];
        }
        
        collectKeys(current, prefix, result);
        return result;
    }
    
    std::string longestPrefixOf(const std::string& query) const {
        TrieNode* current = root;
        std::string prefix;
        std::string longest;
        
        for (char ch : query) {
            if (current->children.find(ch) == current->children.end()) {
                break;
            }
            current = current->children[ch];
            prefix += ch;
            if (current->isEnd) {
                longest = prefix;
            }
        }
        
        return longest;
    }
    
    size_t size() const {
        return count;
    }
    
    bool empty() const {
        return count == 0;
    }
    
    void clear() {
        clearHelper(root);
        root = new TrieNode();
        count = 0;
    }

private:
    void collectKeys(TrieNode* node, const std::string& prefix, std::vector<std::string>& result) const {
        if (node->isEnd) {
            result.push_back(prefix);
        }
        
        for (const auto& pair : node->children) {
            collectKeys(pair.second, prefix + pair.first, result);
        }
    }
};
