#include <vector>
#include <set>
#include <map>
#include <algorithm>
#include <cmath>

using namespace std;

// Generates a sequence based on the given parameters and analyzes it
map<string, unsigned int> analyze_sequence(unsigned int N, unsigned int S, 
                                          unsigned int P, unsigned int Q) {
    const unsigned int m = 1 << 31;
    vector<unsigned int> sequence(N);
    map<string, unsigned int> analysis;

    // Generate sequence
    sequence[0] = S % m;
    for (unsigned int i = 1; i < N; i++) {
        sequence[i] = (sequence[i-1] * P + Q) % m;
    }

    // Basic analysis
    set<unsigned int> unique_values(sequence.begin(), sequence.end());
    analysis["unique_count"] = unique_values.size();
    analysis["sequence_length"] = N;

    // Detect cycle
    unsigned int cycle_length = 0;
    map<unsigned int, unsigned int> value_positions;
    for (unsigned int i = 0; i < sequence.size(); i++) {
        if (value_positions.count(sequence[i])) {
            cycle_length = i - value_positions[sequence[i]];
            break;
        }
        value_positions[sequence[i]] = i;
    }
    analysis["cycle_length"] = cycle_length;

    // Calculate statistics
    unsigned int min_val = *min_element(sequence.begin(), sequence.end());
    unsigned int max_val = *max_element(sequence.begin(), sequence.end());
    analysis["min_value"] = min_val;
    analysis["max_value"] = max_val;
    analysis["value_range"] = max_val - min_val;

    // Count even/odd numbers
    unsigned int even_count = count_if(sequence.begin(), sequence.end(), 
                                     [](unsigned int x) { return x % 2 == 0; });
    analysis["even_count"] = even_count;
    analysis["odd_count"] = N - even_count;

    return analysis;
}
