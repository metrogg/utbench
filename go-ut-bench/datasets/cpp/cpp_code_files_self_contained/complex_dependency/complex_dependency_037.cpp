#include <string>
#include <vector>
#include <algorithm>
#include <cctype>

using namespace std;

// Enhanced string processor that performs multiple operations
vector<string> process_string(const string& input) {
    vector<string> results;
    
    // 1. Original functionality - append greeting
    string with_greeting = input + "Hi";
    results.push_back("Greeting: " + with_greeting);
    
    // 2. Reverse the string
    string reversed = input;
    reverse(reversed.begin(), reversed.end());
    results.push_back("Reversed: " + reversed);
    
    // 3. Convert to uppercase
    string upper = input;
    transform(upper.begin(), upper.end(), upper.begin(), ::toupper);
    results.push_back("Uppercase: " + upper);
    
    // 4. Count vowels
    int vowel_count = 0;
    for (char c : input) {
        if (tolower(c) == 'a' || tolower(c) == 'e' || tolower(c) == 'i' || 
            tolower(c) == 'o' || tolower(c) == 'u') {
            vowel_count++;
        }
    }
    results.push_back("Vowel count: " + to_string(vowel_count));
    
    // 5. Check if palindrome
    string clean_input;
    for (char c : input) {
        if (isalpha(c)) {
            clean_input += tolower(c);
        }
    }
    string clean_reversed = clean_input;
    reverse(clean_reversed.begin(), clean_reversed.end());
    bool is_palindrome = (clean_input == clean_reversed);
    results.push_back("Is palindrome: " + string(is_palindrome ? "true" : "false"));
    
    return results;
}
