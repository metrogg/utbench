#include <vector>
#include <map>
#include <bitset>
#include <iomanip>
#include <climits>

using namespace std;

// Function to count the number of set bits (1s) in an integer
int count_set_bits(unsigned int n) {
    int count = 0;
    while (n) {
        count += n & 1;
        n >>= 1;
    }
    return count;
}

// Function to reverse the bits of an integer
unsigned int reverse_bits(unsigned int n) {
    unsigned int reversed = 0;
    for (int i = 0; i < sizeof(n) * CHAR_BIT; i++) {
        reversed = (reversed << 1) | (n & 1);
        n >>= 1;
    }
    return reversed;
}

// Function to check if a number is a power of two
bool is_power_of_two(unsigned int n) {
    return n && !(n & (n - 1));
}

// Function to find the position of the most significant set bit
int msb_position(unsigned int n) {
    if (n == 0) return -1;
    int pos = 0;
    while (n >>= 1) {
        pos++;
    }
    return pos;
}

// Function to perform multiple bit operations and return results
map<string, string> perform_bit_operations(unsigned int n) {
    map<string, string> results;
    
    // Binary representation
    results["binary"] = bitset<32>(n).to_string();
    
    // Count of set bits
    results["set_bits"] = to_string(count_set_bits(n));
    
    // Reversed bits
    results["reversed"] = bitset<32>(reverse_bits(n)).to_string();
    
    // Is power of two
    results["is_power_of_two"] = is_power_of_two(n) ? "true" : "false";
    
    // MSB position
    results["msb_position"] = to_string(msb_position(n));
    
    return results;
}
