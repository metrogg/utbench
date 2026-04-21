#include <vector>
#include <climits>
#include <algorithm>
using namespace std;

struct SubarrayResult {
    int max_sum;
    int start_index;
    int end_index;
    vector<int> subarray_elements;
};

// Enhanced function to find maximum subarray with additional information
SubarrayResult find_max_subarray(const vector<int>& nums) {
    if (nums.empty()) {
        return {0, -1, -1, {}};
    }

    int max_so_far = INT_MIN;
    int max_ending_here = 0;
    int start = 0, end = 0;
    int temp_start = 0;

    for (int i = 0; i < nums.size(); i++) {
        max_ending_here += nums[i];

        if (max_so_far < max_ending_here) {
            max_so_far = max_ending_here;
            start = temp_start;
            end = i;
        }

        if (max_ending_here < 0) {
            max_ending_here = 0;
            temp_start = i + 1;
        }
    }

    // Extract the subarray elements
    vector<int> subarray(nums.begin() + start, nums.begin() + end + 1);

    return {max_so_far, start, end, subarray};
}

// Additional function to find minimum subarray (for comparison)
SubarrayResult find_min_subarray(const vector<int>& nums) {
    if (nums.empty()) {
        return {0, -1, -1, {}};
    }

    int min_so_far = INT_MAX;
    int min_ending_here = 0;
    int start = 0, end = 0;
    int temp_start = 0;

    for (int i = 0; i < nums.size(); i++) {
        min_ending_here += nums[i];

        if (min_so_far > min_ending_here) {
            min_so_far = min_ending_here;
            start = temp_start;
            end = i;
        }

        if (min_ending_here > 0) {
            min_ending_here = 0;
            temp_start = i + 1;
        }
    }

    vector<int> subarray(nums.begin() + start, nums.begin() + end + 1);
    return {min_so_far, start, end, subarray};
}
