#include <vector>
#include <string>
#include <algorithm>
#include <ctime>
#include <stdexcept>
#include <map>
#include <iomanip>

using namespace std;

class User {
public:
    User() : fname(""), lname(""), uid_(""), pwd_(""), loginCount(0), lastLogin(0), isActive(false) {}

    User(string firstname, string lastname, string user_id, string password, int lcount = 0)
        : fname(firstname), lname(lastname), uid_(user_id), pwd_(password), 
          loginCount(lcount), lastLogin(0), isActive(true) {}

    // Getters
    string uid() const { return uid_; }
    string pwd() const { return pwd_; }
    int hours() const { return loginCount; }
    time_t getLastLogin() const { return lastLogin; }
    bool active() const { return isActive; }
    string fullName() const { return fname + " " + lname; }

    // Authentication methods
    bool login(string pwd) {
        if (!isActive) return false;
        if (pwd == pwd_) {
            loginCount++;
            lastLogin = time(nullptr);
            return true;
        }
        return false;
    }

    void deactivate() { isActive = false; }
    void activate() { isActive = true; }
    void changePassword(string oldPwd, string newPwd) {
        if (oldPwd == pwd_) {
            pwd_ = newPwd;
        } else {
            throw invalid_argument("Incorrect current password");
        }
    }

    // Static methods for user management
    static vector<User> makeUsersFrom(vector<string> userData) {
        vector<User> users;
        for (const auto& data : userData) {
            vector<string> fields = splitString(data, ',');
            if (fields.size() >= 5) {
                users.emplace_back(
                    fields[0], fields[1], fields[2], fields[3], stoi(fields[4])
                );
            }
        }
        return users;
    }

    static map<string, string> generatePasswordReport(const vector<User>& users) {
        map<string, string> report;
        for (const auto& user : users) {
            string strength = "weak";
            if (user.pwd_.length() >= 12) {
                strength = "strong";
            } else if (user.pwd_.length() >= 8) {
                strength = "medium";
            }
            report[user.uid()] = strength;
        }
        return report;
    }

    // Operator overloads
    friend bool operator<(const User& lhs, const User& rhs) {
        return lhs.uid_ < rhs.uid_;
    }

    friend ostream& operator<<(ostream& os, const User& user) {
        os << "User: " << user.fname << " " << user.lname 
           << " (" << user.uid_ << "), Logins: " << user.loginCount;
        if (!user.isActive) os << " [INACTIVE]";
        return os;
    }

private:
    string fname, lname;
    string uid_, pwd_;
    int loginCount;
    time_t lastLogin;
    bool isActive;

    static vector<string> splitString(const string& s, char delimiter) {
        vector<string> tokens;
        string token;
        istringstream tokenStream(s);
        while (getline(tokenStream, token, delimiter)) {
            tokens.push_back(token);
        }
        return tokens;
    }
};

class UserManager {
public:
    void addUser(const User& user) {
        if (users.find(user.uid()) != users.end()) {
            throw invalid_argument("User ID already exists");
        }
        users[user.uid()] = user;
    }

    User* getUser(const string& uid) {
        auto it = users.find(uid);
        return it != users.end() ? &it->second : nullptr;
    }

    void removeUser(const string& uid) {
        users.erase(uid);
    }

    vector<User> getActiveUsers() const {
        vector<User> active;
        for (const auto& pair : users) {
            if (pair.second.active()) {
                active.push_back(pair.second);
            }
        }
        return active;
    }

    void printAllUsers() const {
        cout << "=== User List ===" << endl;
        for (const auto& pair : users) {
            cout << pair.second << endl;
        }
    }

private:
    map<string, User> users;
};
