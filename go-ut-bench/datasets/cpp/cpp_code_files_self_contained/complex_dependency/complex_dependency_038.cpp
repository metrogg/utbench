#include <vector>
#include <algorithm>
#include <climits>

using namespace std;

class LongestPathFinder {
private:
    vector<vector<int>> matrix;
    vector<vector<int>> dp;
    int rows, cols;
    
    // Directions: up, down, left, right
    const vector<pair<int, int>> directions = {{-1, 0}, {1, 0}, {0, -1}, {0, 1}};
    
    // DFS with memoization
    int dfs(int i, int j) {
        if (dp[i][j] != -1) return dp[i][j];
        
        int max_path = 1;
        for (const auto& dir : directions) {
            int ni = i + dir.first;
            int nj = j + dir.second;
            
            if (ni >= 0 && ni < rows && nj >= 0 && nj < cols && 
                matrix[ni][nj] > matrix[i][j]) {
                max_path = max(max_path, 1 + dfs(ni, nj));
            }
        }
        
        dp[i][j] = max_path;
        return max_path;
    }

public:
    LongestPathFinder(const vector<vector<int>>& input_matrix) 
        : matrix(input_matrix), rows(input_matrix.size()) {
        if (rows == 0) cols = 0;
        else cols = input_matrix[0].size();
        
        dp.resize(rows, vector<int>(cols, -1));
    }
    
    int findLongestIncreasingPath() {
        if (rows == 0 || cols == 0) return 0;
        
        int result = 1;
        for (int i = 0; i < rows; ++i) {
            for (int j = 0; j < cols; ++j) {
                result = max(result, dfs(i, j));
            }
        }
        return result;
    }
    
    // Additional functionality: get the actual path sequence
    vector<int> getLongestPathSequence() {
        if (rows == 0 || cols == 0) return {};
        
        int max_len = findLongestIncreasingPath();
        vector<int> path;
        
        // Find starting position of the longest path
        int start_i = -1, start_j = -1;
        for (int i = 0; i < rows; ++i) {
            for (int j = 0; j < cols; ++j) {
                if (dp[i][j] == max_len) {
                    start_i = i;
                    start_j = j;
                    break;
                }
            }
            if (start_i != -1) break;
        }
        
        // Reconstruct the path
        int current_i = start_i, current_j = start_j;
        path.push_back(matrix[current_i][current_j]);
        
        while (true) {
            bool found_next = false;
            for (const auto& dir : directions) {
                int ni = current_i + dir.first;
                int nj = current_j + dir.second;
                
                if (ni >= 0 && ni < rows && nj >= 0 && nj < cols && 
                    matrix[ni][nj] > matrix[current_i][current_j] && 
                    dp[ni][nj] == dp[current_i][current_j] - 1) {
                    current_i = ni;
                    current_j = nj;
                    path.push_back(matrix[current_i][current_j]);
                    found_next = true;
                    break;
                }
            }
            if (!found_next) break;
        }
        
        return path;
    }
};
