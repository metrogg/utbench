#include <string>

using namespace std;

int add(int a, int b) {
    return a + b;
}

int multiply(int a, int b) {
    return a * b;
}

bool isPositive(int n) {
    return n > 0;
}

string greet(const string& name) {
    if (name.empty()) {
        return "Hello, World!";
    }
    return "Hello, " + name + "!";
}