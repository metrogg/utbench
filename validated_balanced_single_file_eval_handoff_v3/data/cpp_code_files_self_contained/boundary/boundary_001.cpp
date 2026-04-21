#include <vector>
#include <string>
#include <algorithm>
#include <cmath>
#include <stdexcept>

using namespace std;

// Converts a decimal number to binary and returns as string
string decimalToBinary(int decimal) {
    if (decimal == 0) return "0";
    
    bool isNegative = false;
    if (decimal < 0) {
        isNegative = true;
        decimal = abs(decimal);
    }
    
    string binary;
    while (decimal > 0) {
        binary = to_string(decimal % 2) + binary;
        decimal /= 2;
    }
    
    if (isNegative) {
        binary = "-" + binary;
    }
    
    return binary;
}

// Converts binary string to decimal number
int binaryToDecimal(const string& binary) {
    if (binary.empty()) {
        throw invalid_argument("Empty binary string");
    }
    
    bool isNegative = false;
    string processed = binary;
    
    if (binary[0] == '-') {
        isNegative = true;
        processed = binary.substr(1);
    }
    
    int decimal = 0;
    int power = 0;
    
    for (int i = processed.length() - 1; i >= 0; --i) {
        char c = processed[i];
        if (c != '0' && c != '1') {
            throw invalid_argument("Invalid binary digit: " + string(1, c));
        }
        
        int digit = c - '0';
        decimal += digit * pow(2, power);
        power++;
    }
    
    return isNegative ? -decimal : decimal;
}

// Validates if a string is a valid binary number
bool isValidBinary(const string& binary) {
    if (binary.empty()) return false;
    
    size_t start = 0;
    if (binary[0] == '-') {
        if (binary.length() == 1) return false;
        start = 1;
    }
    
    for (size_t i = start; i < binary.length(); ++i) {
        if (binary[i] != '0' && binary[i] != '1') {
            return false;
        }
    }
    
    return true;
}
