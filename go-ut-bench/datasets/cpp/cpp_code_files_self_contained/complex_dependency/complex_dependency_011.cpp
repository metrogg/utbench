#include <vector>
#include <unordered_map>
#include <algorithm>
#include <climits>

using namespace std;

class ConstraintSolver {
private:
    int n; // Number of variables
    int m; // Number of constraints
    vector<vector<int>> constraints; // Constraints matrix
    vector<bool> solution; // Current solution
    bool solution_found;

    // Check if current partial solution satisfies all constraints
    bool is_valid(const vector<bool>& partial, int depth) {
        for (const auto& constraint : constraints) {
            bool satisfied = false;
            for (int i = 0; i <= depth; i++) {
                if (partial[i] && constraint[i] == 1) {
                    satisfied = true;
                    break;
                }
                if (!partial[i] && constraint[i] == -1) {
                    satisfied = true;
                    break;
                }
            }
            if (!satisfied) {
                // Check remaining variables (not yet assigned)
                for (int i = depth+1; i < n; i++) {
                    if (constraint[i] != 0) {
                        satisfied = true; // Potential to be satisfied later
                        break;
                    }
                }
                if (!satisfied) return false;
            }
        }
        return true;
    }

    // Recursive backtracking with pruning
    void backtrack(int var_index, int selected_count, int max_selected) {
        if (solution_found) return;
        if (selected_count > max_selected) return;
        if (!is_valid(solution, var_index-1)) return;
        
        if (var_index == n) {
            if (is_valid(solution, n-1)) {
                solution_found = true;
            }
            return;
        }

        // Try including this variable
        solution[var_index] = true;
        backtrack(var_index+1, selected_count+1, max_selected);
        
        if (solution_found) return;
        
        // Try excluding this variable
        solution[var_index] = false;
        backtrack(var_index+1, selected_count, max_selected);
    }

public:
    ConstraintSolver(int vars, int cons) : n(vars), m(cons), 
        constraints(cons, vector<int>(vars, 0)), solution(vars, false), solution_found(false) {}

    void add_constraint(int cons_idx, int var, int value) {
        constraints[cons_idx][var] = value;
    }

    // Solve with minimal number of selected variables
    pair<bool, vector<bool>> solve() {
        for (int k = 1; k <= n; k++) {
            solution_found = false;
            backtrack(0, 0, k);
            if (solution_found) {
                return {true, solution};
            }
        }
        return {false, vector<bool>(n, false)};
    }

    // Alternative solving method using heuristic
    pair<bool, vector<bool>> solve_heuristic() {
        vector<int> var_scores(n, 0);
        
        // Score variables based on how many constraints they satisfy
        for (int i = 0; i < n; i++) {
            for (const auto& constraint : constraints) {
                if (constraint[i] == 1 || constraint[i] == -1) {
                    var_scores[i]++;
                }
            }
        }
        
        // Sort variables by score (descending)
        vector<int> var_order(n);
        for (int i = 0; i < n; i++) var_order[i] = i;
        sort(var_order.begin(), var_order.end(), [&var_scores](int a, int b) {
            return var_scores[a] > var_scores[b];
        });
        
        // Greedy selection
        vector<bool> greedy_sol(n, false);
        vector<bool> constraint_satisfied(m, false);
        
        for (int var : var_order) {
            bool needed = false;
            for (int c = 0; c < m; c++) {
                if (!constraint_satisfied[c] && constraints[c][var] != 0) {
                    needed = true;
                    break;
                }
            }
            
            if (needed) {
                greedy_sol[var] = true;
                for (int c = 0; c < m; c++) {
                    if ((greedy_sol[var] && constraints[c][var] == 1) ||
                        (!greedy_sol[var] && constraints[c][var] == -1)) {
                        constraint_satisfied[c] = true;
                    }
                }
            }
        }
        
        // Verify solution
        bool valid = true;
        for (int c = 0; c < m; c++) {
            bool sat = false;
            for (int var = 0; var < n; var++) {
                if ((greedy_sol[var] && constraints[c][var] == 1) ||
                    (!greedy_sol[var] && constraints[c][var] == -1)) {
                    sat = true;
                    break;
                }
            }
            if (!sat) {
                valid = false;
                break;
            }
        }
        
        return {valid, greedy_sol};
    }
};
