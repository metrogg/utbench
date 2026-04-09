#include <string>
#include <vector>
#include <map>
#include <stdexcept>
#include <ctime>
#include <cstring>
#include <memory>
#include <sstream>

using namespace std;

// Enhanced MCC Message Types
enum class MCCMessageType {
    REGISTER,
    HEARTBEAT,
    TASK_REQUEST,
    TASK_RESULT,
    RESOURCE_QUERY,
    RESOURCE_RESPONSE,
    ERROR
};

// Enhanced MCC Message Base Class
class MCCMessage {
protected:
    MCCMessageType mType;
    time_t mTimestamp;
    string mAgentId;
    string mContent;

public:
    MCCMessage(MCCMessageType type, const string& agentId, const string& content = "")
        : mType(type), mTimestamp(time(nullptr)), mAgentId(agentId), mContent(content) {}

    virtual ~MCCMessage() = default;

    MCCMessageType type() const { return mType; }
    time_t timestamp() const { return mTimestamp; }
    const string& agentId() const { return mAgentId; }
    const string& content() const { return mContent; }

    virtual string serialize() const {
        stringstream ss;
        ss << static_cast<int>(mType) << "|"
           << mTimestamp << "|"
           << mAgentId.length() << "|" << mAgentId << "|"
           << mContent.length() << "|" << mContent;
        return ss.str();
    }

    virtual bool deserialize(const string& data) {
        try {
            size_t pos = 0;
            size_t next_pos = data.find('|', pos);
            int type = stoi(data.substr(pos, next_pos - pos));
            mType = static_cast<MCCMessageType>(type);
            pos = next_pos + 1;

            next_pos = data.find('|', pos);
            mTimestamp = stol(data.substr(pos, next_pos - pos));
            pos = next_pos + 1;

            next_pos = data.find('|', pos);
            size_t agentIdLen = stoul(data.substr(pos, next_pos - pos));
            pos = next_pos + 1;
            mAgentId = data.substr(pos, agentIdLen);
            pos += agentIdLen + 1;

            next_pos = data.find('|', pos);
            size_t contentLen = stoul(data.substr(pos, next_pos - pos));
            pos = next_pos + 1;
            mContent = data.substr(pos, contentLen);

            return true;
        } catch (...) {
            return false;
        }
    }

    virtual string toString() const {
        string typeStr;
        switch (mType) {
            case MCCMessageType::REGISTER: typeStr = "REGISTER"; break;
            case MCCMessageType::HEARTBEAT: typeStr = "HEARTBEAT"; break;
            case MCCMessageType::TASK_REQUEST: typeStr = "TASK_REQUEST"; break;
            case MCCMessageType::TASK_RESULT: typeStr = "TASK_RESULT"; break;
            case MCCMessageType::RESOURCE_QUERY: typeStr = "RESOURCE_QUERY"; break;
            case MCCMessageType::RESOURCE_RESPONSE: typeStr = "RESOURCE_RESPONSE"; break;
            case MCCMessageType::ERROR: typeStr = "ERROR"; break;
        }
        return "Type: " + typeStr + ", Agent: " + mAgentId + ", Time: " + to_string(mTimestamp) + ", Content: " + mContent;
    }
};

// MCC Server with enhanced functionality
class MCCServer {
private:
    map<int, string> mConnectedAgents; // fd -> agentId
    map<string, int> mAgentToFd;      // agentId -> fd
    map<string, time_t> mLastHeartbeat; // agentId -> last heartbeat time

public:
    virtual ~MCCServer() = default;

