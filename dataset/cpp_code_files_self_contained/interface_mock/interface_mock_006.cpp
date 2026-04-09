#include <vector>
#include <cmath>
#include <algorithm>
#include <limits>

using namespace std;

// Enhanced Vector3 class with more operations
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

    float Dot(const Vector3& other) const {
        return x * other.x + y * other.y + z * other.z;
    }

    Vector3 Cross(const Vector3& other) const {
        return Vector3(
            y * other.z - z * other.y,
            z * other.x - x * other.z,
            x * other.y - y * other.x
        );
    }

    float Magnitude() const {
        return sqrt(x*x + y*y + z*z);
    }

    Vector3 Normalized() const {
        float mag = Magnitude();
        if (mag > 0) return *this * (1.0f / mag);
        return *this;
    }
};

// Collision Shape Base Class
class CollisionShape {
public:
    Vector3 pos;
    Vector3 velocity;
    float restitution = 0.8f; // Bounciness factor

    virtual ~CollisionShape() = default;
    virtual Vector3 GetFarthestPointInDirection(const Vector3& direction) const = 0;
};

// Sphere Collision Shape
class CollisionSphere : public CollisionShape {
public:
    float radius;

    CollisionSphere(const Vector3& p, float r) : radius(r) {
        pos = p;
    }

    Vector3 GetFarthestPointInDirection(const Vector3& direction) const override {
        return pos + direction.Normalized() * radius;
    }
};

// Axis-Aligned Bounding Box
class CollisionAABB : public CollisionShape {
public:
    Vector3 halfdims;

    CollisionAABB(const Vector3& p, const Vector3& hd) : halfdims(hd) {
        pos = p;
    }

    Vector3 GetFarthestPointInDirection(const Vector3& direction) const override {
        return Vector3(
            pos.x + (direction.x > 0 ? halfdims.x : -halfdims.x),
            pos.y + (direction.y > 0 ? halfdims.y : -halfdims.y),
            pos.z + (direction.z > 0 ? halfdims.z : -halfdims.z)
        );
    }
};

// Collision Result
class CollisionData {
public:
    Vector3 point;
    Vector3 normal;
    float penetration;
    CollisionShape* shapeA;
    CollisionShape* shapeB;

    void ResolveCollision() const {
        // Calculate relative velocity
        Vector3 relativeVel = shapeB->velocity - shapeA->velocity;
        float velAlongNormal = relativeVel.Dot(normal);

        // Do not resolve if objects are separating
        if (velAlongNormal > 0) return;

        // Calculate impulse scalar
        float e = min(shapeA->restitution, shapeB->restitution);
        float j = -(1 + e) * velAlongNormal;
        
        // Apply impulse
        Vector3 impulse = normal * j;
        shapeA->velocity = shapeA->velocity - impulse * 0.5f;
        shapeB->velocity = shapeB->velocity + impulse * 0.5f;

        // Positional correction to prevent sinking
        const float percent = 0.2f;
        const float slop = 0.01f;
        Vector3 correction = normal * (max(penetration - slop, 0.0f) / 
                             (1.0f / shapeA->restitution + 1.0f / shapeB->restitution)) * percent;
        shapeA->pos = shapeA->pos - correction * (1.0f / shapeA->restitution);
        shapeB->pos = shapeB->pos + correction * (1.0f / shapeB->restitution);
    }
};

// Collision Detection Functions
bool SphereSphereCollision(const CollisionSphere& a, const CollisionSphere& b, CollisionData& data) {
    Vector3 diff = b.pos - a.pos;
    float dist = diff.Magnitude();
    float radiusSum = a.radius + b.radius;

    if (dist < radiusSum) {
        data.penetration = radiusSum - dist;
        data.normal = diff.Normalized();
        data.point = a.pos + data.normal * a.radius;
        return true;
    }
    return false;
}

