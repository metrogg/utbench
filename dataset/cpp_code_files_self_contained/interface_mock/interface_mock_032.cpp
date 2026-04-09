#include <vector>
#include <chrono>
#include <string>
#include <map>
#include <iomanip>
#include <stdexcept>

using namespace std;
using namespace chrono;

class AdvancedTimer {
public:
    typedef void (*Callback)(double, void*);
    
    // Timer modes
    enum class TimerMode {
        HIGH_RESOLUTION,
        STEADY_CLOCK,
        SYSTEM_CLOCK
    };

    // Measurement units
    enum class TimeUnit {
        NANOSECONDS,
        MICROSECONDS,
        MILLISECONDS,
        SECONDS
    };

    // Constructor with optional callback and context
    AdvancedTimer(Callback pfn = nullptr, void* context = nullptr, 
                 TimerMode mode = TimerMode::HIGH_RESOLUTION,
                 TimeUnit unit = TimeUnit::MILLISECONDS)
        : m_callback(pfn), m_context(context), m_mode(mode), m_unit(unit) {
        Reset();
    }

    // Reset the timer
    void Reset() {
        m_isRunning = true;
        m_startTime = getCurrentTime();
        m_laps.clear();
        m_totalPausedTime = 0;
        m_pauseStartTime = 0;
    }

    // Stop the timer and return elapsed time
    double Stop() {
        if (!m_isRunning) return GetLastDelta();
        
        m_endTime = getCurrentTime();
        m_isRunning = false;
        m_delta = calculateDuration(m_startTime, m_endTime) - m_totalPausedTime;
        
        if (m_callback) {
            m_callback(m_delta, m_context);
        }
        return m_delta;
    }

    // Record a lap time
    double Lap() {
        if (!m_isRunning) return GetLastDelta();
        
        auto current = getCurrentTime();
        double lapTime = calculateDuration(m_lastLapTime ? m_lastLapTime : m_startTime, current) - 
                        (m_isPaused ? 0 : (current - m_pauseStartTime));
        
        m_laps.push_back(lapTime);
        m_lastLapTime = current;
        return lapTime;
    }

    // Pause the timer
    void Pause() {
        if (!m_isRunning || m_isPaused) return;
        
        m_pauseStartTime = getCurrentTime();
        m_isPaused = true;
    }

    // Resume the timer
    void Resume() {
        if (!m_isRunning || !m_isPaused) return;
        
        auto current = getCurrentTime();
        m_totalPausedTime += calculateDuration(m_pauseStartTime, current);
        m_isPaused = false;
    }

    // Get elapsed time without stopping
    double GetElapsed() const {
        if (!m_isRunning) return m_delta;
        
        auto current = getCurrentTime();
        return calculateDuration(m_startTime, current) - 
               (m_isPaused ? 0 : (current - m_pauseStartTime)) - 
               m_totalPausedTime;
    }

    // Get last measured delta
    double GetLastDelta() const {
        return m_delta;
    }

    // Get all lap times
    const vector<double>& GetLaps() const {
        return m_laps;
    }

    // Set the callback function
    void SetCallback(Callback pfn) {
        m_callback = pfn;
    }

    // Set time unit for output
    void SetTimeUnit(TimeUnit unit) {
        m_unit = unit;
    }

    // Get current timer statistics
    map<string, double> GetStatistics() const {
        map<string, double> stats;
        stats["last"] = m_delta;
        stats["elapsed"] = GetElapsed();
        
        if (!m_laps.empty()) {
            double sum = 0;
            double min = m_laps[0];
            double max = m_laps[0];
            
            for (double lap : m_laps) {
                sum += lap;
                if (lap < min) min = lap;
                if (lap > max) max = lap;
            }
            
            stats["laps_count"] = m_laps.size();
            stats["laps_total"] = sum;
            stats["laps_avg"] = sum / m_laps.size();
            stats["laps_min"] = min;
            stats["laps_max"] = max;
        }
        
        return stats;
    }

private:
    Callback m_callback;
    void* m_context;
    TimerMode m_mode;
    TimeUnit m_unit;
    
    bool m_isRunning = false;
    bool m_isPaused = false;
    double m_startTime = 0;
    double m_endTime = 0;
    double m_lastLapTime = 0;
    double m_pauseStartTime = 0;
    double m_totalPausedTime = 0;
    double m_delta = 0;
    vector<double> m_laps;

    // Get current time based on selected mode
    double getCurrentTime() const {
        switch (m_mode) {
            case TimerMode::HIGH_RESOLUTION:
                return duration_cast<nanoseconds>(high_resolution_clock::now().time_since_epoch()).count() / 1e9;
            case TimerMode::STEADY_CLOCK:
                return duration_cast<nanoseconds>(steady_clock::now().time_since_epoch()).count() / 1e9;
            case TimerMode::SYSTEM_CLOCK:
                return duration_cast<nanoseconds>(system_clock::now().time_since_epoch()).count() / 1e9;
            default:
                throw runtime_error("Invalid timer mode");
        }
    }

    // Calculate duration between two timestamps in selected unit
    double calculateDuration(double start, double end) const {
        double seconds = end - start;
        
        switch (m_unit) {
            case TimeUnit::NANOSECONDS: return seconds * 1e9;
            case TimeUnit::MICROSECONDS: return seconds * 1e6;
            case TimeUnit::MILLISECONDS: return seconds * 1e3;
            case TimeUnit::SECONDS: return seconds;
            default: return seconds;
        }
    }
};
