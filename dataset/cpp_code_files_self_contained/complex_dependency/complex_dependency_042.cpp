#include <vector>
#include <algorithm>
#include <cmath>
#include <map>
#include <iomanip>

using namespace std;

// Function to calculate comprehensive statistics for a vector of integers
map<string, double> calculate_statistics(const vector<int>& numbers) {
    map<string, double> stats;
    
    if (numbers.empty()) {
        stats["error"] = 1;
        return stats;
    }

    // Basic calculations
    int sum = 0;
    long long product = 1;
    int min_val = numbers[0];
    int max_val = numbers[0];

    for (int num : numbers) {
        sum += num;
        product *= num;
        if (num < min_val) min_val = num;
        if (num > max_val) max_val = num;
    }

    double average = static_cast<double>(sum) / numbers.size();

    // Variance calculation
    double variance = 0;
    for (int num : numbers) {
        variance += pow(num - average, 2);
    }
    variance /= numbers.size();

    // Median calculation
    vector<int> sorted_numbers = numbers;
    sort(sorted_numbers.begin(), sorted_numbers.end());
    double median;
    if (sorted_numbers.size() % 2 == 0) {
        median = (sorted_numbers[sorted_numbers.size()/2 - 1] + 
                 sorted_numbers[sorted_numbers.size()/2]) / 2.0;
    } else {
        median = sorted_numbers[sorted_numbers.size()/2];
    }

    // Store all statistics
    stats["sum"] = sum;
    stats["product"] = product;
    stats["average"] = average;
    stats["min"] = min_val;
    stats["max"] = max_val;
    stats["range"] = max_val - min_val;
    stats["variance"] = variance;
    stats["std_dev"] = sqrt(variance);
    stats["median"] = median;

    return stats;
}
