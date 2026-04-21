#include <vector>
#include <map>
#include <algorithm>
#include <climits>
#include <cmath>
#include <iomanip>

using namespace std;

// Enhanced function to analyze an array of integers
map<string, double> analyze_array(const vector<int>& arr) {
    map<string, double> results;
    
    if (arr.empty()) {
        results["error"] = -1;
        return results;
    }

    // Basic statistics
    int sum = 0;
    int min_val = INT_MAX;
    int max_val = INT_MIN;
    map<int, int> frequency_map;

    for (int num : arr) {
        sum += num;
        if (num < min_val) min_val = num;
        if (num > max_val) max_val = num;
        frequency_map[num]++;
    }

    double mean = static_cast<double>(sum) / arr.size();
    results["sum"] = sum;
    results["mean"] = mean;
    results["min"] = min_val;
    results["max"] = max_val;

    // Calculate median
    vector<int> sorted_arr = arr;
    sort(sorted_arr.begin(), sorted_arr.end());
    if (sorted_arr.size() % 2 == 1) {
        results["median"] = sorted_arr[sorted_arr.size() / 2];
    } else {
        results["median"] = (sorted_arr[sorted_arr.size() / 2 - 1] + 
                            sorted_arr[sorted_arr.size() / 2]) / 2.0;
    }

    // Calculate mode(s)
    int max_freq = 0;
    vector<int> modes;
    for (const auto& pair : frequency_map) {
        if (pair.second > max_freq) {
            max_freq = pair.second;
            modes.clear();
            modes.push_back(pair.first);
        } else if (pair.second == max_freq) {
            modes.push_back(pair.first);
        }
    }
    results["mode_count"] = modes.size();
    if (modes.size() == 1) {
        results["mode"] = modes[0];
    }

    // Calculate standard deviation
    double variance = 0.0;
    for (int num : arr) {
        variance += pow(num - mean, 2);
    }
    variance /= arr.size();
    results["std_dev"] = sqrt(variance);

    // Frequency analysis
    results["unique_count"] = frequency_map.size();
    results["most_frequent"] = max_freq;

    return results;
}
