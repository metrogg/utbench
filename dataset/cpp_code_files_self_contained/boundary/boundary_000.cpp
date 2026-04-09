#include <string>
#include <vector>
#include <algorithm>
#include <stdexcept>

using namespace std;

class BigInt {
private:
    string number;
    bool isNegative;

    // Helper function to remove leading zeros
    void normalize() {
        size_t nonZeroIndex = number.find_first_not_of('0');
        if (nonZeroIndex == string::npos) {
            number = "0";
            isNegative = false;
        } else {
            number = number.substr(nonZeroIndex);
        }
        if (number == "0") isNegative = false;
    }

    // Helper function for string validation
    bool isValidNumber(const string& num) const {
        if (num.empty()) return false;
        size_t start = 0;
        if (num[0] == '-') {
            if (num.length() == 1) return false;
            start = 1;
        }
        return all_of(num.begin() + start, num.end(), [](char c) {
            return isdigit(c);
        });
    }

    // Helper function for addition of absolute values
    string addAbsolute(const string& a, const string& b) const {
        string result;
        int i = a.length() - 1;
        int j = b.length() - 1;
        int carry = 0;

        while (i >= 0 || j >= 0 || carry) {
            int digitA = (i >= 0) ? (a[i--] - '0') : 0;
            int digitB = (j >= 0) ? (b[j--] - '0') : 0;
            int sum = digitA + digitB + carry;
            carry = sum / 10;
            result.push_back((sum % 10) + '0');
        }

        reverse(result.begin(), result.end());
        return result;
    }

    // Helper function for subtraction of absolute values (a >= b)
    string subtractAbsolute(string a, string b) const {
        string result;
        int i = a.length() - 1;
        int j = b.length() - 1;
        int borrow = 0;

        while (i >= 0) {
            int digitA = (a[i--] - '0') - borrow;
            int digitB = (j >= 0) ? (b[j--] - '0') : 0;
            borrow = 0;

            if (digitA < digitB) {
                digitA += 10;
                borrow = 1;
            }

            result.push_back((digitA - digitB) + '0');
        }

        reverse(result.begin(), result.end());
        return result;
    }

    // Comparison of absolute values
    int compareAbsolute(const string& a, const string& b) const {
        if (a.length() != b.length()) {
            return (a.length() > b.length()) ? 1 : -1;
        }
        return a.compare(b);
    }

public:
    // Constructors
    BigInt() : number("0"), isNegative(false) {}
    BigInt(const string& num) {
        if (!isValidNumber(num)) {
            throw invalid_argument("Invalid number string");
        }
        
        if (num[0] == '-') {
            isNegative = true;
            number = num.substr(1);
        } else {
            isNegative = false;
            number = num;
        }
        normalize();
    }
    BigInt(long long num) : BigInt(to_string(num)) {}

    // Arithmetic operations
    BigInt operator+(const BigInt& other) const {
        if (isNegative == other.isNegative) {
            string sum = addAbsolute(number, other.number);
            return BigInt(isNegative ? "-" + sum : sum);
        } else {
            int cmp = compareAbsolute(number, other.number);
            if (cmp == 0) return BigInt(0);
            
            if (cmp > 0) {
                string diff = subtractAbsolute(number, other.number);
                return BigInt(isNegative ? "-" + diff : diff);
            } else {
                string diff = subtractAbsolute(other.number, number);
                return BigInt(other.isNegative ? "-" + diff : diff);
            }
        }
    }

    BigInt operator-(const BigInt& other) const {
        if (isNegative != other.isNegative) {
            string sum = addAbsolute(number, other.number);
            return BigInt(isNegative ? "-" + sum : sum);
        } else {
            int cmp = compareAbsolute(number, other.number);
            if (cmp == 0) return BigInt(0);
            
            if (cmp > 0) {
                string diff = subtractAbsolute(number, other.number);
                return BigInt(isNegative ? "-" + diff : diff);
            } else {
                string diff = subtractAbsolute(other.number, number);
                return BigInt(!isNegative ? "-" + diff : diff);
            }
        }
    }

    // Comparison operators
    bool operator==(const BigInt& other) const {
        return (number == other.number) && (isNegative == other.isNegative);
    }

    bool operator<(const BigInt& other) const {
        if (isNegative != other.isNegative) {
            return isNegative;
        }
        if (isNegative) {
            return compareAbsolute(number, other.number) > 0;
        }
        return compareAbsolute(number, other.number) < 0;
    }

    bool operator>(const BigInt& other) const {
        return other < *this;
    }

    bool operator<=(const BigInt& other) const {
        return !(*this > other);
    }

    bool operator>=(const BigInt& other) const {
        return !(*this < other);
    }

    // String representation
    string toString() const {
        return isNegative ? "-" + number : number;
    }

    // Get number of digits
    size_t digitCount() const {
        return number.length();
    }

    // Multiplication (added for more complexity)
    BigInt operator*(const BigInt& other) const {
        if (number == "0" || other.number == "0") {
            return BigInt(0);
        }

        string result(number.length() + other.number.length(), '0');
        for (int i = number.length() - 1; i >= 0; i--) {
            int carry = 0;
            for (int j = other.number.length() - 1; j >= 0; j--) {
                int product = (number[i] - '0') * (other.number[j] - '0') + 
                             (result[i + j + 1] - '0') + carry;
                carry = product / 10;
                result[i + j + 1] = (product % 10) + '0';
            }
            result[i] += carry;
        }

        BigInt res(result);
        res.isNegative = isNegative != other.isNegative;
        res.normalize();
        return res;
    }

    // Power operation (added for more complexity)
    BigInt pow(unsigned int exponent) const {
        if (exponent == 0) return BigInt(1);
        BigInt result(1);
        BigInt base = *this;
        while (exponent > 0) {
            if (exponent % 2 == 1) {
                result = result * base;
            }
            base = base * base;
            exponent /= 2;
        }
        return result;
    }
};
