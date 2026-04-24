#include <vector>
#include <cmath>
#include <stdexcept>

using namespace std;

// Simplified Vector4 class for 3D coordinates
class Vector4 {
public:
    float x, y, z, w;

    Vector4() : x(0), y(0), z(0), w(1) {}
    Vector4(float x, float y, float z, float w = 1) : x(x), y(y), z(z), w(w) {}

    Vector4 operator+(const Vector4& other) const {
        return Vector4(x + other.x, y + other.y, z + other.z, w + other.w);
    }

    Vector4 operator-(const Vector4& other) const {
        return Vector4(x - other.x, y - other.y, z - other.z, w - other.w);
    }

    Vector4 operator*(float scalar) const {
        return Vector4(x * scalar, y * scalar, z * scalar, w * scalar);
    }

    float dot(const Vector4& other) const {
        return x * other.x + y * other.y + z * other.z + w * other.w;
    }

    Vector4 cross(const Vector4& other) const {
        return Vector4(
            y * other.z - z * other.y,
            z * other.x - x * other.z,
            x * other.y - y * other.x,
            0
        );
    }

    float magnitude() const {
        return sqrt(x*x + y*y + z*z);
    }

    Vector4 normalized() const {
        float mag = magnitude();
        if (mag == 0) return *this;
        return *this * (1.0f / mag);
    }
};

// Simplified Matrix4x4 class for transformations
class Matrix4x4 {
public:
    float m[4][4];

    Matrix4x4() {
        for (int i = 0; i < 4; ++i)
            for (int j = 0; j < 4; ++j)
                m[i][j] = (i == j) ? 1.0f : 0.0f;
    }

    static Matrix4x4 perspective(float fov, float aspect, float near, float far) {
        Matrix4x4 result;
        float tanHalfFov = tan(fov / 2);

        result.m[0][0] = 1.0f / (aspect * tanHalfFov);
        result.m[1][1] = 1.0f / tanHalfFov;
        result.m[2][2] = -(far + near) / (far - near);
        result.m[2][3] = -1.0f;
        result.m[3][2] = -(2.0f * far * near) / (far - near);
        result.m[3][3] = 0.0f;

        return result;
    }

    static Matrix4x4 translation(float x, float y, float z) {
        Matrix4x4 result;
        result.m[3][0] = x;
        result.m[3][1] = y;
        result.m[3][2] = z;
        return result;
    }

    static Matrix4x4 rotationX(float angle) {
        Matrix4x4 result;
        float c = cos(angle);
        float s = sin(angle);
        result.m[1][1] = c;
        result.m[1][2] = s;
        result.m[2][1] = -s;
        result.m[2][2] = c;
        return result;
    }

    static Matrix4x4 rotationY(float angle) {
        Matrix4x4 result;
        float c = cos(angle);
        float s = sin(angle);
        result.m[0][0] = c;
        result.m[0][2] = -s;
        result.m[2][0] = s;
        result.m[2][2] = c;
        return result;
    }

    static Matrix4x4 rotationZ(float angle) {
        Matrix4x4 result;
        float c = cos(angle);
        float s = sin(angle);
        result.m[0][0] = c;
        result.m[0][1] = s;
        result.m[1][0] = -s;
        result.m[1][1] = c;
        return result;
    }

    Vector4 operator*(const Vector4& v) const {
        Vector4 result;
        result.x = m[0][0] * v.x + m[1][0] * v.y + m[2][0] * v.z + m[3][0] * v.w;
        result.y = m[0][1] * v.x + m[1][1] * v.y + m[2][1] * v.z + m[3][1] * v.w;
        result.z = m[0][2] * v.x + m[1][2] * v.y + m[2][2] * v.z + m[3][2] * v.w;
        result.w = m[0][3] * v.x + m[1][3] * v.y + m[2][3] * v.z + m[3][3] * v.w;
        return result;
    }

    Matrix4x4 operator*(const Matrix4x4& other) const {
        Matrix4x4 result;
        for (int i = 0; i < 4; ++i) {
            for (int j = 0; j < 4; ++j) {
                result.m[i][j] = 0;
                for (int k = 0; k < 4; ++k) {
                    result.m[i][j] += m[i][k] * other.m[k][j];
                }
            }
        }
        return result;
    }
};

