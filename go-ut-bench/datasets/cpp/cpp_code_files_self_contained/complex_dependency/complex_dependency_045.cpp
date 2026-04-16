#include <vector>
#include <tuple>
#include <set>
#include <thread>
#include <mutex>
#include <algorithm>
#include <stdexcept>

using namespace std;

class EnhancedVM {
private:
    vector<tuple<char, int>> program;
    set<tuple<int, int>> results;
    mutex results_mutex;
    vector<thread> threads;
    bool debug_mode;

    struct ExecutionState {
        set<int> visited_pcs;
        bool can_modify;
        int pc;
        int acc;
        int modifications_made;
    };

    void execute(ExecutionState state) {
        while (true) {
            // Check for infinite loop
            if (state.visited_pcs.count(state.pc)) {
                add_result(EXIT_FAILURE, state.acc);
                return;
            }

            // Check for successful termination
            if (state.pc >= program.size()) {
                add_result(EXIT_SUCCESS, state.acc);
                return;
            }

            state.visited_pcs.insert(state.pc);
            auto [opcode, arg] = program[state.pc];

            switch(opcode) {
                case 'n': {  // nop
                    if (state.can_modify && state.modifications_made < 2) {
                        // Branch with modified instruction (nop->jmp)
                        auto branch_state = state;
                        branch_state.can_modify = false;
                        branch_state.modifications_made++;
                        branch_state.pc += arg;
                        thread t(&EnhancedVM::execute, this, branch_state);
                        threads.push_back(move(t));
                    }
                    state.pc++;
                    break;
                }
                case 'a': {  // acc
                    state.acc += arg;
                    state.pc++;
                    break;
                }
                case 'j': {  // jmp
                    if (state.can_modify && state.modifications_made < 2) {
                        // Branch with modified instruction (jmp->nop)
                        auto branch_state = state;
                        branch_state.can_modify = false;
                        branch_state.modifications_made++;
                        branch_state.pc++;
                        thread t(&EnhancedVM::execute, this, branch_state);
                        threads.push_back(move(t));
                    }
                    state.pc += arg;
                    break;
                }
                default:
                    throw runtime_error("Invalid opcode");
            }
        }
    }

    void add_result(int status, int acc) {
        lock_guard<mutex> lock(results_mutex);
        results.insert({status, acc});
        if (debug_mode) {
            cout << "Result added: status=" << status << ", acc=" << acc << endl;
        }
    }

public:
    EnhancedVM(const vector<tuple<char, int>>& program, bool debug = false) 
        : program(program), debug_mode(debug) {}

    set<tuple<int, int>> run() {
        ExecutionState init_state{
            {},  // visited_pcs
            true,  // can_modify
            0,  // pc
            0,  // acc
            0   // modifications_made
        };

        execute(init_state);

        // Wait for all threads to complete
        for (auto& t : threads) {
            if (t.joinable()) {
                t.join();
            }
        }

        return results;
    }

    static vector<tuple<char, int>> parse_program(const vector<string>& instructions) {
        vector<tuple<char, int>> program;
        for (const auto& instr : instructions) {
            size_t space_pos = instr.find(' ');
            string op = instr.substr(0, space_pos);
            int arg = stoi(instr.substr(space_pos + 1));
            
            char opcode;
            if (op == "nop") opcode = 'n';
            else if (op == "acc") opcode = 'a';
            else if (op == "jmp") opcode = 'j';
            else throw runtime_error("Unknown instruction: " + op);
            
            program.emplace_back(opcode, arg);
        }
        return program;
    }
};
