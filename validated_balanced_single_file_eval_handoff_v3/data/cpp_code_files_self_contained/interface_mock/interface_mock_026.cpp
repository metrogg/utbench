#include <vector>
#include <map>
#include <cmath>
#include <algorithm>

using namespace std;

struct PlayerState {
    float position_x;
    float position_y;
    bool alive;
    bool shooting;
    float speed;
    int health;
    int score;
    vector<pair<float, float>> projectiles;
};

class AdvancedPlayer {
private:
    PlayerState state;
    float boundary_left;
    float boundary_right;
    float boundary_top;
    float boundary_bottom;
    float projectile_speed;
    int max_health;

public:
    AdvancedPlayer(float start_x, float start_y, 
                  float left_bound, float right_bound,
                  float top_bound, float bottom_bound,
                  int initial_health = 100)
        : boundary_left(left_bound), boundary_right(right_bound),
          boundary_top(top_bound), boundary_bottom(bottom_bound),
          projectile_speed(500.0f), max_health(initial_health) {
        state.position_x = start_x;
        state.position_y = start_y;
        state.alive = true;
        state.shooting = false;
        state.speed = 0.0f;
        state.health = initial_health;
        state.score = 0;
    }

    void processMovement(float delta_time, bool move_left, bool move_right, bool move_up, bool move_down) {
        if (!state.alive) return;

        float move_x = 0.0f;
        float move_y = 0.0f;

        if (move_left) move_x -= 300.0f;
        if (move_right) move_x += 300.0f;
        if (move_up) move_y -= 300.0f;
        if (move_down) move_y += 300.0f;

        // Normalize diagonal movement
        if (move_x != 0 && move_y != 0) {
            move_x *= 0.7071f; // 1/sqrt(2)
            move_y *= 0.7071f;
        }

        state.position_x += move_x * delta_time;
        state.position_y += move_y * delta_time;

        // Apply boundaries
        state.position_x = max(boundary_left, min(state.position_x, boundary_right));
        state.position_y = max(boundary_top, min(state.position_y, boundary_bottom));

        state.speed = sqrt(move_x * move_x + move_y * move_y);
    }

    void processShooting(float delta_time, bool shoot_command) {
        if (!state.alive) return;

        if (shoot_command && !state.shooting) {
            // Add new projectile at player position
            state.projectiles.emplace_back(state.position_x, state.position_y - 20.0f);
            state.shooting = true;
        }
        else if (!shoot_command) {
            state.shooting = false;
        }

        // Update existing projectiles
        for (auto& proj : state.projectiles) {
            proj.second -= projectile_speed * delta_time;
        }

        // Remove projectiles that are out of bounds
        state.projectiles.erase(
            remove_if(state.projectiles.begin(), state.projectiles.end(),
                [this](const pair<float, float>& p) {
                    return p.second < boundary_top;
                }),
            state.projectiles.end()
        );
    }

    void takeDamage(int damage) {
        if (!state.alive) return;
        
        state.health -= damage;
        if (state.health <= 0) {
            state.health = 0;
            state.alive = false;
        }
    }

    void addScore(int points) {
        state.score += points;
    }

    void heal(int amount) {
        if (!state.alive) return;
        state.health = min(max_health, state.health + amount);
    }

    const PlayerState& getState() const {
        return state;
    }

    vector<pair<float, float>> getActiveProjectiles() const {
        return state.projectiles;
    }
};
