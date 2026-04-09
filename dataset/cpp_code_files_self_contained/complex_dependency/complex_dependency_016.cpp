#include <map>
#include <queue>
#include <vector>
#include <algorithm>

using namespace std;

struct Node {
    int value;
    Node* left;
    Node* right;
    int level;  // Track node level
};

struct TreeData {
    int parent;
    int child;
    char direction;
};

// Enhanced node creation with level tracking
Node* createNode(int value, int level = 0) {
    Node* node = new Node;
    node->value = value;
    node->left = nullptr;
    node->right = nullptr;
    node->level = level;
    return node;
}

// Build tree with level tracking and additional validation
Node* buildTree(const vector<TreeData>& data) {
    if (data.empty()) return nullptr;

    Node* root = nullptr;
    map<int, Node*> nodeMap;

    for (const auto& entry : data) {
        Node* parent;

        if (nodeMap.find(entry.parent) == nodeMap.end()) {
            parent = createNode(entry.parent);
            nodeMap[entry.parent] = parent;
            if (!root) root = parent;
        } else {
            parent = nodeMap[entry.parent];
        }

        Node* child = createNode(entry.child, parent->level + 1);
        if (entry.direction == 'L') {
            if (parent->left) {
                cerr << "Warning: Overwriting left child of node " << parent->value << endl;
            }
            parent->left = child;
        } else {
            if (parent->right) {
                cerr << "Warning: Overwriting right child of node " << parent->value << endl;
            }
            parent->right = child;
        }
        nodeMap[entry.child] = child;
    }

    return root;
}

// Check if all leaves are at same level (BFS approach)
bool areLeavesAtSameLevel(Node* root) {
    if (!root) return true;

    queue<Node*> q;
    q.push(root);
    int leafLevel = -1;

    while (!q.empty()) {
        Node* current = q.front();
        q.pop();

        if (!current->left && !current->right) {
            if (leafLevel == -1) {
                leafLevel = current->level;
            } else if (current->level != leafLevel) {
                return false;
            }
        }

        if (current->left) q.push(current->left);
        if (current->right) q.push(current->right);
    }

    return true;
}

// Additional function to get tree height
int getTreeHeight(Node* root) {
    if (!root) return 0;
    return 1 + max(getTreeHeight(root->left), getTreeHeight(root->right));
}

// Function to check if tree is balanced
bool isTreeBalanced(Node* root) {
    if (!root) return true;
    
    int leftHeight = getTreeHeight(root->left);
    int rightHeight = getTreeHeight(root->right);
    
    return abs(leftHeight - rightHeight) <= 1 && 
           isTreeBalanced(root->left) && 
           isTreeBalanced(root->right);
}

// Enhanced solve function with more detailed output
void analyzeTree(Node* root) {
    if (!root) {
        cout << "Tree is empty" << endl;
        return;
    }

    bool sameLevel = areLeavesAtSameLevel(root);
    bool balanced = isTreeBalanced(root);
    int height = getTreeHeight(root);

    cout << "Tree Analysis Results:" << endl;
    cout << "1. All leaves at same level: " << (sameLevel ? "Yes" : "No") << endl;
    cout << "2. Tree is balanced: " << (balanced ? "Yes" : "No") << endl;
    cout << "3. Tree height: " << height << endl;
}