// Enhanced 3D Sprite class with additional functionality
class Sprite3D {
private:
    Vector4 vertices[4];  // Clockwise order: TL, BL, BR, TR
    Vector4 colors[4];
    Matrix4x4 transform;
    Matrix4x4 perspective;

public:
    Sprite3D() {
        // Default unit square
        vertices[0] = Vector4(-0.5f, 0.5f, 0);  // TL
        vertices[1] = Vector4(-0.5f, -0.5f, 0); // BL
        vertices[2] = Vector4(0.5f, -0.5f, 0);  // BR
        vertices[3] = Vector4(0.5f, 0.5f, 0);   // TR
        
        // Default white color
        for (int i = 0; i < 4; ++i) {
            colors[i] = Vector4(1, 1, 1, 1);
        }
        
        // Default perspective (45 deg FOV, 1:1 aspect, near=0.1, far=100)
        perspective = Matrix4x4::perspective(3.14159f/4, 1.0f, 0.1f, 100.0f);
    }

    Sprite3D(Vector4 v1, Vector4 v2, Vector4 v3, Vector4 v4) {
        vertices[0] = v1;
        vertices[1] = v2;
        vertices[2] = v3;
        vertices[3] = v4;
        
        for (int i = 0; i < 4; ++i) {
            colors[i] = Vector4(1, 1, 1, 1);
        }
        
        perspective = Matrix4x4::perspective(3.14159f/4, 1.0f, 0.1f, 100.0f);
    }

    void setColors(Vector4 c1, Vector4 c2, Vector4 c3, Vector4 c4) {
        colors[0] = c1;
        colors[1] = c2;
        colors[2] = c3;
        colors[3] = c4;
    }

    void moveTo(float x, float y, float z) {
        transform = Matrix4x4::translation(x, y, z);
    }

    void moveDelta(float dx, float dy, float dz) {
        transform = Matrix4x4::translation(dx, dy, dz) * transform;
    }

    void rotate(float xAngle, float yAngle, float zAngle) {
        Matrix4x4 rotX = Matrix4x4::rotationX(xAngle);
        Matrix4x4 rotY = Matrix4x4::rotationY(yAngle);
        Matrix4x4 rotZ = Matrix4x4::rotationZ(zAngle);
        transform = rotZ * rotY * rotX * transform;
    }

    void setPerspective(Matrix4x4 newPerspective) {
        perspective = newPerspective;
    }

    Vector4 getCenter() const {
        Vector4 center;
        for (int i = 0; i < 4; ++i) {
            center = center + vertices[i];
        }
        return center * 0.25f;
    }

    vector<Vector4> getTransformedVertices() const {
        vector<Vector4> result;
        for (int i = 0; i < 4; ++i) {
            Vector4 transformed = transform * vertices[i];
            result.push_back(transformed);
        }
        return result;
    }

    vector<Vector4> getScreenVertices() const {
        vector<Vector4> result;
        vector<Vector4> worldVerts = getTransformedVertices();
        
        for (const auto& vert : worldVerts) {
            Vector4 screenVert = perspective * vert;
            // Perspective divide
            if (screenVert.w != 0) {
                screenVert.x /= screenVert.w;
                screenVert.y /= screenVert.w;
                screenVert.z /= screenVert.w;
            }
            result.push_back(screenVert);
        }
        
        return result;
    }

    float calculateArea() const {
        vector<Vector4> verts = getTransformedVertices();
        Vector4 v1 = verts[1] - verts[0];
        Vector4 v2 = verts[3] - verts[0];
        Vector4 crossProd = v1.cross(v2);
        return crossProd.magnitude() / 2.0f;
    }

    bool isFacingCamera() const {
        vector<Vector4> verts = getTransformedVertices();
        Vector4 normal = (verts[1] - verts[0]).cross(verts[3] - verts[0]);
        Vector4 viewDir = Vector4(0, 0, -1, 0); // Looking down negative Z
        return normal.dot(viewDir) < 0;
    }
};
