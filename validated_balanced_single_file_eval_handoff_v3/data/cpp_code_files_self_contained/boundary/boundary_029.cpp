#include <vector>
#include <string>
#include <map>
#include <algorithm>
#include <cctype>

using namespace std;

// Enhanced string analyzer with multiple string operations
map<string, int> analyze_string(const string& str) {
    map<string, int> results;
    
    // Basic length calculation
    results["length"] = str.length();
    
    // Count uppercase letters
    results["uppercase"] = count_if(str.begin(), str.end(), [](char c) {
        return isupper(c);
    });
    
    // Count lowercase letters
    results["lowercase"] = count_if(str.begin(), str.end(), [](char c) {
        return islower(c);
    });
    
    // Count digits
    results["digits"] = count_if(str.begin(), str.end(), [](char c) {
        return isdigit(c);
    });
    
    // Count whitespace characters
    results["whitespace"] = count_if(str.begin(), str.end(), [](char c) {
        return isspace(c);
    });
    
    // Count vowels
    results["vowels"] = count_if(str.begin(), str.end(), [](char c) {
        c = tolower(c);
        return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u';
    });
    
    // Count consonants
    results["consonants"] = count_if(str.begin(), str.end(), [](char c) {
        c = tolower(c);
        return isalpha(c) && !(c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u');
    });
    
    // Count special characters
    results["special"] = count_if(str.begin(), str.end(), [](char c) {
        return !isalnum(c) && !isspace(c);
    });
    
    return results;
}

// Alternative string length counter (similar to original but improved)
int custom_strlen(const char* str) {
    if (str == nullptr) return 0;
    int length = 0;
    while (*str != '\0') {
        length++;
        str++;
    }
    return length;
}
