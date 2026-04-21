#include <vector>
#include <string>
#include <stdexcept>

class Logger {
public:
    virtual void log(const std::string& message) = 0;
    virtual std::vector<std::string> getLogs() = 0;
    virtual ~Logger() {}
};

class ConsoleLogger : public Logger {
private:
    std::vector<std::string> logs;
    
public:
    ConsoleLogger() : logs() {}
    
    void log(const std::string& message) override {
        logs.push_back(message);
    }
    
    std::vector<std::string> getLogs() override {
        return logs;
    }
};

class FileLogger : public Logger {
private:
    std::string filename;
    std::vector<std::string> logs;
    
public:
    FileLogger(const std::string& name) : filename(name), logs() {}
    
    void log(const std::string& message) override {
        logs.push_back("[" + filename + "] " + message);
    }
    
    std::vector<std::string> getLogs() override {
        return logs;
    }
};

class Application {
private:
    Logger* logger;
    
public:
    Application(Logger* log) : logger(log) {}
    
    std::string run(const std::string& taskName) {
        logger->log("Starting task: " + taskName);
        std::string result = "Task " + taskName + " completed";
        logger->log(result);
        return result;
    }
    
    std::vector<std::string> getHistory() {
        return logger->getLogs();
    }
};