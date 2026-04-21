#include <string>
#include <vector>
#include <algorithm>
#include <cctype>
#include <map>

using namespace std;

// Enhanced string processing function with multiple cleaning options
map<string, string> process_string(const string& input, bool trim=true, bool remove_duplicates=true, 
                                  bool to_lowercase=false, bool remove_punctuation=false) {
    string result = input;
    
    // Convert to lowercase if requested
    if (to_lowercase) {
        transform(result.begin(), result.end(), result.begin(), 
                 [](unsigned char c){ return tolower(c); });
    }
    
    // Remove punctuation if requested
    if (remove_punctuation) {
        result.erase(remove_if(result.begin(), result.end(), 
                    [](unsigned char c){ return ispunct(c); }), result.end());
    }
    
    // Remove duplicate spaces if requested
    if (remove_duplicates) {
        string::size_type pos = 0;
        while ((pos = result.find("  ", pos)) != string::npos) {
            result.replace(pos, 2, " ");
        }
    }
    
    // Trim leading and trailing spaces if requested
    if (trim) {
        // Trim leading spaces
        while (!result.empty() && result[0] == ' ') {
            result.erase(0, 1);
        }
        // Trim trailing spaces
        while (!result.empty() && result.back() == ' ') {
            result.pop_back();
        }
    }
    
    // Calculate statistics
    map<string, string> output;
    output["processed"] = result;
    output["original_length"] = to_string(input.length());
    output["processed_length"] = to_string(result.length());
    output["space_reduction"] = to_string(input.length() - result.length());
    
    // Word count (simple implementation)
    int word_count = 0;
    bool in_word = false;
    for (char c : result) {
        if (isspace(c)) {
            in_word = false;
        } else if (!in_word) {
            in_word = true;
            word_count++;
        }
    }
    output["word_count"] = to_string(word_count);
    
    return output;
}
