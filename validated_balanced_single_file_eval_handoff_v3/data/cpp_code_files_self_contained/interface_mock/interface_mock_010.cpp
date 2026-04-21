#include <vector>
#include <cmath>
#include <map>
#include <stdexcept>

using namespace std;

enum class PlatformType { Thick, Thin };
enum class CollisionType { None, Good, Bad };

struct Vector2f {
    float x, y;
    Vector2f(float x = 0, float y = 0) : x(x), y(y) {}
    Vector2f operator+(const Vector2f& other) const { return Vector2f(x + other.x, y + other.y); }
    Vector2f operator-(const Vector2f& other) const { return Vector2f(x - other.x, y - other.y); }
    Vector2f operator/(float divisor) const { return Vector2f(x / divisor, y / divisor); }
    float length() const { return sqrt(x*x + y*y); }
};

class Player {
public:
    Vector2f position;
    Vector2f size;
    bool isGrounded;

    Player(Vector2f pos = Vector2f(), Vector2f sz = Vector2f(20, 40)) 
        : position(pos), size(sz), isGrounded(false) {}

    Vector2f getPosition() const { return position; }
    Vector2f getSize() const { return size; }
    void move(Vector2f offset) { position = position + offset; }
    void land() { isGrounded = true; }
};

class Platform {
    Vector2f size;
    Vector2f position;
    PlatformType type;
    const float COLLISION_TOLERANCE = 5.f;

public:
    Platform(PlatformType t, Vector2f pos) : type(t), position(pos) {
        size = getPlatformSize(type);
    }

    CollisionType collision(Player& player) const {
        const Vector2f playerPos = player.getPosition();
        const Vector2f playerHalf = player.getSize() / 2.f;
        const Vector2f thisHalf = size / 2.f;

        const Vector2f delta = playerPos - position;
        const float intersectX = abs(delta.x) - (playerHalf.x + thisHalf.x);
        const float intersectY = abs(delta.y) - (playerHalf.y + thisHalf.y);

        if (intersectX < 0.f && intersectY < 0.f) {
            if (intersectX > intersectY && delta.y < 0.f && intersectY < -COLLISION_TOLERANCE) {
                player.move(Vector2f(0.f, intersectY));
                player.land();
                return CollisionType::Good;
            }
            return CollisionType::Bad;
        }
        return CollisionType::None;
    }

    Vector2f getPlatformSize(PlatformType type) const {
        static const map<PlatformType, Vector2f> sizes = {
            {PlatformType::Thick, Vector2f(30.f, 30.f)},
            {PlatformType::Thin, Vector2f(30.f, 5.f)}
        };
        return sizes.at(type);
    }

    Vector2f getPosition() const { return position; }
    Vector2f getSize() const { return size; }
    PlatformType getType() const { return type; }
};

class PhysicsEngine {
    vector<Platform> platforms;
    Player player;
    float gravity = 9.8f;

public:
    PhysicsEngine(const vector<Platform>& plats, const Player& pl) 
        : platforms(plats), player(pl) {}

    void update(float deltaTime) {
        // Apply gravity if not grounded
        if (!player.isGrounded) {
            player.move(Vector2f(0, gravity * deltaTime));
        }

        // Reset grounded state before collision checks
        player.isGrounded = false;

        // Check collisions with all platforms
        for (const auto& platform : platforms) {
            platform.collision(player);
        }
    }

    Player getPlayer() const { return player; }
};
