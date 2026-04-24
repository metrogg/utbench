#include <vector>
#include <string>
#include <map>
#include <algorithm>
#include <stdexcept>
#include <memory>

using namespace std;

// Database connection class with enhanced functionality
class DatabaseConnection {
private:
    string name;
    bool isOpen;
    int connectionId;
    static int nextId;

public:
    DatabaseConnection(const string& name) : name(name), isOpen(false) {
        connectionId = ++nextId;
    }

    bool open() {
        if (isOpen) return false;
        isOpen = true;
        return true;
    }

    bool close() {
        if (!isOpen) return false;
        isOpen = false;
        return true;
    }

    bool executeQuery(const string& query) {
        if (!isOpen) return false;
        // Simulate query execution
        return true;
    }

    string getName() const { return name; }
    int getId() const { return connectionId; }
    bool isConnectionOpen() const { return isOpen; }
};

int DatabaseConnection::nextId = 0;

// Database connection manager with advanced features
class DatabaseManager {
private:
    vector<shared_ptr<DatabaseConnection>> connections;
    map<string, shared_ptr<DatabaseConnection>> connectionMap;

public:
    // Find connection by name and return its index or -1 if not found
    int findConnection(const string& name) const {
        auto it = connectionMap.find(name);
        if (it == connectionMap.end()) return -1;
        
        for (size_t i = 0; i < connections.size(); ++i) {
            if (connections[i]->getName() == name) {
                return static_cast<int>(i);
            }
        }
        return -1;
    }

    // Add a new connection
    bool addConnection(const string& name) {
        if (connectionMap.find(name) != connectionMap.end()) return false;
        
        auto conn = make_shared<DatabaseConnection>(name);
        connections.push_back(conn);
        connectionMap[name] = conn;
        return true;
    }

    // Open a connection by name
    bool openConnection(const string& name) {
        auto it = connectionMap.find(name);
        if (it == connectionMap.end()) return false;
        return it->second->open();
    }

    // Close a connection by name
    bool closeConnection(const string& name) {
        auto it = connectionMap.find(name);
        if (it == connectionMap.end()) return false;
        return it->second->close();
    }

    // Get connection status
    string getConnectionStatus(const string& name) const {
        auto it = connectionMap.find(name);
        if (it == connectionMap.end()) return "Not Found";
        return it->second->isConnectionOpen() ? "Open" : "Closed";
    }

    // Execute query on a connection
    bool executeQuery(const string& connName, const string& query) {
        auto it = connectionMap.find(connName);
        if (it == connectionMap.end()) return false;
        return it->second->executeQuery(query);
    }

    // Get all connection names
    vector<string> getAllConnectionNames() const {
        vector<string> names;
        for (const auto& conn : connections) {
            names.push_back(conn->getName());
        }
        return names;
    }
};
