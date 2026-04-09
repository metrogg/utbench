#include <vector>
#include <algorithm>
#include <cassert>
#include <map>
#include <set>

using namespace std;

// Enhanced Segment Tree Node with lazy propagation
struct SegmentTreeNode {
    int sum;
    int lazy;  // For lazy propagation
    int left, right;
    SegmentTreeNode *left_child, *right_child;
    
    SegmentTreeNode() : sum(0), lazy(0), left(0), right(0), 
                       left_child(nullptr), right_child(nullptr) {}
};

class EnhancedSegmentTree {
private:
    SegmentTreeNode* root;
    
    // Build the segment tree recursively
    SegmentTreeNode* build(const vector<int>& points, int left_idx, int right_idx) {
        SegmentTreeNode* node = new SegmentTreeNode();
        node->left = points[left_idx];
        node->right = points[right_idx];
        
        if (left_idx != right_idx) {
            int mid = left_idx + (right_idx - left_idx) / 2;
            node->left_child = build(points, left_idx, mid);
            node->right_child = build(points, mid + 1, right_idx);
        }
        
        return node;
    }
    
    // Push down lazy updates to children
    void push_down(SegmentTreeNode* node) {
        if (node->lazy != 0 && node->left_child) {
            node->left_child->sum += node->lazy;
            node->left_child->lazy += node->lazy;
            
            node->right_child->sum += node->lazy;
            node->right_child->lazy += node->lazy;
            
            node->lazy = 0;
        }
    }
    
    // Range update with lazy propagation
    void range_add(SegmentTreeNode* node, int l, int r, int val) {
        if (node->right < l || node->left > r) return;
        
        if (l <= node->left && node->right <= r) {
            node->sum += val;
            node->lazy += val;
            return;
        }
        
        push_down(node);
        range_add(node->left_child, l, r, val);
        range_add(node->right_child, l, r, val);
    }
    
    // Point query with lazy propagation
    int point_query(SegmentTreeNode* node, int point) {
        if (node->left > point || node->right < point) return 0;
        
        if (node->left == node->right) {
            return node->sum;
        }
        
        push_down(node);
        if (point <= node->left_child->right) {
            return point_query(node->left_child, point);
        } else {
            return point_query(node->right_child, point);
        }
    }
    
public:
    EnhancedSegmentTree(const vector<int>& points) {
        if (!points.empty()) {
            root = build(points, 0, points.size() - 1);
        } else {
            root = nullptr;
        }
    }
    
    // Add value to all points in range [l, r]
    void add_range(int l, int r, int val = 1) {
        if (root) range_add(root, l, r, val);
    }
    
    // Query the value at a specific point
    int query_point(int point) {
        if (root) return point_query(root, point);
        return 0;
    }
    
    // Batch query multiple points
    vector<int> query_points(const vector<int>& points) {
        vector<int> results;
        for (int point : points) {
            results.push_back(query_point(point));
        }
        return results;
    }
};

// Function to process flower blooms and person queries
vector<int> process_flower_blooms(const vector<vector<int>>& flowers, 
                                 const vector<int>& persons) {
    // Collect all unique points
    set<int> unique_points;
    for (const auto& flower : flowers) {
        unique_points.insert(flower[0]);
        unique_points.insert(flower[1]);
    }
    for (int person : persons) {
        unique_points.insert(person);
    }
    
    vector<int> sorted_points(unique_points.begin(), unique_points.end());
    EnhancedSegmentTree st(sorted_points);
    
    // Apply flower bloom intervals
    for (const auto& flower : flowers) {
        st.add_range(flower[0], flower[1]);
    }
    
    // Query all person points
    return st.query_points(persons);
}
