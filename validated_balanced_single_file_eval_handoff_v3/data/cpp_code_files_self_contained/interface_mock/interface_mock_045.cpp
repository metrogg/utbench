#include <vector>
#include <map>
#include <cmath>
#include <stdexcept>

using namespace std;

class GameCharacter {
private:
    int x, y;
    int speed;
    int size;
    int screenWidth, screenHeight;

public:
    GameCharacter(int startX, int startY, int characterSize, int speed, int screenW, int screenH)
        : x(startX), y(startY), size(characterSize), speed(speed), 
          screenWidth(screenW), screenHeight(screenH) {}

    void move(char direction) {
        switch(direction) {
            case 'L': x -= speed; break;
            case 'R': x += speed; break;
            case 'U': y -= speed; break;
            case 'D': y += speed; break;
        }
        clampPosition();
    }

    void clampPosition() {
        x = max(0, min(x, screenWidth - size));
        y = max(0, min(y, screenHeight - size));
    }

    bool checkCollision(int objX, int objY, int objSize) const {
        return (x < objX + objSize &&
                x + size > objX &&
                y < objY + objSize &&
                y + size > objY);
    }

    void setPosition(int newX, int newY) {
        x = newX;
        y = newY;
        clampPosition();
    }

    map<string, int> getPosition() const {
        return {{"x", x}, {"y", y}};
    }

    vector<vector<int>> getBoundingBox() const {
        return {
            {x, y},
            {x + size, y},
            {x + size, y + size},
            {x, y + size}
        };
    }

    double distanceTo(int targetX, int targetY) const {
        int dx = (x + size/2) - targetX;
        int dy = (y + size/2) - targetY;
        return sqrt(dx*dx + dy*dy);
    }
};

class GameWorld {
private:
    int width, height;
    vector<vector<int>> obstacles;
    vector<vector<int>> collectibles;

public:
    GameWorld(int w, int h) : width(w), height(h) {}

    void addObstacle(int x, int y, int size) {
        obstacles.push_back({x, y, size});
    }

    void addCollectible(int x, int y, int size) {
        collectibles.push_back({x, y, size});
    }

    bool checkObstacleCollision(int x, int y, int size) const {
        for (const auto& obs : obstacles) {
            if (x < obs[0] + obs[2] &&
                x + size > obs[0] &&
                y < obs[1] + obs[2] &&
                y + size > obs[1]) {
                return true;
            }
        }
        return false;
    }

    int checkCollectibleCollision(int x, int y, int size) {
        for (size_t i = 0; i < collectibles.size(); ++i) {
            const auto& col = collectibles[i];
            if (x < col[0] + col[2] &&
                x + size > col[0] &&
                y < col[1] + col[2] &&
                y + size > col[1]) {
                int points = col[2] * 10; // Bigger collectibles give more points
                collectibles.erase(collectibles.begin() + i);
                return points;
            }
        }
        return 0;
    }

    vector<vector<int>> getWorldGrid(int cellSize) const {
        vector<vector<int>> grid;
        for (int y = 0; y < height; y += cellSize) {
            for (int x = 0; x < width; x += cellSize) {
                grid.push_back({x, y});
            }
        }
        return grid;
    }
};
