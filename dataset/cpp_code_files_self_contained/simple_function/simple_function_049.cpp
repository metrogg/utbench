#include <vector>
#include <string>
#include <algorithm>
#include <tuple>

using namespace std;

// Enhanced LCS function that returns multiple results
tuple<string, int, pair<int, int>, pair<int, int>> findLongestCommonSubstring(
    const string& str1, const string& str2) {
    
    int m = str1.length();
    int n = str2.length();
    vector<vector<int>> dp(m + 1, vector<int>(n + 1, 0));
    int max_length = 0;
    int end_pos_str1 = 0;
    int end_pos_str2 = 0;

    for (int i = 1; i <= m; i++) {
        for (int j = 1; j <= n; j++) {
            if (str1[i-1] == str2[j-1]) {
                dp[i][j] = dp[i-1][j-1] + 1;
                if (dp[i][j] > max_length) {
                    max_length = dp[i][j];
                    end_pos_str1 = i;
                    end_pos_str2 = j;
                }
            }
        }
    }

    if (max_length == 0) {
        return make_tuple("", 0, make_pair(-1, -1), make_pair(-1, -1));
    }

    int start_pos_str1 = end_pos_str1 - max_length;
    int start_pos_str2 = end_pos_str2 - max_length;
    string lcs = str1.substr(start_pos_str1, max_length);

    return make_tuple(
        lcs,
        max_length,
        make_pair(start_pos_str1, end_pos_str1 - 1),
        make_pair(start_pos_str2, end_pos_str2 - 1)
    );
}

// Helper function to print the result in a readable format
void printLCSResult(const tuple<string, int, pair<int, int>, pair<int, int>>& result) {
    auto [lcs, length, pos1, pos2] = result;
    
    cout << "Longest Common Substring: \"" << lcs << "\"" << endl;
    cout << "Length: " << length << endl;
    cout << "Positions in String 1: [" << pos1.first << ", " << pos1.second << "]" << endl;
    cout << "Positions in String 2: [" << pos2.first << ", " << pos2.second << "]" << endl;
    cout << "----------------------------------" << endl;
}
