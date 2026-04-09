#include <cmath>
#include <vector>
#include <stdexcept>

using namespace std;

class Quaternion {
private:
    double x, y, z, w;

public:
    // Constructors
    Quaternion() : x(0), y(0), z(0), w(1) {}
    Quaternion(double x, double y, double z, double w) : x(x), y(y), z(z), w(w) {}
    
    // Core operations
    Quaternion operator+(const Quaternion& q) const {
        return Quaternion(x + q.x, y + q.y, z + q.z, w + q.w);
    }
    
    Quaternion operator-(const Quaternion& q) const {
        return Quaternion(x - q.x, y - q.y, z - q.z, w - q.w);
    }
    
    Quaternion operator*(const Quaternion& q) const {
        return Quaternion(
            w * q.x + x * q.w + y * q.z - z * q.y,
            w * q.y - x * q.z + y * q.w + z * q.x,
            w * q.z + x * q.y - y * q.x + z * q.w,
            w * q.w - x * q.x - y * q.y - z * q.z
        );
    }
    
    Quaternion operator*(double scalar) const {
        return Quaternion(x * scalar, y * scalar, z * scalar, w * scalar);
    }
    
    Quaternion conjugate() const {
        return Quaternion(-x, -y, -z, w);
    }
    
    double norm() const {
        return x*x + y*y + z*z + w*w;
    }
    
    Quaternion inverse() const {
        double n = norm();
        if (n == 0) throw runtime_error("Cannot invert zero quaternion");
        return conjugate() * (1.0 / n);
    }
    
    // Advanced operations
    static Quaternion fromAxisAngle(double ax, double ay, double az, double angle) {
        double halfAngle = angle / 2;
        double sinHalf = sin(halfAngle);
        return Quaternion(
            ax * sinHalf,
            ay * sinHalf,
            az * sinHalf,
            cos(halfAngle)
        );
    }
    
    static Quaternion slerp(const Quaternion& q1, const Quaternion& q2, double t) {
        double dot = q1.x*q2.x + q1.y*q2.y + q1.z*q2.z + q1.w*q2.w;
        
        // Ensure shortest path
        Quaternion q2adj = (dot < 0) ? q2 * -1 : q2;
        dot = abs(dot);
        
        const double DOT_THRESHOLD = 0.9995;
        if (dot > DOT_THRESHOLD) {
            // Linear interpolation for very close quaternions
            Quaternion result = q1 + (q2adj - q1) * t;
            return result * (1.0 / sqrt(result.norm()));
        }
        
        double theta0 = acos(dot);
        double theta = theta0 * t;
        double sinTheta = sin(theta);
        double sinTheta0 = sin(theta0);
        
        double s1 = cos(theta) - dot * sinTheta / sinTheta0;
        double s2 = sinTheta / sinTheta0;
        
        return (q1 * s1) + (q2adj * s2);
    }
    
    // Rotation operations
    Quaternion rotate(const Quaternion& point) const {
        return (*this) * point * this->inverse();
    }
    
    vector<double> rotateVector(double vx, double vy, double vz) const {
        Quaternion p(vx, vy, vz, 0);
        Quaternion result = (*this) * p * this->inverse();
        return {result.x, result.y, result.z};
    }
    
    // Conversion to other representations
    vector<double> toEulerAngles() const {
        // Roll (x-axis rotation)
        double sinr_cosp = 2 * (w * x + y * z);
        double cosr_cosp = 1 - 2 * (x * x + y * y);
        double roll = atan2(sinr_cosp, cosr_cosp);

        // Pitch (y-axis rotation)
        double sinp = 2 * (w * y - z * x);
        double pitch;
        if (abs(sinp) >= 1)
            pitch = copysign(M_PI / 2, sinp); // Use 90 degrees if out of range
        else
            pitch = asin(sinp);

        // Yaw (z-axis rotation)
        double siny_cosp = 2 * (w * z + x * y);
        double cosy_cosp = 1 - 2 * (y * y + z * z);
        double yaw = atan2(siny_cosp, cosy_cosp);

        return {roll, pitch, yaw};
    }
    
    // Getters
    vector<double> getComponents() const { return {x, y, z, w}; }
};
