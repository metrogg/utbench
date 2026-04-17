#include <vector>
#include <cmath>
#include <map>
#include <stdexcept>

using namespace std;

// Simulated Nunchuk data structure
struct NunchukData {
    int16_t joystickX;
    int16_t joystickY;
    bool buttonC;
    bool buttonZ;
    int16_t accelX;
    int16_t accelY;
    int16_t accelZ;
};

// Enhanced Nunchuk controller class
class NunchukController {
private:
    int16_t joystickThreshold;
    vector<pair<int16_t, int16_t>> movementHistory;
    static const int HISTORY_SIZE = 10;
    
public:
    NunchukController(int16_t threshold = 5) : joystickThreshold(threshold) {}
    
    // Simulate reading nunchuk data (in real implementation, this would read from hardware)
    bool readData(NunchukData& data) {
        // In a real implementation, this would read from actual hardware
        // For simulation purposes, we'll assume it always succeeds
        return true;
    }
    
    // Process nunchuk data with enhanced features
    map<string, int> processMovement(const NunchukData& data) {
        map<string, int> result;
        
        // Calculate movement magnitude
        int16_t magnitude = sqrt(data.joystickX * data.joystickX + data.joystickY * data.joystickY);
        
        // Determine movement direction (8-way)
        string direction = "neutral";
        if (abs(data.joystickX) > joystickThreshold || abs(data.joystickY) > joystickThreshold) {
            if (data.joystickY < -joystickThreshold) {
                if (data.joystickX < -joystickThreshold) direction = "up-left";
                else if (data.joystickX > joystickThreshold) direction = "up-right";
                else direction = "up";
            }
            else if (data.joystickY > joystickThreshold) {
                if (data.joystickX < -joystickThreshold) direction = "down-left";
                else if (data.joystickX > joystickThreshold) direction = "down-right";
                else direction = "down";
            }
            else {
                if (data.joystickX < -joystickThreshold) direction = "left";
                else if (data.joystickX > joystickThreshold) direction = "right";
            }
        }
        
        // Track movement history for gesture recognition
        if (movementHistory.size() >= HISTORY_SIZE) {
            movementHistory.erase(movementHistory.begin());
        }
        movementHistory.push_back({data.joystickX, data.joystickY});
        
        // Calculate movement velocity (simple difference from last position)
        int velocity = 0;
        if (movementHistory.size() > 1) {
            auto last = movementHistory.back();
            auto prev = *(movementHistory.end() - 2);
            velocity = sqrt(pow(last.first - prev.first, 2) + pow(last.second - prev.second, 2));
        }
        
        // Store results
        result["magnitude"] = magnitude;
        result["velocity"] = velocity;
        result["button_c"] = data.buttonC ? 1 : 0;
        result["button_z"] = data.buttonZ ? 1 : 0;
        result["accel_x"] = data.accelX;
        result["accel_y"] = data.accelY;
        result["accel_z"] = data.accelZ;
        
        // For direction, we'll use a separate output since it's a string
        // In C++ we can't mix types in a map, so we'll handle it separately
        
        return result;
    }
    
    // Detect if a gesture pattern matches (simple implementation)
    bool detectGesture(const vector<pair<int16_t, int16_t>>& pattern) {
        if (movementHistory.size() < pattern.size()) return false;
        
        for (size_t i = 0; i < pattern.size(); i++) {
            size_t historyIndex = movementHistory.size() - pattern.size() + i;
            auto historyPos = movementHistory[historyIndex];
            auto patternPos = pattern[i];
            
            if (abs(historyPos.first - patternPos.first) > joystickThreshold * 2 ||
                abs(historyPos.second - patternPos.second) > joystickThreshold * 2) {
                return false;
            }
        }
        return true;
    }
};
