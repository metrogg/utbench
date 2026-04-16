#include <vector>
#include <algorithm>
#include <climits>
using namespace std;

const int MAX_DIMENSIONS = 5;
const int MAX_NODES = 200000;

struct Point {
    int dimensions[MAX_DIMENSIONS];
    
    bool operator<(const Point& other) const {
        for (int i = 0; i < MAX_DIMENSIONS; ++i) {
            if (dimensions[i] != other.dimensions[i]) {
                return dimensions[i] < other.dimensions[i];
            }
        }
        return false;
    }
};

struct SegmentTreeNode {
    int left, right;
    Point max_point;
    int count;
    Point min_point;
};

class MultiDimensionalSegmentTree {
private:
    vector<SegmentTreeNode> tree;
    vector<Point> points;
    
    void build(int node, int left, int right) {
        tree[node].left = left;
        tree[node].right = right;
        
        if (left == right) {
            tree[node].max_point = points[left];
            tree[node].min_point = points[left];
            tree[node].count = 1;
            return;
        }
        
        int mid = (left + right) / 2;
        build(2 * node, left, mid);
        build(2 * node + 1, mid + 1, right);
        
        // Update max and min points
        for (int i = 0; i < MAX_DIMENSIONS; ++i) {
            tree[node].max_point.dimensions[i] = max(
                tree[2 * node].max_point.dimensions[i],
                tree[2 * node + 1].max_point.dimensions[i]
            );
            tree[node].min_point.dimensions[i] = min(
                tree[2 * node].min_point.dimensions[i],
                tree[2 * node + 1].min_point.dimensions[i]
            );
        }
        tree[node].count = tree[2 * node].count + tree[2 * node + 1].count;
    }
    
    int query(int node, const Point& upper_bounds) {
        // Check if all dimensions of max_point are <= upper_bounds
        bool all_leq = true;
        for (int i = 0; i < MAX_DIMENSIONS; ++i) {
            if (tree[node].max_point.dimensions[i] > upper_bounds.dimensions[i]) {
                all_leq = false;
                break;
            }
        }
        
        if (all_leq) {
            return tree[node].count;
        }
        
        // Check if any dimension of min_point is > upper_bounds
        bool any_gt = false;
        for (int i = 0; i < MAX_DIMENSIONS; ++i) {
            if (tree[node].min_point.dimensions[i] > upper_bounds.dimensions[i]) {
                any_gt = true;
                break;
            }
        }
        
        if (any_gt) {
            return 0;
        }
        
        if (tree[node].left == tree[node].right) {
            return 0;
        }
        
        return query(2 * node, upper_bounds) + query(2 * node + 1, upper_bounds);
    }

public:
    MultiDimensionalSegmentTree(const vector<Point>& input_points) : points(input_points) {
        if (points.empty()) {
            throw invalid_argument("Input points cannot be empty");
        }
        
        tree.resize(4 * points.size() + 10);
        build(1, 0, points.size() - 1);
    }
    
    int count_points_in_hyperrectangle(const Point& upper_bounds) {
        return query(1, upper_bounds);
    }
    
    void update_point(int index, const Point& new_point) {
        if (index < 0 || index >= points.size()) {
            throw out_of_range("Index out of bounds");
        }
        
        points[index] = new_point;
        update(1, index, new_point);
    }

private:
    void update(int node, int index, const Point& new_point) {
        if (tree[node].left == tree[node].right) {
            tree[node].max_point = new_point;
            tree[node].min_point = new_point;
            return;
        }
        
        int mid = (tree[node].left + tree[node].right) / 2;
        if (index <= mid) {
            update(2 * node, index, new_point);
        } else {
            update(2 * node + 1, index, new_point);
        }
        
        // Update max and min points
        for (int i = 0; i < MAX_DIMENSIONS; ++i) {
            tree[node].max_point.dimensions[i] = max(
                tree[2 * node].max_point.dimensions[i],
                tree[2 * node + 1].max_point.dimensions[i]
            );
            tree[node].min_point.dimensions[i] = min(
                tree[2 * node].min_point.dimensions[i],
                tree[2 * node + 1].min_point.dimensions[i]
            );
        }
    }
};
