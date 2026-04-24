#include <vector>
#include <functional>
#include <unordered_set>

using namespace std;

/**
 * Enhanced version of the original solve function with additional features:
 * 1. Supports tracking the solution path
 * 2. Includes memoization for better performance
 * 3. Provides multiple solution finding capability
 * 4. Has configurable depth limit for large inputs
 */
vector<vector<int>> solve_enhanced(const vector<vector<int>>& vectors, int target, 
                                  int max_solutions = 1, int depth_limit = 1000) {
    // Memoization cache to store intermediate results
    struct MemoKey {
        int index;
        int current_xor;
        bool operator==(const MemoKey& other) const {
            return index == other.index && current_xor == other.current_xor;
        }
    };
    
    struct MemoKeyHash {
        size_t operator()(const MemoKey& k) const {
            return hash<int>()(k.index) ^ hash<int>()(k.current_xor);
        }
    };
    
    unordered_map<MemoKey, bool, MemoKeyHash> memo;
    vector<vector<int>> solutions;
    vector<int> current_path;
    
    // Recursive lambda with memoization
    function<bool(int, int, bool)> solver = [&](int index, int current, bool find_all) {
        if (index == 0) {
            if (current == target) {
                if (!current_path.empty()) {
                    solutions.push_back(current_path);
                }
                return true;
            }
            return false;
        }
        
        // Check memoization cache
        MemoKey key{index, current};
        if (memo.find(key) != memo.end()) {
            return memo[key];
        }
        
        // Depth limit check
        if (vectors.size() - index > depth_limit) {
            return false;
        }
        
        bool found = false;
        for (int num : vectors[index - 1]) {
            current_path.push_back(num);
            if (solver(index - 1, current ^ num, find_all)) {
                found = true;
                if (!find_all) {
                    current_path.pop_back();
                    break;
                }
            }
            current_path.pop_back();
            
            // Early exit if we've found enough solutions
            if (find_all && solutions.size() >= max_solutions) {
                break;
            }
        }
        
        memo[key] = found;
        return found;
    };
    
    solver(vectors.size(), 0, max_solutions > 1);
    return solutions;
}

/**
 * Simplified version similar to original solve function
 */
bool solve_basic(const vector<vector<int>>& vectors, int target) {
    if (vectors.empty()) return target == 0;
    
    function<bool(int, int)> solver = [&](int index, int current) {
        if (index == 0) {
            return current == target;
        }
        for (int num : vectors[index - 1]) {
            if (solver(index - 1, current ^ num)) {
                return true;
            }
        }
        return false;
    };
    
    return solver(vectors.size(), 0);
}
