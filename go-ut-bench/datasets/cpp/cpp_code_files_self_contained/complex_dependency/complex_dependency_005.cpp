#include <vector>
#include <string>
#include <climits>
#include <algorithm>
#include <unordered_map>

using namespace std;

const int INF = INT_MAX;

class BracketBalancer {
private:
    string sequence;
    vector<vector<int>> dp;
    vector<vector<int>> split_point;
    unordered_map<char, char> bracket_pairs = {
        {'(', ')'}, {'[', ']'}, {'{', '}'}, 
        {'<', '>'}, /* Add more bracket types if needed */
    };

    bool is_opening(char c) {
        return bracket_pairs.find(c) != bracket_pairs.end();
    }

    bool is_closing(char c) {
        for (auto& pair : bracket_pairs) {
            if (pair.second == c) return true;
        }
        return false;
    }

    bool is_matching_pair(char open, char close) {
        auto it = bracket_pairs.find(open);
        return it != bracket_pairs.end() && it->second == close;
    }

public:
    BracketBalancer(const string& s) : sequence(s) {
        int n = sequence.size();
        dp.resize(n, vector<int>(n, INF));
        split_point.resize(n, vector<int>(n, -1));
    }

    int find_min_additions() {
        int n = sequence.size();
        if (n == 0) return 0;

        for (int i = 0; i < n; ++i) {
            dp[i][i] = 1;  // Single bracket needs a pair
        }

        for (int length = 2; length <= n; ++length) {
            for (int i = 0; i + length - 1 < n; ++i) {
                int j = i + length - 1;

                // Case 1: The current brackets match
                if (is_matching_pair(sequence[i], sequence[j])) {
                    if (length == 2) {
                        dp[i][j] = 0;
                    } else {
                        dp[i][j] = dp[i+1][j-1];
                    }
                }

                // Case 2: Try all possible splits
                for (int k = i; k < j; ++k) {
                    if (dp[i][j] > dp[i][k] + dp[k+1][j]) {
                        dp[i][j] = dp[i][k] + dp[k+1][j];
                        split_point[i][j] = k;
                    }
                }
            }
        }

        return dp[0][n-1];
    }

    string get_balanced_sequence() {
        int n = sequence.size();
        if (n == 0) return "";

        find_min_additions();  // Ensure DP table is filled
        return reconstruct(0, n-1);
    }

private:
    string reconstruct(int i, int j) {
        if (i > j) return "";
        if (i == j) {
            if (is_opening(sequence[i])) {
                return string(1, sequence[i]) + bracket_pairs[sequence[i]];
            } else if (is_closing(sequence[j])) {
                for (auto& pair : bracket_pairs) {
                    if (pair.second == sequence[j]) {
                        return string(1, pair.first) + sequence[j];
                    }
                }
            }
            return "";  // Not a bracket character
        }

        if (is_matching_pair(sequence[i], sequence[j]) && 
            (i+1 == j || dp[i][j] == dp[i+1][j-1])) {
            return sequence[i] + reconstruct(i+1, j-1) + sequence[j];
        }

        int k = split_point[i][j];
        if (k != -1) {
            return reconstruct(i, k) + reconstruct(k+1, j);
        }

        return "";  // Shouldn't reach here for valid inputs
    }
};
