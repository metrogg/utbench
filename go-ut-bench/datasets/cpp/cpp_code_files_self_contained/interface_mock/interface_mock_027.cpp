#include <vector>
#include <cmath>
#include <algorithm>
#include <climits>
#include <map>
#include <string>

using namespace std;

class GhostEscapeAnalyzer {
private:
    struct Position {
        int x;
        int y;
        Position(int x, int y) : x(x), y(y) {}
    };

    Position player;
    Position target;
    vector<Position> ghosts;
    vector<Position> obstacles;

    int manhattanDistance(const Position& a, const Position& b) {
        return abs(a.x - b.x) + abs(a.y - b.y);
    }

    bool isPositionBlocked(int x, int y) {
        for (const auto& obs : obstacles) {
            if (obs.x == x && obs.y == y) return true;
        }
        return false;
    }

public:
    GhostEscapeAnalyzer(int playerX, int playerY, int targetX, int targetY) 
        : player(playerX, playerY), target(targetX, targetY) {}

    void addGhost(int x, int y) {
        ghosts.emplace_back(x, y);
    }

    void addObstacle(int x, int y) {
        obstacles.emplace_back(x, y);
    }

    // Enhanced escape analysis with multiple factors
    map<string, bool> analyzeEscape() {
        map<string, bool> results;
        
        // Basic escape check
        int playerDistance = manhattanDistance(player, target);
        bool canEscape = true;
        
        for (const auto& ghost : ghosts) {
            int ghostDistance = manhattanDistance(ghost, target);
            if (ghostDistance <= playerDistance) {
                canEscape = false;
                break;
            }
        }
        results["basic_escape"] = canEscape;

        // Check if path is blocked by obstacles
        bool pathBlocked = false;
        int steps = playerDistance;
        Position current = player;
        
        for (int i = 0; i < steps; ++i) {
            if (current.x < target.x) current.x++;
            else if (current.x > target.x) current.x--;
            else if (current.y < target.y) current.y++;
            else if (current.y > target.y) current.y--;
            
            if (isPositionBlocked(current.x, current.y)) {
                pathBlocked = true;
                break;
            }
        }
        results["path_clear"] = !pathBlocked;

        // Check if any ghost can intercept
        bool canIntercept = false;
        for (const auto& ghost : ghosts) {
            int ghostToPlayer = manhattanDistance(ghost, player);
            int ghostToTarget = manhattanDistance(ghost, target);
            
            if (ghostToPlayer <= playerDistance && ghostToTarget <= playerDistance) {
                canIntercept = true;
                break;
            }
        }
        results["intercept_possible"] = canIntercept;

        // Final escape possibility considering all factors
        results["can_escape"] = canEscape && !pathBlocked && !canIntercept;

        return results;
    }

    // Find safest path distance considering ghosts
    int findSafestPathDistance() {
        if (player.x == target.x && player.y == target.y) return 0;

        int minDistance = manhattanDistance(player, target);
        int maxGhostThreat = 0;

        for (const auto& ghost : ghosts) {
            int ghostDist = manhattanDistance(ghost, player);
            if (ghostDist < minDistance) {
                maxGhostThreat = max(maxGhostThreat, minDistance - ghostDist);
            }
        }

        return minDistance + maxGhostThreat;
    }
};
