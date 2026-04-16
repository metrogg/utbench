#include <vector>
#include <algorithm>
#include <climits>

using namespace std;

const int INF = INT_MAX;

struct SegmentNode {
    int max_add;    // Maximum value to add
    int min_set;    // Minimum value to set
    int lazy_max;   // Lazy propagation for max operations
    int lazy_min;   // Lazy propagation for min operations
    
    SegmentNode() : max_add(-INF), min_set(INF), lazy_max(-INF), lazy_min(INF) {}
};

class AdvancedWallBuilder {
private:
    vector<SegmentNode> segment_tree;
    int size;
    
    void push_lazy(int node, int node_left, int node_right) {
        if (node_left == node_right) return;
        
        // Propagate max operations
        if (segment_tree[node].lazy_max != -INF) {
            int left_child = 2 * node;
            int right_child = 2 * node + 1;
            
            segment_tree[left_child].max_add = max(segment_tree[left_child].max_add, segment_tree[node].lazy_max);
            segment_tree[left_child].lazy_max = max(segment_tree[left_child].lazy_max, segment_tree[node].lazy_max);
            
            segment_tree[right_child].max_add = max(segment_tree[right_child].max_add, segment_tree[node].lazy_max);
            segment_tree[right_child].lazy_max = max(segment_tree[right_child].lazy_max, segment_tree[node].lazy_max);
            
            segment_tree[node].lazy_max = -INF;
        }
        
        // Propagate min operations
        if (segment_tree[node].lazy_min != INF) {
            int left_child = 2 * node;
            int right_child = 2 * node + 1;
            
            segment_tree[left_child].min_set = min(segment_tree[left_child].min_set, segment_tree[node].lazy_min);
            segment_tree[left_child].lazy_min = min(segment_tree[left_child].lazy_min, segment_tree[node].lazy_min);
            
            segment_tree[right_child].min_set = min(segment_tree[right_child].min_set, segment_tree[node].lazy_min);
            segment_tree[right_child].lazy_min = min(segment_tree[right_child].lazy_min, segment_tree[node].lazy_min);
            
            segment_tree[node].lazy_min = INF;
        }
    }
    
public:
    AdvancedWallBuilder(int n) {
        size = 1;
        while (size < n) size <<= 1;
        segment_tree.resize(2 * size);
    }
    
    void update_range(int node, int node_left, int node_right, int left, int right, int op, int value) {
        push_lazy(node, node_left, node_right);
        
        if (node_right < left || node_left > right) return;
        
        if (left <= node_left && node_right <= right) {
            if (op == 1) {  // Max operation (ensure at least value)
                segment_tree[node].max_add = max(segment_tree[node].max_add, value);
                segment_tree[node].lazy_max = max(segment_tree[node].lazy_max, value);
            } else {  // Min operation (ensure at most value)
                segment_tree[node].min_set = min(segment_tree[node].min_set, value);
                segment_tree[node].lazy_min = min(segment_tree[node].lazy_min, value);
            }
            push_lazy(node, node_left, node_right);
            return;
        }
        
        int mid = (node_left + node_right) / 2;
        update_range(2 * node, node_left, mid, left, right, op, value);
        update_range(2 * node + 1, mid + 1, node_right, left, right, op, value);
    }
    
    vector<int> get_final_heights(int n) {
        vector<int> result(n);
        for (int i = 0; i < n; ++i) {
            int node = size + i;
            int current_max = -INF;
            int current_min = INF;
            
            // Traverse up the tree to collect all operations affecting this position
            while (node >= 1) {
                current_max = max(current_max, segment_tree[node].max_add);
                current_min = min(current_min, segment_tree[node].min_set);
                node /= 2;
            }
            
            // The final height is constrained by both max and min operations
            if (current_max == -INF && current_min == INF) {
                result[i] = 0;  // Default height if no operations
            } else if (current_max == -INF) {
                result[i] = current_min;
            } else if (current_min == INF) {
                result[i] = current_max;
            } else {
                result[i] = min(current_max, current_min);
            }
        }
        return result;
    }
};

vector<int> build_advanced_wall(int n, const vector<int>& ops, 
                               const vector<int>& lefts, 
                               const vector<int>& rights, 
                               const vector<int>& heights) {
    AdvancedWallBuilder builder(n);
    for (size_t i = 0; i < ops.size(); ++i) {
        builder.update_range(1, 0, n - 1, lefts[i], rights[i], ops[i], heights[i]);
    }
    return builder.get_final_heights(n);
}
