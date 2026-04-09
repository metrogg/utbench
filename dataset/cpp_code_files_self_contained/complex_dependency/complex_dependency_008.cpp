#include <vector>
#include <unordered_map>
#include <climits>
using namespace std;

const int N_MAX = 10000;

class BinaryTreeLCA {
private:
    unordered_map<int, int> indices;
    vector<int> pre_order;
    vector<int> in_order;
    int position_U, position_V;

    // Helper function to build tree from preorder and inorder traversals
    struct TreeNode {
        int val;
        TreeNode* left;
        TreeNode* right;
        TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    };

    TreeNode* buildTree(int pre_start, int in_start, int in_end) {
        if (in_start > in_end) return nullptr;
        
        TreeNode* root = new TreeNode(pre_order[pre_start]);
        int in_index = indices[root->val];
        int left_size = in_index - in_start;
        
        root->left = buildTree(pre_start + 1, in_start, in_index - 1);
        root->right = buildTree(pre_start + left_size + 1, in_index + 1, in_end);
        
        return root;
    }

    TreeNode* findLCA(TreeNode* root, int u, int v) {
        if (!root) return nullptr;
        if (root->val == u || root->val == v) return root;
        
        TreeNode* left = findLCA(root->left, u, v);
        TreeNode* right = findLCA(root->right, u, v);
        
        if (left && right) return root;
        return left ? left : right;
    }

public:
    void setTraversals(const vector<int>& pre, const vector<int>& in) {
        pre_order = pre;
        in_order = in;
        indices.clear();
        for (int i = 0; i < in_order.size(); ++i) {
            indices[in_order[i]] = i;
        }
    }

    string findLCAInfo(int U, int V) {
        if (indices.find(U) == indices.end() && indices.find(V) == indices.end()) {
            return "ERROR: " + to_string(U) + " and " + to_string(V) + " are not found.";
        }
        if (indices.find(U) == indices.end()) {
            return "ERROR: " + to_string(U) + " is not found.";
        }
        if (indices.find(V) == indices.end()) {
            return "ERROR: " + to_string(V) + " is not found.";
        }

        TreeNode* root = buildTree(0, 0, in_order.size() - 1);
        TreeNode* lca = findLCA(root, U, V);

        if (!lca) return "ERROR: No LCA found.";
        if (lca->val == U) {
            return to_string(U) + " is an ancestor of " + to_string(V) + ".";
        }
        if (lca->val == V) {
            return to_string(V) + " is an ancestor of " + to_string(U) + ".";
        }
        return "LCA of " + to_string(U) + " and " + to_string(V) + " is " + to_string(lca->val) + ".";
    }
};
