#include <vector>
#include <cmath>
#include <memory>
#include <unordered_map>

// Simple 2D vector class to replace the original Vector2
class Vector2 {
public:
    float x, y;

    Vector2(float x = 0.0f, float y = 0.0f) : x(x), y(y) {}

    Vector2 operator+(const Vector2& other) const {
        return Vector2(x + other.x, y + other.y);
    }

    Vector2 operator-(const Vector2& other) const {
        return Vector2(x - other.x, y - other.y);
    }

    Vector2 operator*(float scalar) const {
        return Vector2(x * scalar, y * scalar);
    }

    float magnitude() const {
        return std::sqrt(x * x + y * y);
    }

    Vector2 normalized() const {
        float mag = magnitude();
        if (mag > 0) {
            return Vector2(x / mag, y / mag);
        }
        return Vector2(0, 0);
    }

    static float distance(const Vector2& a, const Vector2& b) {
        return (a - b).magnitude();
    }
};

// Simplified Blackboard class
class Blackboard {
private:
    std::unordered_map<std::string, float> floatValues;
    std::unordered_map<std::string, Vector2> vectorValues;
    std::unordered_map<std::string, bool> boolValues;

public:
    void setFloat(const std::string& key, float value) { floatValues[key] = value; }
    float getFloat(const std::string& key, float defaultValue = 0.0f) const {
        auto it = floatValues.find(key);
        return it != floatValues.end() ? it->second : defaultValue;
    }

    void setVector(const std::string& key, const Vector2& value) { vectorValues[key] = value; }
    Vector2 getVector(const std::string& key, const Vector2& defaultValue = Vector2()) const {
        auto it = vectorValues.find(key);
        return it != vectorValues.end() ? it->second : defaultValue;
    }

    void setBool(const std::string& key, bool value) { boolValues[key] = value; }
    bool getBool(const std::string& key, bool defaultValue = false) const {
        auto it = boolValues.find(key);
        return it != boolValues.end() ? it->second : defaultValue;
    }
};

// Forward declaration
class SteeringBehaviourManager;

// Enhanced SteeringBehaviour class with additional functionality
class SteeringBehaviour {
protected:
    Blackboard* blackboard = nullptr;
    float weight = 1.0f;
    bool enabled = true;

public:
    SteeringBehaviourManager* sbm = nullptr;

    SteeringBehaviour() = default;
    virtual ~SteeringBehaviour() = default;

    // Core update function
    virtual Vector2 update() = 0;

    // Initialization function
    virtual void initialise(Blackboard* bb) {
        blackboard = bb;
    }

    // New: Enable/disable behavior
    virtual void setEnabled(bool state) {
        enabled = state;
    }

    // New: Get behavior weight
    virtual float getWeight() const {
        return weight;
    }

    // New: Set behavior weight
    virtual void setWeight(float newWeight) {
        weight = newWeight;
    }

    // New: Get behavior name (for debugging)
    virtual std::string getName() const {
        return "BaseSteeringBehavior";
    }

    // New: Reset behavior state
    virtual void reset() {}
};

// Concrete implementation of a seek behavior
class SeekBehavior : public SteeringBehaviour {
private:
    Vector2 target;

public:
    SeekBehavior(const Vector2& targetPos = Vector2()) : target(targetPos) {}

    Vector2 update() override {
        if (!enabled) return Vector2();

        Vector2 currentPos = blackboard->getVector("position");
        Vector2 desiredVelocity = (target - currentPos).normalized() * blackboard->getFloat("max_speed", 5.0f);
        Vector2 currentVelocity = blackboard->getVector("velocity");
        Vector2 steering = (desiredVelocity - currentVelocity) * weight;

        return steering;
    }

    void setTarget(const Vector2& newTarget) {
        target = newTarget;
    }

    std::string getName() const override {
        return "SeekBehavior";
    }
};

// Concrete implementation of a flee behavior
class FleeBehavior : public SteeringBehaviour {
private:
    Vector2 threatPosition;
    float panicDistance = 10.0f;

public:
    FleeBehavior(const Vector2& threatPos = Vector2(), float distance = 10.0f) 
        : threatPosition(threatPos), panicDistance(distance) {}

    Vector2 update() override {
        if (!enabled) return Vector2();

        Vector2 currentPos = blackboard->getVector("position");
        float distance = Vector2::distance(currentPos, threatPosition);

        if (distance > panicDistance) return Vector2();

        Vector2 desiredVelocity = (currentPos - threatPosition).normalized() * blackboard->getFloat("max_speed", 5.0f);
        Vector2 currentVelocity = blackboard->getVector("velocity");
        Vector2 steering = (desiredVelocity - currentVelocity) * weight;

        return steering;
    }

    void setThreat(const Vector2& newThreat) {
        threatPosition = newThreat;
    }

    std::string getName() const override {
        return "FleeBehavior";
    }
};
