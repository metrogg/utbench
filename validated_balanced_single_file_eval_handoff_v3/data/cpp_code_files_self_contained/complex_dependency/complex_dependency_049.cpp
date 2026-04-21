#include <vector>
#include <string>
#include <cctype>
#include <map>

using namespace std;

// Enhanced character analysis function that provides multiple character properties
map<string, string> analyze_character(char ch) {
    map<string, string> result;
    
    // Determine character case
    if (islower(ch)) {
        result["case"] = "lowercase";
    } else if (isupper(ch)) {
        result["case"] = "UPPERCASE";
    } else {
        result["case"] = "Invalid";
    }
    
    // Determine character type
    if (isalpha(ch)) {
        result["type"] = "Alphabetic";
    } else if (isdigit(ch)) {
        result["type"] = "Numeric";
    } else if (isspace(ch)) {
        result["type"] = "Whitespace";
    } else if (ispunct(ch)) {
        result["type"] = "Punctuation";
    } else {
        result["type"] = "Other";
    }
    
    // Get ASCII value
    result["ascii"] = to_string(static_cast<int>(ch));
    
    // Determine if character is printable
    result["printable"] = isprint(ch) ? "Yes" : "No";
    
    // Determine if character is control character
    result["control"] = iscntrl(ch) ? "Yes" : "No";
    
    return result;
}

// Function to analyze multiple characters in a string
vector<map<string, string>> analyze_string(const string& str) {
    vector<map<string, string>> results;
    for (char ch : str) {
        results.push_back(analyze_character(ch));
    }
    return results;
}
