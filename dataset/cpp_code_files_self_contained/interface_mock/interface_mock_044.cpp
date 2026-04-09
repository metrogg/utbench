#include <vector>
#include <cmath>
#include <array>
#include <iomanip>

using namespace std;

// Structure to hold robot pose (position and orientation)
struct RobotPose {
    double x;       // x-coordinate in meters
    double y;       // y-coordinate in meters
    double theta;   // orientation in radians
    
    vector<double> getAsVector() const {
        return {x, y, theta};
    }
};

// Extended Kalman Filter for robot localization
class RobotLocalization {
private:
    RobotPose pose;
    array<array<double, 3>, 3> covariance;
    
    // Helper function to keep angle between -PI and PI
    void normalizeAngle(double& angle) {
        angle -= 2*M_PI * floor(angle/(2*M_PI));  // 0..2*PI
        if (angle > M_PI) angle -= 2*M_PI;        // -PI..+PI
    }
    
    // Initialize covariance matrix
    void initCovariance() {
        covariance = {{
            {0.1, 0.0, 0.0},
            {0.0, 0.1, 0.0},
            {0.0, 0.0, 0.1}
        }};
    }
    
public:
    // Constructor with initial pose
    RobotLocalization(const vector<double>& initial_pose) {
        pose.x = initial_pose[0];
        pose.y = initial_pose[1];
        pose.theta = initial_pose[2];
        normalizeAngle(pose.theta);
        initCovariance();
    }
    
    // Get current pose
    vector<double> getPose() const {
        return pose.getAsVector();
    }
    
    // Prediction step of Kalman filter
    void predict(double v, double w, double dt) {
        // Time interval in seconds
        const double dt_sec = dt * 0.001;  // assuming dt is in milliseconds
        
        // Current orientation
        const double theta = pose.theta;
        const double cos_theta = cos(theta);
        const double sin_theta = sin(theta);
        
        // Calculate distance and angle change
        const double dist = v * dt_sec;
        const double gamma = w * dt_sec;
        
        // Update pose prediction
        pose.x += dist * cos_theta;
        pose.y += dist * sin_theta;
        pose.theta += gamma;
        normalizeAngle(pose.theta);
        
        // Jacobian matrix A (df/dx)
        array<array<double, 3>, 3> A = {{
            {1, 0, -dist * sin_theta},
            {0, 1, dist * cos_theta},
            {0, 0, 1}
        }};
        
        // Update covariance: P_new = A * P * A^T
        array<array<double, 3>, 3> P_new = {{{0}}};
        for (int i = 0; i < 3; ++i) {
            for (int j = 0; j < 3; ++j) {
                for (int k = 0; k < 3; ++k) {
                    for (int l = 0; l < 3; ++l) {
                        P_new[i][j] += A[i][k] * covariance[k][l] * A[j][l];
                    }
                }
            }
        }
        
        // Process noise (B * Q * B^T)
        // Simplified noise model compared to original
        const double motion_noise = 0.05;
        for (int i = 0; i < 3; ++i) {
            for (int j = 0; j < 3; ++j) {
                if (i == j) {
                    P_new[i][j] += motion_noise * motion_noise;
                }
            }
        }
        
        // Update covariance matrix
        covariance = P_new;
        
        // Ensure positive definiteness
        for (int i = 0; i < 3; ++i) {
            if (covariance[i][i] < 0.00001) covariance[i][i] = 0.00001;
        }
    }
    
    // Update step of Kalman filter with landmark measurement
    void updateLandmark(double measured_distance, double measured_angle, 
                       double landmark_x, double landmark_y) {
        // Predicted measurement
        const double dx = landmark_x - pose.x;
        const double dy = landmark_y - pose.y;
        double predicted_angle = atan2(dy, dx) - pose.theta;
        normalizeAngle(predicted_angle);
        
        // Measurement innovation (error)
        double error = measured_angle - predicted_angle;
        normalizeAngle(error);
        
        // Jacobian H (dh/dx)
        const double r2 = dx*dx + dy*dy;
        array<double, 3> H = {
            dy / r2,
            -dx / r2,
            -1
        };
        
        // Innovation covariance: S = H * P * H^T + R
        array<double, 3> PHt = {0};
        for (int i = 0; i < 3; ++i) {
            for (int k = 0; k < 3; ++k) {
                PHt[i] += covariance[i][k] * H[k];
            }
        }
        
        double innov_cov = 0;
        for (int k = 0; k < 3; ++k) {
            innov_cov += H[k] * PHt[k];
        }
        innov_cov += 0.01;  // Measurement noise (simplified)
        
        // Kalman gain: K = P * H^T * S^-1
        array<double, 3> K;
        for (int i = 0; i < 3; ++i) {
            K[i] = PHt[i] / innov_cov;
        }
        
        // Update state
        pose.x += K[0] * error;
        pose.y += K[1] * error;
        pose.theta += K[2] * error;
        normalizeAngle(pose.theta);
        
        // Update covariance: P = (I - K*H) * P
        for (int i = 0; i < 3; ++i) {
            for (int j = 0; j < 3; ++j) {
                covariance[i][j] -= K[i] * PHt[j];
            }
        }
        
        // Ensure positive definiteness
        for (int i = 0; i < 3; ++i) {
            if (covariance[i][i] < 0.00001) covariance[i][i] = 0.00001;
        }
    }
};
