#include <string>
#include <vector>
#include <map>
#include <algorithm>

using namespace std;

enum class ErrorSeverity {
    INFO,
    WARNING,
    ERROR,
    CRITICAL
};

enum class ErrorCategory {
    SYNTAX,
    SEMANTIC,
    RUNTIME,
    SYSTEM,
    NETWORK
};

class Error {
private:
    int error_num;
    string msg;
    int line;
    ErrorSeverity severity;
    ErrorCategory category;
    vector<string> context;

public:
    Error(int num, string _msg, int _line, 
          ErrorSeverity _severity = ErrorSeverity::ERROR, 
          ErrorCategory _category = ErrorCategory::SYNTAX)
        : error_num(num), msg(_msg), line(_line), 
          severity(_severity), category(_category) {}

    // Add contextual information to the error
    void add_context(const string& ctx) {
        context.push_back(ctx);
    }

    // Comparison by line number then by severity
    bool operator<(const Error& e) const {
        if (line != e.line) return line < e.line;
        return severity < e.severity;
    }

    // Equality comparison
    bool operator==(const Error& e) const {
        return error_num == e.error_num && line == e.line;
    }

    // Get formatted error message
    string format() const {
        string result = "[line: " + to_string(line) + 
                       ", error: " + to_string(error_num) + 
                       ", severity: " + severity_to_string() + 
                       ", category: " + category_to_string() + "] " + msg;
        
        if (!context.empty()) {
            result += "\nContext:";
            for (const auto& ctx : context) {
                result += "\n  - " + ctx;
            }
        }
        return result;
    }

    // Getters
    int get_line() const { return line; }
    ErrorSeverity get_severity() const { return severity; }
    ErrorCategory get_category() const { return category; }

private:
    string severity_to_string() const {
        switch(severity) {
            case ErrorSeverity::INFO: return "INFO";
            case ErrorSeverity::WARNING: return "WARNING";
            case ErrorSeverity::ERROR: return "ERROR";
            case ErrorSeverity::CRITICAL: return "CRITICAL";
            default: return "UNKNOWN";
        }
    }

    string category_to_string() const {
        switch(category) {
            case ErrorCategory::SYNTAX: return "SYNTAX";
            case ErrorCategory::SEMANTIC: return "SEMANTIC";
            case ErrorCategory::RUNTIME: return "RUNTIME";
            case ErrorCategory::SYSTEM: return "SYSTEM";
            case ErrorCategory::NETWORK: return "NETWORK";
            default: return "UNKNOWN";
        }
    }
};

class ErrorManager {
private:
    vector<Error> errors;
    map<int, string> error_descriptions;

public:
    ErrorManager() {
        // Initialize some common error descriptions
        error_descriptions = {
            {100, "Syntax error"},
            {101, "Missing semicolon"},
            {200, "Type mismatch"},
            {300, "Division by zero"},
            {400, "File not found"},
            {500, "Network timeout"}
        };
    }

    // Add a new error
    void add_error(const Error& err) {
        errors.push_back(err);
    }

    // Get all errors sorted by line number and severity
    vector<Error> get_sorted_errors() const {
        vector<Error> sorted = errors;
        sort(sorted.begin(), sorted.end());
        return sorted;
    }

    // Get errors filtered by severity
    vector<Error> get_errors_by_severity(ErrorSeverity severity) const {
        vector<Error> filtered;
        copy_if(errors.begin(), errors.end(), back_inserter(filtered),
               [severity](const Error& e) { return e.get_severity() == severity; });
        return filtered;
    }

    // Get the standard description for an error number
    string get_error_description(int error_num) const {
        auto it = error_descriptions.find(error_num);
        return it != error_descriptions.end() ? it->second : "Unknown error";
    }

    // Check if there are any errors of ERROR or higher severity
    bool has_serious_errors() const {
        return any_of(errors.begin(), errors.end(), 
                     [](const Error& e) { 
                         return e.get_severity() >= ErrorSeverity::ERROR; 
                     });
    }
};

// Overloaded output operator for Error class
ostream& operator<<(ostream& out, const Error& e) {
    out << e.format();
    return out;
}
