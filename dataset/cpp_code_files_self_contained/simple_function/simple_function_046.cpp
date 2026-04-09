#include <vector>
#include <map>
#include <string>
#include <algorithm>

using namespace std;

// Function to analyze digits in a number and return detailed statistics
map<string, vector<int>> analyze_digits(int number) {
    map<string, vector<int>> result;
    vector<int> digits;
    string num_str = to_string(number);
    
    // Extract each digit
    for (char c : num_str) {
        digits.push_back(c - '0');
    }
    
    // Store all digits
    result["digits"] = digits;
    
    // Find unique digits
    vector<int> unique_digits = digits;
    sort(unique_digits.begin(), unique_digits.end());
    unique_digits.erase(unique(unique_digits.begin(), unique_digits.end()), unique_digits.end());
    result["unique_digits"] = unique_digits;
    
    // Find repeated digits
    vector<int> repeated_digits;
    map<int, int> digit_count;
    for (int d : digits) {
        digit_count[d]++;
    }
    for (auto& pair : digit_count) {
        if (pair.second > 1) {
            repeated_digits.push_back(pair.first);
        }
    }
    result["repeated_digits"] = repeated_digits;
    
    // Check if all digits are unique
    result["all_unique"] = {(digits.size() == unique_digits.size()) ? 1 : 0};
    
    return result;
}

// Function to check if any digits are repeated (similar to original logic)
bool has_repeated_digits(int number) {
    string num_str = to_string(number);
    for (size_t i = 0; i < num_str.size(); ++i) {
        for (size_t j = i + 1; j < num_str.size(); ++j) {
            if (num_str[i] == num_str[j]) {
                return true;
            }
        }
    }
    return false;
}
