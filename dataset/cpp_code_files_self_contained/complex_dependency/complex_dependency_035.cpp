#include <string>
#include <vector>
#include <algorithm>
#include <map>

using namespace std;

// Enhanced string pattern matching with multiple algorithms and statistics
map<string, int> pattern_matcher(const string& text, const string& pattern, 
                                bool case_sensitive = true, bool count_overlaps = true) {
    map<string, int> results;
    
    if (pattern.empty()) {
        results["error"] = -1;
        return results;
    }
    
    string modified_text = text;
    string modified_pattern = pattern;
    
    if (!case_sensitive) {
        transform(text.begin(), text.end(), modified_text.begin(), ::tolower);
        transform(pattern.begin(), pattern.end(), modified_pattern.begin(), ::tolower);
    }
    
    // Naive matching (original algorithm)
    int naive_count = 0;
    int pattern_len = modified_pattern.length();
    int text_len = modified_text.length();
    
    for (int i = 0; i + pattern_len <= text_len; ) {
        bool found = true;
        for (int j = 0; j < pattern_len; j++) {
            if (modified_text[i + j] != modified_pattern[j]) {
                found = false;
                break;
            }
        }
        if (found) {
            naive_count++;
            i += count_overlaps ? 1 : pattern_len;
        } else {
            i++;
        }
    }
    
    // Additional statistics
    int first_occurrence = -1;
    int last_occurrence = -1;
    vector<int> positions;
    
    for (int i = 0; i + pattern_len <= text_len; i++) {
        bool found = true;
        for (int j = 0; j < pattern_len; j++) {
            if (modified_text[i + j] != modified_pattern[j]) {
                found = false;
                break;
            }
        }
        if (found) {
            if (first_occurrence == -1) first_occurrence = i;
            last_occurrence = i;
            positions.push_back(i);
        }
    }
    
    // Store all results
    results["total_matches"] = naive_count;
    results["first_position"] = first_occurrence;
    results["last_position"] = last_occurrence;
    results["pattern_length"] = pattern_len;
    results["text_length"] = text_len;
    
    return results;
}
