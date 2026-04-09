#include <string>
#include <vector>
#include <map>
#include <memory>

using namespace std;

class Action {
public:
    Action(const string& name, const string& type) : name(name), type(type) {}
    string get_name() const { return name; }
    string get_type() const { return type; }
private:
    string name;
    string type;
};

class StatusEffect {
protected:
    size_t duration;
    string name;
    string description;
    vector<string> immunities;
    map<string, string> actionInteractions;

public:
    StatusEffect(size_t duration, const string& name) 
        : duration(duration), name(name), description("No description available") {}

    virtual ~StatusEffect() = default;

    size_t get_duration() const { return duration; }
    string get_name() const { return name; }
    string get_description() const { return description; }
    
    void set_duration(size_t time) { duration = time; }
    void set_description(const string& desc) { description = desc; }
    
    void add_immunity(const string& actionType) {
        immunities.push_back(actionType);
    }
    
    void add_interaction(const string& actionType, const string& result) {
        actionInteractions[actionType] = result;
    }

    virtual void apply_effect() {
        if (duration > 0) {
            cout << "Applying " << name << " effect. Duration left: " << duration << endl;
            duration--;
        }
    }

    virtual bool check_immunity(const Action& action) const {
        for (const auto& immunity : immunities) {
            if (action.get_type() == immunity) {
                return true;
            }
        }
        return false;
    }

    virtual string can_perform(const Action& action) const {
        auto it = actionInteractions.find(action.get_type());
        if (it != actionInteractions.end()) {
            return it->second;
        }
        return "No special interaction";
    }

    virtual unique_ptr<StatusEffect> clone() const = 0;
};

class PoisonEffect : public StatusEffect {
public:
    PoisonEffect(size_t duration) : StatusEffect(duration, "Poison") {
        description = "Causes damage over time";
        add_immunity("Poison");
        add_interaction("Antidote", "Cured");
        add_interaction("Heal", "Effect reduced");
    }

    void apply_effect() override {
        StatusEffect::apply_effect();
        if (duration > 0) {
            cout << "Taking 5 poison damage!" << endl;
        }
    }

    unique_ptr<StatusEffect> clone() const override {
        return make_unique<PoisonEffect>(*this);
    }
};

class BurnEffect : public StatusEffect {
public:
    BurnEffect(size_t duration) : StatusEffect(duration, "Burn") {
        description = "Causes fire damage over time";
        add_immunity("Fire");
        add_interaction("Water", "Effect removed");
    }

    void apply_effect() override {
        StatusEffect::apply_effect();
        if (duration > 0) {
            cout << "Taking 3 burn damage!" << endl;
        }
    }

    unique_ptr<StatusEffect> clone() const override {
        return make_unique<BurnEffect>(*this);
    }
};

class FreezeEffect : public StatusEffect {
public:
    FreezeEffect(size_t duration) : StatusEffect(duration, "Freeze") {
        description = "Prevents movement and actions";
        add_interaction("Fire", "Effect removed");
        add_interaction("Physical", "No effect while frozen");
    }

    void apply_effect() override {
        StatusEffect::apply_effect();
        if (duration > 0) {
            cout << "Cannot move or act!" << endl;
        }
    }

    unique_ptr<StatusEffect> clone() const override {
        return make_unique<FreezeEffect>(*this);
    }
};
