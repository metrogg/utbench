#include <vector>
#include <algorithm>
#include <climits>

using namespace std;

class SegmentTree {
private:
    struct Node {
        int l, r;
        long long sum, min_val, max_val;
        long long lazy;
        Node *left, *right;
        
        Node(int s, int e) : l(s), r(e), sum(0), min_val(0), max_val(0), 
                            lazy(0), left(nullptr), right(nullptr) {}
    };

    Node* root;
    
    void push_down(Node* node) {
        if (node->lazy != 0 && node->left) {
            node->left->sum += node->lazy * (node->left->r - node->left->l + 1);
            node->left->min_val += node->lazy;
            node->left->max_val += node->lazy;
            node->left->lazy += node->lazy;
            
            node->right->sum += node->lazy * (node->right->r - node->right->l + 1);
            node->right->min_val += node->lazy;
            node->right->max_val += node->lazy;
            node->right->lazy += node->lazy;
            
            node->lazy = 0;
        }
    }
    
    void build(Node* node, const vector<long long>& data, int s, int e) {
        if (s == e) {
            node->sum = node->min_val = node->max_val = data[s];
            return;
        }
        
        int mid = (s + e) / 2;
        node->left = new Node(s, mid);
        node->right = new Node(mid + 1, e);
        
        build(node->left, data, s, mid);
        build(node->right, data, mid + 1, e);
        
        node->sum = node->left->sum + node->right->sum;
        node->min_val = min(node->left->min_val, node->right->min_val);
        node->max_val = max(node->left->max_val, node->right->max_val);
    }
    
    void update_range(Node* node, int l, int r, long long val) {
        if (node->r < l || node->l > r) return;
        
        if (l <= node->l && node->r <= r) {
            node->sum += val * (node->r - node->l + 1);
            node->min_val += val;
            node->max_val += val;
            node->lazy += val;
            return;
        }
        
        push_down(node);
        update_range(node->left, l, r, val);
        update_range(node->right, l, r, val);
        
        node->sum = node->left->sum + node->right->sum;
        node->min_val = min(node->left->min_val, node->right->min_val);
        node->max_val = max(node->left->max_val, node->right->max_val);
    }
    
    long long query_sum(Node* node, int l, int r) {
        if (node->r < l || node->l > r) return 0;
        if (l <= node->l && node->r <= r) return node->sum;
        
        push_down(node);
        return query_sum(node->left, l, r) + query_sum(node->right, l, r);
    }
    
    long long query_min(Node* node, int l, int r) {
        if (node->r < l || node->l > r) return LLONG_MAX;
        if (l <= node->l && node->r <= r) return node->min_val;
        
        push_down(node);
        return min(query_min(node->left, l, r), query_min(node->right, l, r));
    }
    
    long long query_max(Node* node, int l, int r) {
        if (node->r < l || node->l > r) return LLONG_MIN;
        if (l <= node->l && node->r <= r) return node->max_val;
        
        push_down(node);
        return max(query_max(node->left, l, r), query_max(node->right, l, r));
    }

public:
    SegmentTree(const vector<long long>& data) {
        if (data.empty()) {
            root = nullptr;
            return;
        }
        root = new Node(0, data.size() - 1);
        build(root, data, 0, data.size() - 1);
    }
    
    ~SegmentTree() {
        // In a real implementation, we'd need a proper destructor
        // This is simplified for the benchmark
    }
    
    void range_add(int l, int r, long long val) {
        if (root) update_range(root, l, r, val);
    }
    
    long long range_sum(int l, int r) {
        return root ? query_sum(root, l, r) : 0;
    }
    
    long long range_min(int l, int r) {
        return root ? query_min(root, l, r) : LLONG_MAX;
    }
    
    long long range_max(int l, int r) {
        return root ? query_max(root, l, r) : LLONG_MIN;
    }
};
