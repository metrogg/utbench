#include <string>
#include <cmath>
#include <vector>
#include <stdexcept>

using namespace std;

class Vector {
public:
    float x, y, z;

    Vector(float x=0, float y=0, float z=0) : x(x), y(y), z(z) {}

    Vector operator+(const Vector& other) const {
        return Vector(x + other.x, y + other.y, z + other.z);
    }

    Vector operator-(const Vector& other) const {
        return Vector(x - other.x, y - other.y, z - other.z);
    }

    Vector operator*(float scalar) const {
        return Vector(x * scalar, y * scalar, z * scalar);
    }

    float dot(const Vector& other) const {
        return x * other.x + y * other.y + z * other.z;
    }

    Vector cross(const Vector& other) const {
        return Vector(
            y * other.z - z * other.y,
            z * other.x - x * other.z,
            x * other.y - y * other.x
        );
    }

    float magnitude() const {
        return sqrt(x*x + y*y + z*z);
    }

    Vector normalized() const {
        float mag = magnitude();
        if (mag == 0) return Vector();
        return Vector(x/mag, y/mag, z/mag);
    }

    void print() const {
        cout << "(" << x << ", " << y << ", " << z << ")";
    }
};

class AABB {
public:
    Vector min, max;

    AABB(Vector min=Vector(), Vector max=Vector()) : min(min), max(max) {}

    bool contains(const Vector& point) const {
        return point.x >= min.x && point.x <= max.x &&
               point.y >= min.y && point.y <= max.y &&
               point.z >= min.z && point.z <= max.z;
    }

    void print() const {
        cout << "AABB[Min: ";
        min.print();
        cout << ", Max: ";
        max.print();
        cout << "]";
    }
};

class Camera {
private:
    string name;
    Vector pos;
    Vector rotator; // in degrees
    Vector forwardVector;
    Vector upVector;
    Vector rightVector;
    
    float fov;      // in degrees
    float zFar;
    float zNear;
    float width;
    float height;

    void updateVectors() {
        // Convert rotation from degrees to radians
        float pitch = rotator.x * M_PI / 180.0f;
        float yaw = rotator.y * M_PI / 180.0f;
        float roll = rotator.z * M_PI / 180.0f;

        // Calculate forward vector
        forwardVector = Vector(
            cos(yaw) * cos(pitch),
            sin(pitch),
            sin(yaw) * cos(pitch)
        ).normalized();

        // Calculate right vector
        Vector worldUp(0, 1, 0);
        rightVector = forwardVector.cross(worldUp).normalized();

        // Calculate up vector
        upVector = rightVector.cross(forwardVector).normalized();

        // Apply roll if needed
        if (roll != 0) {
            // Rotate up and right vectors around forward axis
            float cosRoll = cos(roll);
            float sinRoll = sin(roll);
            
            Vector newRight = rightVector * cosRoll + upVector * sinRoll;
            upVector = rightVector * -sinRoll + upVector * cosRoll;
            rightVector = newRight;
        }
    }

public:
    Camera(const string& name="DefaultCamera", 
           const Vector& pos=Vector(), 
           const Vector& rotator=Vector(),
           float fov=60.0f, 
           float zFar=1000.0f, 
           float zNear=0.1f,
           float width=800.0f, 
           float height=600.0f) {
        Init(name, pos, rotator, fov, zFar, zNear, width, height);
    }

    void Init(const string& srcName, 
              const Vector& srcPos, 
              const Vector& srcRotator,
              float srcFov, 
              float srcZFar, 
              float srcZNear,
              float srcWidth, 
              float srcHeight) {
        name = srcName;
        pos = srcPos;
        rotator = srcRotator;
        fov = srcFov;
        zFar = srcZFar;
        zNear = srcZNear;
        width = srcWidth;
        height = srcHeight;
        updateVectors();
    }

    void SetPos(const Vector& src) {
        pos = src;
    }

