#include <mutex>
#include <vector>
#include <chrono>
#include <iomanip>
#include <sstream>
#include <memory>
#include <map>

using namespace std;

// Log level enumeration
enum class LogLevel {
    TRACE,
    DEBUG,
    INFO,
    WARNING,
    ERROR,
    CRITICAL
};

// Convert log level to string
string logLevelToString(LogLevel level) {
    static const map<LogLevel, string> levelStrings = {
        {LogLevel::TRACE, "TRACE"},
        {LogLevel::DEBUG, "DEBUG"},
        {LogLevel::INFO, "INFO"},
        {LogLevel::WARNING, "WARNING"},
        {LogLevel::ERROR, "ERROR"},
        {LogLevel::CRITICAL, "CRITICAL"}
    };
    return levelStrings.at(level);
}

// Message structure containing log information
struct LogMessage {
    LogLevel level;
    string content;
    chrono::system_clock::time_point timestamp;
    string source; // Optional source information
};

// Base formatter interface
class LogFormatter {
public:
    virtual ~LogFormatter() = default;
    virtual string format(const LogMessage& msg) = 0;
};

// Default formatter implementation
class DefaultFormatter : public LogFormatter {
public:
    string format(const LogMessage& msg) override {
        ostringstream oss;
        auto time = chrono::system_clock::to_time_t(msg.timestamp);
        oss << put_time(localtime(&time), "%Y-%m-%d %H:%M:%S") << " ["
            << logLevelToString(msg.level) << "] ";
        if (!msg.source.empty()) {
            oss << "(" << msg.source << ") ";
        }
        oss << msg.content << "\n";
        return oss.str();
    }
};

// Base channel interface
class LogChannel {
public:
    virtual ~LogChannel() = default;
    virtual void log(const LogMessage& msg) = 0;
    virtual void flush() = 0;
    virtual void setMinLevel(LogLevel level) = 0;
    virtual void setFormatter(unique_ptr<LogFormatter> formatter) = 0;
};

// Thread-safe ostream channel implementation
class OStreamChannel : public LogChannel {
public:
    OStreamChannel(ostream& os, LogLevel minLevel = LogLevel::INFO, 
                  bool forceFlush = false, 
                  unique_ptr<LogFormatter> formatter = make_unique<DefaultFormatter>())
        : os_(os), minLevel_(minLevel), forceFlush_(forceFlush), 
          formatter_(move(formatter)) {}

    void log(const LogMessage& msg) override {
        if (msg.level < minLevel_) return;
        
        lock_guard<mutex> lock(mutex_);
        os_ << formatter_->format(msg);
        if (forceFlush_) {
            os_.flush();
        }
    }

    void flush() override {
        lock_guard<mutex> lock(mutex_);
        os_.flush();
    }

    void setMinLevel(LogLevel level) override {
        lock_guard<mutex> lock(mutex_);
        minLevel_ = level;
    }

    void setFormatter(unique_ptr<LogFormatter> formatter) override {
        lock_guard<mutex> lock(mutex_);
        formatter_ = move(formatter);
    }

private:
    ostream& os_;
    LogLevel minLevel_;
    bool forceFlush_;
    unique_ptr<LogFormatter> formatter_;
    mutex mutex_;
};

// Simple logger class that uses channels
class Logger {
public:
    void addChannel(unique_ptr<LogChannel> channel) {
        channels_.push_back(move(channel));
    }

    void log(LogLevel level, const string& message, const string& source = "") {
        LogMessage msg{level, message, chrono::system_clock::now(), source};
        for (auto& channel : channels_) {
            channel->log(msg);
        }
    }

    void flush() {
        for (auto& channel : channels_) {
            channel->flush();
        }
    }

private:
    vector<unique_ptr<LogChannel>> channels_;
};
