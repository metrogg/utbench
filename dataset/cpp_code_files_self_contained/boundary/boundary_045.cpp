#include <vector>
#include <map>
#include <cmath>
#include <algorithm>

using namespace std;

// Function to calculate various series properties
map<string, double> calculate_series_properties(const vector<int>& series) {
    map<string, double> results;
    
    if (series.empty()) {
        results["error"] = 1;
        return results;
    }
    
    // Basic calculations
    int sum = 0;
    int product = 1;
    int min_val = series[0];
    int max_val = series[0];
    double harmonic_sum = 0.0;
    double geometric_sum = 0.0;
    
    for (int num : series) {
        sum += num;
        product *= num;
        min_val = min(min_val, num);
        max_val = max(max_val, num);
        harmonic_sum += 1.0 / num;
        geometric_sum += log(num);
    }
    
    // Store results
    results["sum"] = sum;
    results["product"] = product;
    results["min"] = min_val;
    results["max"] = max_val;
    results["arithmetic_mean"] = static_cast<double>(sum) / series.size();
    results["harmonic_mean"] = series.size() / harmonic_sum;
    results["geometric_mean"] = exp(geometric_sum / series.size());
    
    // Calculate variance and standard deviation
    double mean = results["arithmetic_mean"];
    double variance = 0.0;
    for (int num : series) {
        variance += pow(num - mean, 2);
    }
    variance /= series.size();
    results["variance"] = variance;
    results["std_dev"] = sqrt(variance);
    
    return results;
}

// Function to generate different types of series
vector<int> generate_series(const string& series_type, int n) {
    vector<int> series;
    
    if (n <= 0) return series;
    
    if (series_type == "natural") {
        for (int i = 1; i <= n; i++) {
            series.push_back(i);
        }
    }
    else if (series_type == "even") {
        for (int i = 1; i <= n; i++) {
            series.push_back(2 * i);
        }
    }
    else if (series_type == "odd") {
        for (int i = 0; i < n; i++) {
            series.push_back(2 * i + 1);
        }
    }
    else if (series_type == "square") {
        for (int i = 1; i <= n; i++) {
            series.push_back(i * i);
        }
    }
    else if (series_type == "fibonacci") {
        if (n >= 1) series.push_back(1);
        if (n >= 2) series.push_back(1);
        for (int i = 2; i < n; i++) {
            series.push_back(series[i-1] + series[i-2]);
        }
    }
    
    return series;
}
