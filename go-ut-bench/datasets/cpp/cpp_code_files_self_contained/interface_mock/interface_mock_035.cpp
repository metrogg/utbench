#include <vector>
#include <algorithm>
#include <sstream>
#include <map>
#include <iomanip>

using namespace std;

struct BridgeCrossingResult {
    int total_time;
    vector<string> steps;
    map<string, int> time_breakdown;
};

BridgeCrossingResult optimized_bridge_crossing(vector<int>& people) {
    BridgeCrossingResult result;
    result.total_time = 0;
    ostringstream oss;
    
    if (people.empty()) return result;
    
    sort(people.begin(), people.end());
    
    while (!people.empty()) {
        int n = people.size();
        
        if (n == 1) {
            oss << people[0] << " crosses alone\n";
            result.total_time += people[0];
            result.time_breakdown["single_crossing"] = people[0];
            people.clear();
        } 
        else if (n == 2) {
            oss << people[0] << " and " << people[1] << " cross together\n";
            result.total_time += max(people[0], people[1]);
            result.time_breakdown["pair_crossing"] = max(people[0], people[1]);
            people.clear();
        } 
        else if (n == 3) {
            // Special case for 3 people
            oss << people[0] << " and " << people[1] << " cross\n";
            oss << people[0] << " returns\n";
            oss << people[0] << " and " << people[2] << " cross\n";
            
            int time = people[0] + people[1] + people[2];
            result.total_time += time;
            result.time_breakdown["first_pair"] = people[1];
            result.time_breakdown["return"] = people[0];
            result.time_breakdown["second_pair"] = people[2];
            people.clear();
        } 
        else {
            int A1 = people[0];
            int A2 = people[1];
            int An_1 = people[n-2];
            int An = people[n-1];
            
            if (2*A2 < (A1 + An_1)) {
                // Strategy 1: Fastest two as runners
                oss << A1 << " and " << A2 << " cross\n";
                oss << A1 << " returns\n";
                oss << An_1 << " and " << An << " cross\n";
                oss << A2 << " returns\n";
                
                int time = A2 + A1 + An + A2;
                result.total_time += time;
                result.time_breakdown["fastest_pair"] += A2;
                result.time_breakdown["fastest_return"] += A1;
                result.time_breakdown["slowest_pair"] += An;
                result.time_breakdown["runner_return"] += A2;
                
                people.pop_back();
                people.pop_back();
            } 
            else {
                // Strategy 2: Fastest as runner
                oss << A1 << " and " << An << " cross\n";
                oss << A1 << " returns\n";
                oss << A1 << " and " << An_1 << " cross\n";
                oss << A1 << " returns\n";
                
                int time = An + A1 + An_1 + A1;
                result.total_time += time;
                result.time_breakdown["slowest_cross"] += An;
                result.time_breakdown["fastest_return"] += A1;
                result.time_breakdown["second_slowest_cross"] += An_1;
                result.time_breakdown["final_return"] += A1;
                
                people.pop_back();
                people.pop_back();
            }
        }
    }
    
    string steps = oss.str();
    steps = steps.substr(0, steps.size() - 1); // Remove last newline
    result.steps = vector<string>();
    istringstream iss(steps);
    string line;
    while (getline(iss, line)) {
        result.steps.push_back(line);
    }
    
    return result;
}
