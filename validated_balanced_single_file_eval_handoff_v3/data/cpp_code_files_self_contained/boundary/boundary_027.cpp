#include <string>
#include <vector>
#include <sstream>
#include <algorithm>
#include <iomanip>

using namespace std;

struct TextStats {
    int space_count;
    int word_count;
    int line_count;
    int char_count;
    string squeezed_text;
};

TextStats process_text(const string& input_text) {
    TextStats stats = {0, 0, 0, 0, ""};
    bool at_spaces = true;
    bool new_line = true;
    bool in_word = false;
    stringstream ss(input_text);
    string line;
    
    while (getline(ss, line)) {
        stats.line_count++;
        new_line = true;
        
        for (char c : line) {
            stats.char_count++;
            
            if (c == ' ') {
                if (at_spaces) {
                    stats.space_count++;
                } else {
                    if (in_word) {
                        stats.word_count++;
                        in_word = false;
                    }
                    at_spaces = true;
                    stats.squeezed_text += ' ';
                }
            } else {
                if (new_line) {
                    if (stats.space_count > 0) {
                        stringstream count_ss;
                        count_ss << setw(2) << stats.space_count;
                        stats.squeezed_text += count_ss.str() + " ";
                        stats.space_count = 0;
                    }
                    new_line = false;
                }
                
                if (at_spaces) {
                    at_spaces = false;
                    in_word = true;
                }
                stats.squeezed_text += c;
            }
        }
        
        // Handle word at end of line
        if (in_word) {
            stats.word_count++;
            in_word = false;
        }
        
        stats.squeezed_text += '\n';
        at_spaces = true;
    }
    
    return stats;
}
