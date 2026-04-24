#include <vector>
#include <algorithm>
#include <map>
#include <utility>

using namespace std;

// Enhanced binary search that finds all occurrences of target in a sorted vector
// Returns a pair representing the first and last positions of the target
// If target is not found, returns (-1, -1)
pair<int, int> binarySearchRange(const vector<int>& nums, int target) {
    if (nums.empty()) return make_pair(-1, -1);
    
    // Find first occurrence
    int left = 0, right = nums.size() - 1;
    int first = -1;
    while (left <= right) {
        int mid = left + (right - left) / 2;
        if (nums[mid] >= target) {
            right = mid - 1;
        } else {
            left = mid + 1;
        }
        if (nums[mid] == target) first = mid;
    }
    
    // Find last occurrence
    left = 0, right = nums.size() - 1;
    int last = -1;
    while (left <= right) {
        int mid = left + (right - left) / 2;
        if (nums[mid] <= target) {
            left = mid + 1;
        } else {
            right = mid - 1;
        }
        if (nums[mid] == target) last = mid;
    }
    
    return make_pair(first, last);
}

// Counts the number of occurrences of target in the sorted vector
int countOccurrences(const vector<int>& nums, int target) {
    auto range = binarySearchRange(nums, target);
    if (range.first == -1) return 0;
    return range.second - range.first + 1;
}

// Finds the insertion position for target to maintain sorted order
int findInsertPosition(const vector<int>& nums, int target) {
    int left = 0, right = nums.size() - 1;
    while (left <= right) {
        int mid = left + (right - left) / 2;
        if (nums[mid] < target) {
            left = mid + 1;
        } else {
            right = mid - 1;
        }
    }
    return left;
}

// Checks if target exists in the vector
bool contains(const vector<int>& nums, int target) {
    auto range = binarySearchRange(nums, target);
    return range.first != -1;
}
