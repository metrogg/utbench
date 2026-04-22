#include <vector>
#include <string>
#include <map>
#include <algorithm>

using namespace std;

// Enhanced state representation with more properties
struct CharacterState {
    bool has_wood;
    bool has_axe;
    bool has_food;
    bool has_money;
    int energy;
    int wood_count;

    bool operator==(const CharacterState& other) const {
        return has_wood == other.has_wood &&
               has_axe == other.has_axe &&
               has_food == other.has_food &&
               has_money == other.has_money &&
               energy == other.energy &&
               wood_count == other.wood_count;
    }
};

// Base Action class template
template<typename T>
class Action {
public:
    virtual ~Action() = default;
    virtual bool can_run(const T& state) = 0;
    virtual void plan_effects(T& state) = 0;
    virtual bool execute(T& state) = 0;
    virtual int cost() const { return 1; }
    virtual string name() const = 0;
};

// Enhanced Goal class with multiple conditions
template<typename T>
class Goal {
public:
    virtual ~Goal() = default;
    virtual int distance_to(const T& state) const = 0;
    virtual bool is_satisfied(const T& state) const = 0;
    virtual string description() const = 0;
};

// Planner implementation
template<typename T>
class Planner {
public:
    vector<Action<T>*> find_plan(T start_state, Goal<T>& goal, vector<Action<T>*> actions, int max_depth = 10) {
        vector<Action<T>*> plan;
        if (goal.is_satisfied(start_state)) {
            return plan;
        }

        if (!find_plan_recursive(start_state, goal, actions, plan, max_depth)) {
            plan.clear();
        }
        return plan;
    }

private:
    bool find_plan_recursive(T current_state, Goal<T>& goal, vector<Action<T>*>& actions, 
                           vector<Action<T>*>& plan, int depth) {
        if (depth <= 0) return false;
        if (goal.is_satisfied(current_state)) return true;

        for (auto* action : actions) {
            if (!action->can_run(current_state)) continue;

            T new_state = current_state;
            action->plan_effects(new_state);

            plan.push_back(action);
            if (find_plan_recursive(new_state, goal, actions, plan, depth - 1)) {
                return true;
            }
            plan.pop_back();
        }
        return false;
    }
};

// Concrete Actions
class CutWood : public Action<CharacterState> {
public:
    bool can_run(const CharacterState& state) override {
        return state.has_axe && state.energy > 10;
    }

    void plan_effects(CharacterState& state) override {
        state.has_wood = true;
        state.wood_count += 1;
        state.energy -= 15;
    }

    bool execute(CharacterState& state) override {
        if (!can_run(state)) return false;
        plan_effects(state);
        return true;
    }

    string name() const override { return "Cut Wood"; }
};

class GrabAxe : public Action<CharacterState> {
public:
    bool can_run(const CharacterState& state) override {
        return !state.has_axe;
    }

    void plan_effects(CharacterState& state) override {
        state.has_axe = true;
    }

    bool execute(CharacterState& state) override {
        if (!can_run(state)) return false;
        plan_effects(state);
        return true;
    }

    string name() const override { return "Grab Axe"; }
};

class BuyFood : public Action<CharacterState> {
public:
    bool can_run(const CharacterState& state) override {
        return state.has_money && !state.has_food;
    }

    void plan_effects(CharacterState& state) override {
        state.has_food = true;
        state.has_money = false;
    }

    bool execute(CharacterState& state) override {
        if (!can_run(state)) return false;
        plan_effects(state);
        return true;
    }

    string name() const override { return "Buy Food"; }
};

class EatFood : public Action<CharacterState> {
public:
    bool can_run(const CharacterState& state) override {
        return state.has_food;
    }

    void plan_effects(CharacterState& state) override {
        state.has_food = false;
        state.energy = min(100, state.energy + 30);
    }

    bool execute(CharacterState& state) override {
        if (!can_run(state)) return false;
        plan_effects(state);
        return true;
    }

    string name() const override { return "Eat Food"; }
};

// Concrete Goals
class GatherWoodGoal : public Goal<CharacterState> {
public:
    int distance_to(const CharacterState& state) const override {
        if (state.wood_count >= 3) return 0;
        return 3 - state.wood_count;
    }

    bool is_satisfied(const CharacterState& state) const override {
        return state.wood_count >= 3;
    }

    string description() const override {
        return "Gather at least 3 pieces of wood";
    }
};

class SurviveGoal : public Goal<CharacterState> {
public:
    int distance_to(const CharacterState& state) const override {
        if (state.energy > 20) return 0;
        return 21 - state.energy;
    }

    bool is_satisfied(const CharacterState& state) const override {
        return state.energy > 20;
    }

    string description() const override {
        return "Maintain energy above 20";
    }
};