    // Process incoming message
    virtual void processMessage(int fd, const string& rawMessage) {
        unique_ptr<MCCMessage> message = parseMessage(rawMessage);
        if (!message) {
            sendError(fd, "Invalid message format");
            return;
        }

        switch (message->type()) {
            case MCCMessageType::REGISTER:
                handleRegister(fd, *message);
                break;
            case MCCMessageType::HEARTBEAT:
                handleHeartbeat(fd, *message);
                break;
            case MCCMessageType::TASK_REQUEST:
                handleTaskRequest(fd, *message);
                break;
            case MCCMessageType::TASK_RESULT:
                handleTaskResult(fd, *message);
                break;
            case MCCMessageType::RESOURCE_QUERY:
                handleResourceQuery(fd, *message);
                break;
            case MCCMessageType::RESOURCE_RESPONSE:
                handleResourceResponse(fd, *message);
                break;
            case MCCMessageType::ERROR:
                handleError(fd, *message);
                break;
        }
    }

    // Connection management
    virtual void onNewConnection(int fd) {
        cout << "New connection established: " << fd << endl;
    }

    virtual void onDisconnect(int fd) {
        auto it = mConnectedAgents.find(fd);
        if (it != mConnectedAgents.end()) {
            string agentId = it->second;
            mAgentToFd.erase(agentId);
            mLastHeartbeat.erase(agentId);
            mConnectedAgents.erase(fd);
            cout << "Agent disconnected: " << agentId << " (FD: " << fd << ")" << endl;
        } else {
            cout << "Unknown connection closed: " << fd << endl;
        }
    }

protected:
    unique_ptr<MCCMessage> parseMessage(const string& rawMessage) {
        unique_ptr<MCCMessage> message(new MCCMessage(MCCMessageType::ERROR, ""));
        if (message->deserialize(rawMessage)) {
            return message;
        }
        return nullptr;
    }

    void sendMessage(int fd, const MCCMessage& message) {
        string serialized = message.serialize();
        cout << "Sending message to FD " << fd << ": " << message.toString() << endl;
    }

    void sendError(int fd, const string& errorMsg) {
        MCCMessage error(MCCMessageType::ERROR, "", errorMsg);
        sendMessage(fd, error);
    }

    // Message handlers
    virtual void handleRegister(int fd, const MCCMessage& message) {
        if (mAgentToFd.find(message.agentId()) != mAgentToFd.end()) {
            sendError(fd, "Agent already registered");
            return;
        }

        mConnectedAgents[fd] = message.agentId();
        mAgentToFd[message.agentId()] = fd;
        mLastHeartbeat[message.agentId()] = time(nullptr);

        MCCMessage response(MCCMessageType::REGISTER, "server", "Registration successful");
        sendMessage(fd, response);
        cout << "Agent registered: " << message.agentId() << " (FD: " << fd << ")" << endl;
    }

    virtual void handleHeartbeat(int fd, const MCCMessage& message) {
        auto it = mConnectedAgents.find(fd);
        if (it == mConnectedAgents.end() || it->second != message.agentId()) {
            sendError(fd, "Unauthorized heartbeat");
            return;
        }

        mLastHeartbeat[message.agentId()] = time(nullptr);
        cout << "Heartbeat received from: " << message.agentId() << endl;
    }

    virtual void handleTaskRequest(int fd, const MCCMessage& message) {
        cout << "Task request from " << message.agentId() << ": " << message.content() << endl;
        // In a real implementation, this would dispatch the task to appropriate resources
        MCCMessage response(MCCMessageType::TASK_REQUEST, "server", "Task queued");
        sendMessage(fd, response);
    }

    virtual void handleTaskResult(int fd, const MCCMessage& message) {
        cout << "Task result from " << message.agentId() << ": " << message.content() << endl;
    }

    virtual void handleResourceQuery(int fd, const MCCMessage& message) {
        cout << "Resource query from " << message.agentId() << endl;
        // In a real implementation, this would query available resources
        MCCMessage response(MCCMessageType::RESOURCE_RESPONSE, "server", "CPU:4,Memory:8GB,Disk:100GB");
        sendMessage(fd, response);
    }

    virtual void handleResourceResponse(int fd, const MCCMessage& message) {
        cout << "Resource response from " << message.agentId() << ": " << message.content() << endl;
    }

    virtual void handleError(int fd, const MCCMessage& message) {
        cout << "Error received from " << message.agentId() << ": " << message.content() << endl;
    }
};
