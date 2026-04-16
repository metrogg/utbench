#include <vector>
#include <algorithm>
#include <cassert>
#include <string>
#include <map>
#include <numeric>

struct BWTResult {
    std::string transformed;
    int primary_index;
    std::vector<int> suffix_array;
};

// Structure for BWT processing
typedef struct {
    uint64_t u, v;
} pair64_t;

// BWT Encoding function
BWTResult bwt_encode(const std::string& input) {
    if (input.empty()) {
        return {"", -1, {}};
    }

    // Add sentinel character if not present
    std::string s = input;
    if (s.back() != '$') {
        s += '$';
    }

    // Generate all rotations
    std::vector<std::string> rotations;
    for (size_t i = 0; i < s.size(); ++i) {
        rotations.push_back(s.substr(i) + s.substr(0, i));
    }

    // Sort rotations lexicographically
    std::sort(rotations.begin(), rotations.end());

    // Extract last column and find primary index
    std::string transformed;
    int primary_index = -1;
    std::vector<int> suffix_array;
    
    for (size_t i = 0; i < rotations.size(); ++i) {
        transformed += rotations[i].back();
        if (rotations[i] == s) {
            primary_index = i;
        }
        // Store the original string index for suffix array
        suffix_array.push_back(rotations.size() - rotations[i].find('$') - 1);
    }

    return {transformed, primary_index, suffix_array};
}

// BWT Decoding function
std::string bwt_decode(const std::string& transformed, int primary_index) {
    if (transformed.empty() || primary_index < 0 || primary_index >= transformed.size()) {
        return "";
    }

    // Create the table for inversion
    std::vector<std::string> table(transformed.size(), "");
    for (size_t i = 0; i < transformed.size(); ++i) {
        // Prepend each character to the corresponding string
        for (size_t j = 0; j < transformed.size(); ++j) {
            table[j] = transformed[j] + table[j];
        }
        // Sort the table
        std::sort(table.begin(), table.end());
    }

    // The original string is at the primary index
    return table[primary_index];
}

// Optimized BWT Encoding using suffix array (more efficient)
BWTResult bwt_encode_optimized(const std::string& input) {
    if (input.empty()) {
        return {"", -1, {}};
    }

    std::string s = input;
    if (s.back() != '$') {
        s += '$';
    }

    // Create suffix array
    std::vector<int> suffix_array(s.size());
    std::iota(suffix_array.begin(), suffix_array.end(), 0);
    
    // Sort suffixes based on their string values
    std::sort(suffix_array.begin(), suffix_array.end(),
        [&s](int a, int b) -> bool {
            return s.compare(a, std::string::npos, s, b, std::string::npos) < 0;
        });

    // Build BWT string
    std::string transformed;
    int primary_index = -1;
    for (size_t i = 0; i < suffix_array.size(); ++i) {
        int pos = suffix_array[i];
        if (pos == 0) {
            primary_index = i;
            transformed += s.back();
        } else {
            transformed += s[pos - 1];
        }
    }

    return {transformed, primary_index, suffix_array};
}
