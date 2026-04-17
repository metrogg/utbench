#include <string>
#include <vector>
#include <set>
#include <algorithm>
#include <cmath>
#include <map>

using namespace std;

// Enhanced number manipulation utilities
namespace NumberUtils {
    int toNumber(const string &str) {
        int num = 0;
        for (char c : str) {
            if (isdigit(c)) {
                num = num * 10 + (c - '0');
            }
        }
        return num;
    }

    string toString(int num) {
        if (num == 0) return "0";
        string str;
        while (num > 0) {
            str += '0' + (num % 10);
            num /= 10;
        }
        reverse(str.begin(), str.end());
        return str;
    }

    vector<int> getAllRotations(int num) {
        string s = toString(num);
        set<int> rotations;
        
        for (size_t i = 1; i < s.length(); i++) {
            string rotated = s.substr(i) + s.substr(0, i);
            rotations.insert(toNumber(rotated));
        }
        
        return vector<int>(rotations.begin(), rotations.end());
    }

    bool isPrime(int n) {
        if (n <= 1) return false;
        if (n == 2) return true;
        if (n % 2 == 0) return false;
        for (int i = 3; i * i <= n; i += 2) {
            if (n % i == 0) return false;
        }
        return true;
    }
}

// Enhanced number analyzer with statistics
map<string, int> analyzeNumberRange(int A, int B) {
    map<string, int> results;
    set<int> countedPairs;
    int totalRecycled = 0;
    int totalPrimes = 0;
    int totalPalindromes = 0;
    int maxRotationCount = 0;
    set<int> uniquePrimes;
    set<int> uniquePalindromes;

    for (int n = A; n <= B; n++) {
        // Check for prime numbers
        if (NumberUtils::isPrime(n)) {
            totalPrimes++;
            uniquePrimes.insert(n);
        }

        // Check for palindromes
        string s = NumberUtils::toString(n);
        string rev = s;
        reverse(rev.begin(), rev.end());
        if (s == rev) {
            totalPalindromes++;
            uniquePalindromes.insert(n);
        }

        // Count recycled pairs
        vector<int> rotations = NumberUtils::getAllRotations(n);
        maxRotationCount = max(maxRotationCount, (int)rotations.size());

        for (int rotated : rotations) {
            if (rotated > n && rotated <= B) {
                pair<int, int> p = make_pair(n, rotated);
                if (countedPairs.find(rotated) == countedPairs.end()) {
                    totalRecycled++;
                    countedPairs.insert(rotated);
                }
            }
        }
    }

    results["total_numbers"] = B - A + 1;
    results["recycled_pairs"] = totalRecycled;
    results["prime_numbers"] = totalPrimes;
    results["unique_primes"] = uniquePrimes.size();
    results["palindromes"] = totalPalindromes;
    results["unique_palindromes"] = uniquePalindromes.size();
    results["max_rotations"] = maxRotationCount;

    return results;
}
