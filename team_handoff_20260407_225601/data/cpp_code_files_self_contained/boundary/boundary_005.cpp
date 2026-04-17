#include <vector>
#include <string>
#include <cmath>
#include <stdexcept>
#include <map>

using namespace std;

// Enhanced binary number converter with multiple representations
map<string, string> binary_converter(const string& binary_str) {
    // Validate input
    if (binary_str.empty()) {
        throw invalid_argument("Error: Empty binary string");
    }
    if (binary_str.length() > 32) {
        throw invalid_argument("Error: Binary string too long (max 32 bits)");
    }
    for (char c : binary_str) {
        if (c != '0' && c != '1') {
            throw invalid_argument("Error: Invalid binary digit '" + string(1, c) + "'");
        }
    }

    map<string, string> results;
    int n = binary_str.length();
    
    // Unsigned representation
    unsigned long unsigned_val = 0;
    for (int i = 0; i < n; i++) {
        unsigned_val = (unsigned_val << 1) | (binary_str[i] - '0');
    }
    results["unsigned"] = to_string(unsigned_val);

    // Signed magnitude representation
    if (binary_str[0] == '0') {
        results["signed_magnitude"] = "+" + to_string(unsigned_val);
    } else {
        unsigned long magnitude = 0;
        for (int i = 1; i < n; i++) {
            magnitude = (magnitude << 1) | (binary_str[i] - '0');
        }
        results["signed_magnitude"] = "-" + to_string(magnitude);
    }

    // Ones' complement representation
    if (binary_str[0] == '0') {
        results["ones_complement"] = "+" + to_string(unsigned_val);
    } else {
        unsigned long ones_comp = 0;
        string flipped = binary_str;
        for (int i = 0; i < n; i++) {
            flipped[i] = (binary_str[i] == '0') ? '1' : '0';
        }
        for (int i = 1; i < n; i++) {
            ones_comp = (ones_comp << 1) | (flipped[i] - '0');
        }
        results["ones_complement"] = "-" + to_string(ones_comp);
    }

    // Two's complement representation
    if (binary_str[0] == '0') {
        results["twos_complement"] = "+" + to_string(unsigned_val);
    } else {
        unsigned long twos_comp = 0;
        string flipped = binary_str;
        bool carry = true;
        // Flip bits and add 1
        for (int i = n-1; i >= 0; i--) {
            if (carry) {
                if (flipped[i] == '1') {
                    flipped[i] = '0';
                } else {
                    flipped[i] = '1';
                    carry = false;
                }
            } else if (i != 0) { // Don't flip the sign bit for ones' complement
                flipped[i] = (binary_str[i] == '0') ? '1' : '0';
            }
        }
        for (int i = 1; i < n; i++) {
            twos_comp = (twos_comp << 1) | (flipped[i] - '0');
        }
        results["twos_complement"] = "-" + to_string(twos_comp + (carry ? 1 : 0));
    }

    // Hexadecimal representation
    string hex_str;
    unsigned_val = 0;
    for (char c : binary_str) {
        unsigned_val = (unsigned_val << 1) | (c - '0');
    }
    char hex_digits[] = "0123456789ABCDEF";
    while (unsigned_val > 0) {
        hex_str = hex_digits[unsigned_val % 16] + hex_str;
        unsigned_val /= 16;
    }
    results["hexadecimal"] = "0x" + (hex_str.empty() ? "0" : hex_str);

    return results;
}
