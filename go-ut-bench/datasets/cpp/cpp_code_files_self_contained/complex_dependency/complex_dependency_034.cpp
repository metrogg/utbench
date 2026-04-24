#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <cctype>

using namespace std;

// Enhanced string analysis function that provides multiple string metrics
map<string, int> analyze_string(const string& str) {
    map<string, int> results;
    
    // Basic length metrics
    results["length"] = str.length();
    results["bytes"] = str.size();
    
    // Character type counts
    int uppercase = 0, lowercase = 0, digits = 0, whitespace = 0, special = 0;
    for (char c : str) {
        if (isupper(c)) uppercase++;
        else if (islower(c)) lowercase++;
        else if (isdigit(c)) digits++;
        else if (isspace(c)) whitespace++;
        else special++;
    }
    
    results["uppercase"] = uppercase;
    results["lowercase"] = lowercase;
    results["digits"] = digits;
    results["whitespace"] = whitespace;
    results["special"] = special;
    
    // Word count (simplified)
    int word_count = 0;
    bool in_word = false;
    for (char c : str) {
        if (isalpha(c)) {
            if (!in_word) {
                word_count++;
                in_word = true;
            }
        } else {
            in_word = false;
        }
    }
    results["words"] = word_count;
    
    // Palindrome check
    string lower_str = str;
    transform(lower_str.begin(), lower_str.end(), lower_str.begin(), ::tolower);
    string reversed_str(lower_str.rbegin(), lower_str.rend());
    results["is_palindrome"] = (lower_str == reversed_str);
    
    return results;
}
