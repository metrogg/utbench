#include <vector>
#include <map>
#include <utility>
#include <queue>
#include <iomanip>

using namespace std;

class GridPathFinder {
private:
    const int MOD = 1000000007;
    const int dirs[4][2] = {{1,0}, {-1,0}, {0,1}, {0,-1}};

public:
    // Enhanced DP solution with path tracking capability
    int findPaths(int m, int n, int N, int i, int j, bool trackPaths = false) {
        if (N == 0) return 0;
        
        // 3D DP array: dp[step][row][col]
        vector<vector<vector<int>>> dp(N+1, vector<vector<int>>(m, vector<int>(n, 0)));
        
        // Optional: Store paths for visualization (only for small N due to memory)
        map<pair<int, int>, vector<vector<pair<int, int>>>> pathMap;
        
        for (int step = 1; step <= N; ++step) {
            for (int row = 0; row < m; ++row) {
                for (int col = 0; col < n; ++col) {
                    for (const auto& d : dirs) {
                        int new_row = row + d[0];
                        int new_col = col + d[1];
                        
                        if (new_row < 0 || new_row >= m || new_col < 0 || new_col >= n) {
                            dp[step][row][col] = (dp[step][row][col] + 1) % MOD;
                            
                            if (trackPaths && step <= 3) { // Limit to small steps for demo
                                vector<pair<int, int>> path = {{row, col}};
                                pathMap[{row, col}].push_back(path);
                            }
                        } else {
                            dp[step][row][col] = (dp[step][row][col] + dp[step-1][new_row][new_col]) % MOD;
                            
                            if (trackPaths && step <= 3) {
                                for (auto& p : pathMap[{new_row, new_col}]) {
                                    if (p.size() == step-1) {
                                        vector<pair<int, int>> new_path = p;
                                        new_path.push_back({row, col});
                                        pathMap[{row, col}].push_back(new_path);
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        
        if (trackPaths && N <= 3) {
            cout << "Possible paths for visualization (N=" << N << "):" << endl;
            for (const auto& path : pathMap[{i, j}]) {
                if (path.size() == N) {
                    cout << "Path: ";
                    for (const auto& p : path) {
                        cout << "(" << p.first << "," << p.second << ") ";
                    }
                    cout << "-> exit" << endl;
                }
            }
        }
        
        return dp[N][i][j];
    }
    
    // BFS solution for comparison (not efficient for large N)
    int findPathsBFS(int m, int n, int N, int i, int j) {
        if (N == 0) return 0;
        
        queue<tuple<int, int, int>> q;
        q.push({i, j, 0});
        int count = 0;
        
        while (!q.empty()) {
            auto [x, y, steps] = q.front();
            q.pop();
            
            if (steps > N) continue;
            
            if (x < 0 || x >= m || y < 0 || y >= n) {
                count = (count + 1) % MOD;
                continue;
            }
            
            if (steps == N) continue;
            
            for (const auto& d : dirs) {
                q.push({x + d[0], y + d[1], steps + 1});
            }
        }
        
        return count;
    }
};
