#include <vector>
#include <cmath>
#include <map>
#include <stdexcept>

using namespace std;

// Simplified vector3 class to replace glm::vec3
class Vector3 {
public:
    float x, y, z;

    Vector3(float x = 0, float y = 0, float z = 0) : x(x), y(y), z(z) {}

    Vector3 operator+(const Vector3& other) const {
        return Vector3(x + other.x, y + other.y, z + other.z);
    }

    Vector3 operator-(const Vector3& other) const {
        return Vector3(x - other.x, y - other.y, z - other.z);
    }

    Vector3 operator*(float scalar) const {
        return Vector3(x * scalar, y * scalar, z * scalar);
    }

    float magnitude() const {
        return sqrt(x*x + y*y + z*z);
    }

    Vector3 normalized() const {
        float mag = magnitude();
        if (mag > 0) {
            return *this * (1.0f / mag);
        }
        return *this;
    }
};

// Simulated input system to replace GLFW
class InputSystem {
private:
    map<int, bool> keyStates;

public:
    void setKeyState(int key, bool state) {
        keyStates[key] = state;
    }

    bool getKey(int key) const {
        auto it = keyStates.find(key);
        return it != keyStates.end() ? it->second : false;
    }
};

// Enhanced Player class with physics and collision
class Player {
private:
    Vector3 position;
    Vector3 velocity;
    Vector3 acceleration;
    bool isGrounded;
    float moveSpeed;
    float jumpForce;
    float gravity;
    float groundY;

    vector<Vector3> obstacles;

public:
    Player() : position(0, 0, 0), velocity(0, 0, 0), acceleration(0, 0, 0),
               isGrounded(false), moveSpeed(5.0f), jumpForce(7.0f),
               gravity(-9.8f), groundY(0.0f) {
        // Add some sample obstacles
        obstacles.push_back(Vector3(2, 0, 2));
        obstacles.push_back(Vector3(-2, 0, -2));
    }

    void Update(float deltaTime, const InputSystem& input) {
        // Process input
        Vector3 inputMovement(0, 0, 0);
        if (input.getKey(1)) inputMovement.x -= 1; // LEFT
        if (input.getKey(2)) inputMovement.x += 1; // RIGHT
        if (input.getKey(3)) inputMovement.z -= 1; // UP
        if (input.getKey(4)) inputMovement.z += 1; // DOWN

        // Normalize input if moving diagonally
        if (inputMovement.magnitude() > 0) {
            inputMovement = inputMovement.normalized();
        }

        // Apply movement
        acceleration.x = inputMovement.x * moveSpeed;
        acceleration.z = inputMovement.z * moveSpeed;

        // Handle jumping
        if (input.getKey(5) && isGrounded) { // SPACE
            velocity.y = jumpForce;
            isGrounded = false;
        }

        // Apply gravity
        acceleration.y = gravity;

        // Update velocity and position
        velocity = velocity + acceleration * deltaTime;
        Vector3 newPosition = position + velocity * deltaTime;

        // Check for collisions with obstacles
        bool collided = false;
        for (const auto& obstacle : obstacles) {
            float distance = (newPosition - obstacle).magnitude();
            if (distance < 1.5f) { // Simple collision radius
                collided = true;
                break;
            }
        }

        // Update position if no collision
        if (!collided) {
            position = newPosition;
        } else {
            velocity = Vector3(0, 0, 0); // Stop on collision
        }

        // Ground check
        if (position.y <= groundY) {
            position.y = groundY;
            velocity.y = 0;
            isGrounded = true;
        }
    }

    Vector3 GetPosition() const { return position; }
    Vector3 GetVelocity() const { return velocity; }
    bool IsGrounded() const { return isGrounded; }

    void AddObstacle(const Vector3& obstacle) {
        obstacles.push_back(obstacle);
    }
};
