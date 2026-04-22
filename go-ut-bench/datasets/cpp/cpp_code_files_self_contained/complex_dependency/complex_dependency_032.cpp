#include <iomanip>
#include <string>
#include <vector>
#include <map>
#include <cctype>

using namespace std;

const int COL_WIDTH = 10;

// Enhanced function to analyze a character and return all its properties
map<string, int> analyze_character(char c) {
    map<string, int> properties;
    
    properties["isalnum"] = isalnum(c) ? 1 : 0;
    properties["isalpha"] = isalpha(c) ? 1 : 0;
    properties["iscntrl"] = iscntrl(c) ? 1 : 0;
    properties["isdigit"] = isdigit(c) ? 1 : 0;
    properties["isgraph"] = isgraph(c) ? 1 : 0;
    properties["islower"] = islower(c) ? 1 : 0;
    properties["isprint"] = isprint(c) ? 1 : 0;
    properties["ispunct"] = ispunct(c) ? 1 : 0;
    properties["isspace"] = isspace(c) ? 1 : 0;
    properties["isupper"] = isupper(c) ? 1 : 0;
    properties["isxdigit"] = isxdigit(c) ? 1 : 0;
    
    // Additional derived properties
    properties["to_upper"] = toupper(c);
    properties["to_lower"] = tolower(c);
    properties["ascii_value"] = static_cast<int>(c);
    
    return properties;
}

// Function to analyze a string and return statistics about its characters
map<string, int> analyze_string(const string& s) {
    map<string, int> stats;
    
    stats["total_chars"] = s.length();
    stats["alnum_chars"] = 0;
    stats["alpha_chars"] = 0;
    stats["digit_chars"] = 0;
    stats["lower_chars"] = 0;
    stats["upper_chars"] = 0;
    stats["space_chars"] = 0;
    stats["punct_chars"] = 0;
    stats["cntrl_chars"] = 0;
    
    for (char c : s) {
        if (isalnum(c)) stats["alnum_chars"]++;
        if (isalpha(c)) stats["alpha_chars"]++;
        if (isdigit(c)) stats["digit_chars"]++;
        if (islower(c)) stats["lower_chars"]++;
        if (isupper(c)) stats["upper_chars"]++;
        if (isspace(c)) stats["space_chars"]++;
        if (ispunct(c)) stats["punct_chars"]++;
        if (iscntrl(c)) stats["cntrl_chars"]++;
    }
    
    return stats;
}
