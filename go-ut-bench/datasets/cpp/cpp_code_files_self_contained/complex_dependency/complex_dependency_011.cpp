#include <vector>
#include <stack>
#include <queue>
using namespace std;

// TreeNode structure with additional methods
struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    
    // Helper method to create a complete binary tree from vector
    static TreeNode* createTree(const vector<int>& values, int index = 0) {
        if (index >= values.size() || values[index] == -1) {
            return nullptr;
        }
        TreeNode* node = new TreeNode(values[index]);
        node->left = createTree(values, 2 * index + 1);
        node->right = createTree(values, 2 * index + 2);
        return node;
    }
    
    // Helper method to delete tree and free memory
    static void deleteTree(TreeNode* root) {
        if (!root) return;
        deleteTree(root->left);
        deleteTree(root->right);
        delete root;
    }
};

// Enhanced spiral traversal with multiple output options
vector<vector<int>> spiralTraversal(TreeNode* root, bool includeLevelMarkers = false) {
    vector<vector<int>> result;
    if (!root) return result;
    
    stack<TreeNode*> currentLevel;
    stack<TreeNode*> nextLevel;
    currentLevel.push(root);
    
    bool leftToRight = true;
    int level = 0;
    
    while (!currentLevel.empty()) {
        vector<int> currentLevelValues;
        int levelSize = currentLevel.size();
        
        for (int i = 0; i < levelSize; i++) {
            TreeNode* node = currentLevel.top();
            currentLevel.pop();
            currentLevelValues.push_back(node->val);
            
            if (leftToRight) {
                if (node->left) nextLevel.push(node->left);
                if (node->right) nextLevel.push(node->right);
            } else {
                if (node->right) nextLevel.push(node->right);
                if (node->left) nextLevel.push(node->left);
            }
        }
        
        if (includeLevelMarkers) {
            result.push_back({level});
        }
        result.push_back(currentLevelValues);
        
        leftToRight = !leftToRight;
        swap(currentLevel, nextLevel);
        level++;
    }
    
    return result;
}

// Additional tree analysis functions
int treeHeight(TreeNode* root) {
    if (!root) return 0;
    return 1 + max(treeHeight(root->left), treeHeight(root->right));
}

bool isPerfectTree(TreeNode* root) {
    if (!root) return true;
    queue<TreeNode*> q;
    q.push(root);
    bool mustBeLeaf = false;
    
    while (!q.empty()) {
        TreeNode* node = q.front();
        q.pop();
        
        if (mustBeLeaf && (node->left || node->right)) {
            return false;
        }
        
        if (node->left && node->right) {
            q.push(node->left);
            q.push(node->right);
        } else if (node->left || node->right) {
            mustBeLeaf = true;
            if (node->left) q.push(node->left);
            if (node->right) q.push(node->right);
        } else {
            mustBeLeaf = true;
        }
    }
    return true;
}
