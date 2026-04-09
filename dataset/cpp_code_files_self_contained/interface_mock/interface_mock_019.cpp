#include <unordered_map>
#include <string>
#include <functional>
#include <vector>
#include <memory>

using namespace std;

class Button {
public:
    enum class State { RELEASED, PRESSED, HELD, JUST_RELEASED };

    Button() : m_state(State::RELEASED), m_framesHeld(0) {}

    void updateFrame() {
        switch (m_state) {
            case State::PRESSED:
                m_state = State::HELD;
                m_framesHeld++;
                break;
            case State::HELD:
                m_framesHeld++;
                break;
            case State::JUST_RELEASED:
                m_state = State::RELEASED;
                m_framesHeld = 0;
                break;
            default:
                break;
        }
    }

    void press() {
        if (m_state == State::RELEASED) {
            m_state = State::PRESSED;
            m_framesHeld = 1;
        }
    }

    void release() {
        if (m_state == State::PRESSED || m_state == State::HELD) {
            m_state = State::JUST_RELEASED;
        }
    }

    State getState() const { return m_state; }
    int getFramesHeld() const { return m_framesHeld; }
    bool isPressed() const { return m_state == State::PRESSED; }
    bool isHeld() const { return m_state == State::HELD; }
    bool isReleased() const { return m_state == State::RELEASED; }

private:
    State m_state;
    int m_framesHeld;
};

class InputManager {
public:
    InputManager() = default;

    void updateAll() {
        for (auto& [name, button] : m_buttons) {
            button.updateFrame();
        }
    }

    void addButton(const string& name) {
        m_buttons.emplace(name, Button());
    }

    void pressButton(const string& name) {
        if (m_buttons.find(name) != m_buttons.end()) {
            m_buttons[name].press();
        }
    }

    void releaseButton(const string& name) {
        if (m_buttons.find(name) != m_buttons.end()) {
            m_buttons[name].release();
        }
    }

    Button::State getButtonState(const string& name) const {
        if (m_buttons.find(name) != m_buttons.end()) {
            return m_buttons.at(name).getState();
        }
        return Button::State::RELEASED;
    }

    int getButtonFramesHeld(const string& name) const {
        if (m_buttons.find(name) != m_buttons.end()) {
            return m_buttons.at(name).getFramesHeld();
        }
        return 0;
    }

    void addButtonCombo(const vector<string>& buttons, const string& comboName, 
                       function<void()> callback) {
        m_combos[comboName] = {buttons, callback, false};
    }

    void checkCombos() {
        for (auto& [comboName, combo] : m_combos) {
            bool allPressed = true;
            for (const auto& btnName : combo.buttons) {
                if (getButtonState(btnName) != Button::State::PRESSED && 
                    getButtonState(btnName) != Button::State::HELD) {
                    allPressed = false;
                    break;
                }
            }

            if (allPressed && !combo.wasTriggered) {
                combo.callback();
                combo.wasTriggered = true;
            } else if (!allPressed) {
                combo.wasTriggered = false;
            }
        }
    }

private:
    struct Combo {
        vector<string> buttons;
        function<void()> callback;
        bool wasTriggered;
    };

    unordered_map<string, Button> m_buttons;
    unordered_map<string, Combo> m_combos;
};
