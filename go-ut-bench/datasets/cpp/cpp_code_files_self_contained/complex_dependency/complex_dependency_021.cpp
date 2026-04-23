#include <vector>
#include <string>
#include <unordered_set>
#include <algorithm>
#include <cmath>

using namespace std;

// Function to check if a number contains all digits 0-9
bool containsAllDigits(long long number) {
    unordered_set<int> digits;
    if (number == 0) {
        digits.insert(0);
    } else {
        while (number > 0) {
            digits.insert(number % 10);
            number /= 10;
        }
    }
    return digits.size() == 10;
}

// Function to find the smallest multiple that contains all digits 0-9
// Returns -1 if no such multiple exists (though theoretically it always should)
long long findSmallestCompleteMultiple(int n) {
    if (n == 0) return -1; // Special case for zero
    
    unordered_set<int> seenDigits;
    long long currentMultiple;
    int maxIterations = 1000; // Safety limit to prevent infinite loops
    
    for (int i = 1; i <= maxIterations; i++) {
        currentMultiple = static_cast<long long>(n) * i;
        long long temp = currentMultiple;
        
        // Handle zero separately
        if (temp == 0) {
            seenDigits.insert(0);
        } else {
            while (temp > 0) {
                seenDigits.insert(temp % 10);
                temp /= 10;
            }
        }
        
        if (seenDigits.size() == 10) {
            return currentMultiple;
        }
    }
    
    return -1; // If we didn't find all digits in maxIterations
}

// Enhanced function that returns more detailed information
vector<long long> analyzeNumberProperties(int n) {
    vector<long long> result;
    
    // 1. Smallest multiple with all digits
    long long completeMultiple = findSmallestCompleteMultiple(n);
    result.push_back(completeMultiple);
    
    // 2. Number of multiples needed to find all digits
    if (n == 0) {
        result.push_back(-1);
    } else {
        int count = 0;
        unordered_set<int> seenDigits;
        for (int i = 1; i <= 1000; i++) {
            count++;
            long long temp = static_cast<long long>(n) * i;
            if (temp == 0) {
                seenDigits.insert(0);
            } else {
                while (temp > 0) {
                    seenDigits.insert(temp % 10);
                    temp /= 10;
                }
            }
            if (seenDigits.size() == 10) break;
        }
        result.push_back(seenDigits.size() == 10 ? count : -1);
    }
    
    // 3. First multiple that adds a new digit
    if (n == 0) {
        result.push_back(0);
    } else {
        unordered_set<int> seenDigits;
        for (int i = 1; i <= 1000; i++) {
            long long temp = static_cast<long long>(n) * i;
            int initialSize = seenDigits.size();
            if (temp == 0) {
                seenDigits.insert(0);
            } else {
                while (temp > 0) {
                    seenDigits.insert(temp % 10);
                    temp /= 10;
                }
            }
            if (seenDigits.size() > initialSize) {
                result.push_back(i);
                break;
            }
        }
    }
    
    return result;
}
