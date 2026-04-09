#include <vector>
#include <algorithm>
#include <string>
#include <map>
#include <stdexcept>

using namespace std;

// Enhanced character sorter with multiple sorting strategies and validation
class CharacterSorter {
private:
    int char_count;
    vector<char> characters;
    map<char, int> char_weights;

    // Initialize character weights (simulating some external comparison logic)
    void initialize_weights() {
        for (int i = 0; i < char_count; ++i) {
            char_weights[characters[i]] = (i % 2 == 0) ? i * 2 : i * 3;
        }
    }

public:
    CharacterSorter(int n) : char_count(n) {
        if (n <= 0 || n > 26) {
            throw invalid_argument("Character count must be between 1 and 26");
        }
        
        // Initialize characters from 'A' to 'A'+n-1
        for (int i = 0; i < n; ++i) {
            characters.push_back('A' + i);
        }
        initialize_weights();
    }

    // Sort by ascending order (A, B, C...)
    string sort_ascending() {
        sort(characters.begin(), characters.end());
        return string(characters.begin(), characters.end());
    }

    // Sort by descending order (Z, Y, X...)
    string sort_descending() {
        sort(characters.begin(), characters.end(), greater<char>());
        return string(characters.begin(), characters.end());
    }

    // Sort by custom weight (defined in initialize_weights)
    string sort_by_weight() {
        sort(characters.begin(), characters.end(), 
            [this](char a, char b) { return char_weights[a] < char_weights[b]; });
        return string(characters.begin(), characters.end());
    }

    // Sort by alternating pattern (A, C, E... B, D, F...)
    string sort_alternating() {
        vector<char> even_chars, odd_chars;
        for (int i = 0; i < char_count; ++i) {
            if (i % 2 == 0) {
                even_chars.push_back(characters[i]);
            } else {
                odd_chars.push_back(characters[i]);
            }
        }
        sort(even_chars.begin(), even_chars.end());
        sort(odd_chars.begin(), odd_chars.end());
        
        characters.clear();
        characters.insert(characters.end(), even_chars.begin(), even_chars.end());
        characters.insert(characters.end(), odd_chars.begin(), odd_chars.end());
        
        return string(characters.begin(), characters.end());
    }

    // Get original character sequence
    string get_original() const {
        string result;
        for (int i = 0; i < char_count; ++i) {
            result += 'A' + i;
        }
        return result;
    }
};
