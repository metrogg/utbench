#include <stdexcept>

int clamp(int value, int minVal, int maxVal) {
    if (minVal > maxVal) {
        throw std::invalid_argument("minVal must be less than or equal to maxVal");
    }
    if (value < minVal) {
        return minVal;
    }
    if (value > maxVal) {
        return maxVal;
    }
    return value;
}