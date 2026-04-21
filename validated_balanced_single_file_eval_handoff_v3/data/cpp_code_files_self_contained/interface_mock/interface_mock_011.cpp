#include <vector>
#include <cmath>
#include <algorithm>
#include <stdexcept>
#include <map>

using namespace std;

struct Vector2 {
    float x, y;
    
    Vector2(float x = 0, float y = 0) : x(x), y(y) {}
    
    Vector2 operator+(const Vector2& other) const {
        return Vector2(x + other.x, y + other.y);
    }
    
    Vector2 operator-(const Vector2& other) const {
        return Vector2(x - other.x, y - other.y);
    }
    
    float magnitude() const {
        return sqrt(x*x + y*y);
    }
    
    Vector2 normalized() const {
        float mag = magnitude();
        if (mag > 0) {
            return Vector2(x/mag, y/mag);
        }
        return Vector2(0, 0);
    }
};

enum class ColliderType {
    CIRCLE,
    RECTANGLE,
    POLYGON
};

class Collider {
private:
    ColliderType _type;
    Vector2 _position;
    Vector2 _size; // For rectangle
    float _radius; // For circle
    vector<Vector2> _vertices; // For polygon
    
public:
    Collider(ColliderType type, Vector2 position = Vector2(), 
             Vector2 size = Vector2(1,1), float radius = 1.0f, 
             vector<Vector2> vertices = {}) 
        : _type(type), _position(position), _size(size), 
          _radius(radius), _vertices(vertices) {}
    
    void SetPosition(Vector2 position) {
        _position = position;
    }
    
    Vector2 GetPosition() const {
        return _position;
    }
    
    Vector2 clamp(Vector2 position, Vector2 min, Vector2 max) const {
        Vector2 result;
        result.x = std::fmax(min.x, std::fmin(position.x, max.x));
        result.y = std::fmax(min.y, std::fmin(position.y, max.y));
        return result;
    }
    
    bool CheckCollision(const Collider& other) const {
        if (_type == ColliderType::CIRCLE && other._type == ColliderType::CIRCLE) {
            return CircleCircleCollision(other);
        } else if (_type == ColliderType::RECTANGLE && other._type == ColliderType::RECTANGLE) {
            return RectRectCollision(other);
        } else if (_type == ColliderType::CIRCLE && other._type == ColliderType::RECTANGLE) {
            return CircleRectCollision(other);
        } else if (_type == ColliderType::RECTANGLE && other._type == ColliderType::CIRCLE) {
            return other.CircleRectCollision(*this);
        }
        return false;
    }
    
    Vector2 GetCollisionNormal(const Collider& other) const {
        if (_type == ColliderType::CIRCLE && other._type == ColliderType::CIRCLE) {
            Vector2 direction = other._position - _position;
            return direction.normalized();
        } else if (_type == ColliderType::RECTANGLE && other._type == ColliderType::CIRCLE) {
            return GetRectCircleCollisionNormal(other);
        }
        return Vector2(0, 0);
    }
    
private:
    bool CircleCircleCollision(const Collider& other) const {
        float distance = (_position - other._position).magnitude();
        return distance < (_radius + other._radius);
    }
    
    bool RectRectCollision(const Collider& other) const {
        bool xOverlap = abs(_position.x - other._position.x) < (_size.x + other._size.x)/2;
        bool yOverlap = abs(_position.y - other._position.y) < (_size.y + other._size.y)/2;
        return xOverlap && yOverlap;
    }
    
    bool CircleRectCollision(const Collider& rect) const {
        Vector2 closestPoint = clamp(_position, 
                                    Vector2(rect._position.x - rect._size.x/2, rect._position.y - rect._size.y/2),
                                    Vector2(rect._position.x + rect._size.x/2, rect._position.y + rect._size.y/2));
        float distance = (_position - closestPoint).magnitude();
        return distance < _radius;
    }
    
    Vector2 GetRectCircleCollisionNormal(const Collider& circle) const {
        Vector2 closestPoint = clamp(circle._position, 
                                    Vector2(_position.x - _size.x/2, _position.y - _size.y/2),
                                    Vector2(_position.x + _size.x/2, _position.y + _size.y/2));
        Vector2 normal = (circle._position - closestPoint).normalized();
        return normal;
    }
};
