#include <vector>
#include <string>
#include <sstream>
#include <algorithm>
#include <numeric>
#include <climits>
#include <unordered_map>

using namespace std;

// Extended GCD function that returns both GCD and coefficients for Bézout's identity
struct ExtendedGCDResult {
    int gcd;
    int x;
    int y;
};

ExtendedGCDResult extendedGCD(int a, int b) {
    if (b == 0) {
        return {a, 1, 0};
    }
    auto result = extendedGCD(b, a % b);
    return {result.gcd, result.y, result.x - (a / b) * result.y};
}

// Function to calculate GCD of a vector of numbers using multiple methods
unordered_map<string, int> calculateGCDs(const vector<int>& numbers) {
    unordered_map<string, int> results;
    
    if (numbers.empty()) {
        return results;
    }
    
    // 1. Pairwise GCD (original method)
    int pairwise_gcd = 0;
    for (size_t i = 0; i < numbers.size(); ++i) {
        for (size_t j = i + 1; j < numbers.size(); ++j) {
            pairwise_gcd = max(pairwise_gcd, extendedGCD(numbers[i], numbers[j]).gcd);
        }
    }
    results["pairwise"] = pairwise_gcd;
    
    // 2. Sequential GCD (gcd of all numbers)
    int sequential_gcd = numbers[0];
    for (size_t i = 1; i < numbers.size(); ++i) {
        sequential_gcd = extendedGCD(sequential_gcd, numbers[i]).gcd;
        if (sequential_gcd == 1) break; // GCD can't be smaller than 1
    }
    results["sequential"] = sequential_gcd;
    
    // 3. Using C++ built-in GCD (C++17) for comparison
    int builtin_gcd = numbers[0];
    for (size_t i = 1; i < numbers.size(); ++i) {
        builtin_gcd = gcd(builtin_gcd, numbers[i]);
        if (builtin_gcd == 1) break;
    }
    results["builtin"] = builtin_gcd;
    
    return results;
}

// Function to parse input string into vector of integers
vector<int> parseInput(const string& input) {
    vector<int> numbers;
    istringstream iss(input);
    int num;
    while (iss >> num) {
        numbers.push_back(num);
    }
    return numbers;
}
