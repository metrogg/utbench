#include <vector>
#include <unordered_set>
#include <algorithm>
#include <string>

using namespace std;

class AlphabetPathFinder {
private:
    vector<vector<char>> matrix;
    vector<vector<bool>> visited;
    unordered_set<char> used_chars;
    int rows, cols;
    int max_length;
    
    void dfs(int x, int y, int current_length) {
        if (x < 0 || x >= rows || y < 0 || y >= cols || 
            visited[x][y] || used_chars.count(matrix[x][y])) {
            return;
        }
        
        visited[x][y] = true;
        used_chars.insert(matrix[x][y]);
        current_length++;
        
        if (current_length > max_length) {
            max_length = current_length;
        }
        
        // Explore all four directions
        dfs(x + 1, y, current_length);
        dfs(x - 1, y, current_length);
        dfs(x, y + 1, current_length);
        dfs(x, y - 1, current_length);
        
        // Backtrack
        visited[x][y] = false;
        used_chars.erase(matrix[x][y]);
    }
    
public:
    AlphabetPathFinder(const vector<vector<char>>& input_matrix) 
        : matrix(input_matrix), 
          rows(input_matrix.size()),
          cols(rows > 0 ? input_matrix[0].size() : 0),
          max_length(0) {
        if (rows > 0 && cols > 0) {
            visited = vector<vector<bool>>(rows, vector<bool>(cols, false));
        }
    }
    
    int findLongestUniquePath() {
        if (rows == 0 || cols == 0) return 0;
        
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                dfs(i, j, 0);
            }
        }
        
        return max_length;
    }
    
    vector<vector<char>> getMatrix() const {
        return matrix;
    }
};
