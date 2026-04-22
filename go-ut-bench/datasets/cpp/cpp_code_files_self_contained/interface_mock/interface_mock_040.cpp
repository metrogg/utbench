#include <string>
#include <vector>
#include <algorithm>
#include <map>
#include <cctype>

using namespace std;

// Processes a single river-water pair with enhanced functionality
map<string, string> process_river_water_pair(const string& river, const string& water) {
    map<string, string> result;
    
    // Process river name
    string processed_river = river;
    for (char &c : processed_river) {
        if (c == '.' || c == '-' || c == '/' || c == ' ') {
            c = '_';
        } else if (isupper(c)) {
            c = tolower(c);
        }
    }
    
    // Process water name with additional transformations
    string processed_water = water;
    for (char &c : processed_water) {
        if (c == '.' || c == '-' || c == '/' || c == ' ') {
            c = '_';
        } else if (isupper(c)) {
            c = tolower(c);
        }
    }
    
    // Remove consecutive underscores
    auto new_end = unique(processed_river.begin(), processed_river.end(), 
        [](char a, char b) { return a == '_' && b == '_'; });
    processed_river.erase(new_end, processed_river.end());
    
    new_end = unique(processed_water.begin(), processed_water.end(), 
        [](char a, char b) { return a == '_' && b == '_'; });
    processed_water.erase(new_end, processed_water.end());
    
    // Trim underscores from start and end
    if (!processed_river.empty() && processed_river.front() == '_') {
        processed_river.erase(0, 1);
    }
    if (!processed_river.empty() && processed_river.back() == '_') {
        processed_river.pop_back();
    }
    
    if (!processed_water.empty() && processed_water.front() == '_') {
        processed_water.erase(0, 1);
    }
    if (!processed_water.empty() && processed_water.back() == '_') {
        processed_water.pop_back();
    }
    
    result["river"] = processed_river;
    result["water"] = processed_water;
    result["combined"] = processed_river + " -> " + processed_water;
    
    return result;
}

// Processes multiple river-water pairs
vector<map<string, string>> process_river_water_pairs(const vector<pair<string, string>>& pairs) {
    vector<map<string, string>> results;
    
    for (const auto& pair : pairs) {
        results.push_back(process_river_water_pair(pair.first, pair.second));
    }
    
    return results;
}
