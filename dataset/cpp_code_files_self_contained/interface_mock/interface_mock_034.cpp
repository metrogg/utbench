#include <vector>
#include <cmath>
#include <array>

using namespace std;

class Matrix4x4 {
public:
    array<array<float, 4>, 4> data;

    Matrix4x4() {
        // Initialize as identity matrix
        for (int i = 0; i < 4; ++i) {
            for (int j = 0; j < 4; ++j) {
                data[i][j] = (i == j) ? 1.0f : 0.0f;
            }
        }
    }

    static Matrix4x4 createTranslation(float tx, float ty, float tz) {
        Matrix4x4 m;
        m.data[0][3] = tx;
        m.data[1][3] = ty;
        m.data[2][3] = tz;
        return m;
    }

    static Matrix4x4 createScale(float sx, float sy, float sz) {
        Matrix4x4 m;
        m.data[0][0] = sx;
        m.data[1][1] = sy;
        m.data[2][2] = sz;
        return m;
    }

    static Matrix4x4 createRotationX(float angle) {
        float rad = angle * M_PI / 180.0f;
        float c = cos(rad);
        float s = sin(rad);
        
        Matrix4x4 m;
        m.data[1][1] = c;
        m.data[1][2] = -s;
        m.data[2][1] = s;
        m.data[2][2] = c;
        return m;
    }

    static Matrix4x4 createRotationY(float angle) {
        float rad = angle * M_PI / 180.0f;
        float c = cos(rad);
        float s = sin(rad);
        
        Matrix4x4 m;
        m.data[0][0] = c;
        m.data[0][2] = s;
        m.data[2][0] = -s;
        m.data[2][2] = c;
        return m;
    }

    static Matrix4x4 createRotationZ(float angle) {
        float rad = angle * M_PI / 180.0f;
        float c = cos(rad);
        float s = sin(rad);
        
        Matrix4x4 m;
        m.data[0][0] = c;
        m.data[0][1] = -s;
        m.data[1][0] = s;
        m.data[1][1] = c;
        return m;
    }

    Matrix4x4 multiply(const Matrix4x4& other) const {
        Matrix4x4 result;
        for (int i = 0; i < 4; ++i) {
            for (int j = 0; j < 4; ++j) {
                result.data[i][j] = 0;
                for (int k = 0; k < 4; ++k) {
                    result.data[i][j] += data[i][k] * other.data[k][j];
                }
            }
        }
        return result;
    }

    array<float, 4> transformPoint(const array<float, 4>& point) const {
        array<float, 4> result = {0, 0, 0, 0};
        for (int i = 0; i < 4; ++i) {
            for (int j = 0; j < 4; ++j) {
                result[i] += data[i][j] * point[j];
            }
        }
        return result;
    }
};

class Camera {
private:
    array<float, 3> position;
    array<float, 3> lookAt;
    array<float, 3> upVector;
    float fov;
    float aspectRatio;
    float nearPlane;
    float farPlane;

public:
    Camera(array<float, 3> pos, array<float, 3> look, array<float, 3> up, 
           float fovDegrees, float aspect, float near, float far)
        : position(pos), lookAt(look), upVector(up), fov(fovDegrees), 
          aspectRatio(aspect), nearPlane(near), farPlane(far) {}

    Matrix4x4 getViewMatrix() const {
        array<float, 3> zAxis = {
            position[0] - lookAt[0],
            position[1] - lookAt[1],
            position[2] - lookAt[2]
        };
        normalize(zAxis);

        array<float, 3> xAxis = crossProduct(upVector, zAxis);
        normalize(xAxis);

        array<float, 3> yAxis = crossProduct(zAxis, xAxis);

        Matrix4x4 view;
        view.data[0][0] = xAxis[0];
        view.data[1][0] = xAxis[1];
        view.data[2][0] = xAxis[2];
        view.data[0][1] = yAxis[0];
        view.data[1][1] = yAxis[1];
        view.data[2][1] = yAxis[2];
        view.data[0][2] = zAxis[0];
        view.data[1][2] = zAxis[1];
        view.data[2][2] = zAxis[2];
        view.data[0][3] = -dotProduct(xAxis, position);
        view.data[1][3] = -dotProduct(yAxis, position);
        view.data[2][3] = -dotProduct(zAxis, position);

        return view;
    }

    Matrix4x4 getProjectionMatrix() const {
        float f = 1.0f / tan(fov * M_PI / 360.0f);
        float rangeInv = 1.0f / (nearPlane - farPlane);

        Matrix4x4 proj;
        proj.data[0][0] = f / aspectRatio;
        proj.data[1][1] = f;
        proj.data[2][2] = (nearPlane + farPlane) * rangeInv;
        proj.data[2][3] = 2 * nearPlane * farPlane * rangeInv;
        proj.data[3][2] = -1.0f;
        proj.data[3][3] = 0.0f;

        return proj;
    }

private:
    static float dotProduct(const array<float, 3>& a, const array<float, 3>& b) {
        return a[0]*b[0] + a[1]*b[1] + a[2]*b[2];
    }

    static array<float, 3> crossProduct(const array<float, 3>& a, const array<float, 3>& b) {
        return {
            a[1]*b[2] - a[2]*b[1],
            a[2]*b[0] - a[0]*b[2],
            a[0]*b[1] - a[1]*b[0]
        };
    }

    static void normalize(array<float, 3>& v) {
        float length = sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2]);
        if (length > 0) {
            v[0] /= length;
            v[1] /= length;
            v[2] /= length;
        }
    }
};

array<float, 3> applyTransformations(
    const array<float, 3>& point,
    const vector<Matrix4x4>& transformations,
    const Camera& camera
) {
    // Convert 3D point to homogeneous coordinates
    array<float, 4> homogeneousPoint = {point[0], point[1], point[2], 1.0f};
    
    // Apply all transformations
    for (const auto& transform : transformations) {
        homogeneousPoint = transform.transformPoint(homogeneousPoint);
    }
    
    // Apply view transformation
    Matrix4x4 viewMatrix = camera.getViewMatrix();
    homogeneousPoint = viewMatrix.transformPoint(homogeneousPoint);
    
    // Apply projection
    Matrix4x4 projMatrix = camera.getProjectionMatrix();
    homogeneousPoint = projMatrix.transformPoint(homogeneousPoint);
    
    // Perspective division
    if (homogeneousPoint[3] != 0.0f) {
        homogeneousPoint[0] /= homogeneousPoint[3];
        homogeneousPoint[1] /= homogeneousPoint[3];
        homogeneousPoint[2] /= homogeneousPoint[3];
    }
    
    // Return as 3D point (discarding w component)
    return {homogeneousPoint[0], homogeneousPoint[1], homogeneousPoint[2]};
}
