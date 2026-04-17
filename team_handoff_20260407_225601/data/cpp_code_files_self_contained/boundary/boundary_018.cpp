#include <string>
#include <algorithm>
#include <vector>
#include <stdexcept>

using namespace std;

// Function to add two binary strings with additional features
string addBinaryNumbers(const string& a, const string& b, bool validate_input = true) {
    // Input validation
    if (validate_input) {
        auto is_binary = [](const string& s) {
            return all_of(s.begin(), s.end(), [](char c) { return c == '0' || c == '1'; });
        };
        
        if (!is_binary(a) || !is_binary(b)) {
            throw invalid_argument("Input strings must contain only '0' and '1' characters");
        }
    }

    int i = a.size() - 1, j = b.size() - 1;
    int carry = 0;
    string result;

    while (i >= 0 || j >= 0 || carry > 0) {
        int digit_a = (i >= 0) ? a[i--] - '0' : 0;
        int digit_b = (j >= 0) ? b[j--] - '0' : 0;
        
        int sum = digit_a + digit_b + carry;
        result.push_back((sum % 2) + '0');
        carry = sum / 2;
    }

    reverse(result.begin(), result.end());
    
    // Remove leading zeros (optional)
    size_t first_non_zero = result.find_first_not_of('0');
    if (first_non_zero != string::npos) {
        result = result.substr(first_non_zero);
    } else {
        result = "0"; // All zeros case
    }

    return result;
}

// Extended function to add multiple binary strings
string addMultipleBinary(const vector<string>& binaries) {
    if (binaries.empty()) return "0";
    
    string result = binaries[0];
    for (size_t i = 1; i < binaries.size(); ++i) {
        result = addBinaryNumbers(result, binaries[i]);
    }
    
    return result;
}

// Function to convert binary string to decimal
unsigned long long binaryToDecimal(const string& binary) {
    unsigned long long decimal = 0;
    for (char c : binary) {
        decimal = (decimal << 1) | (c - '0');
    }
    return decimal;
}
