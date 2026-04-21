#include <vector>
#include <string>
#include <algorithm>
#include <climits>
#include <map>

using namespace std;

struct Pattern {
    vector<char> chars;
    vector<int> counts;
};

Pattern extract_pattern(const string& str) {
    Pattern pat;
    if (str.empty()) return pat;
    
    char current = str[0];
    int count = 1;
    
    for (size_t i = 1; i < str.size(); i++) {
        if (str[i] == current) {
            count++;
        } else {
            pat.chars.push_back(current);
            pat.counts.push_back(count);
            current = str[i];
            count = 1;
        }
    }
    pat.chars.push_back(current);
    pat.counts.push_back(count);
    
    return pat;
}

int calculate_min_operations(const vector<Pattern>& patterns) {
    if (patterns.empty()) return -1;
    
    // Check if all patterns have the same character sequence
    Pattern first = patterns[0];
    for (const auto& pat : patterns) {
        if (pat.chars != first.chars) {
            return -1;
        }
    }
    
    int total_operations = 0;
    size_t pattern_length = first.chars.size();
    
    for (size_t i = 0; i < pattern_length; i++) {
        vector<int> counts;
        for (const auto& pat : patterns) {
            counts.push_back(pat.counts[i]);
        }
        
        // Find median to minimize operations
        sort(counts.begin(), counts.end());
        int median = counts[counts.size() / 2];
        
        int operations = 0;
        for (int cnt : counts) {
            operations += abs(cnt - median);
        }
        
        total_operations += operations;
    }
    
    return total_operations;
}

string pattern_matching_solver(const vector<string>& strings) {
    if (strings.empty()) return "Invalid input: no strings provided";
    
    vector<Pattern> patterns;
    for (const auto& str : strings) {
        patterns.push_back(extract_pattern(str));
    }
    
    int result = calculate_min_operations(patterns);
    if (result == -1) {
        return "Fegla Won";
    } else {
        return to_string(result);
    }
}
