#include <vector>
#include <algorithm>
#include <map>
#include <utility>

using namespace std;

// Enhanced function to find wiggle subsequence with additional statistics
map<string, int> analyzeWiggleSequence(const vector<int>& nums) {
    map<string, int> result;
    
    if(nums.empty()) {
        result["max_length"] = 0;
        result["increasing_count"] = 0;
        result["decreasing_count"] = 0;
        result["peak_count"] = 0;
        result["valley_count"] = 0;
        return result;
    }
    
    int p = 1, n = 1; // p for positive difference, n for negative difference
    int increasing = 0, decreasing = 0;
    int peaks = 0, valleys = 0;
    
    for(int i = 1; i < nums.size(); i++) {
        if(nums[i] > nums[i-1]) {
            p = max(p, n + 1);
            increasing++;
            // Check if this is a valley to peak transition
            if(i > 1 && nums[i-1] < nums[i-2]) {
                peaks++;
            }
        }
        else if(nums[i] < nums[i-1]) {
            n = max(n, p + 1);
            decreasing++;
            // Check if this is a peak to valley transition
            if(i > 1 && nums[i-1] > nums[i-2]) {
                valleys++;
            }
        }
    }
    
    result["max_length"] = max(p, n);
    result["increasing_count"] = increasing;
    result["decreasing_count"] = decreasing;
    result["peak_count"] = peaks;
    result["valley_count"] = valleys;
    
    return result;
}

// Additional function to extract the actual wiggle subsequence
vector<int> getWiggleSubsequence(const vector<int>& nums) {
    if(nums.empty()) return {};
    
    vector<int> subsequence;
    subsequence.push_back(nums[0]);
    int direction = 0; // 0: undetermined, 1: increasing, -1: decreasing
    
    for(int i = 1; i < nums.size(); i++) {
        if(nums[i] > nums[i-1]) {
            if(direction <= 0) {
                subsequence.push_back(nums[i]);
                direction = 1;
            }
        }
        else if(nums[i] < nums[i-1]) {
            if(direction >= 0) {
                subsequence.push_back(nums[i]);
                direction = -1;
            }
        }
    }
    
    return subsequence;
}
