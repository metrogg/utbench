#include <string>
#include <vector>
#include <map>
#include <algorithm>
using namespace std;

// Quaternion multiplication rules
string quaternion_mult(const string& a, const string& b) {
    // Handle signs
    int sign = 1;
    string x = a, y = b;
    
    if (x[0] == '-') {
        sign *= -1;
        x = x.substr(1);
    }
    if (y[0] == '-') {
        sign *= -1;
        y = y.substr(1);
    }
    
    // Multiplication table
    static const map<pair<char, char>, string> mult_table = {
        {{'1', '1'}, "1"}, {{'1', 'i'}, "i"}, {{'1', 'j'}, "j"}, {{'1', 'k'}, "k"},
        {{'i', '1'}, "i"}, {{'i', 'i'}, "-1"}, {{'i', 'j'}, "k"}, {{'i', 'k'}, "-j"},
        {{'j', '1'}, "j"}, {{'j', 'i'}, "-k"}, {{'j', 'j'}, "-1"}, {{'j', 'k'}, "i"},
        {{'k', '1'}, "k"}, {{'k', 'i'}, "j"}, {{'k', 'j'}, "-i"}, {{'k', 'k'}, "-1"}
    };
    
    string result = mult_table.at({x[0], y[0]});
    if (result[0] == '-') {
        sign *= -1;
        result = result.substr(1);
    }
    
    return (sign == -1) ? "-" + result : result;
}

// Fast exponentiation for quaternions
string quaternion_pow(const string& base, long long exp) {
    if (exp == 0) return "1";
    if (exp == 1) return base;
    
    string half = quaternion_pow(base, exp / 2);
    string result = quaternion_mult(half, half);
    
    if (exp % 2 == 1) {
        result = quaternion_mult(result, base);
    }
    
    return result;
}

// Check if the string can be split into i, j, k products
bool can_split_ijk(const string& s, long long repeat) {
    // First check if the total product is -1
    string total = "1";
    for (char c : s) {
        total = quaternion_mult(total, string(1, c));
    }
    total = quaternion_pow(total, repeat);
    if (total != "-1") return false;
    
    // Now check if we can find i and j in reasonable repeats
    const int max_check = min(10LL, repeat);
    string extended;
    for (int i = 0; i < max_check; ++i) {
        extended += s;
    }
    
    // Find i
    string current = "1";
    size_t pos = 0;
    while (pos < extended.size() && current != "i") {
        current = quaternion_mult(current, string(1, extended[pos]));
        pos++;
    }
    if (current != "i") return false;
    
    // Find j
    current = "1";
    while (pos < extended.size() && current != "j") {
        current = quaternion_mult(current, string(1, extended[pos]));
        pos++;
    }
    if (current != "j") return false;
    
    return true;
}
