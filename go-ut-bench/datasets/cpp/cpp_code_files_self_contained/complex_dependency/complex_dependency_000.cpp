#include <unordered_map>
#include <vector>
#include <queue>
#include <algorithm>
using namespace std;

struct Node {
    int data;
    Node* left;
    Node* right;
    Node(int val) : data(val), left(nullptr), right(nullptr) {}
};

// Recursive solution for Largest Independent Set (LIS)
int lisRecursive(Node* root) {
    if (!root) return 0;
    
    // Case 1: Exclude current node
    int exclude = lisRecursive(root->left) + lisRecursive(root->right);
    
    // Case 2: Include current node
    int include = 1;
    if (root->left) {
        include += lisRecursive(root->left->left) + lisRecursive(root->left->right);
    }
    if (root->right) {
        include += lisRecursive(root->right->left) + lisRecursive(root->right->right);
    }
    
    return max(exclude, include);
}

// Memoized solution using dynamic programming
int lisMemoized(Node* root, unordered_map<Node*, int>& dp) {
    if (!root) return 0;
    if (dp.find(root) != dp.end()) return dp[root];
    
    int exclude = lisMemoized(root->left, dp) + lisMemoized(root->right, dp);
    
    int include = 1;
    if (root->left) {
        include += lisMemoized(root->left->left, dp) + lisMemoized(root->left->right, dp);
    }
    if (root->right) {
        include += lisMemoized(root->right->left, dp) + lisMemoized(root->right->right, dp);
    }
    
    return dp[root] = max(exclude, include);
}

// Iterative solution using post-order traversal
int lisIterative(Node* root) {
    if (!root) return 0;
    
    unordered_map<Node*, int> dp;
    vector<Node*> postOrder;
    queue<Node*> q;
    q.push(root);
    
    // Generate post-order traversal using BFS and stack
    while (!q.empty()) {
        Node* current = q.front();
        q.pop();
        postOrder.push_back(current);
        if (current->right) q.push(current->right);
        if (current->left) q.push(current->left);
    }
    reverse(postOrder.begin(), postOrder.end());
    
    for (Node* node : postOrder) {
        int exclude = 0;
        if (node->left) exclude += dp[node->left];
        if (node->right) exclude += dp[node->right];
        
        int include = 1;
        if (node->left) {
            if (node->left->left) include += dp[node->left->left];
            if (node->left->right) include += dp[node->left->right];
        }
        if (node->right) {
            if (node->right->left) include += dp[node->right->left];
            if (node->right->right) include += dp[node->right->right];
        }
        
        dp[node] = max(exclude, include);
    }
    
    return dp[root];
}

// Helper function to build a binary tree from level order traversal
Node* buildTree(const vector<int>& nodes) {
    if (nodes.empty() || nodes[0] == -1) return nullptr;
    
    Node* root = new Node(nodes[0]);
    queue<Node*> q;
    q.push(root);
    
    int i = 1;
    while (!q.empty() && i < nodes.size()) {
        Node* current = q.front();
        q.pop();
        
        if (i < nodes.size() && nodes[i] != -1) {
            current->left = new Node(nodes[i]);
            q.push(current->left);
        }
        i++;
        
        if (i < nodes.size() && nodes[i] != -1) {
            current->right = new Node(nodes[i]);
            q.push(current->right);
        }
        i++;
    }
    
    return root;
}

// Helper function to delete the tree and free memory
void deleteTree(Node* root) {
    if (!root) return;
    deleteTree(root->left);
    deleteTree(root->right);
    delete root;
}
