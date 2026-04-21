#include <vector>
#include <algorithm>
#include <map>
#include <climits>

using namespace std;

class SpatialAnalyzer {
private:
    vector<int> fenwickTree;
    int maxCoord;

    void update(int index, int delta) {
        for (; index <= maxCoord; index += index & -index) {
            fenwickTree[index] += delta;
        }
    }

    int query(int index) {
        int sum = 0;
        for (; index > 0; index -= index & -index) {
            sum += fenwickTree[index];
        }
        return sum;
    }

public:
    SpatialAnalyzer(int maxCoordinate) : maxCoord(maxCoordinate) {
        fenwickTree.resize(maxCoord + 2, 0);
    }

    // Process points and return level counts
    map<int, int> countLevels(const vector<pair<int, int>>& points) {
        map<int, int> levelCounts;
        for (const auto& point : points) {
            int x = point.first + 1; // 1-based indexing
            int currentLevel = query(x);
            levelCounts[currentLevel]++;
            update(x, 1);
        }
        return levelCounts;
    }

    // Range query: count points in [x1, x2] × [y1, y2]
    int rangeQuery(const vector<pair<int, int>>& points, 
                  int x1, int x2, int y1, int y2) {
        // Reset Fenwick Tree
        fill(fenwickTree.begin(), fenwickTree.end(), 0);
        
        // Sort points by y-coordinate for efficient processing
        vector<pair<int, int>> sortedPoints = points;
        sort(sortedPoints.begin(), sortedPoints.end(), 
            [](const auto& a, const auto& b) { return a.second < b.second; });

        int count = 0;
        for (const auto& point : sortedPoints) {
            if (point.second >= y1 && point.second <= y2) {
                int x = point.first + 1;
                if (x >= x1 + 1 && x <= x2 + 1) {
                    count++;
                }
            }
        }
        return count;
    }

    // Find points dominating a given point (x', y' >= x, y)
    int countDominatingPoints(const vector<pair<int, int>>& points, int x, int y) {
        // Reset Fenwick Tree
        fill(fenwickTree.begin(), fenwickTree.end(), 0);
        
        // Sort points by y-coordinate in descending order
        vector<pair<int, int>> sortedPoints = points;
        sort(sortedPoints.begin(), sortedPoints.end(), 
            [](const auto& a, const auto& b) { return a.second > b.second; });

        int count = 0;
        for (const auto& point : sortedPoints) {
            if (point.second >= y) {
                int currentX = point.first + 1;
                if (currentX >= x + 1) {
                    count += query(currentX);
                    update(currentX, 1);
                }
            } else {
                break;
            }
        }
        return count;
    }
};
