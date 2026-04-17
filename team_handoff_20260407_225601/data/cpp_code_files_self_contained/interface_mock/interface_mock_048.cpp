#include <vector>
#include <map>
#include <string>
#include <stdexcept>

using namespace std;

class PlayerMovementSimulator {
private:
    int xAxisPos;
    int yAxisPos;
    int xMapMovement;
    int yMapMovement;
    int yMovementM;
    int yMovementP;
    int xMovementM;
    int xMovementP;
    const int LEVEL_WIDTH = 720;
    const int LEVEL_HEIGHT = 510;

public:
    PlayerMovementSimulator(int xPos, int yPos, int xMap, int yMap) {
        xAxisPos = xPos;
        yAxisPos = yPos;
        xMapMovement = xMap;
        yMapMovement = yMap;
        yMovementM = 0;
        yMovementP = 0;
        xMovementM = 0;
        xMovementP = 0;
    }

    // Simulate movement based on keyboard input
    void simulateMovement(char direction) {
        switch(direction) {
            case 'U': // Up
                if (yAxisPos != 0) {
                    if (yAxisPos >= LEVEL_HEIGHT / 2 && yAxisPos <= 300 && yMapMovement > 0) {
                        yMovementP -= 1;
                        yMovementM -= 1;
                        if (yMovementP == -5) {
                            yAxisPos -= 1;
                            yMovementP = 0;
                        }
                        if (yMovementM == -3) {
                            yMapMovement -= 1;
                            yMovementM = 0;
                        }
                    } else {
                        yAxisPos -= 1;
                    }
                }
                break;
                
            case 'D': // Down
                if (yAxisPos != LEVEL_HEIGHT) {
                    if (yAxisPos >= LEVEL_HEIGHT / 2 && yAxisPos <= 350 && yMapMovement != 45) {
                        yMovementP += 1;
                        yMovementM += 1;
                        if (yMovementP == 5) {
                            yAxisPos += 1;
                            yMovementP = 0;
                        }
                        if (yMovementM == 3) {
                            yMapMovement += 1;
                            yMovementM = 0;
                        }
                    } else {
                        yAxisPos += 1;
                    }
                }
                break;
                
            case 'L': // Left
                if (xAxisPos != 20) {
                    if (xAxisPos >= LEVEL_WIDTH / 2 && xAxisPos <= 400 && xMapMovement > 0) {
                        xMovementP -= 1;
                        xMovementM -= 1;
                        if (xMovementP == -5) {
                            xAxisPos -= 1;
                            xMovementP = 0;
                        }
                        if (xMovementM == -3) {
                            xMapMovement -= 1;
                            xMovementM = 0;
                        }
                    } else {
                        xAxisPos -= 1;
                    }
                }
                break;
                
            case 'R': // Right
                if (xAxisPos != LEVEL_WIDTH) {
                    if (xAxisPos >= LEVEL_WIDTH / 2 && xAxisPos <= 400 && xMapMovement != 39) {
                        xMovementP += 1;
                        xMovementM += 1;
                        if (xMovementP == 5) {
                            xAxisPos += 1;
                            xMovementP = 0;
                        }
                        if (xMovementM == 3) {
                            xMapMovement += 1;
                            xMovementM = 0;
                        }
                    } else {
                        xAxisPos += 1;
                    }
                }
                break;
                
            default:
                throw invalid_argument("Invalid direction. Use U, D, L, or R.");
        }
    }

    // Get current position and map data
    map<string, int> getPositionData() {
        return {
            {"xPos", xAxisPos},
            {"yPos", yAxisPos},
            {"xMap", xMapMovement},
            {"yMap", yMapMovement}
        };
    }

    // Simulate multiple movements from a string of commands
    void executeMovementSequence(const string& commands) {
        for (char cmd : commands) {
            simulateMovement(cmd);
        }
    }
};
