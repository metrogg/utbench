#include <cmath>
#include <vector>
#include <iomanip>
#include <stdexcept>

using namespace std;

class Vector3 {
public:
    float x, y, z;

    Vector3() : x(0), y(0), z(0) {}
    Vector3(float x, float y, float z) : x(x), y(y), z(z) {}

    float length() const {
        return sqrt(x*x + y*y + z*z);
    }

    Vector3 normalize() const {
        float len = length();
        if (len == 0) return Vector3();
        return Vector3(x/len, y/len, z/len);
    }

    Vector3 operator+(const Vector3& other) const {
        return Vector3(x + other.x, y + other.y, z + other.z);
    }

    Vector3 operator-(const Vector3& other) const {
        return Vector3(x - other.x, y - other.y, z - other.z);
    }

    Vector3 operator*(float scalar) const {
        return Vector3(x * scalar, y * scalar, z * scalar);
    }

    float dot(const Vector3& other) const {
        return x*other.x + y*other.y + z*other.z;
    }

    Vector3 cross(const Vector3& other) const {
        return Vector3(
            y*other.z - z*other.y,
            z*other.x - x*other.z,
            x*other.y - y*other.x
        );
    }
};

class LineSegment3D {
private:
    Vector3 point1, point2;

public:
    LineSegment3D() : point1(), point2() {}
    LineSegment3D(const Vector3& p1, const Vector3& p2) : point1(p1), point2(p2) {}
    LineSegment3D(float p1x, float p1y, float p1z, float p2x, float p2y, float p2z) 
        : point1(p1x, p1y, p1z), point2(p2x, p2y, p2z) {}

    float length() const {
        Vector3 diff = point2 - point1;
        return diff.length();
    }

    Vector3 direction() const {
        return (point2 - point1).normalize();
    }

    Vector3 midpoint() const {
        return Vector3(
            (point1.x + point2.x) / 2,
            (point1.y + point2.y) / 2,
            (point1.z + point2.z) / 2
        );
    }

    bool contains_point(const Vector3& point, float tolerance = 0.0001f) const {
        Vector3 line_vec = point2 - point1;
        Vector3 point_vec = point - point1;
        
        // Check if point is on the line
        Vector3 cross_product = line_vec.cross(point_vec);
        if (cross_product.length() > tolerance) {
            return false;
        }
        
        // Check if point is between the endpoints
        float dot_product = line_vec.dot(point_vec);
        if (dot_product < -tolerance || dot_product > line_vec.dot(line_vec) + tolerance) {
            return false;
        }
        
        return true;
    }

    float distance_to_point(const Vector3& point) const {
        Vector3 line_vec = point2 - point1;
        Vector3 point_vec = point - point1;
        
        float line_length = line_vec.length();
        if (line_length == 0) return (point - point1).length();
        
        Vector3 line_dir = line_vec.normalize();
        float projection_length = point_vec.dot(line_dir);
        
        if (projection_length < 0) return (point - point1).length();
        if (projection_length > line_length) return (point - point2).length();
        
        Vector3 projection = point1 + line_dir * projection_length;
        return (point - projection).length();
    }

    pair<bool, Vector3> find_intersection(const LineSegment3D& other, float tolerance = 0.0001f) const {
        Vector3 p1 = point1, p2 = point2;
        Vector3 p3 = other.point1, p4 = other.point2;
        
        Vector3 p13 = p1 - p3;
        Vector3 p43 = p4 - p3;
        Vector3 p21 = p2 - p1;
        
        if (abs(p43.x) < tolerance && abs(p43.y) < tolerance && abs(p43.z) < tolerance) {
            return {false, Vector3()}; // Second segment is a point
        }
        
        if (abs(p21.x) < tolerance && abs(p21.y) < tolerance && abs(p21.z) < tolerance) {
            return {false, Vector3()}; // First segment is a point
        }
        
        float d1343 = p13.dot(p43);
        float d4321 = p43.dot(p21);
        float d1321 = p13.dot(p21);
        float d4343 = p43.dot(p43);
        float d2121 = p21.dot(p21);
        
        float denom = d2121 * d4343 - d4321 * d4321;
        if (abs(denom) < tolerance) {
            return {false, Vector3()}; // Lines are parallel
        }
        
        float numer = d1343 * d4321 - d1321 * d4343;
        float mua = numer / denom;
        float mub = (d1343 + d4321 * mua) / d4343;
        
        Vector3 intersection = p1 + p21 * mua;
        
        // Check if intersection point is within both segments
        bool within_segment1 = (mua >= -tolerance && mua <= 1.0f + tolerance);
        bool within_segment2 = (mub >= -tolerance && mub <= 1.0f + tolerance);
        
        return {within_segment1 && within_segment2, intersection};
    }

    Vector3 get_point1() const { return point1; }
    Vector3 get_point2() const { return point2; }
};
