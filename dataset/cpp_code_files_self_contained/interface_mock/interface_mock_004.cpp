#include <vector>
#include <string>
#include <algorithm>
#include <stdexcept>
#include <memory>

using namespace std;

// Define Role types as an enum for better type safety
enum class RoleType {
    STUDENT,
    TEACHER,
    ADMIN,
    GUEST,
    RESEARCHER
};

// Role class with enhanced functionality
class Role {
public:
    Role(RoleType type, const string& description = "") 
        : type(type), description(description) {}

    RoleType getType() const { return type; }
    string getDescription() const { return description; }
    void setDescription(const string& desc) { description = desc; }

    virtual string getRoleDetails() const {
        return "Role: " + roleTypeToString(type) + 
               ", Description: " + description;
    }

    static string roleTypeToString(RoleType type) {
        switch(type) {
            case RoleType::STUDENT: return "Student";
            case RoleType::TEACHER: return "Teacher";
            case RoleType::ADMIN: return "Administrator";
            case RoleType::GUEST: return "Guest";
            case RoleType::RESEARCHER: return "Researcher";
            default: return "Unknown";
        }
    }

private:
    RoleType type;
    string description;
};

// Enhanced Person class with validation and additional features
class Person {
public:
    Person(const string& cnp, const string& firstName, 
           const string& lastName, const string& email)
        : cnp(validateCNP(cnp)), 
          firstName(validateName(firstName, "First name")),
          lastName(validateName(lastName, "Last name")),
          email(validateEmail(email)) {}

    // Getters
    string getCNP() const { return cnp; }
    string getFirstName() const { return firstName; }
    string getLastName() const { return lastName; }
    string getEmail() const { return email; }
    string getFullName() const { return firstName + " " + lastName; }

    // Setters with validation
    void setCNP(const string& newCNP) { cnp = validateCNP(newCNP); }
    void setFirstName(const string& newName) { firstName = validateName(newName, "First name"); }
    void setLastName(const string& newName) { lastName = validateName(newName, "Last name"); }
    void setEmail(const string& newEmail) { email = validateEmail(newEmail); }

    // Role management
    void addRole(unique_ptr<Role> role) {
        roles.push_back(move(role));
    }

    bool hasRole(RoleType type) const {
        return any_of(roles.begin(), roles.end(), 
            [type](const unique_ptr<Role>& r) { return r->getType() == type; });
    }

    vector<string> getRoleDetails() const {
        vector<string> details;
        for (const auto& role : roles) {
            details.push_back(role->getRoleDetails());
        }
        return details;
    }

    // Display information
    void displayInfo() const {
        cout << "Person Information:\n";
        cout << "CNP: " << cnp << "\n";
        cout << "Name: " << getFullName() << "\n";
        cout << "Email: " << email << "\n";
        cout << "Roles:\n";
        for (const auto& detail : getRoleDetails()) {
            cout << "- " << detail << "\n";
        }
    }

private:
    string cnp;
    string firstName;
    string lastName;
    string email;
    vector<unique_ptr<Role>> roles;

    // Validation functions
    static string validateCNP(const string& cnp) {
        if (cnp.length() != 13 || !all_of(cnp.begin(), cnp.end(), ::isdigit)) {
            throw invalid_argument("CNP must be 13 digits");
        }
        return cnp;
    }

    static string validateName(const string& name, const string& fieldName) {
        if (name.empty()) {
            throw invalid_argument(fieldName + " cannot be empty");
        }
        if (any_of(name.begin(), name.end(), ::isdigit)) {
            throw invalid_argument(fieldName + " cannot contain digits");
        }
        return name;
    }

    static string validateEmail(const string& email) {
        if (email.find('@') == string::npos || email.find('.') == string::npos) {
            throw invalid_argument("Invalid email format");
        }
        return email;
    }
};