bool SphereAABBCollision(const CollisionSphere& sphere, const CollisionAABB& aabb, CollisionData& data) {
    Vector3 closestPoint = sphere.pos;
    closestPoint.x = max(aabb.pos.x - aabb.halfdims.x, min(sphere.pos.x, aabb.pos.x + aabb.halfdims.x));
    closestPoint.y = max(aabb.pos.y - aabb.halfdims.y, min(sphere.pos.y, aabb.pos.y + aabb.halfdims.y));
    closestPoint.z = max(aabb.pos.z - aabb.halfdims.z, min(sphere.pos.z, aabb.pos.z + aabb.halfdims.z));

    Vector3 diff = closestPoint - sphere.pos;
    float dist = diff.Magnitude();

    if (dist < sphere.radius) {
        data.penetration = sphere.radius - dist;
        data.normal = diff.Normalized();
        data.point = closestPoint;
        return true;
    }
    return false;
}

// GJK (Gilbert-Johnson-Keerthi) algorithm for convex shapes
bool GJKCollision(const CollisionShape& a, const CollisionShape& b, CollisionData& data) {
    Vector3 direction(1, 0, 0); // Initial search direction
    Vector3 simplex[4];
    int simplexSize = 0;

    Vector3 supportA = a.GetFarthestPointInDirection(direction);
    Vector3 supportB = b.GetFarthestPointInDirection(direction * -1.0f);
    Vector3 support = supportA - supportB;

    simplex[0] = support;
    simplexSize = 1;

    direction = support * -1.0f;

    while (true) {
        supportA = a.GetFarthestPointInDirection(direction);
        supportB = b.GetFarthestPointInDirection(direction * -1.0f);
        support = supportA - supportB;

        if (support.Dot(direction) <= 0) {
            return false; // No collision
        }

        simplex[simplexSize] = support;
        simplexSize++;

        // Handle simplex
        if (simplexSize == 2) {
            // Line case
            Vector3 a = simplex[1];
            Vector3 b = simplex[0];
            Vector3 ab = b - a;
            Vector3 ao = a * -1.0f;

            if (ab.Dot(ao) >= 0) {
                direction = ab.Cross(ao).Cross(ab);
            } else {
                simplex[0] = a;
                simplexSize = 1;
                direction = ao;
            }
        } else if (simplexSize == 3) {
            // Triangle case
            Vector3 a = simplex[2];
            Vector3 b = simplex[1];
            Vector3 c = simplex[0];
            Vector3 ab = b - a;
            Vector3 ac = c - a;
            Vector3 ao = a * -1.0f;
            Vector3 abc = ab.Cross(ac);

            if (abc.Cross(ac).Dot(ao) >= 0) {
                if (ac.Dot(ao) >= 0) {
                    simplex[0] = a;
                    simplex[1] = c;
                    simplexSize = 2;
                    direction = ac.Cross(ao).Cross(ac);
                } else {
                    goto checkAB;
                }
            } else {
                checkAB:
                if (ab.Cross(abc).Dot(ao) >= 0) {
                    simplex[0] = a;
                    simplex[1] = b;
                    simplexSize = 2;
                    direction = ab.Cross(ao).Cross(ab);
                } else {
                    if (abc.Dot(ao) >= 0) {
                        direction = abc;
                    } else {
                        simplex[0] = a;
                        simplex[1] = b;
                        simplex[2] = c;
                        simplexSize = 3;
                        direction = abc * -1.0f;
                    }
                }
            }
        } else {
            // Tetrahedron case
            Vector3 a = simplex[3];
            Vector3 b = simplex[2];
            Vector3 c = simplex[1];
            Vector3 d = simplex[0];
            Vector3 ab = b - a;
            Vector3 ac = c - a;
            Vector3 ad = d - a;
            Vector3 ao = a * -1.0f;
            Vector3 abc = ab.Cross(ac);
            Vector3 acd = ac.Cross(ad);
            Vector3 adb = ad.Cross(ab);

            if (abc.Dot(ao) >= 0) {
                simplex[0] = a;
                simplex[1] = b;
                simplex[2] = c;
                simplexSize = 3;
                direction = abc;
                continue;
            }

            if (acd.Dot(ao) >= 0) {
                simplex[0] = a;
                simplex[1] = c;
                simplex[2] = d;
                simplexSize = 3;
                direction = acd;
                continue;
            }

            if (adb.Dot(ao) >= 0) {
                simplex[0] = a;
                simplex[1] = d;
                simplex[2] = b;
                simplexSize = 3;
                direction = adb;
                continue;
            }

            // Origin is inside the tetrahedron - collision!
            return true;
        }
    }
}
