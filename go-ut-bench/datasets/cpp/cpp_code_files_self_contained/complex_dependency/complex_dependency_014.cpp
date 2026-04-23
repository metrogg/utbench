#include <vector>
#include <string>
#include <map>
#include <algorithm>

using namespace std;

// Enhanced pattern matching function that can detect multiple pattern types
map<string, bool> analyze_string_patterns(const string& str) {
    map<string, bool> results;
    
    // Initialize pattern flags
    bool abc_pattern = false;       // Original ABC pattern (A>=B<=C)
    bool palindrome_pattern = false; // Any palindrome sequence of length >=3
    bool increasing_pattern = false; // Strictly increasing sequence (ABC where A<B<C)
    bool decreasing_pattern = false; // Strictly decreasing sequence (CBA where C<B<A)
    bool zigzag_pattern = false;    // Zigzag pattern (A<B>C<D>E etc.)
    
    // Variables for pattern tracking
    int n = str.length();
    vector<int> counts(26, 0);      // Count of each character
    char prev_char = '\0';
    int current_streak = 1;
    
    // Check for all patterns in a single pass
    for (int i = 0; i < n; ++i) {
        char c = str[i];
        counts[c - 'a']++;
        
        // Check for ABC pattern (original functionality)
        if (i >= 2) {
            if (str[i-2] + 1 == str[i-1] && str[i-1] + 1 == str[i]) {
                int a = counts[str[i-2] - 'a'];
                int b = counts[str[i-1] - 'a'];
                int c = counts[str[i] - 'a'];
                if (a >= b && c >= b) {
                    abc_pattern = true;
                }
            }
        }
        
        // Check for palindrome pattern
        if (i >= 2 && str[i] == str[i-2]) {
            palindrome_pattern = true;
        }
        
        // Check for increasing pattern
        if (i >= 2 && str[i-2] + 1 == str[i-1] && str[i-1] + 1 == str[i]) {
            increasing_pattern = true;
        }
        
        // Check for decreasing pattern
        if (i >= 2 && str[i-2] - 1 == str[i-1] && str[i-1] - 1 == str[i]) {
            decreasing_pattern = true;
        }
        
        // Check for zigzag pattern (needs at least 3 characters)
        if (i >= 2) {
            if ((str[i-2] < str[i-1] && str[i-1] > str[i]) || 
                (str[i-2] > str[i-1] && str[i-1] < str[i])) {
                zigzag_pattern = true;
            }
        }
    }
    
    // Store results
    results["abc_pattern"] = abc_pattern;
    results["palindrome"] = palindrome_pattern;
    results["increasing"] = increasing_pattern;
    results["decreasing"] = decreasing_pattern;
    results["zigzag"] = zigzag_pattern;
    
    return results;
}
