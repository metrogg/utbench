#include <vector>
#include <cmath>
#include <cstring>
#include <stdexcept>

#define CAP_PI 3.14159265358979323846f

struct Vector3 {
    float x, y, z;
    Vector3(float x = 0, float y = 0, float z = 0) : x(x), y(y), z(z) {}
};

class AffineMatrix {
private:
    float m[16];

public:
    // Constructors
    AffineMatrix() { setIdentity(); }
    AffineMatrix(const AffineMatrix& other) { memcpy(m, other.m, sizeof(m)); }
    
    // Basic matrix operations
    void setIdentity() {
        memset(m, 0, sizeof(m));
        m[0] = m[5] = m[10] = m[15] = 1.0f;
    }
    
    void setZero() {
        memset(m, 0, sizeof(m));
        m[15] = 1.0f;
    }
    
    void setScale(float s) {
        setZero();
        m[0] = m[5] = m[10] = s;
    }
    
    void setScale(float sx, float sy, float sz) {
        setZero();
        m[0] = sx;
        m[5] = sy;
        m[10] = sz;
        m[15] = 1.0f;
    }
    
    void setTranslation(float x, float y, float z) {
        setIdentity();
        m[12] = x;
        m[13] = y;
        m[14] = z;
    }
    
    void setRotationX(float degrees) {
        float radians = degrees * CAP_PI / 180.0f;
        float c = cosf(radians);
        float s = sinf(radians);
        
        setIdentity();
        m[5] = c;  m[9] = -s;
        m[6] = s;  m[10] = c;
    }
    
    void setRotationY(float degrees) {
        float radians = degrees * CAP_PI / 180.0f;
        float c = cosf(radians);
        float s = sinf(radians);
        
        setIdentity();
        m[0] = c;  m[8] = s;
        m[2] = -s; m[10] = c;
    }
    
    void setRotationZ(float degrees) {
        float radians = degrees * CAP_PI / 180.0f;
        float c = cosf(radians);
        float s = sinf(radians);
        
        setIdentity();
        m[0] = c;  m[4] = -s;
        m[1] = s;  m[5] = c;
    }
    
    // Matrix operations
    void multiply(const AffineMatrix& other) {
        float result[16] = {0};
        
        for (int i = 0; i < 4; ++i) {
            for (int j = 0; j < 4; ++j) {
                for (int k = 0; k < 4; ++k) {
                    result[i*4 + j] += m[i*4 + k] * other.m[k*4 + j];
                }
            }
        }
        
        memcpy(m, result, sizeof(m));
    }
    
    bool invert() {
        float inv[16], det;
        
        inv[0] = m[5] * m[10] * m[15] - m[5] * m[11] * m[14] - m[9] * m[6] * m[15] 
               + m[9] * m[7] * m[14] + m[13] * m[6] * m[11] - m[13] * m[7] * m[10];
        
        inv[4] = -m[4] * m[10] * m[15] + m[4] * m[11] * m[14] + m[8] * m[6] * m[15] 
                - m[8] * m[7] * m[14] - m[12] * m[6] * m[11] + m[12] * m[7] * m[10];
        
        inv[8] = m[4] * m[9] * m[15] - m[4] * m[11] * m[13] - m[8] * m[5] * m[15] 
               + m[8] * m[7] * m[13] + m[12] * m[5] * m[11] - m[12] * m[7] * m[9];
        
        inv[12] = -m[4] * m[9] * m[14] + m[4] * m[10] * m[13] + m[8] * m[5] * m[14] 
                - m[8] * m[6] * m[13] - m[12] * m[5] * m[10] + m[12] * m[6] * m[9];
        
        inv[1] = -m[1] * m[10] * m[15] + m[1] * m[11] * m[14] + m[9] * m[2] * m[15] 
               - m[9] * m[3] * m[14] - m[13] * m[2] * m[11] + m[13] * m[3] * m[10];
        
        inv[5] = m[0] * m[10] * m[15] - m[0] * m[11] * m[14] - m[8] * m[2] * m[15] 
               + m[8] * m[3] * m[14] + m[12] * m[2] * m[11] - m[12] * m[3] * m[10];
        
        inv[9] = -m[0] * m[9] * m[15] + m[0] * m[11] * m[13] + m[8] * m[1] * m[15] 
               - m[8] * m[3] * m[13] - m[12] * m[1] * m[11] + m[12] * m[3] * m[9];
        
        inv[13] = m[0] * m[9] * m[14] - m[0] * m[10] * m[13] - m[8] * m[1] * m[14] 
                + m[8] * m[2] * m[13] + m[12] * m[1] * m[10] - m[12] * m[2] * m[9];
        
        inv[2] = m[1] * m[6] * m[15] - m[1] * m[7] * m[14] - m[5] * m[2] * m[15] 
               + m[5] * m[3] * m[14] + m[13] * m[2] * m[7] - m[13] * m[3] * m[6];
        
        inv[6] = -m[0] * m[6] * m[15] + m[0] * m[7] * m[14] + m[4] * m[2] * m[15] 
               - m[4] * m[3] * m[14] - m[12] * m[2] * m[7] + m[12] * m[3] * m[6];
        
        inv[10] = m[0] * m[5] * m[15] - m[0] * m[7] * m[13] - m[4] * m[1] * m[15] 
                + m[4] * m[3] * m[13] + m[12] * m[1] * m[7] - m[12] * m[3] * m[5];
        
        inv[14] = -m[0] * m[5] * m[14] + m[0] * m[6] * m[13] + m[4] * m[1] * m[14] 
                - m[4] * m[2] * m[13] - m[12] * m[1] * m[6] + m[12] * m[2] * m[5];
        
        inv[3] = -m[1] * m[6] * m[11] + m[1] * m[7] * m[10] + m[5] * m[2] * m[11] 
               - m[5] * m[3] * m[10] - m[9] * m[2] * m[7] + m[9] * m[3] * m[6];
        
        inv[7] = m[0] * m[6] * m[11] - m[0] * m[7] * m[10] - m[4] * m[2] * m[11] 
               + m[4] * m[3] * m[10] + m[8] * m[2] * m[7] - m[8] * m[3] * m[6];
        
        inv[11] = -m[0] * m[5] * m[11] + m[0] * m[7] * m[9] + m[4] * m[1] * m[11] 
                - m[4] * m[3] * m[9] - m[8] * m[1] * m[7] + m[8] * m[3] * m[5];
        
        inv[15] = m[0] * m[5] * m[10] - m[0] * m[6] * m[9] - m[4] * m[1] * m[10] 
                + m[4] * m[2] * m[9] + m[8] * m[1] * m[6] - m[8] * m[2] * m[5];
        
        det = m[0] * inv[0] + m[1] * inv[4] + m[2] * inv[8] + m[3] * inv[12];
        
        if (det == 0) return false;
        
        det = 1.0f / det;
        
        for (int i = 0; i < 16; i++)
            m[i] = inv[i] * det;
        
        return true;
    }
    
