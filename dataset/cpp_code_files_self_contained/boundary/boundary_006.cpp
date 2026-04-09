#include <stdexcept>
#include <string>
#include <iomanip>
#include <sstream>

class Time {
private:
    int hours;
    int minutes;
    
    void normalize() {
        hours += minutes / 60;
        minutes = minutes % 60;
        hours = hours % 24;
    }

public:
    Time() : hours(0), minutes(0) {}
    Time(int h, int m = 0) : hours(h), minutes(m) {
        if (h < 0 || m < 0) {
            throw std::invalid_argument("Time values cannot be negative");
        }
        normalize();
    }
    
    // Arithmetic operations
    Time operator+(const Time &t) const {
        return Time(hours + t.hours, minutes + t.minutes);
    }
    
    Time operator-(const Time &t) const {
        int totalMinutes = (hours * 60 + minutes) - (t.hours * 60 + t.minutes);
        if (totalMinutes < 0) {
            throw std::runtime_error("Resulting time cannot be negative");
        }
        return Time(0, totalMinutes);
    }
    
    Time operator*(double factor) const {
        if (factor <= 0) {
            throw std::invalid_argument("Multiplication factor must be positive");
        }
        int totalMinutes = static_cast<int>((hours * 60 + minutes) * factor);
        return Time(0, totalMinutes);
    }
    
    // Comparison operators
    bool operator==(const Time &t) const {
        return hours == t.hours && minutes == t.minutes;
    }
    
    bool operator<(const Time &t) const {
        if (hours == t.hours) {
            return minutes < t.minutes;
        }
        return hours < t.hours;
    }
    
    // Compound assignment operators
    Time& operator+=(const Time &t) {
        hours += t.hours;
        minutes += t.minutes;
        normalize();
        return *this;
    }
    
    Time& operator-=(const Time &t) {
        *this = *this - t;
        return *this;
    }
    
    // Conversion to string
    std::string toString() const {
        std::ostringstream oss;
        oss << std::setw(2) << std::setfill('0') << hours << ":"
            << std::setw(2) << std::setfill('0') << minutes;
        return oss.str();
    }
    
    // Getters
    int getHours() const { return hours; }
    int getMinutes() const { return minutes; }
    
    // Time manipulation
    void addMinutes(int m) {
        if (m < 0) {
            throw std::invalid_argument("Minutes to add cannot be negative");
        }
        minutes += m;
        normalize();
    }
    
    void addHours(int h) {
        if (h < 0) {
            throw std::invalid_argument("Hours to add cannot be negative");
        }
        hours += h;
        normalize();
    }
    
    void reset(int h = 0, int m = 0) {
        if (h < 0 || m < 0) {
            throw std::invalid_argument("Time values cannot be negative");
        }
        hours = h;
        minutes = m;
        normalize();
    }
    
    // Time difference in minutes
    int differenceInMinutes(const Time &t) const {
        return (hours * 60 + minutes) - (t.hours * 60 + t.minutes);
    }
    
    // Time formatting
    std::string format12Hour() const {
        std::ostringstream oss;
        int displayHours = hours % 12;
        if (displayHours == 0) displayHours = 12;
        oss << std::setw(2) << std::setfill('0') << displayHours << ":"
            << std::setw(2) << std::setfill('0') << minutes
            << (hours < 12 ? " AM" : " PM");
        return oss.str();
    }
};
