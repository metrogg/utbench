#include <vector>
#include <algorithm>
#include <climits>
#include <map>

using namespace std;

// Enhanced function to solve the ribbon cutting problem with additional statistics
map<string, int> solve_ribbon_cutting(int n, const vector<int>& pieces) {
    if (pieces.empty()) {
        throw invalid_argument("Error: No piece lengths provided");
    }

    // Initialize DP array
    vector<int> dp(n + 1, INT_MIN);
    dp[0] = 0;  // Base case: 0 pieces for length 0

    // Mark valid pieces
    for (int piece : pieces) {
        if (piece <= n) {
            dp[piece] = 1;
        }
    }

    // Fill DP table
    for (int i = 1; i <= n; ++i) {
        for (int piece : pieces) {
            if (i >= piece && dp[i - piece] != INT_MIN) {
                dp[i] = max(dp[i], dp[i - piece] + 1);
            }
        }
    }

    // Calculate additional statistics
    int max_pieces = dp[n];
    int min_pieces = INT_MAX;
    vector<int> possible_lengths;
    int impossible_lengths = 0;

    for (int i = 1; i <= n; ++i) {
        if (dp[i] != INT_MIN) {
            min_pieces = min(min_pieces, dp[i]);
            possible_lengths.push_back(i);
        } else {
            impossible_lengths++;
        }
    }

    // Prepare results
    map<string, int> results;
    results["max_pieces"] = (max_pieces == INT_MIN) ? -1 : max_pieces;
    results["min_pieces_for_possible"] = (min_pieces == INT_MAX) ? -1 : min_pieces;
    results["possible_lengths_count"] = possible_lengths.size();
    results["impossible_lengths_count"] = impossible_lengths;
    results["largest_possible"] = possible_lengths.empty() ? -1 : possible_lengths.back();
    results["smallest_possible"] = possible_lengths.empty() ? -1 : possible_lengths.front();

    return results;
}
