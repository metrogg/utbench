#include <algorithm>
#include <utility>
#include <vector>
#include <queue>
#include <climits>

using namespace std;

// TreeNode structure
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
};

// Enhanced balance check with additional tree metrics
pair<bool, int> checkTreeBalance(TreeNode* root) {
    if (root == nullptr) {
        return make_pair(true, 0);
    }
    
    auto left = checkTreeBalance(root->left);
    auto right = checkTreeBalance(root->right);
    
    int height = 1 + max(left.second, right.second);
    bool balanced = abs(left.second - right.second) <= 1 && left.first && right.first;
    
    return make_pair(balanced, height);
}

// Additional tree analysis functions
struct TreeAnalysis {
    bool isBalanced;
    int height;
    int minValue;
    int maxValue;
    bool isBST;
};

TreeAnalysis analyzeTree(TreeNode* root) {
    if (root == nullptr) {
        return {true, 0, INT_MAX, INT_MIN, true};
    }
    
    auto left = analyzeTree(root->left);
    auto right = analyzeTree(root->right);
    
    TreeAnalysis result;
    result.height = 1 + max(left.height, right.height);
    result.isBalanced = abs(left.height - right.height) <= 1 && 
                        left.isBalanced && right.isBalanced;
    result.minValue = min({root->val, left.minValue, right.minValue});
    result.maxValue = max({root->val, left.maxValue, right.maxValue});
    result.isBST = (root->val > left.maxValue) && 
                  (root->val < right.minValue) && 
                  left.isBST && right.isBST;
    
    return result;
}

// Helper function to build tree from level order traversal
TreeNode* buildTree(const vector<int>& nodes) {
    if (nodes.empty() || nodes[0] == -1) return nullptr;
    
    TreeNode* root = new TreeNode(nodes[0]);
    queue<TreeNode*> q;
    q.push(root);
    
    int i = 1;
    while (!q.empty() && i < nodes.size()) {
        TreeNode* current = q.front();
        q.pop();
        
        if (i < nodes.size() && nodes[i] != -1) {
            current->left = new TreeNode(nodes[i]);
            q.push(current->left);
        }
        i++;
        
        if (i < nodes.size() && nodes[i] != -1) {
            current->right = new TreeNode(nodes[i]);
            q.push(current->right);
        }
        i++;
    }
    
    return root;
}
