#include <vector>
#include <string>
#include <algorithm>
#include <cctype>

using namespace std;

// Trim whitespace from the left side of a string (in-place)
string& ltrim_inplace(string& s) {
    s.erase(s.begin(), find_if(s.begin(), s.end(), [](int ch) {
        return !isspace(ch);
    }));
    return s;
}

// Trim whitespace from the right side of a string (in-place)
string& rtrim_inplace(string& s) {
    s.erase(find_if(s.rbegin(), s.rend(), [](int ch) {
        return !isspace(ch);
    }).base(), s.end());
    return s;
}

// Trim whitespace from both sides of a string (in-place)
string& trim_inplace(string& s) {
    ltrim_inplace(s);
    rtrim_inplace(s);
    return s;
}

// Trim whitespace from both sides of a string (returns new string)
string trim(const string& s) {
    string result = s;
    trim_inplace(result);
    return result;
}

// Check if a string starts with a given prefix
bool startswith(const string& str, const string& prefix) {
    if (prefix.length() > str.length()) return false;
    return equal(prefix.begin(), prefix.end(), str.begin());
}

// Check if a string ends with a given suffix
bool endswith(const string& str, const string& suffix) {
    if (suffix.length() > str.length()) return false;
    return equal(suffix.rbegin(), suffix.rend(), str.rbegin());
}

// Replace all occurrences of a substring in a string (in-place)
void replace_all_inplace(string& str, const string& from, const string& to) {
    if (from.empty()) return;
    
    size_t start_pos = 0;
    while ((start_pos = str.find(from, start_pos)) != string::npos) {
        str.replace(start_pos, from.length(), to);
        start_pos += to.length();
    }
}

// Split a string into tokens based on whitespace
void split(const string& str, vector<string>& tokens) {
    tokens.clear();
    size_t start = 0;
    size_t end = str.find_first_of(" \t\n\r");
    
    while (end != string::npos) {
        if (end != start) { // Only add non-empty tokens
            tokens.push_back(str.substr(start, end - start));
        }
        start = end + 1;
        end = str.find_first_of(" \t\n\r", start);
    }
    
    // Add the last token
    if (start < str.length()) {
        tokens.push_back(str.substr(start));
    }
}

// Join a container of strings with optional delimiter
string join(const vector<string>& elements, const string& delimiter = "") {
    if (elements.empty()) return "";
    
    string result;
    for (size_t i = 0; i < elements.size(); ++i) {
        if (i != 0) {
            result += delimiter;
        }
        result += elements[i];
    }
    return result;
}

// Join strings using iterators with optional delimiter
string join(vector<string>::const_iterator first, 
            vector<string>::const_iterator last, 
            const string& delimiter = "") {
    if (first == last) return "";
    
    string result;
    for (auto it = first; it != last; ++it) {
        if (it != first) {
            result += delimiter;
        }
        result += *it;
    }
    return result;
}
