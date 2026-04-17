#include <string>
#include <vector>
#include <algorithm>
#include <map>

using namespace std;

// Function to check if two strings are rotated versions of each other
bool isRotatedString(const string& a, const string& b) {
    if (a.length() != b.length()) return false;
    if (a.empty() && b.empty()) return true;
    string doubled = a + a;
    return doubled.find(b) != string::npos;
}

// Enhanced function to find all possible rotation matches and their positions
map<int, string> findAllRotations(const string& a, const string& b) {
    map<int, string> results;
    
    if (a.length() != b.length()) {
        results[-1] = "Strings are of different lengths";
        return results;
    }
    
    if (a.empty() || b.empty()) {
        results[-1] = "Empty string input";
        return results;
    }
    
    string doubled = a + a;
    size_t pos = doubled.find(b);
    int rotation_count = 0;
    
    while (pos != string::npos && pos < a.length()) {
        int rotation = static_cast<int>(pos);
        results[rotation] = "Rotation found at position " + to_string(rotation);
        pos = doubled.find(b, pos + 1);
        rotation_count++;
    }
    
    if (rotation_count == 0) {
        results[-1] = "No rotation found";
    }
    
    return results;
}

// Function to check if strings are rotated by exactly k positions
bool isKRotatedString(const string& a, const string& b, int k) {
    if (a.length() != b.length()) return false;
    if (k < 0 || k >= a.length()) return false;
    
    string rotated;
    rotated = a.substr(k) + a.substr(0, k);
    return rotated == b;
}
