#include <vector>
#include <cmath>
#include <algorithm>
#include <limits>

using namespace std;

class Vector3D {
private:
    double x, y, z;

public:
    Vector3D() : x(0), y(0), z(0) {}
    Vector3D(double x, double y, double z) : x(x), y(y), z(z) {}

    double getX() const { return x; }
    double getY() const { return y; }
    double getZ() const { return z; }

    Vector3D operator+(const Vector3D& other) const {
        return Vector3D(x + other.x, y + other.y, z + other.z);
    }

    Vector3D operator-(const Vector3D& other) const {
        return Vector3D(x - other.x, y - other.y, z - other.z);
    }

    Vector3D operator*(double scalar) const {
        return Vector3D(x * scalar, y * scalar, z * scalar);
    }

    double dot(const Vector3D& other) const {
        return x * other.x + y * other.y + z * other.z;
    }

    Vector3D cross(const Vector3D& other) const {
        return Vector3D(
            y * other.z - z * other.y,
            z * other.x - x * other.z,
            x * other.y - y * other.x
        );
    }

    double length() const {
        return sqrt(x * x + y * y + z * z);
    }

    Vector3D normalize() const {
        double len = length();
        if (len > 0) {
            return Vector3D(x / len, y / len, z / len);
        }
        return *this;
    }
};

class Cube {
private:
    Vector3D minCorner;
    Vector3D maxCorner;

    struct IntersectionResult {
        bool intersects;
        double distance;
        Vector3D point;
        Vector3D normal;
    };

public:
    Cube() {}
    Cube(const Vector3D& minCorner, const Vector3D& maxCorner)
        : minCorner(minCorner), maxCorner(maxCorner) {}

    // Enhanced ray-cube intersection with normal calculation
    IntersectionResult intersect(const Vector3D& rayOrigin, const Vector3D& rayDirection) const {
        IntersectionResult result;
        result.intersects = false;
        result.distance = numeric_limits<double>::max();

        Vector3D invDir(
            1.0 / rayDirection.getX(),
            1.0 / rayDirection.getY(),
            1.0 / rayDirection.getZ()
        );

        double t1 = (minCorner.getX() - rayOrigin.getX()) * invDir.getX();
        double t2 = (maxCorner.getX() - rayOrigin.getX()) * invDir.getX();
        double t3 = (minCorner.getY() - rayOrigin.getY()) * invDir.getY();
        double t4 = (maxCorner.getY() - rayOrigin.getY()) * invDir.getY();
        double t5 = (minCorner.getZ() - rayOrigin.getZ()) * invDir.getZ();
        double t6 = (maxCorner.getZ() - rayOrigin.getZ()) * invDir.getZ();

        double tmin = max(max(min(t1, t2), min(t3, t4)), min(t5, t6));
        double tmax = min(min(max(t1, t2), max(t3, t4)), max(t5, t6));

        // If tmax < 0, ray is intersecting but the cube is behind the origin
        if (tmax < 0) {
            return result;
        }

        // If tmin > tmax, ray doesn't intersect the cube
        if (tmin > tmax) {
            return result;
        }

        result.intersects = true;
        result.distance = tmin;
        result.point = rayOrigin + rayDirection * tmin;

        // Calculate normal at intersection point
        Vector3D center = (minCorner + maxCorner) * 0.5;
        Vector3D pointToCenter = result.point - center;
        Vector3D halfSize = (maxCorner - minCorner) * 0.5;

        // Find which face was hit by comparing distances
        double minDist = numeric_limits<double>::max();
        Vector3D normal;

        // Check X faces
        double dist = abs(halfSize.getX() - abs(pointToCenter.getX()));
        if (dist < minDist) {
            minDist = dist;
            normal = Vector3D(pointToCenter.getX() > 0 ? 1 : -1, 0, 0);
        }

        // Check Y faces
        dist = abs(halfSize.getY() - abs(pointToCenter.getY()));
        if (dist < minDist) {
            minDist = dist;
            normal = Vector3D(0, pointToCenter.getY() > 0 ? 1 : -1, 0);
        }

        // Check Z faces
        dist = abs(halfSize.getZ() - abs(pointToCenter.getZ()));
        if (dist < minDist) {
            normal = Vector3D(0, 0, pointToCenter.getZ() > 0 ? 1 : -1);
        }

        result.normal = normal.normalize();
        return result;
    }

    // Calculate volume of the cube
    double volume() const {
        Vector3D size = maxCorner - minCorner;
        return size.getX() * size.getY() * size.getZ();
    }

    // Check if a point is inside the cube
    bool containsPoint(const Vector3D& point) const {
        return (point.getX() >= minCorner.getX() && point.getX() <= maxCorner.getX() &&
                point.getY() >= minCorner.getY() && point.getY() <= maxCorner.getY() &&
                point.getZ() >= minCorner.getZ() && point.getZ() <= maxCorner.getZ());
    }
};
