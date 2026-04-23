#include <string>
#include <vector>
#include <algorithm>
#include <map>

using namespace std;

// Function to perform advanced string manipulations and analysis
map<string, string> analyze_string(const string& input_str, int substring_start = 0, int substring_length = -1) {
    map<string, string> result;
    
    // 1. Basic string info
    result["original"] = input_str;
    result["length"] = to_string(input_str.length());
    
    // 2. Substring extraction with bounds checking
    string substring;
    if (substring_start < 0 || substring_start >= input_str.length()) {
        substring = "Invalid start position";
    } else {
        int actual_length = (substring_length == -1) ? input_str.length() - substring_start : substring_length;
        actual_length = min(actual_length, (int)input_str.length() - substring_start);
        substring = input_str.substr(substring_start, actual_length);
    }
    result["substring"] = substring;
    
    // 3. String reversal
    string reversed(input_str.rbegin(), input_str.rend());
    result["reversed"] = reversed;
    
    // 4. Character frequency analysis
    map<char, int> freq_map;
    for (char c : input_str) {
        freq_map[c]++;
    }
    string freq_str;
    for (const auto& pair : freq_map) {
        freq_str += string(1, pair.first) + ":" + to_string(pair.second) + " ";
    }
    result["char_frequency"] = freq_str;
    
    // 5. Case conversion
    string upper, lower;
    transform(input_str.begin(), input_str.end(), back_inserter(upper), ::toupper);
    transform(input_str.begin(), input_str.end(), back_inserter(lower), ::tolower);
    result["uppercase"] = upper;
    result["lowercase"] = lower;
    
    // 6. Check if palindrome
    string is_palindrome = (input_str == reversed) ? "true" : "false";
    result["is_palindrome"] = is_palindrome;
    
    // 7. First and last characters
    if (!input_str.empty()) {
        result["first_char"] = string(1, input_str.front());
        result["last_char"] = string(1, input_str.back());
    } else {
        result["first_char"] = "N/A";
        result["last_char"] = "N/A";
    }
    
    return result;
}

// Helper function to generate repeated character string
string generate_repeated_string(char c, int count) {
    if (count <= 0) return "";
    return string(count, c);
}
