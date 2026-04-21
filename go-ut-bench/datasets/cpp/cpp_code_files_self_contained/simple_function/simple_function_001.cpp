#include <stdexcept>
#include <string>
#include <algorithm>

std::string reverseString(const std::string& input) {
    if (input.empty()) {
        return "";
    }
    std::string result = input;
    std::reverse(result.begin(), result.end());
    return result;
}

bool isPalindrome(const std::string& input) {
    if (input.empty()) {
        return false;
    }
    std::string reversed = reverseString(input);
    return input == reversed;
}