    void SetRotator(const Vector& src) {
        rotator = src;
        updateVectors();
    }

    Vector GetRightSideVector() const { return rightVector; }
    Vector GetLeftSideVector() const { return rightVector * -1.0f; }
    Vector GetForwardVector() const { return forwardVector; }
    Vector GetUpVector() const { return upVector; }
    Vector GetRotator() const { return rotator; }
    Vector GetPos() const { return pos; }

    AABB GetViewAABB() const {
        // Calculate view frustum corners in world space
        float aspect = width / height;
        float tanFov = tan(fov * 0.5f * M_PI / 180.0f);
        
        Vector nearCenter = pos + forwardVector * zNear;
        Vector farCenter = pos + forwardVector * zFar;

        // Near plane dimensions
        float nearHeight = 2.0f * tanFov * zNear;
        float nearWidth = nearHeight * aspect;

        // Far plane dimensions
        float farHeight = 2.0f * tanFov * zFar;
        float farWidth = farHeight * aspect;

        // Calculate corners
        Vector nearTopLeft = nearCenter + upVector * (nearHeight * 0.5f) - rightVector * (nearWidth * 0.5f);
        Vector nearTopRight = nearCenter + upVector * (nearHeight * 0.5f) + rightVector * (nearWidth * 0.5f);
        Vector nearBottomLeft = nearCenter - upVector * (nearHeight * 0.5f) - rightVector * (nearWidth * 0.5f);
        Vector nearBottomRight = nearCenter - upVector * (nearHeight * 0.5f) + rightVector * (nearWidth * 0.5f);

        Vector farTopLeft = farCenter + upVector * (farHeight * 0.5f) - rightVector * (farWidth * 0.5f);
        Vector farTopRight = farCenter + upVector * (farHeight * 0.5f) + rightVector * (farWidth * 0.5f);
        Vector farBottomLeft = farCenter - upVector * (farHeight * 0.5f) - rightVector * (farWidth * 0.5f);
        Vector farBottomRight = farCenter - upVector * (farHeight * 0.5f) + rightVector * (farWidth * 0.5f);

        // Find min and max coordinates for AABB
        Vector min = pos;
        Vector max = pos;

        auto updateMinMax = [&](const Vector& point) {
            min.x = fmin(min.x, point.x);
            min.y = fmin(min.y, point.y);
            min.z = fmin(min.z, point.z);
            max.x = fmax(max.x, point.x);
            max.y = fmax(max.y, point.y);
            max.z = fmax(max.z, point.z);
        };

        updateMinMax(nearTopLeft);
        updateMinMax(nearTopRight);
        updateMinMax(nearBottomLeft);
        updateMinMax(nearBottomRight);
        updateMinMax(farTopLeft);
        updateMinMax(farTopRight);
        updateMinMax(farBottomLeft);
        updateMinMax(farBottomRight);

        return AABB(min, max);
    }

    bool IsAABBInside(const AABB& aabb) const {
        // Simple AABB vs frustum check (simplified for demo)
        AABB viewAABB = GetViewAABB();
        return !(aabb.max.x < viewAABB.min.x || aabb.min.x > viewAABB.max.x ||
                 aabb.max.y < viewAABB.min.y || aabb.min.y > viewAABB.max.y ||
                 aabb.max.z < viewAABB.min.z || aabb.min.z > viewAABB.max.z);
    }

    void printInfo() const {
        cout << "Camera: " << name << endl;
        cout << "Position: "; pos.print(); cout << endl;
        cout << "Rotation: "; rotator.print(); cout << endl;
        cout << "FOV: " << fov << " degrees" << endl;
        cout << "Clipping: Near=" << zNear << ", Far=" << zFar << endl;
        cout << "Resolution: " << width << "x" << height << endl;
        cout << "Forward: "; forwardVector.print(); cout << endl;
        cout << "Up: "; upVector.print(); cout << endl;
        cout << "Right: "; rightVector.print(); cout << endl;
    }
};
