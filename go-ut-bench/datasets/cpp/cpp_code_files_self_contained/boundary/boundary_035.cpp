#include <vector>
#include <map>
#include <cmath>
#include <stdexcept>

using namespace std;

class StepperController {
private:
    // Motor parameters
    int dirPin;
    int enablePin;
    int stepPin;
    int microSteps;
    int stepsPerRev;
    double gearRatio;
    
    // Movement parameters
    double currentPosition;  // in revolutions
    double targetPosition;
    double currentSpeed;     // in rev/s
    double targetSpeed;
    double acceleration;     // in rev/s²
    bool isEnabled;
    bool isMoving;
    bool isContinuous;
    
public:
    StepperController(int dirPin, int enablePin, int stepPin, 
                     int microSteps = 1, int stepsPerRev = 200, double gearRatio = 1.0)
        : dirPin(dirPin), enablePin(enablePin), stepPin(stepPin),
          microSteps(microSteps), stepsPerRev(stepsPerRev), gearRatio(gearRatio),
          currentPosition(0), targetPosition(0), currentSpeed(0), targetSpeed(0),
          acceleration(1.0), isEnabled(false), isMoving(false), isContinuous(false) {}
    
    void enable() {
        isEnabled = true;
        // Simulate enabling outputs
    }
    
    void disable() {
        isEnabled = false;
        isMoving = false;
        isContinuous = false;
        currentSpeed = 0;
        // Simulate disabling outputs
    }
    
    void setCurrentPosition(double newPos) {
        if (!isMoving) {
            currentPosition = newPos;
        }
    }
    
    double getCurrentPosition() const {
        return currentPosition;
    }
    
    void setSpeed(double speedRevPerSec) {
        targetSpeed = speedRevPerSec;
        if (isContinuous) {
            // Update continuous movement speed
            currentSpeed = targetSpeed;
        }
    }
    
    double getCurrentSpeed() const {
        return currentSpeed;
    }
    
    void setAcceleration(double accelRevPerSecSq) {
        acceleration = accelRevPerSecSq;
    }
    
    void move(double revolutions) {
        if (!isEnabled) return;
        
        isContinuous = false;
        targetPosition = currentPosition + revolutions;
        isMoving = true;
        
        // Simulate acceleration
        double distance = abs(targetPosition - currentPosition);
        double timeToAccel = abs(targetSpeed) / acceleration;
        double accelDistance = 0.5 * acceleration * timeToAccel * timeToAccel;
        
        if (accelDistance > distance / 2) {
            // Triangular speed profile
            double maxSpeed = sqrt(acceleration * distance);
            targetSpeed = (revolutions > 0) ? maxSpeed : -maxSpeed;
        }
    }
    
    void moveTo(double position) {
        move(position - currentPosition);
    }
    
    void runForward() {
        if (!isEnabled) return;
        
        isContinuous = true;
        isMoving = true;
        targetSpeed = abs(targetSpeed);  // Ensure positive speed
        currentSpeed = targetSpeed;
    }
    
    void runBackward() {
        if (!isEnabled) return;
        
        isContinuous = true;
        isMoving = true;
        targetSpeed = -abs(targetSpeed);  // Ensure negative speed
        currentSpeed = targetSpeed;
    }
    
    void stopMove() {
        if (isMoving) {
            if (isContinuous) {
                // For continuous movement, just stop
                isContinuous = false;
                currentSpeed = 0;
                isMoving = false;
            } else {
                // For targeted movement, calculate deceleration distance
                double decelTime = abs(currentSpeed) / acceleration;
                double decelDistance = 0.5 * currentSpeed * decelTime;
                
                // Update target position to stopping point
                if (currentSpeed > 0) {
                    targetPosition = min(targetPosition, currentPosition + decelDistance);
                } else {
                    targetPosition = max(targetPosition, currentPosition - decelDistance);
                }
            }
        }
    }
    
    void forceStop() {
        currentSpeed = 0;
        isMoving = false;
        isContinuous = false;
    }
    
    void update(double deltaTime) {
        if (!isMoving || !isEnabled) return;
        
        if (isContinuous) {
            // Continuous movement - just update position based on speed
            currentPosition += currentSpeed * deltaTime;
        } else {
            // Targeted movement with acceleration/deceleration
            double distanceToTarget = targetPosition - currentPosition;
            double direction = (distanceToTarget > 0) ? 1.0 : -1.0;
            
            // Calculate stopping distance with current speed
            double stoppingDistance = (currentSpeed * currentSpeed) / (2 * acceleration);
            
            if (abs(distanceToTarget) <= stoppingDistance) {
                // Need to decelerate
                double speedChange = -direction * acceleration * deltaTime;
                if ((currentSpeed > 0 && currentSpeed + speedChange < 0) || 
                    (currentSpeed < 0 && currentSpeed + speedChange > 0)) {
                    // Would overshoot - stop exactly at target
                    currentSpeed = 0;
                    currentPosition = targetPosition;
                    isMoving = false;
                } else {
                    currentSpeed += speedChange;
                    currentPosition += currentSpeed * deltaTime;
                }
            } else {
                // Can accelerate or maintain speed
                double speedChange = direction * acceleration * deltaTime;
                if (abs(currentSpeed + speedChange) > abs(targetSpeed)) {
                    // Would exceed target speed - cap at target
                    currentSpeed = direction * abs(targetSpeed);
                } else {
                    currentSpeed += speedChange;
                }
                currentPosition += currentSpeed * deltaTime;
            }
        }
    }
    
    map<string, double> getStatus() const {
        return {
            {"position", currentPosition},
            {"speed", currentSpeed},
            {"target_position", targetPosition},
            {"target_speed", targetSpeed},
            {"acceleration", acceleration},
            {"is_enabled", isEnabled},
            {"is_moving", isMoving},
            {"is_continuous", isContinuous}
        };
    }
};
