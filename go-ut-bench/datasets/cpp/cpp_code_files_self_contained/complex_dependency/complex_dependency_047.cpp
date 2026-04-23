#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <cctype>

using namespace std;

// Enhanced string analysis function that provides multiple string metrics
map<string, int> analyze_string(const string& str) {
    map<string, int> metrics;
    
    // Basic string properties
    metrics["length"] = str.length();
    metrics["size"] = str.size();
    
    // Character type counts
    metrics["uppercase"] = 0;
    metrics["lowercase"] = 0;
    metrics["digits"] = 0;
    metrics["whitespace"] = 0;
    metrics["special"] = 0;
    
    // Word count (simplified)
    metrics["words"] = 0;
    bool in_word = false;
    
    for (char c : str) {
        if (isupper(c)) {
            metrics["uppercase"]++;
        } else if (islower(c)) {
            metrics["lowercase"]++;
        } else if (isdigit(c)) {
            metrics["digits"]++;
        } else if (isspace(c)) {
            metrics["whitespace"]++;
            in_word = false;
            continue;
        } else {
            metrics["special"]++;
        }
        
        // Simple word counting
        if (!isspace(c) && !in_word) {
            metrics["words"]++;
            in_word = true;
        }
    }
    
    // String palindrome check
    string lower_str = str;
    transform(lower_str.begin(), lower_str.end(), lower_str.begin(), ::tolower);
    string reversed_str = lower_str;
    reverse(reversed_str.begin(), reversed_str.end());
    metrics["is_palindrome"] = (lower_str == reversed_str) ? 1 : 0;
    
    return metrics;
}
