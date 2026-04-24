#include <string>
#include <vector>
#include <cctype>
#include <algorithm>
#include <map>

using namespace std;

// Enhanced string processing function that performs multiple operations
map<string, string> process_string(const string& input) {
    map<string, string> results;
    
    // 1. Original string
    results["original"] = input;
    
    // 2. Reversed string
    string reversed = input;
    reverse(reversed.begin(), reversed.end());
    results["reversed"] = reversed;
    
    // 3. Uppercase version
    string upper = input;
    transform(upper.begin(), upper.end(), upper.begin(), ::toupper);
    results["uppercase"] = upper;
    
    // 4. Lowercase version
    string lower = input;
    transform(lower.begin(), lower.end(), lower.begin(), ::tolower);
    results["lowercase"] = lower;
    
    // 5. Alternating case version
    string alternate = input;
    for (size_t i = 0; i < alternate.size(); ++i) {
        if (i % 2 == 0) {
            alternate[i] = toupper(alternate[i]);
        } else {
            alternate[i] = tolower(alternate[i]);
        }
    }
    results["alternating_case"] = alternate;
    
    // 6. Word count (assuming words are separated by spaces)
    size_t word_count = 0;
    bool in_word = false;
    for (char c : input) {
        if (isalpha(c) && !in_word) {
            in_word = true;
            word_count++;
        } else if (isspace(c)) {
            in_word = false;
        }
    }
    results["word_count"] = to_string(word_count);
    
    // 7. Character count (excluding spaces)
    size_t char_count = count_if(input.begin(), input.end(), 
                                [](char c) { return !isspace(c); });
    results["non_space_chars"] = to_string(char_count);
    
    return results;
}
