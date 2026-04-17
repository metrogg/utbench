#include <vector>
#include <cmath>
#include <algorithm>
#include <stdexcept>
#include <iomanip>

using namespace std;

// Represents a grayscale image as a 2D vector of pixel values (0-255)
class GrayscaleImage {
private:
    vector<vector<int>> pixels;
    int width;
    int height;

public:
    GrayscaleImage(int w, int h) : width(w), height(h) {
        pixels = vector<vector<int>>(h, vector<int>(w, 0));
    }

    GrayscaleImage(const vector<vector<int>>& pixel_data) {
        if (pixel_data.empty() || pixel_data[0].empty()) {
            throw invalid_argument("Image cannot be empty");
        }
        height = pixel_data.size();
        width = pixel_data[0].size();
        pixels = pixel_data;
    }

    int get_width() const { return width; }
    int get_height() const { return height; }
    int get_pixel(int y, int x) const { return pixels[y][x]; }
    void set_pixel(int y, int x, int value) { pixels[y][x] = value; }

    // Applies Sobel operator for edge detection
    void detect_edges() {
        vector<vector<int>> output(height, vector<int>(width, 0));
        
        // Sobel kernels
        const int sobel_x[3][3] = {{-1, 0, 1}, {-2, 0, 2}, {-1, 0, 1}};
        const int sobel_y[3][3] = {{-1, -2, -1}, {0, 0, 0}, {1, 2, 1}};

        for (int y = 1; y < height - 1; y++) {
            for (int x = 1; x < width - 1; x++) {
                int gx = 0, gy = 0;

                // Apply Sobel operator
                for (int ky = -1; ky <= 1; ky++) {
                    for (int kx = -1; kx <= 1; kx++) {
                        int pixel = pixels[y + ky][x + kx];
                        gx += pixel * sobel_x[ky + 1][kx + 1];
                        gy += pixel * sobel_y[ky + 1][kx + 1];
                    }
                }

                // Calculate gradient magnitude
                int magnitude = static_cast<int>(sqrt(gx * gx + gy * gy));
                output[y][x] = min(255, magnitude); // Clamp to 255
            }
        }

        // Copy edges back to original image
        for (int y = 1; y < height - 1; y++) {
            for (int x = 1; x < width - 1; x++) {
                pixels[y][x] = output[y][x];
            }
        }
    }

    // Prints the image matrix (for testing purposes)
    void print_image() const {
        for (const auto& row : pixels) {
            for (int pixel : row) {
                cout << setw(4) << pixel;
            }
            cout << endl;
        }
    }
};
