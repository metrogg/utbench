#include <vector>
#include <cmath>
#include <map>
#include <iomanip>

using namespace std;

const float PI = 3.14159265358979323846f;
const float TWO_PI = 2.0f * PI;

class WaveGenerator {
private:
    float currentAngle = 0.0f;
    float angleIncrement = 0.0f;
    float level = 1.0f;
    float tailOff = 0.0f;
    bool isPlaying = false;
    float frequency = 440.0f;
    int sampleRate = 44100;
    string waveType = "sine";

public:
    void startNote(float freq, float vel) {
        frequency = freq;
        currentAngle = 0.0f;
        angleIncrement = frequency / sampleRate * TWO_PI;
        tailOff = 0.0f;
        level = vel;
        isPlaying = true;
    }

    void stopNote(bool allowTailOff) {
        if (allowTailOff) {
            if (tailOff == 0.0f)
                tailOff = 1.0f;
        } else {
            level = 0.0f;
            currentAngle = 0.0f;
            isPlaying = false;
        }
    }

    void setWaveType(const string& type) {
        waveType = type;
    }

    void setSampleRate(int rate) {
        sampleRate = rate;
    }

    float generateSample() {
        if (!isPlaying) return 0.0f;

        float value = 0.0f;
        
        if (currentAngle > TWO_PI) {
            currentAngle -= TWO_PI;
        }

        if (waveType == "sine") {
            value = sin(currentAngle) * level;
        } 
        else if (waveType == "square") {
            value = (sin(currentAngle) > 0) ? level : -level;
        } 
        else if (waveType == "sawtooth") {
            value = (1.0f / PI * currentAngle - 1.0f) * level;
        } 
        else if (waveType == "triangle") {
            if (currentAngle < PI) {
                value = (-1.0f + 2.0f * currentAngle / PI) * level;
            } else {
                value = (3.0f - 2.0f * currentAngle / PI) * level;
            }
        }

        if (tailOff > 0.0f) {
            value *= tailOff;
            tailOff *= 0.99f;
            if (tailOff <= 0.05f) {
                isPlaying = false;
            }
        }

        currentAngle += angleIncrement;
        return value;
    }

    vector<float> generateBlock(int numSamples) {
        vector<float> samples;
        samples.reserve(numSamples);
        for (int i = 0; i < numSamples; ++i) {
            samples.push_back(generateSample());
        }
        return samples;
    }
};
