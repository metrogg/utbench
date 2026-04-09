#include <vector>
#include <string>
#include <algorithm>
#include <map>
#include <set>

using namespace std;

class AdvancedStringAnalyzer {
public:
    // Find the longest common substring using dynamic programming
    string findLongestCommonSubstring(const string& str1, const string& str2) {
        int m = str1.length();
        int n = str2.length();
        vector<vector<int>> dp(m + 1, vector<int>(n + 1, 0));
        int max_length = 0;
        int end_pos = 0;

        for (int i = 1; i <= m; ++i) {
            for (int j = 1; j <= n; ++j) {
                if (str1[i-1] == str2[j-1]) {
                    dp[i][j] = dp[i-1][j-1] + 1;
                    if (dp[i][j] > max_length) {
                        max_length = dp[i][j];
                        end_pos = i - 1;
                    }
                }
            }
        }

        if (max_length == 0) return "";
        return str1.substr(end_pos - max_length + 1, max_length);
    }

    // Find all common substrings of maximum length
    vector<string> findAllLongestCommonSubstrings(const string& str1, const string& str2) {
        int m = str1.length();
        int n = str2.length();
        vector<vector<int>> dp(m + 1, vector<int>(n + 1, 0));
        int max_length = 0;
        set<string> result_set;

        for (int i = 1; i <= m; ++i) {
            for (int j = 1; j <= n; ++j) {
                if (str1[i-1] == str2[j-1]) {
                    dp[i][j] = dp[i-1][j-1] + 1;
                    if (dp[i][j] > max_length) {
                        max_length = dp[i][j];
                        result_set.clear();
                        result_set.insert(str1.substr(i - max_length, max_length));
                    } else if (dp[i][j] == max_length) {
                        result_set.insert(str1.substr(i - max_length, max_length));
                    }
                }
            }
        }

        return vector<string>(result_set.begin(), result_set.end());
    }

    // Find the longest common subsequence (different from substring)
    string findLongestCommonSubsequence(const string& str1, const string& str2) {
        int m = str1.length();
        int n = str2.length();
        vector<vector<int>> dp(m + 1, vector<int>(n + 1, 0));

        for (int i = 1; i <= m; ++i) {
            for (int j = 1; j <= n; ++j) {
                if (str1[i-1] == str2[j-1]) {
                    dp[i][j] = dp[i-1][j-1] + 1;
                } else {
                    dp[i][j] = max(dp[i-1][j], dp[i][j-1]);
                }
            }
        }

        // Reconstruct the subsequence
        string lcs;
        int i = m, j = n;
        while (i > 0 && j > 0) {
            if (str1[i-1] == str2[j-1]) {
                lcs = str1[i-1] + lcs;
                i--;
                j--;
            } else if (dp[i-1][j] > dp[i][j-1]) {
                i--;
            } else {
                j--;
            }
        }

        return lcs;
    }

    // Find common substrings with length at least k
    vector<string> findCommonSubstringsWithMinLength(const string& str1, const string& str2, int k) {
        vector<string> result;
        int m = str1.length();
        int n = str2.length();
        vector<vector<int>> dp(m + 1, vector<int>(n + 1, 0));

        for (int i = 1; i <= m; ++i) {
            for (int j = 1; j <= n; ++j) {
                if (str1[i-1] == str2[j-1]) {
                    dp[i][j] = dp[i-1][j-1] + 1;
                    if (dp[i][j] >= k) {
                        string candidate = str1.substr(i - dp[i][j], dp[i][j]);
                        // Avoid duplicates
                        if (find(result.begin(), result.end(), candidate) == result.end()) {
                            result.push_back(candidate);
                        }
                    }
                }
            }
        }

        return result;
    }
};
