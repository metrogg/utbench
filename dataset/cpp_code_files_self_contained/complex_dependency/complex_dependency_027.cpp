#include <string>
#include <vector>
#include <cctype>
#include <algorithm>
#include <map>
#include <sstream>
#include <iomanip>

using namespace std;

// Case-insensitive string comparison (returns 0 if equal, <0 if s1<s2, >0 if s1>s2)
int strCaseCmp(const string& s1, const string& s2) {
    auto it1 = s1.begin();
    auto it2 = s2.begin();
    
    while (it1 != s1.end() && it2 != s2.end()) {
        char c1 = tolower(*it1);
        char c2 = tolower(*it2);
        if (c1 != c2) {
            return c1 - c2;
        }
        ++it1;
        ++it2;
    }
    
    if (it1 == s1.end() && it2 == s2.end()) return 0;
    if (it1 == s1.end()) return -1;
    return 1;
}

// Convert string to lowercase
string toLower(const string& s) {
    string result;
    transform(s.begin(), s.end(), back_inserter(result), 
              [](unsigned char c) { return tolower(c); });
    return result;
}

// Convert string to uppercase
string toUpper(const string& s) {
    string result;
    transform(s.begin(), s.end(), back_inserter(result), 
              [](unsigned char c) { return toupper(c); });
    return result;
}

// Check if string is numeric
bool isNumeric(const string& s) {
    if (s.empty()) return false;
    
    size_t start = 0;
    if (s[0] == '-') {
        if (s.size() == 1) return false;
        start = 1;
    }
    
    bool decimalFound = false;
    for (size_t i = start; i < s.size(); ++i) {
        if (s[i] == '.') {
            if (decimalFound) return false;
            decimalFound = true;
        } else if (!isdigit(s[i])) {
            return false;
        }
    }
    return true;
}

// Count occurrences of a character in a string
size_t countChar(const string& s, char c) {
    return count_if(s.begin(), s.end(), 
                   [c](char ch) { return tolower(ch) == tolower(c); });
}

// Enhanced string analysis function
map<string, string> analyzeString(const string& s) {
    map<string, string> result;
    
    result["original"] = s;
    result["lowercase"] = toLower(s);
    result["uppercase"] = toUpper(s);
    result["length"] = to_string(s.length());
    
    size_t letters = count_if(s.begin(), s.end(), 
                            [](char c) { return isalpha(c); });
    size_t digits = count_if(s.begin(), s.end(), 
                           [](char c) { return isdigit(c); });
    size_t spaces = count_if(s.begin(), s.end(), 
                           [](char c) { return isspace(c); });
    size_t others = s.length() - letters - digits - spaces;
    
    result["letters"] = to_string(letters);
    result["digits"] = to_string(digits);
    result["spaces"] = to_string(spaces);
    result["others"] = to_string(others);
    result["is_numeric"] = isNumeric(s) ? "true" : "false";
    
    return result;
}

// Case-insensitive string search
bool containsIgnoreCase(const string& haystack, const string& needle) {
    auto it = search(haystack.begin(), haystack.end(),
                    needle.begin(), needle.end(),
                    [](char ch1, char ch2) { 
                        return tolower(ch1) == tolower(ch2); 
                    });
    return it != haystack.end();
}
