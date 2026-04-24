#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <cctype>

using namespace std;

// Function to analyze a string and return various statistics
map<string, int> analyze_string(const string& input) {
    map<string, int> stats;
    
    // Basic length statistics
    stats["length"] = input.length();
    stats["words"] = 0;
    stats["uppercase"] = 0;
    stats["lowercase"] = 0;
    stats["digits"] = 0;
    stats["spaces"] = 0;
    stats["vowels"] = 0;
    stats["consonants"] = 0;
    
    if (input.empty()) {
        return stats;
    }
    
    bool in_word = false;
    const string vowels = "aeiouAEIOU";
    
    for (char c : input) {
        if (isupper(c)) {
            stats["uppercase"]++;
        } else if (islower(c)) {
            stats["lowercase"]++;
        }
        
        if (isdigit(c)) {
            stats["digits"]++;
        } else if (isspace(c)) {
            stats["spaces"]++;
            in_word = false;
        } else if (!in_word) {
            stats["words"]++;
            in_word = true;
        }
        
        if (isalpha(c)) {
            if (vowels.find(c) != string::npos) {
                stats["vowels"]++;
            } else {
                stats["consonants"]++;
            }
        }
    }
    
    return stats;
}

// Function to process and transform a string
string process_string(const string& input) {
    if (input.empty()) {
        return input;
    }
    
    string result;
    bool capitalize_next = true;
    
    for (char c : input) {
        if (capitalize_next && isalpha(c)) {
            result += toupper(c);
            capitalize_next = false;
        } else {
            result += tolower(c);
        }
        
        if (c == '.' || c == '!' || c == '?') {
            capitalize_next = true;
        }
    }
    
    return result;
}
