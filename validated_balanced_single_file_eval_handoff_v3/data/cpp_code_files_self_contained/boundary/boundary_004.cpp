#include <vector>
#include <string>
#include <iomanip>
#include <stdexcept>
#include <cmath>

using namespace std;

struct Time {
    int hours;
    int minutes;
    int seconds;

    // Constructor with validation
    Time(int h = 0, int m = 0, int s = 0) {
        if (h < 0 || m < 0 || s < 0 || m >= 60 || s >= 60) {
            throw invalid_argument("Invalid time values");
        }
        hours = h;
        minutes = m;
        seconds = s;
    }
};

// Check if time1 is later than time2
bool isLater(const Time &time1, const Time &time2) {
    if (time1.hours > time2.hours) return true;
    if (time1.hours < time2.hours) return false;
    
    if (time1.minutes > time2.minutes) return true;
    if (time1.minutes < time2.minutes) return false;
    
    return time1.seconds > time2.seconds;
}

// Convert Time to total seconds
int timeToSeconds(const Time &time) {
    return time.hours * 3600 + time.minutes * 60 + time.seconds;
}

// Convert seconds to Time with validation
Time secondsToTime(int totalSeconds) {
    if (totalSeconds < 0) {
        throw invalid_argument("Negative seconds not allowed");
    }
    
    // Handle overflow beyond 24 hours
    totalSeconds %= 86400; // 24*60*60
    
    int h = totalSeconds / 3600;
    int remaining = totalSeconds % 3600;
    int m = remaining / 60;
    int s = remaining % 60;
    
    return Time(h, m, s);
}

// Add seconds to a time and return new time
Time addSeconds(const Time &time, int seconds) {
    return secondsToTime(timeToSeconds(time) + seconds);
}

// Calculate time difference (time1 - time2)
Time timeDifference(const Time &time1, const Time &time2) {
    int diff = timeToSeconds(time1) - timeToSeconds(time2);
    if (diff < 0) diff += 86400; // Handle negative difference (next day)
    return secondsToTime(diff);
}

// Format time as HH:MM:SS
string formatTime(const Time &time) {
    char buffer[9];
    snprintf(buffer, sizeof(buffer), "%02d:%02d:%02d", 
             time.hours, time.minutes, time.seconds);
    return string(buffer);
}

// Calculate average of multiple times
Time averageTime(const vector<Time> &times) {
    if (times.empty()) {
        throw invalid_argument("Empty time list");
    }
    
    long total = 0;
    for (const auto &t : times) {
        total += timeToSeconds(t);
    }
    
    return secondsToTime(round(static_cast<double>(total) / times.size()));
}
