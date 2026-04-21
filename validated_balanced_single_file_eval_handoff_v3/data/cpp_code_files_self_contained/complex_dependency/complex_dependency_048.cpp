#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <cctype>

using namespace std;

// Function to perform advanced string comparison and analysis
map<string, string> analyze_strings(const string& str1, const string& str2) {
    map<string, string> result;
    
    // Basic comparisons
    result["string1"] = str1;
    result["string2"] = str2;
    result["length1"] = to_string(str1.length());
    result["length2"] = to_string(str2.length());
    result["lexical_compare"] = (str1 == str2) ? "equal" : (str1 > str2) ? "greater" : "less";
    
    // Character analysis
    result["last_char_diff"] = to_string(str1.empty() || str2.empty() ? 0 : str1.back() - str2.back());
    
    // Case sensitivity analysis
    bool case_diff = false;
    if (str1.length() == str2.length()) {
        for (size_t i =  0; i < str1.length(); ++i) {
            if (tolower(str1[i]) != tolower(str2[i])) {
                case_diff = true;
                break;
            }
        }
    }
    result["case_sensitive_equal"] = (!case_diff && str1.length() == str2.length()) ? "true" : "false";
    
    // Numeric content analysis
    auto count_digits = [](const string& s) {
        return count_if(s.begin(), s.end(), [](char c) { return isdigit(c); });
    };
    result["digits_in_str1"] = to_string(count_digits(str1));
    result["digits_in_str2"] = to_string(count_digits(str2));
    
    // Special characters analysis
    auto count_special = [](const string& s) {
        return count_if(s.begin(), s.end(), [](char c) { return !isalnum(c); });
    };
    result["special_in_str1"] = to_string(count_special(str1));
    result["special_in_str2"] = to_string(count_special(str2));
    
    return result;
}
