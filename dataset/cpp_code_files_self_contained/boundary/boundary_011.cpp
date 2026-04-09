#include <vector>
#include <string>
#include <stdexcept>
#include <algorithm>
#include <cmath>

using namespace std;

// Function to perform mathematical operations on a vector of numbers
vector<double> vector_math_operations(const vector<double>& numbers) {
    if (numbers.empty()) {
        throw invalid_argument("Error: Empty vector provided");
    }
    
    vector<double> results;
    
    // 1. Calculate pairwise differences
    for (size_t i = 0; i < numbers.size(); ++i) {
        for (size_t j = i + 1; j < numbers.size(); ++j) {
            results.push_back(numbers[i] - numbers[j]);
        }
    }
    
    // 2. Calculate absolute differences
    for (size_t i = 0; i < numbers.size(); ++i) {
        for (size_t j = i + 1; j < numbers.size(); ++j) {
            results.push_back(abs(numbers[i] - numbers[j]));
        }
    }
    
    // 3. Calculate normalized differences (avoid division by zero)
    for (size_t i = 0; i < numbers.size(); ++i) {
        for (size_t j = i + 1; j < numbers.size(); ++j) {
            if (numbers[j] != 0) {
                results.push_back((numbers[i] - numbers[j]) / numbers[j]);
            } else {
                results.push_back(NAN); // Not a Number for division by zero
            }
        }
    }
    
    // 4. Calculate vector magnitude of differences
    double sum_of_squares = 0;
    for (size_t i = 0; i < numbers.size(); ++i) {
        for (size_t j = i + 1; j < numbers.size(); ++j) {
            double diff = numbers[i] - numbers[j];
            sum_of_squares += diff * diff;
        }
    }
    results.push_back(sqrt(sum_of_squares));
    
    return results;
}
