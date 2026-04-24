#include <vector>
#include <stdexcept>
#include <cmath>

using namespace std;

class Pixel {
private:
    unsigned char red;
    unsigned char green;
    unsigned char blue;
    unsigned char alpha;

public:
    // Default constructor (black, fully opaque)
    Pixel() : red(0), green(0), blue(0), alpha(255) {}

    // Parameterized constructor
    Pixel(unsigned char r, unsigned char g, unsigned char b, unsigned char a = 255) 
        : red(r), green(g), blue(b), alpha(a) {
        validate();
    }

    // Getters
    unsigned char get_red() const { return red; }
    unsigned char get_green() const { return green; }
    unsigned char get_blue() const { return blue; }
    unsigned char get_alpha() const { return alpha; }

    // Setters with validation
    void set_red(unsigned char value) { red = value; validate(); }
    void set_green(unsigned char value) { green = value; validate(); }
    void set_blue(unsigned char value) { blue = value; validate(); }
    void set_alpha(unsigned char value) { alpha = value; validate(); }

    // Color operations
    Pixel blend(const Pixel& other) const {
        if (alpha == 0) return other;
        if (other.alpha == 0) return *this;
        
        double ratio1 = alpha / 255.0;
        double ratio2 = other.alpha / 255.0 * (1 - ratio1);
        double total_ratio = ratio1 + ratio2;
        
        unsigned char r = static_cast<unsigned char>((red * ratio1 + other.red * ratio2) / total_ratio);
        unsigned char g = static_cast<unsigned char>((green * ratio1 + other.green * ratio2) / total_ratio);
        unsigned char b = static_cast<unsigned char>((blue * ratio1 + other.blue * ratio2) / total_ratio);
        unsigned char a = static_cast<unsigned char>(255 * (ratio1 + ratio2 * (1 - ratio1)));
        
        return Pixel(r, g, b, a);
    }

    double brightness() const {
        return (0.299 * red + 0.587 * green + 0.114 * blue) / 255.0;
    }

    Pixel grayscale() const {
        unsigned char gray = static_cast<unsigned char>(0.299 * red + 0.587 * green + 0.114 * blue);
        return Pixel(gray, gray, gray, alpha);
    }

    Pixel invert() const {
        return Pixel(255 - red, 255 - green, 255 - blue, alpha);
    }

private:
    void validate() {
        // Alpha can be 0-255, but other components should be 0-255 regardless
        if (alpha > 255) {
            throw invalid_argument("Alpha value must be between 0 and 255");
        }
    }
};
