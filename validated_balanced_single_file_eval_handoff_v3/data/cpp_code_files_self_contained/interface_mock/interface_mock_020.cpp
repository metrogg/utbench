#include <cmath>
#include <vector>
#include <stdexcept>
#include <memory>

using namespace std;

class CBall {
public:
    CBall(double x, double y, double vx, double vy, double mass) 
        : x(x), y(y), vx(vx), vy(vy), mass(mass) {}
    
    void updatePosition(double dt) {
        x += vx * dt;
        y += vy * dt;
    }
    
    void applyForce(double fx, double fy, double dt) {
        vx += fx / mass * dt;
        vy += fy / mass * dt;
    }
    
    double getX() const { return x; }
    double getY() const { return y; }
    double getVX() const { return vx; }
    double getVY() const { return vy; }
    double getMass() const { return mass; }

private:
    double x, y;    // Position
    double vx, vy;  // Velocity
    double mass;    // Mass
};

class CTrap {
public:
    CTrap(double x, double y, double f, double r) 
        : x(x), y(y), f(f), r(r) {
        if (f <= 0) throw invalid_argument("Force coefficient must be positive");
        if (r <= 0) throw invalid_argument("Radius must be positive");
    }
    
    // Calculate force vector acting on a ball
    pair<double, double> calculateForce(const CBall& ball) const {
        double dx = x - ball.getX();
        double dy = y - ball.getY();
        double distance_sq = dx*dx + dy*dy;
        
        if (distance_sq > r*r || distance_sq == 0) {
            return {0, 0};
        }
        
        double distance = sqrt(distance_sq);
        double force_magnitude = f / distance_sq;
        
        return {force_magnitude * dx/distance, 
                force_magnitude * dy/distance};
    }
    
    // Apply trap effect to a ball over time dt
    void effect(CBall& ball, double dt) const {
        auto [fx, fy] = calculateForce(ball);
        ball.applyForce(fx, fy, dt);
    }
    
    // Simulate multiple balls in the trap's field
    vector<pair<double, double>> simulateBalls(const vector<CBall>& balls, double total_time, double dt) const {
        vector<pair<double, double>> final_positions;
        vector<CBall> ball_copies = balls;
        
        for (double t = 0; t < total_time; t += dt) {
            for (auto& ball : ball_copies) {
                effect(ball, dt);
                ball.updatePosition(dt);
            }
        }
        
        for (const auto& ball : ball_copies) {
            final_positions.emplace_back(ball.getX(), ball.getY());
        }
        
        return final_positions;
    }

private:
    double x, y;  // Position
    double f;     // Force coefficient
    double r;     // Effective radius
};
