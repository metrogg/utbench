#include <stdexcept>
#include <cmath>

double safeDivide(double numerator, double denominator) {
    if (denominator == 0) {
        throw std::runtime_error("Cannot divide by zero");
    }
    if (numerator == 0) {
        return 0.0;
    }
    double result = numerator / denominator;
    return result;
}

int safeSqrt(int value) {
    if (value < 0) {
        throw std::invalid_argument("Cannot compute square root of negative number");
    }
    if (value == 0) {
        return 0;
    }
    return static_cast<int>(std::sqrt(static_cast<double>(value)));
}