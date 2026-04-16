#include <sstream>
#include <vector>
#include <map>
#include <algorithm>
#include <cctype>

using namespace std;

// Function to process and analyze name strings
map<string, vector<string>> process_names(const string& greeting, const vector<string>& names) {
    map<string, vector<string>> result;
    
    // Process each name
    for (const auto& name : names) {
        vector<string> components;
        stringstream ss(name);
        string part;
        
        // Split name into components
        while (ss >> part) {
            components.push_back(part);
        }
        
        // Generate greetings
        vector<string> greetings;
        if (!components.empty()) {
            // Basic greeting
            greetings.push_back(greeting + ", " + name);
            
            // Formal greeting (last name only)
            if (components.size() > 1) {
                greetings.push_back(greeting + ", Mr./Ms. " + components.back());
            }
            
            // Initials version
            string initials;
            for (const auto& c : components) {
                if (!c.empty()) {
                    initials += toupper(c[0]);
                    initials += ". ";
                }
            }
            if (!initials.empty()) {
                initials.pop_back(); // Remove trailing space
                greetings.push_back(greeting + ", " + initials);
            }
            
            // Reversed name
            string reversed_name;
            for (auto it = components.rbegin(); it != components.rend(); ++it) {
                reversed_name += *it + " ";
            }
            if (!reversed_name.empty()) {
                reversed_name.pop_back(); // Remove trailing space
                greetings.push_back(greeting + ", " + reversed_name);
            }
        }
        
        result[name] = greetings;
    }
    
    return result;
}
