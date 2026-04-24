#include <vector>
#include <cmath>
#include <algorithm>
#include <stdexcept>
#include <map>
#include <array>

using namespace std;

// Simulate a simple 2D texture class for our standalone implementation
class Texture {
private:
    vector<vector<array<float, 4>>> data;
    int width, height;

public:
    Texture(int w, int h) : width(w), height(h) {
        data.resize(height, vector<array<float, 4>>(width, {0.0f, 0.0f, 0.0f, 1.0f}));
    }

    void setPixel(int x, int y, float r, float g, float b, float a = 1.0f) {
        if (x >= 0 && x < width && y >= 0 && y < height) {
            data[y][x] = {r, g, b, a};
        }
    }

    array<float, 4> getPixel(float u, float v) const {
        // Basic bilinear filtering
        u = max(0.0f, min(1.0f, u));
        v = max(0.0f, min(1.0f, v));
        
        float x = u * (width - 1);
        float y = v * (height - 1);
        
        int x0 = static_cast<int>(floor(x));
        int y0 = static_cast<int>(floor(y));
        int x1 = min(x0 + 1, width - 1);
        int y1 = min(y0 + 1, height - 1);
        
        float xf = x - x0;
        float yf = y - y0;
        
        array<float, 4> p00 = data[y0][x0];
        array<float, 4> p10 = data[y0][x1];
        array<float, 4> p01 = data[y1][x0];
        array<float, 4> p11 = data[y1][x1];
        
        array<float, 4> result;
        for (int i = 0; i < 4; ++i) {
            float top = p00[i] * (1 - xf) + p10[i] * xf;
            float bottom = p01[i] * (1 - xf) + p11[i] * xf;
            result[i] = top * (1 - yf) + bottom * yf;
        }
        
        return result;
    }

    static Texture createDefaultLookup() {
        Texture tex(512, 512);
        // Create a simple lookup table with color gradients
        for (int y = 0; y < 512; ++y) {
            for (int x = 0; x < 512; ++x) {
                float r = x / 512.0f;
                float g = y / 512.0f;
                float b = (x + y) / 1024.0f;
                tex.setPixel(x, y, r, g, b);
            }
        }
        return tex;
    }
};

// Enhanced lookup filter with multiple operations
array<float, 4> applyLookupFilter(const array<float, 4>& inputColor, const Texture& lookupTexture, 
                                 float alpha = 1.0f, bool useSepia = false, float intensity = 1.0f) {
    array<float, 4> color = inputColor;
    
    // Optional sepia tone conversion before lookup
    if (useSepia) {
        float r = color[0];
        float g = color[1];
        float b = color[2];
        color[0] = min(1.0f, (r * 0.393f) + (g * 0.769f) + (b * 0.189f));
        color[1] = min(1.0f, (r * 0.349f) + (g * 0.686f) + (b * 0.168f));
        color[2] = min(1.0f, (r * 0.272f) + (g * 0.534f) + (b * 0.131f));
    }

    // Main lookup table operation (similar to original shader logic)
    float blueColor = color[2] * 63.0f;
    
    float quad1_y = floor(floor(blueColor) / 8.0f);
    float quad1_x = floor(blueColor) - (quad1_y * 8.0f);
    
    float quad2_y = floor(ceil(blueColor) / 8.0f);
    float quad2_x = ceil(blueColor) - (quad2_y * 8.0f);
    
    float texPos1_x = (quad1_x * 0.125f) + 0.5f/512.0f + ((0.125f - 1.0f/512.0f) * color[0]);
    float texPos1_y = (quad1_y * 0.125f) + 0.5f/512.0f + ((0.125f - 1.0f/512.0f) * color[1]);
    
    float texPos2_x = (quad2_x * 0.125f) + 0.5f/512.0f + ((0.125f - 1.0f/512.0f) * color[0]);
    float texPos2_y = (quad2_y * 0.125f) + 0.5f/512.0f + ((0.125f - 1.0f/512.0f) * color[1]);
    
    array<float, 4> newColor1 = lookupTexture.getPixel(texPos1_x, texPos1_y);
    array<float, 4> newColor2 = lookupTexture.getPixel(texPos2_x, texPos2_y);
    
    // Blend between the two lookup results
    float blendFactor = blueColor - floor(blueColor);
    array<float, 4> result;
    for (int i = 0; i < 3; ++i) {
        result[i] = newColor1[i] * (1.0f - blendFactor) + newColor2[i] * blendFactor;
    }
    result[3] = alpha;
    
    // Apply intensity
    if (intensity != 1.0f) {
        for (int i = 0; i < 3; ++i) {
            result[i] = color[i] * (1.0f - intensity) + result[i] * intensity;
        }
    }
    
    return result;
}
