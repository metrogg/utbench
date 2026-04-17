#include <vector>
#include <map>
#include <cmath>
#include <algorithm>
using namespace std;

// Checks if a number is divisible by another
bool isDivisible(int a, int b) {
    if (b == 0) return false; // Handle division by zero
    return a % b == 0;
}

// Finds all divisors of a number (except 1 and itself)
vector<int> findProperDivisors(int n) {
    vector<int> divisors;
    if (n <= 1) return divisors;
    
    for (int i = 2; i <= sqrt(n); i++) {
        if (isDivisible(n, i)) {
            if (n/i == i) {
                divisors.push_back(i);
            } else {
                divisors.push_back(i);
                divisors.push_back(n/i);
            }
        }
    }
    sort(divisors.begin(), divisors.end());
    return divisors;
}

// Checks if a number is prime
bool isPrime(int n) {
    if (n <= 1) return false;
    if (n == 2) return true;
    if (n % 2 == 0) return false;
    
    for (int i = 3; i <= sqrt(n); i += 2) {
        if (isDivisible(n, i)) {
            return false;
        }
    }
    return true;
}

// Finds greatest common divisor using Euclidean algorithm
int findGCD(int a, int b) {
    while (b != 0) {
        int temp = b;
        b = a % b;
        a = temp;
    }
    return a;
}

// Finds least common multiple
int findLCM(int a, int b) {
    if (a == 0 || b == 0) return 0;
    return abs(a * b) / findGCD(a, b);
}

// Performs comprehensive number analysis
map<string, vector<int>> analyzeNumber(int n) {
    map<string, vector<int>> result;
    
    // Basic properties
    result["is_prime"] = {isPrime(n) ? 1 : 0};
    result["is_even"] = {isDivisible(n, 2) ? 1 : 0};
    
    // Divisors
    vector<int> divisors = findProperDivisors(n);
    result["proper_divisors"] = divisors;
    
    // Prime factors
    vector<int> primeFactors;
    int temp = n;
    for (int i = 2; i <= temp; i++) {
        while (isDivisible(temp, i) && isPrime(i)) {
            primeFactors.push_back(i);
            temp /= i;
        }
    }
    result["prime_factors"] = primeFactors;
    
    return result;
}
