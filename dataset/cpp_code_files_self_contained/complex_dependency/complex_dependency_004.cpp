#include <vector>
#include <climits>
#include <algorithm>
#include <unordered_map>

using namespace std;

/**
 * Enhanced game solver with additional features:
 * 1. Supports both optimal play and suboptimal play analysis
 * 2. Tracks the actual moves taken by players
 * 3. Handles negative numbers in the array
 * 4. Provides detailed game analysis
 */
class NumberGameSolver {
private:
    vector<int> nums;
    vector<int> prefixSum;
    unordered_map<string, int> dp;
    unordered_map<string, vector<pair<int, bool>>> moveTracker; // true for left, false for right

    int getSum(int l, int r) {
        if (l > r) return 0;
        return prefixSum[r] - (l > 0 ? prefixSum[l-1] : 0);
    }

    int solve(int l, int r, bool isPlayerA, bool trackMoves = false) {
        string key = to_string(l) + "-" + to_string(r) + "-" + to_string(isPlayerA);
        
        if (l > r) return 0;
        if (dp.find(key) != dp.end()) return dp[key];
        if (l == r) {
            if (trackMoves) moveTracker[key] = {{l, true}};
            return isPlayerA ? nums[l] : -nums[l];
        }

        int maxDiff = INT_MIN;
        vector<pair<int, bool>> bestMoves;

        // Try taking from left (1 to all remaining elements)
        for (int i = l; i <= r; i++) {
            int current = getSum(l, i);
            int remaining = solve(i+1, r, !isPlayerA, trackMoves);
            int diff = (isPlayerA ? current + remaining : -current + remaining);
            
            if (diff > maxDiff) {
                maxDiff = diff;
                if (trackMoves) {
                    bestMoves.clear();
                    for (int j = l; j <= i; j++) {
                        bestMoves.emplace_back(j, true);
                    }
                    string remainingKey = to_string(i+1) + "-" + to_string(r) + "-" + to_string(!isPlayerA);
                    if (moveTracker.find(remainingKey) != moveTracker.end()) {
                        bestMoves.insert(bestMoves.end(), moveTracker[remainingKey].begin(), moveTracker[remainingKey].end());
                    }
                }
            }
        }

        // Try taking from right (1 to all remaining elements)
        for (int i = r; i >= l; i--) {
            int current = getSum(i, r);
            int remaining = solve(l, i-1, !isPlayerA, trackMoves);
            int diff = (isPlayerA ? current + remaining : -current + remaining);
            
            if (diff > maxDiff) {
                maxDiff = diff;
                if (trackMoves) {
                    bestMoves.clear();
                    for (int j = r; j >= i; j--) {
                        bestMoves.emplace_back(j, false);
                    }
                    string remainingKey = to_string(l) + "-" + to_string(i-1) + "-" + to_string(!isPlayerA);
                    if (moveTracker.find(remainingKey) != moveTracker.end()) {
                        bestMoves.insert(bestMoves.end(), moveTracker[remainingKey].begin(), moveTracker[remainingKey].end());
                    }
                }
            }
        }

        if (trackMoves) {
            moveTracker[key] = bestMoves;
        }
        dp[key] = maxDiff;
        return maxDiff;
    }

public:
    NumberGameSolver(const vector<int>& numbers) : nums(numbers) {
        prefixSum.resize(nums.size());
        if (!nums.empty()) {
            prefixSum[0] = nums[0];
            for (int i = 1; i < nums.size(); i++) {
                prefixSum[i] = prefixSum[i-1] + nums[i];
            }
        }
    }

    int getOptimalDifference() {
        dp.clear();
        return solve(0, nums.size()-1, true);
    }

    vector<vector<int>> getOptimalMoves() {
        dp.clear();
        moveTracker.clear();
        solve(0, nums.size()-1, true, true);
        
        string initialKey = to_string(0) + "-" + to_string(nums.size()-1) + "-1";
        vector<pair<int, bool>> moves = moveTracker[initialKey];
        
        vector<vector<int>> result(2); // result[0] for player A, result[1] for player B
        bool isPlayerA = true;
        
        for (const auto& move : moves) {
            if (isPlayerA) {
                result[0].push_back(nums[move.first]);
            } else {
                result[1].push_back(nums[move.first]);
            }
            isPlayerA = !isPlayerA;
        }
        
        return result;
    }

    int simulateGame(const vector<bool>& playerAChoices) {
        int left = 0, right = nums.size() - 1;
        int scoreA = 0, scoreB = 0;
        bool isPlayerATurn = true;
        size_t choiceIndex = 0;

        while (left <= right && choiceIndex < playerAChoices.size()) {
            if (isPlayerATurn) {
                bool takeLeft = playerAChoices[choiceIndex++];
                if (takeLeft) {
                    scoreA += nums[left++];
                } else {
                    scoreA += nums[right--];
                }
            } else {
                // Player B plays optimally
                int takeLeftDiff = nums[left] - getSum(left+1, right);
                int takeRightDiff = nums[right] - getSum(left, right-1);
                
                if (takeLeftDiff > takeRightDiff) {
                    scoreB += nums[left++];
                } else {
                    scoreB += nums[right--];
                }
            }
            isPlayerATurn = !isPlayerATurn;
        }

        // Finish remaining moves with optimal play
        while (left <= right) {
            if (isPlayerATurn) {
                int takeLeftDiff = nums[left] - getSum(left+1, right);
                int takeRightDiff = nums[right] - getSum(left, right-1);
                
                if (takeLeftDiff > takeRightDiff) {
                    scoreA += nums[left++];
                } else {
                    scoreA += nums[right--];
                }
            } else {
                int takeLeftDiff = nums[left] - getSum(left+1, right);
                int takeRightDiff = nums[right] - getSum(left, right-1);
                
                if (takeLeftDiff > takeRightDiff) {
                    scoreB += nums[left++];
                } else {
                    scoreB += nums[right--];
                }
            }
            isPlayerATurn = !isPlayerATurn;
        }

        return scoreA - scoreB;
    }
};
