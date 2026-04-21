#include <vector>
#include <cmath>
#include <algorithm>
#include <stdexcept>
#include <map>

using namespace std;

class RobotController {
private:
    // Control parameters
    double line_y;
    double prop_factor;
    double int_factor;
    double diff_factor;
    double int_error;
    double old_error;
    
    // Robot state
    double x, y, theta;
    
    // Simulation parameters
    double simulation_time;
    double time_step;
    
public:
    RobotController(double set_point, double kp, double ki, double kd, double initial_x, double initial_y, double initial_theta)
        : line_y(set_point), prop_factor(kp), int_factor(ki), diff_factor(kd),
          int_error(0), old_error(0), x(initial_x), y(initial_y), theta(initial_theta),
          simulation_time(0), time_step(0.1) {}
    
    // PID controller for line following
    double calculateControlSignal() {
        double current_error = line_y - y;
        int_error += current_error * time_step;
        double derivative = (current_error - old_error) / time_step;
        old_error = current_error;
        
        return prop_factor * current_error + 
               int_factor * int_error + 
               diff_factor * derivative;
    }
    
    // Simulate robot motion based on control signal
    void updateRobotState(double control_signal, double desired_velocity) {
        // Simple kinematic model
        double angular_velocity = control_signal;
        theta += angular_velocity * time_step;
        
        // Update position
        x += desired_velocity * cos(theta) * time_step;
        y += desired_velocity * sin(theta) * time_step;
        
        simulation_time += time_step;
    }
    
    // Run simulation for given duration
    map<double, vector<double>> simulate(double duration, double desired_velocity) {
        map<double, vector<double>> trajectory;
        
        while (simulation_time < duration) {
            double control = calculateControlSignal();
            updateRobotState(control, desired_velocity);
            
            trajectory[simulation_time] = {x, y, theta, control};
        }
        
        return trajectory;
    }
    
    // Reset simulation
    void reset(double initial_x, double initial_y, double initial_theta) {
        x = initial_x;
        y = initial_y;
        theta = initial_theta;
        simulation_time = 0;
        int_error = 0;
        old_error = 0;
    }
};
