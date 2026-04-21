#include <vector>
#include <stdexcept>

class TreeNode {
public:
    int value;
    TreeNode* left;
    TreeNode* right;
    
    TreeNode(int val) : value(val), left(nullptr), right(nullptr) {}
};

class BinarySearchTree {
private:
    TreeNode* root;
    
    TreeNode* insertRecursive(TreeNode* node, int value) {
        if (node == nullptr) {
            return new TreeNode(value);
        }
        if (value < node->value) {
            node->left = insertRecursive(node->left, value);
        } else {
            node->right = insertRecursive(node->right, value);
        }
        return node;
    }
    
    bool searchRecursive(TreeNode* node, int value) {
        if (node == nullptr) {
            return false;
        }
        if (node->value == value) {
            return true;
        }
        if (value < node->value) {
            return searchRecursive(node->left, value);
        }
        return searchRecursive(node->right, value);
    }
    
    void inorderRecursive(TreeNode* node, std::vector<int>& result) {
        if (node != nullptr) {
            inorderRecursive(node->left, result);
            result.push_back(node->value);
            inorderRecursive(node->right, result);
        }
    }
    
public:
    BinarySearchTree() : root(nullptr) {}
    
    void insert(int value) {
        root = insertRecursive(root, value);
    }
    
    bool search(int value) {
        return searchRecursive(root, value);
    }
    
    std::vector<int> inorderTraversal() {
        std::vector<int> result;
        inorderRecursive(root, result);
        return result;
    }
};