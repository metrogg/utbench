#include <stack>
#include <vector>
#include <algorithm>
#include <numeric>
#include <stdexcept>

using namespace std;

// Enhanced stack processor with multiple operations
vector<int> process_stack_operations(const vector<int>& inputs) {
    stack<int> data_stack;
    vector<int> removed_elements;
    
    for (int num : inputs) {
        if (num == 0) {
            if (!data_stack.empty()) {
                removed_elements.push_back(data_stack.top());
                data_stack.pop();
            }
        } 
        else if (num == -1) {
            // Special operation: reverse the stack
            vector<int> temp;
            while (!data_stack.empty()) {
                temp.push_back(data_stack.top());
                data_stack.pop();
            }
            for (int val : temp) {
                data_stack.push(val);
            }
        }
        else if (num == -2) {
            // Special operation: sort the stack
            vector<int> temp;
            while (!data_stack.empty()) {
                temp.push_back(data_stack.top());
                data_stack.pop();
            }
            sort(temp.begin(), temp.end());
            for (int val : temp) {
                data_stack.push(val);
            }
        }
        else {
            data_stack.push(num);
        }
    }
    
    // Convert stack to vector (top to bottom)
    vector<int> result;
    while (!data_stack.empty()) {
        result.push_back(data_stack.top());
        data_stack.pop();
    }
    reverse(result.begin(), result.end()); // To maintain original order
    
    // Add removed elements information at the end
    result.push_back(-999); // Separator
    result.insert(result.end(), removed_elements.begin(), removed_elements.end());
    
    return result;
}

// Function to calculate statistics about the processed stack
vector<int> calculate_stack_stats(const vector<int>& processed_result) {
    if (processed_result.empty()) {
        throw invalid_argument("Empty processed result");
    }
    
    // Find the separator (-999)
    auto sep_pos = find(processed_result.begin(), processed_result.end(), -999);
    if (sep_pos == processed_result.end()) {
        throw invalid_argument("Invalid processed result format");
    }
    
    // Extract main stack and removed elements
    vector<int> main_stack(processed_result.begin(), sep_pos);
    vector<int> removed_elements(sep_pos + 1, processed_result.end());
    
    vector<int> stats;
    
    // 0: Sum of remaining elements
    int sum_remaining = accumulate(main_stack.begin(), main_stack.end(), 0);
    stats.push_back(sum_remaining);
    
    // 1: Count of remaining elements
    stats.push_back(main_stack.size());
    
    // 2: Sum of removed elements
    int sum_removed = accumulate(removed_elements.begin(), removed_elements.end(), 0);
    stats.push_back(sum_removed);
    
    // 3: Count of removed elements
    stats.push_back(removed_elements.size());
    
    // 4: Max in remaining
    if (!main_stack.empty()) {
        stats.push_back(*max_element(main_stack.begin(), main_stack.end()));
    } else {
        stats.push_back(0);
    }
    
    // 5: Min in remaining
    if (!main_stack.empty()) {
        stats.push_back(*min_element(main_stack.begin(), main_stack.end()));
    } else {
        stats.push_back(0);
    }
    
    return stats;
}
