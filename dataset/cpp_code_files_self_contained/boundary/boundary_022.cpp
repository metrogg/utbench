#include <vector>
#include <map>
#include <stdexcept>

using namespace std;

// Enhanced function to count binary numbers with various constraints
map<string, long long> count_binary_numbers(int n) {
    if (n <= 0) {
        throw invalid_argument("Error: n must be positive");
    }
    if (n > 90) {
        throw invalid_argument("Error: n must be <= 90 to prevent overflow");
    }

    vector<long long> num(n + 1, 0);
    vector<long long> zero(n + 1, 0);
    vector<long long> one(n + 1, 0);
    vector<long long> consecutive_zero(n + 1, 0);
    vector<long long> no_three_consecutive(n + 1, 0);

    // Base cases
    num[1] = 1;
    num[2] = 2;
    zero[1] = 1;
    one[1] = 1;
    zero[2] = 2;
    one[2] = 1;
    consecutive_zero[1] = 1;
    consecutive_zero[2] = 2;
    no_three_consecutive[1] = 2;
    no_three_consecutive[2] = 3;

    for (int i = 3; i <= n; i++) {
        // Original problem: no two consecutive 1s
        zero[i] = num[i - 1];
        one[i] = zero[i - 1];
        num[i] = one[i - 1] + zero[i - 1];

        // Additional constraint 1: no more than 2 consecutive 0s
        consecutive_zero[i] = consecutive_zero[i - 1] + consecutive_zero[i - 2];

        // Additional constraint 2: no three consecutive identical digits
        no_three_consecutive[i] = no_three_consecutive[i - 1] + no_three_consecutive[i - 2];
    }

    map<string, long long> results;
    results["no_consecutive_1s"] = num[n];
    results["no_more_than_2_0s"] = consecutive_zero[n];
    results["no_three_consecutive_same"] = no_three_consecutive[n];
    results["ending_with_0"] = zero[n];
    results["ending_with_1"] = one[n];

    return results;
}
