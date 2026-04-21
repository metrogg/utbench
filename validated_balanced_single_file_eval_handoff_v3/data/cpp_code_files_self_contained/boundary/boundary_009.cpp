#include <vector>
#include <algorithm>

using namespace std;

// Enhanced search in rotated sorted array with duplicates
// Returns a pair containing:
// - bool: whether target was found
// - int: number of occurrences of the target
pair<bool, int> searchRotatedArray(const vector<int>& nums, int target) {
    if (nums.empty()) return {false, 0};

    int left = 0, right = nums.size() - 1;
    int count = 0;

    while (left <= right) {
        int mid = left + (right - left) / 2;

        if (nums[mid] == target) {
            // Found target, now count all occurrences
            count = 1;
            
            // Count duplicates to the left
            int left_ptr = mid - 1;
            while (left_ptr >= 0 && nums[left_ptr] == target) {
                count++;
                left_ptr--;
            }
            
            // Count duplicates to the right
            int right_ptr = mid + 1;
            while (right_ptr < nums.size() && nums[right_ptr] == target) {
                count++;
                right_ptr++;
            }
            
            return {true, count};
        }

        // Handle the case where left, mid, and right are equal
        if (nums[left] == nums[mid] && nums[mid] == nums[right]) {
            left++;
            right--;
        }
        // Right side is sorted
        else if (nums[mid] <= nums[right]) {
            if (target > nums[mid] && target <= nums[right]) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        // Left side is sorted
        else {
            if (target >= nums[left] && target < nums[mid]) {
                right = mid - 1;
            } else {
                left = mid + 1;
            }
        }
    }

    return {false, 0};
}

// Helper function to find the rotation pivot index
int findPivot(const vector<int>& nums) {
    int left = 0, right = nums.size() - 1;
    
    while (left < right) {
        int mid = left + (right - left) / 2;
        
        if (nums[mid] > nums[right]) {
            left = mid + 1;
        } else if (nums[mid] < nums[right]) {
            right = mid;
        } else {
            right--;
        }
    }
    
    return left;
}
