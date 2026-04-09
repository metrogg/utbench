#include <queue>
#include <vector>
#include <functional>
#include <climits>

using namespace std;

template <typename T> 
struct BinaryNode {
    T elem;
    BinaryNode *left;
    BinaryNode *right;
    BinaryNode(T d, BinaryNode *l = nullptr, BinaryNode *r = nullptr) 
        : elem(d), left(l), right(r) {};
};

// Enhanced binary tree class with various traversal methods
template <typename T>
class BinaryTree {
private:
    BinaryNode<T>* root;

    // Helper function for destructor
    void clear(BinaryNode<T>* node) {
        if (node) {
            clear(node->left);
            clear(node->right);
            delete node;
        }
    }

public:
    BinaryTree() : root(nullptr) {}
    BinaryTree(BinaryNode<T>* r) : root(r) {}
    ~BinaryTree() { clear(root); }

    // Level-order traversal (BFS)
    vector<T> levelOrderTraversal() {
        vector<T> result;
        if (!root) return result;

        queue<BinaryNode<T>*> q;
        q.push(root);

        while (!q.empty()) {
            BinaryNode<T>* current = q.front();
            q.pop();
            result.push_back(current->elem);

            if (current->left) q.push(current->left);
            if (current->right) q.push(current->right);
        }
        return result;
    }

    // Pre-order traversal (DFS)
    vector<T> preOrderTraversal() {
        vector<T> result;
        function<void(BinaryNode<T>*)> traverse = [&](BinaryNode<T>* node) {
            if (!node) return;
            result.push_back(node->elem);
            traverse(node->left);
            traverse(node->right);
        };
        traverse(root);
        return result;
    }

    // Find maximum depth of the tree
    int maxDepth() {
        function<int(BinaryNode<T>*)> depth = [&](BinaryNode<T>* node) {
            if (!node) return 0;
            return 1 + max(depth(node->left), depth(node->right));
        };
        return depth(root);
    }

    // Check if tree is balanced
    bool isBalanced() {
        function<bool(BinaryNode<T>*, int&)> check = [&](BinaryNode<T>* node, int& height) {
            if (!node) {
                height = 0;
                return true;
            }

            int leftHeight, rightHeight;
            bool leftBalanced = check(node->left, leftHeight);
            bool rightBalanced = check(node->right, rightHeight);
            height = 1 + max(leftHeight, rightHeight);

            return leftBalanced && rightBalanced && abs(leftHeight - rightHeight) <= 1;
        };

        int height;
        return check(root, height);
    }

    // Find the maximum value in the tree
    T findMax() {
        if (!root) throw runtime_error("Tree is empty");

        T maxVal = root->elem;
        queue<BinaryNode<T>*> q;
        q.push(root);

        while (!q.empty()) {
            BinaryNode<T>* current = q.front();
            q.pop();
            if (current->elem > maxVal) maxVal = current->elem;

            if (current->left) q.push(current->left);
            if (current->right) q.push(current->right);
        }
        return maxVal;
    }
};
