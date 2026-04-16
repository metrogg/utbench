#include <vector>
#include <string>
#include <map>
#include <ctime>
#include <stdexcept>

using namespace std;

class Session {
private:
    string id;
    map<string, string> options;
    string last_status;
    uint64_t revision;
    time_t created_at;
    time_t last_accessed;
    bool is_active;

public:
    Session() : revision(0), created_at(time(nullptr)), last_accessed(time(nullptr)), is_active(true) {}
    
    Session(const string &id, const map<string, string> &options, const string &status, uint64_t rev)
        : id(id), options(options), last_status(status), revision(rev),
          created_at(time(nullptr)), last_accessed(time(nullptr)), is_active(true) {}

    // Validate session based on timeout (in seconds)
    bool is_valid(int timeout_seconds = 3600) const {
        if (!is_active) return false;
        time_t now = time(nullptr);
        return (now - last_accessed) <= timeout_seconds;
    }

    // Update session access time
    void touch() {
        last_accessed = time(nullptr);
    }

    // Session upgrade with new options
    void upgrade(const map<string, string> &new_options, const string &new_status) {
        if (!is_active) {
            throw runtime_error("Cannot upgrade inactive session");
        }
        
        for (const auto &opt : new_options) {
            options[opt.first] = opt.second;
        }
        last_status = new_status;
        revision++;
        touch();
    }

    // Get session summary
    map<string, string> get_summary() const {
        return {
            {"id", id},
            {"status", last_status},
            {"revision", to_string(revision)},
            {"active", is_active ? "true" : "false"},
            {"created_at", to_string(created_at)},
            {"last_accessed", to_string(last_accessed)}
        };
    }

    // Getters and setters
    void set_id(const string &new_val) { id = new_val; touch(); }
    void set_options(const map<string, string> &new_val) { options = new_val; touch(); }
    void set_last_status(const string &new_val) { last_status = new_val; touch(); }
    void set_revision(uint64_t new_val) { revision = new_val; touch(); }
    void set_active(bool active) { is_active = active; touch(); }

    string get_id() const { return id; }
    map<string, string> get_options() const { return options; }
    string get_last_status() const { return last_status; }
    uint64_t get_revision() const { return revision; }
    bool get_active() const { return is_active; }
};

class SessionManager {
private:
    map<string, Session> sessions;
    int default_timeout;

public:
    SessionManager(int timeout = 3600) : default_timeout(timeout) {}

    // Create new session
    string create_session(const map<string, string> &options, const string &initial_status = "new") {
        string id = "sess_" + to_string(time(nullptr)) + "_" + to_string(rand() % 1000);
        sessions[id] = Session(id, options, initial_status, 1);
        return id;
    }

    // Get session if valid
    Session* get_session(const string &id) {
        auto it = sessions.find(id);
        if (it == sessions.end() || !it->second.is_valid(default_timeout)) {
            return nullptr;
        }
        it->second.touch();
        return &it->second;
    }

    // Clean up expired sessions
    void cleanup_sessions() {
        vector<string> to_remove;
        for (auto &pair : sessions) {
            if (!pair.second.is_valid(default_timeout)) {
                to_remove.push_back(pair.first);
            }
        }
        for (const string &id : to_remove) {
            sessions.erase(id);
        }
    }

    // Get all active sessions
    vector<map<string, string>> list_active_sessions() {
        vector<map<string, string>> result;
        for (auto &pair : sessions) {
            if (pair.second.is_valid(default_timeout)) {
                result.push_back(pair.second.get_summary());
            }
        }
        return result;
    }
};
