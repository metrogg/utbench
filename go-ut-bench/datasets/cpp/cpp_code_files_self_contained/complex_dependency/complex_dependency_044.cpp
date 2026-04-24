#include <string>
#include <vector>
#include <algorithm>
#include <map>
#include <cctype>

using namespace std;

// Enhanced string manipulation functions
map<string, string> analyze_string(const string& input) {
    map<string, string> result;
    
    // Basic string properties
    result["original"] = input;
    result["length"] = to_string(input.length());
    result["empty"] = input.empty() ? "true" : "false";
    
    // Case manipulation
    string upper = input;
    transform(upper.begin(), upper.end(), upper.begin(), ::toupper);
    result["uppercase"] = upper;
    
    string lower = input;
    transform(lower.begin(), lower.end(), lower.begin(), ::tolower);
    result["lowercase"] = lower;
    
    // Character statistics
    int letters = 0, digits = 0, spaces = 0, others = 0;
    for (char c : input) {
        if (isalpha(c)) letters++;
        else if (isdigit(c)) digits++;
        else if (isspace(c)) spaces++;
        else others++;
    }
    result["letter_count"] = to_string(letters);
    result["digit_count"] = to_string(digits);
    result["space_count"] = to_string(spaces);
    result["other_count"] = to_string(others);
    
    // Palindrome check
    string reversed = input;
    reverse(reversed.begin(), reversed.end());
    result["palindrome"] = (input == reversed) ? "true" : "false";
    
    // Word count (simple implementation)
    int word_count = 0;
    bool in_word = false;
    for (char c : input) {
        if (isalpha(c)) {
            if (!in_word) {
                word_count++;
                in_word = true;
            }
        } else {
            in_word = false;
        }
    }
    result["word_count"] = to_string(word_count);
    
    return result;
}
