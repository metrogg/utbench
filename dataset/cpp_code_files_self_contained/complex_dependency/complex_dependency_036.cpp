#include <vector>
#include <algorithm>
#include <map>
#include <iomanip>

using namespace std;

// Enhanced Fibonacci search with additional statistics and error checking
map<string, long> fibonacci_search(const vector<long>& arr, long value) {
    map<string, long> result;
    result["found_index"] = -1;
    result["comparisons"] = 0;
    result["fib_steps"] = 0;
    
    if (arr.empty()) {
        result["error"] = 1; // Empty array error
        return result;
    }

    // Check if array is sorted
    if (!is_sorted(arr.begin(), arr.end())) {
        result["error"] = 2; // Unsorted array error
        return result;
    }

    long size = arr.size();
    long fib2 = 0; // (m-2)'th Fibonacci number
    long fib1 = 1; // (m-1)'th Fibonacci number
    long fibM = fib2 + fib1; // m'th Fibonacci number
    result["fib_steps"]++;

    // Find the smallest Fibonacci number greater than or equal to size
    while (fibM < size) {
        fib2 = fib1;
        fib1 = fibM;
        fibM = fib2 + fib1;
        result["fib_steps"]++;
    }

    long offset = -1;

    while (fibM > 1) {
        result["comparisons"]++;
        long i = min(offset + fib2, size - 1);

        if (arr[i] < value) {
            result["comparisons"]++;
            fibM = fib1;
            fib1 = fib2;
            fib2 = fibM - fib1;
            offset = i;
        }
        else if (arr[i] > value) {
            result["comparisons"] += 2;
            fibM = fib2;
            fib1 = fib1 - fib2;
            fib2 = fibM - fib1;
        }
        else {
            result["comparisons"] += 2;
            result["found_index"] = i;
            return result;
        }
    }

    if (fib1 && offset + 1 < size && arr[offset + 1] == value) {
        result["comparisons"]++;
        result["found_index"] = offset + 1;
    }

    return result;
}

// Helper function to generate Fibonacci sequence up to n terms
vector<long> generate_fibonacci_sequence(long n) {
    vector<long> sequence;
    if (n <= 0) return sequence;
    
    sequence.push_back(0);
    if (n == 1) return sequence;
    
    sequence.push_back(1);
    for (long i = 2; i < n; ++i) {
        sequence.push_back(sequence[i-1] + sequence[i-2]);
    }
    
    return sequence;
}
