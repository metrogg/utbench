#include <string>
#include <vector>
#include <map>
#include <algorithm>
#include <ctime>
#include <iomanip>

using namespace std;

// Enhanced password manager with multiple security features
map<string, string> password_manager(const vector<string>& passwords, 
                                    const string& input_password, 
                                    bool enable_history_check = true,
                                    bool enable_common_check = true,
                                    int max_attempts = 3) {
    map<string, string> result;
    
    // Common weak passwords to check against
    const vector<string> common_passwords = {
        "password", "123456", "qwerty", "admin", "welcome"
    };
    
    // Check if password is empty
    if (input_password.empty()) {
        result["status"] = "error";
        result["message"] = "Password cannot be empty";
        return result;
    }
    
    // Check password length
    if (input_password.length() < 8) {
        result["status"] = "error";
        result["message"] = "Password must be at least 8 characters";
        return result;
    }
    
    // Check against common passwords if enabled
    if (enable_common_check) {
        if (find(common_passwords.begin(), common_passwords.end(), input_password) != common_passwords.end()) {
            result["status"] = "error";
            result["message"] = "Password is too common";
            return result;
        }
    }
    
    // Check password history if enabled
    if (enable_history_check) {
        if (find(passwords.begin(), passwords.end(), input_password) != passwords.end()) {
            result["status"] = "error";
            result["message"] = "Password was used before";
            return result;
        }
    }
    
    // Check password complexity
    bool has_upper = false, has_lower = false, has_digit = false, has_special = false;
    for (char c : input_password) {
        if (isupper(c)) has_upper = true;
        else if (islower(c)) has_lower = true;
        else if (isdigit(c)) has_digit = true;
        else has_special = true;
    }
    
    int complexity_score = (has_upper ? 1 : 0) + (has_lower ? 1 : 0) + 
                          (has_digit ? 1 : 0) + (has_special ? 1 : 0);
    
    if (complexity_score < 3) {
        result["status"] = "warning";
        result["message"] = "Password is weak (consider using mixed characters)";
        result["strength"] = "weak";
    } else {
        result["status"] = "success";
        result["message"] = "Password is strong";
        result["strength"] = "strong";
    }
    
    result["length"] = to_string(input_password.length());
    result["complexity"] = to_string(complexity_score);
    
    return result;
}
