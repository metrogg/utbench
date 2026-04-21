#include <vector>
#include <stdexcept>
#include <sstream>
#include <string>

class Stack {
private:
    std::vector<int> items;
    
public:
    void push(int item) {
        items.push_back(item);
    }
    
    int pop() {
        if (isEmpty()) {
            throw std::runtime_error("Stack is empty");
        }
        int item = items.back();
        items.pop_back();
        return item;
    }
    
    int peek() {
        if (isEmpty()) {
            throw std::runtime_error("Stack is empty");
        }
        return items.back();
    }
    
    bool isEmpty() {
        return items.empty();
    }
    
    int size() {
        return static_cast<int>(items.size());
    }
};

int evaluatePostfix(const std::string& expression) {
    Stack stack;
    std::istringstream iss(expression);
    std::string token;
    
    while (iss >> token) {
        if (isdigit(token[0])) {
            stack.push(std::stoi(token));
        } else {
            int b = stack.pop();
            int a = stack.pop();
            int result;
            if (token == "+") {
                result = a + b;
            } else if (token == "-") {
                result = a - b;
            } else if (token == "*") {
                result = a * b;
            } else if (token == "/") {
                result = a / b;
            } else {
                throw std::invalid_argument("Unknown operator: " + token);
            }
            stack.push(result);
        }
    }
    
    return stack.pop();
}