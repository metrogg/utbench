#include <cstring>
#include <vector>
#include <algorithm>

using namespace std;

// Enhanced palindrome analyzer with additional features
vector<int> analyze_palindromes(const string& str) {
    int n = str.length();
    if (n == 0) return vector<int>(1, 0);
    
    // DP tables
    vector<vector<bool>> is_palindrome(n, vector<bool>(n, false));
    vector<vector<int>> k_palindrome(n, vector<int>(n, 0));
    vector<int> cnt(n+2, 0); // cnt[k] = number of k-palindromes
    
    // Initialize for length 1 and 2 palindromes
    for (int i = 0; i < n; i++) {
        is_palindrome[i][i] = true;
        k_palindrome[i][i] = 1;
        cnt[1]++;
    }
    
    for (int i = 0; i < n-1; i++) {
        if (str[i] == str[i+1]) {
            is_palindrome[i][i+1] = true;
            k_palindrome[i][i+1] = 1;
            cnt[1]++;
        }
    }
    
    // Fill DP tables for lengths > 2
    for (int len = 2; len < n; len++) {
        for (int i = 0; i + len < n; i++) {
            int j = i + len;
            
            // Check if current substring is a palindrome
            if (str[i] == str[j] && is_palindrome[i+1][j-1]) {
                is_palindrome[i][j] = true;
                
                // Calculate k-palindrome value
                int m = (i + j) / 2;
                if ((j - i + 1) % 2 == 1) m--;
                k_palindrome[i][j] = k_palindrome[i][m] + 1;
                cnt[k_palindrome[i][j]]++;
            }
        }
    }
    
    // Calculate cumulative counts
    for (int k = n; k >= 1; k--) {
        cnt[k] += cnt[k+1];
    }
    
    // Return counts for k=1 to k=n
    return vector<int>(cnt.begin()+1, cnt.begin()+n+1);
}

// Additional function to find all unique palindromic substrings
vector<string> find_unique_palindromes(const string& str) {
    int n = str.length();
    vector<vector<bool>> dp(n, vector<bool>(n, false));
    vector<string> result;
    
    for (int i = 0; i < n; i++) {
        dp[i][i] = true;
        result.push_back(str.substr(i, 1));
    }
    
    for (int i = 0; i < n-1; i++) {
        if (str[i] == str[i+1]) {
            dp[i][i+1] = true;
            result.push_back(str.substr(i, 2));
        }
    }
    
    for (int len = 2; len < n; len++) {
        for (int i = 0; i + len < n; i++) {
            int j = i + len;
            if (str[i] == str[j] && dp[i+1][j-1]) {
                dp[i][j] = true;
                result.push_back(str.substr(i, len+1));
            }
        }
    }
    
    // Remove duplicates
    sort(result.begin(), result.end());
    result.erase(unique(result.begin(), result.end()), result.end());
    
    return result;
}
