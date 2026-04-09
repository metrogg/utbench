#include <vector>
#include <map>
#include <algorithm>
#include <stdexcept>
#include <cmath>
#include <string>

using namespace std;

class Rational {
private:
    long long numerator;
    long long denominator;

    long long gcd(long long a, long long b) const {
        return b == 0 ? a : gcd(b, a % b);
    }

    void normalize() {
        if (denominator < 0) {
            numerator = -numerator;
            denominator = -denominator;
        }
        
        if (numerator == 0) {
            denominator = 1;
        } else {
            long long common_divisor = gcd(abs(numerator), abs(denominator));
            numerator /= common_divisor;
            denominator /= common_divisor;
        }
    }

public:
    Rational(long long num = 0, long long denom = 1) : numerator(num), denominator(denom) {
        if (denominator == 0) {
            throw invalid_argument("Denominator cannot be zero");
        }
        normalize();
    }

    // Arithmetic operations
    Rational operator+(const Rational& other) const {
        return Rational(
            numerator * other.denominator + other.numerator * denominator,
            denominator * other.denominator
        );
    }

    Rational operator-(const Rational& other) const {
        return Rational(
            numerator * other.denominator - other.numerator * denominator,
            denominator * other.denominator
        );
    }

    Rational operator*(const Rational& other) const {
        return Rational(
            numerator * other.numerator,
            denominator * other.denominator
        );
    }

    Rational operator/(const Rational& other) const {
        if (other.numerator == 0) {
            throw domain_error("Division by zero");
        }
        return Rational(
            numerator * other.denominator,
            denominator * other.numerator
        );
    }

    // Comparison operators
    bool operator==(const Rational& other) const {
        return numerator == other.numerator && denominator == other.denominator;
    }

    bool operator<(const Rational& other) const {
        return numerator * other.denominator < other.numerator * denominator;
    }

    bool operator>(const Rational& other) const {
        return numerator * other.denominator > other.numerator * denominator;
    }

    // Conversion functions
    double to_double() const {
        return static_cast<double>(numerator) / denominator;
    }

    // String representation
    string to_string(bool show_parens = false) const {
        string result;
        if (show_parens && numerator < 0) {
            result += "(";
        }

        long long integer_part = numerator / denominator;
        long long remainder = abs(numerator) % denominator;

        if (denominator == 1) {
            result += std::to_string(numerator);
        } else if (abs(numerator) < denominator) {
            result += std::to_string(numerator) + "/" + std::to_string(denominator);
        } else {
            result += std::to_string(integer_part) + " " + 
                     std::to_string(remainder) + "/" + 
                     std::to_string(denominator);
        }

        if (show_parens && numerator < 0) {
            result += ")";
        }

        return result;
    }

    // Getters
    long long get_numerator() const { return numerator; }
    long long get_denominator() const { return denominator; }
};

vector<Rational> parse_rationals(const string& input) {
    vector<Rational> result;
    size_t start = 0;
    size_t end = input.find(' ');
    
    while (end != string::npos) {
        string token = input.substr(start, end - start);
        size_t slash_pos = token.find('/');
        
        if (slash_pos != string::npos) {
            long long num = stoll(token.substr(0, slash_pos));
            long long den = stoll(token.substr(slash_pos + 1));
            result.emplace_back(num, den);
        } else {
            result.emplace_back(stoll(token), 1);
        }
        
        start = end + 1;
        end = input.find(' ', start);
    }
    
    // Add the last token
    string last_token = input.substr(start);
    size_t slash_pos = last_token.find('/');
    if (slash_pos != string::npos) {
        long long num = stoll(last_token.substr(0, slash_pos));
        long long den = stoll(last_token.substr(slash_pos + 1));
        result.emplace_back(num, den);
    } else {
        result.emplace_back(stoll(last_token), 1);
    }
    
    return result;
}