    void transpose() {
        float temp[16];
        for (int i = 0; i < 4; ++i) {
            for (int j = 0; j < 4; ++j) {
                temp[i*4 + j] = m[j*4 + i];
            }
        }
        memcpy(m, temp, sizeof(m));
    }
    
    // Transformations
    Vector3 transformPoint(const Vector3& point) const {
        return Vector3(
            m[0]*point.x + m[4]*point.y + m[8]*point.z + m[12],
            m[1]*point.x + m[5]*point.y + m[9]*point.z + m[13],
            m[2]*point.x + m[6]*point.y + m[10]*point.z + m[14]
        );
    }
    
    Vector3 transformVector(const Vector3& vector) const {
        return Vector3(
            m[0]*vector.x + m[4]*vector.y + m[8]*vector.z,
            m[1]*vector.x + m[5]*vector.y + m[9]*vector.z,
            m[2]*vector.x + m[6]*vector.y + m[10]*vector.z
        );
    }
    
    // Advanced operations
    float getUniformScale() const {
        Vector3 v(1, 0, 0);
        Vector3 transformed = transformVector(v);
        return sqrtf(transformed.x*transformed.x + transformed.y*transformed.y + transformed.z*transformed.z);
    }
    
    void decompose(Vector3& translation, Vector3& rotation, Vector3& scale) const {
        // Extract translation
        translation.x = m[12];
        translation.y = m[13];
        translation.z = m[14];
        
        // Extract scale
        scale.x = sqrtf(m[0]*m[0] + m[1]*m[1] + m[2]*m[2]);
        scale.y = sqrtf(m[4]*m[4] + m[5]*m[5] + m[6]*m[6]);
        scale.z = sqrtf(m[8]*m[8] + m[9]*m[9] + m[10]*m[10]);
        
        // Extract rotation (simplified)
        float sx = 1.0f / scale.x;
        float sy = 1.0f / scale.y;
        float sz = 1.0f / scale.z;
        
        float m00 = m[0] * sx;
        float m01 = m[4] * sy;
        float m02 = m[8] * sz;
        
        float m10 = m[1] * sx;
        float m11 = m[5] * sy;
        float m12 = m[9] * sz;
        
        float m20 = m[2] * sx;
        float m21 = m[6] * sy;
        float m22 = m[10] * sz;
        
        rotation.y = asinf(-m20);
        if (cosf(rotation.y) != 0) {
            rotation.x = atan2f(m21, m22);
            rotation.z = atan2f(m10, m00);
        } else {
            rotation.x = atan2f(-m12, m11);
            rotation.z = 0;
        }
        
        // Convert to degrees
        rotation.x *= 180.0f / CAP_PI;
        rotation.y *= 180.0f / CAP_PI;
        rotation.z *= 180.0f / CAP_PI;
    }
    
    // Utility functions
    void print() const {
        printf("Matrix:\n");
        printf("[%6.2f %6.2f %6.2f %6.2f]\n", m[0], m[4], m[8], m[12]);
        printf("[%6.2f %6.2f %6.2f %6.2f]\n", m[1], m[5], m[9], m[13]);
        printf("[%6.2f %6.2f %6.2f %6.2f]\n", m[2], m[6], m[10], m[14]);
        printf("[%6.2f %6.2f %6.2f %6.2f]\n", m[3], m[7], m[11], m[15]);
    }
    
    // Operator overloading
    AffineMatrix operator*(const AffineMatrix& other) const {
        AffineMatrix result;
        for (int i = 0; i < 4; ++i) {
            for (int j = 0; j < 4; ++j) {
                result.m[i*4 + j] = 0;
                for (int k = 0; k < 4; ++k) {
                    result.m[i*4 + j] += m[i*4 + k] * other.m[k*4 + j];
                }
            }
        }
        return result;
    }
    
    Vector3 operator*(const Vector3& point) const {
        return transformPoint(point);
    }
};
