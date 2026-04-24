#include <unordered_map>
#include <vector>
#include <utility>
#include <string>
#include <cctype>
#include <algorithm>

using namespace std;

// Enhanced password modification rules
const unordered_map<char, string> CHANGE_MAP = {
    {'1', "@"}, {'0', "%"}, {'l', "L"}, {'O', "o"},
    {'a', "4"}, {'e', "3"}, {'i', "1"}, {'o', "0"}, {'s', "$"}
};

// Password strength levels
enum class Strength { WEAK, MEDIUM, STRONG, VERY_STRONG };

// Checks if password needs modification and modifies it
pair<bool, string> check_and_modify_password(string password) {
    bool modified = false;
    
    for (size_t i = 0; i < password.length(); ) {
        if (CHANGE_MAP.count(password[i])) {
            string replacement = CHANGE_MAP.at(password[i]);
            password.replace(i, 1, replacement);
            modified = true;
            i += replacement.length(); // Skip replaced characters
        } else {
            i++;
        }
    }
    
    return {modified, password};
}

// Analyzes password strength
Strength analyze_password_strength(const string& password) {
    bool has_upper = false, has_lower = false;
    bool has_digit = false, has_special = false;
    int categories = 0;
    
    for (char c : password) {
        if (isupper(c)) has_upper = true;
        else if (islower(c)) has_lower = true;
        else if (isdigit(c)) has_digit = true;
        else has_special = true;
    }
    
    categories = has_upper + has_lower + has_digit + has_special;
    
    if (password.length() < 6) return Strength::WEAK;
    if (password.length() < 8 && categories < 3) return Strength::WEAK;
    if (password.length() < 10 && categories < 3) return Strength::MEDIUM;
    if (password.length() >= 12 && categories == 4) return Strength::VERY_STRONG;
    if (password.length() >= 8 && categories >= 3) return Strength::STRONG;
    
    return Strength::MEDIUM;
}

// Gets string representation of strength level
string strength_to_string(Strength s) {
    switch(s) {
        case Strength::WEAK: return "Weak";
        case Strength::MEDIUM: return "Medium";
        case Strength::STRONG: return "Strong";
        case Strength::VERY_STRONG: return "Very Strong";
        default: return "Unknown";
    }
}
