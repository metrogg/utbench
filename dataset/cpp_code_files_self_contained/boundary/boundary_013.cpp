#include <vector>
#include <string>
#include <stdexcept>
#include <map>

using namespace std;

// Function to compute factorial (n!)
unsigned long long factorial(int n) {
    if (n < 0) throw invalid_argument("Factorial of negative number");
    unsigned long long result = 1;
    for (int i = 2; i <= n; ++i) {
        result *= i;
    }
    return result;
}

// Function to generate kth permutation of first n letters
string kth_permutation(int n, unsigned long long k) {
    if (n < 1 || n > 26) throw invalid_argument("n must be between 1 and 26");
    if (k == 0 || k > factorial(n)) throw invalid_argument("k out of range");

    vector<char> letters;
    for (int i = 0; i < n; ++i) {
        letters.push_back('a' + i);
    }

    string result;
    k--; // convert to 0-based index

    for (int i = n; i >= 1; --i) {
        unsigned long long fact = factorial(i - 1);
        int index = k / fact;
        result += letters[index];
        letters.erase(letters.begin() + index);
        k %= fact;
    }

    return result;
}

// Extended function that provides permutation statistics
map<string, string> permutation_stats(int n, unsigned long long k) {
    map<string, string> stats;
    stats["input_n"] = to_string(n);
    stats["input_k"] = to_string(k);
    
    try {
        string perm = kth_permutation(n, k);
        stats["permutation"] = perm;
        stats["is_palindrome"] = (perm == string(perm.rbegin(), perm.rend())) ? "true" : "false";
        
        // Calculate character frequency
        map<char, int> freq;
        for (char c : perm) {
            freq[c]++;
        }
        string freq_str;
        for (auto& p : freq) {
            freq_str += string(1, p.first) + ":" + to_string(p.second) + " ";
        }
        stats["character_frequency"] = freq_str;
        
        // Calculate inversion count (naive O(n^2) method)
        int inversions = 0;
        for (int i = 0; i < perm.size(); ++i) {
            for (int j = i + 1; j < perm.size(); ++j) {
                if (perm[i] > perm[j]) inversions++;
            }
        }
        stats["inversion_count"] = to_string(inversions);
        
    } catch (const invalid_argument& e) {
        stats["error"] = e.what();
    }
    
    return stats;
}
