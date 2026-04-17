#include <vector>
#include <map>
#include <string>
#include <cmath>
#include <stdexcept>

using namespace std;

// Simulated graphics effect types
enum class EffectType {
    PLASMA,
    DISTORT,
    FIRE,
    WAVE,
    PARTICLES,
    NONE
};

// Configuration for a single effect
struct EffectConfig {
    EffectType type;
    float intensity;
    float speed;
    float scale;
    bool enabled;
};

// Graphics system simulator
class GraphicsSystem {
private:
    vector<EffectConfig> effects;
    float time;
    int width;
    int height;

    // Simulate rendering an effect
    map<string, float> renderEffect(const EffectConfig& effect) {
        map<string, float> result;
        float t = time * effect.speed;
        
        switch(effect.type) {
            case EffectType::PLASMA:
                result["red"] = 0.5f + 0.5f * sin(t * 1.7f);
                result["green"] = 0.5f + 0.5f * sin(t * 2.3f);
                result["blue"] = 0.5f + 0.5f * sin(t * 3.1f);
                result["intensity"] = effect.intensity;
                break;
            case EffectType::DISTORT:
                result["x_offset"] = sin(t * 0.7f) * effect.scale;
                result["y_offset"] = cos(t * 0.9f) * effect.scale;
                result["intensity"] = effect.intensity;
                break;
            case EffectType::FIRE:
                result["heat"] = 0.8f + 0.2f * sin(t * 2.0f);
                result["intensity"] = effect.intensity;
                break;
            case EffectType::WAVE:
                result["amplitude"] = effect.scale * 0.5f;
                result["frequency"] = effect.speed * 0.3f;
                result["phase"] = t;
                break;
            case EffectType::PARTICLES:
                result["count"] = 100 * effect.intensity;
                result["speed"] = effect.speed;
                break;
            default:
                throw invalid_argument("Unknown effect type");
        }
        
        return result;
    }

public:
    GraphicsSystem(int w, int h) : width(w), height(h), time(0.0f) {}
    
    // Add an effect to the system
    void addEffect(EffectType type, float intensity = 1.0f, 
                  float speed = 1.0f, float scale = 1.0f) {
        effects.push_back({type, intensity, speed, scale, true});
    }
    
    // Update the graphics system
    void update(float deltaTime) {
        time += deltaTime;
    }
    
    // Render all enabled effects
    vector<map<string, float>> render() {
        vector<map<string, float>> results;
        
        for (const auto& effect : effects) {
            if (effect.enabled) {
                try {
                    auto effectResult = renderEffect(effect);
                    effectResult["type"] = static_cast<float>(effect.type);
                    results.push_back(effectResult);
                } catch (const invalid_argument& e) {
                    cerr << "Error rendering effect: " << e.what() << endl;
                }
            }
        }
        
        return results;
    }
    
    // Get effect count
    int getEffectCount() const {
        return effects.size();
    }
    
    // Toggle effect by index
    void toggleEffect(int index) {
        if (index >= 0 && index < effects.size()) {
            effects[index].enabled = !effects[index].enabled;
        }
    }
    
    // Get system dimensions
    pair<int, int> getDimensions() const {
        return {width, height};
    }
};
