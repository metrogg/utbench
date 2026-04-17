#include <algorithm>
#include <vector>
#include <locale>
#include <map>
#include <cctype>
#include <sstream>

using namespace std;

// Enhanced string processing and comparison functions

// Converts string to lowercase and removes non-alphabetic characters
string normalize_string(const string& s) {
    string result;
    for (char c : s) {
        if (isalpha(c)) {
            result += tolower(c);
        }
    }
    return result;
}

// Custom comparison function that considers:
// 1. Case-insensitive comparison
// 2. Length of normalized string
// 3. Original string as tiebreaker
bool enhanced_string_compare(const string& a, const string& b) {
    string a_norm = normalize_string(a);
    string b_norm = normalize_string(b);
    
    if (a_norm.length() != b_norm.length()) {
        return a_norm.length() < b_norm.length();
    }
    
    if (a_norm != b_norm) {
        return a_norm < b_norm;
    }
    
    return a < b; // Original strings as tiebreaker
}

// Analyzes string vector and returns statistics
map<string, int> analyze_strings(const vector<string>& strings) {
    map<string, int> stats;
    int total_chars = 0;
    map<char, int> char_counts;
    
    for (const auto& s : strings) {
        total_chars += s.length();
        for (char c : s) {
            char c_lower = tolower(c);
            char_counts[c_lower]++;
        }
    }
    
    stats["total_strings"] = strings.size();
    stats["total_chars"] = total_chars;
    stats["unique_chars"] = char_counts.size();
    stats["avg_length"] = strings.empty() ? 0 : total_chars / strings.size();
    
    return stats;
}

// Processes vector of strings: sorts and analyzes them
vector<string> process_strings(vector<string>& strings) {
    // First sort using enhanced comparison
    sort(strings.begin(), strings.end(), enhanced_string_compare);
    
    // Get statistics
    auto stats = analyze_strings(strings);
    
    // Print statistics
    cout << "String Statistics:" << endl;
    cout << "------------------" << endl;
    for (const auto& stat : stats) {
        cout << stat.first << ": " << stat.second << endl;
    }
    cout << "------------------" << endl;
    
    return strings;
}
