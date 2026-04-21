#include <string>
#include <vector>
#include <algorithm>

using namespace std;

// Enhanced string matching function with multiple search capabilities
vector<int> find_string_matches(const string& text, const string& pattern, bool case_sensitive = true, bool all_matches = false) {
    vector<int> positions;
    
    if (pattern.empty()) {
        return positions;
    }
    
    string text_copy = text;
    string pattern_copy = pattern;
    
    if (!case_sensitive) {
        transform(text_copy.begin(), text_copy.end(), text_copy.begin(), ::tolower);
        transform(pattern_copy.begin(), pattern_copy.end(), pattern_copy.begin(), ::tolower);
    }
    
    size_t pos = 0;
    while ((pos = text_copy.find(pattern_copy, pos)) != string::npos) {
        positions.push_back(static_cast<int>(pos));
        pos += pattern_copy.length();
        if (!all_matches) {
            break;
        }
    }
    
    return positions;
}

// Function to calculate the similarity percentage between two strings
double string_similarity(const string& str1, const string& str2) {
    if (str1.empty() && str2.empty()) return 100.0;
    if (str1.empty() || str2.empty()) return 0.0;
    
    int matches = 0;
    int min_len = min(str1.length(), str2.length());
    
    for (int i = 0; i < min_len; ++i) {
        if (str1[i] == str2[i]) {
            matches++;
        }
    }
    
    return (static_cast<double>(matches) * 100) / max(str1.length(), str2.length());
}
