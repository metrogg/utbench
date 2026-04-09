#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
#include <set>

using namespace std;

class StringAnalyzer {
public:
    // Enhanced function to find length of longest substring without repeating characters
    int lengthOfLongestSubstring(const string& s) {
        unordered_map<char, int> charMap;
        int maxLen = 0;
        int start = 0;
        
        for (int end = 0; end < s.size(); end++) {
            char currentChar = s[end];
            
            if (charMap.find(currentChar) != charMap.end() && charMap[currentChar] >= start) {
                start = charMap[currentChar] + 1;
            }
            
            charMap[currentChar] = end;
            maxLen = max(maxLen, end - start + 1);
        }
        
        return maxLen;
    }

    // New function to find all unique longest substrings without repeating characters
    vector<string> getAllLongestSubstrings(const string& s) {
        unordered_map<char, int> charMap;
        int maxLen = 0;
        int start = 0;
        set<string> resultSet;
        
        for (int end = 0; end < s.size(); end++) {
            char currentChar = s[end];
            
            if (charMap.find(currentChar) != charMap.end() && charMap[currentChar] >= start) {
                // Before moving start, check if current substring is longest
                if (end - start == maxLen) {
                    resultSet.insert(s.substr(start, maxLen));
                } else if (end - start > maxLen) {
                    maxLen = end - start;
                    resultSet.clear();
                    resultSet.insert(s.substr(start, maxLen));
                }
                start = charMap[currentChar] + 1;
            }
            
            charMap[currentChar] = end;
        }
        
        // Check the last possible substring
        if (s.size() - start == maxLen) {
            resultSet.insert(s.substr(start, maxLen));
        } else if (s.size() - start > maxLen) {
            maxLen = s.size() - start;
            resultSet.clear();
            resultSet.insert(s.substr(start, maxLen));
        }
        
        return vector<string>(resultSet.begin(), resultSet.end());
    }

    // New function to find the lexicographically smallest longest substring
    string getLexSmallestLongestSubstring(const string& s) {
        vector<string> allSubstrings = getAllLongestSubstrings(s);
        if (allSubstrings.empty()) return "";
        sort(allSubstrings.begin(), allSubstrings.end());
        return allSubstrings[0];
    }
};
