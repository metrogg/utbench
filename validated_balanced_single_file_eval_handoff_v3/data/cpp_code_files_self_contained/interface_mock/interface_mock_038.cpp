#include <vector>
#include <cmath>
#include <map>
#include <string>

using namespace std;

// Game object structure
struct GameObject {
    double x, y;
    double width, height;
    string type;
    bool active;
    int value;
};

// Character structure
struct Character {
    double x, y;
    double width, height;
    int score;
    int lives;
    bool invincible;
    int invincibleTimer;
};

// Collision detection function
bool checkCollision(const GameObject& obj1, const GameObject& obj2) {
    return !(obj1.x > obj2.x + obj2.width ||
             obj1.x + obj1.width < obj2.x ||
             obj1.y > obj2.y + obj2.height ||
             obj1.y + obj1.height < obj2.y);
}

// Enhanced collision system with multiple game object types
map<string, string> handleCollisions(Character& player, vector<GameObject>& objects) {
    map<string, string> results;
    results["score_change"] = "0";
    results["lives_change"] = "0";
    results["powerup"] = "none";
    results["special_effect"] = "none";

    // Convert player to GameObject for collision detection
    GameObject playerObj = {player.x, player.y, player.width, player.height, "player", true, 0};

    for (auto& obj : objects) {
        if (!obj.active) continue;

        if (checkCollision(playerObj, obj)) {
            if (obj.type == "coin") {
                // Coin collection
                player.score += obj.value;
                obj.active = false;
                results["score_change"] = to_string(obj.value);
                results["special_effect"] = "coin_collected";
            }
            else if (obj.type == "bomb" && !player.invincible) {
                // Bomb collision
                player.lives -= 1;
                obj.active = false;
                results["lives_change"] = "-1";
                results["special_effect"] = "explosion";
                player.invincible = true;
                player.invincibleTimer = 180; // 3 seconds at 60fps
            }
            else if (obj.type == "heart") {
                // Heart collection
                player.lives = min(player.lives + 1, 3); // Max 3 lives
                obj.active = false;
                results["lives_change"] = "+1";
                results["special_effect"] = "health_restored";
            }
            else if (obj.type == "powerup") {
                // Powerup collection
                obj.active = false;
                results["powerup"] = obj.type + "_powerup";
                results["special_effect"] = "powerup_collected";
            }
        }
    }

    // Handle invincibility frames
    if (player.invincible) {
        player.invincibleTimer--;
        if (player.invincibleTimer <= 0) {
            player.invincible = false;
        }
    }

    return results;
}
