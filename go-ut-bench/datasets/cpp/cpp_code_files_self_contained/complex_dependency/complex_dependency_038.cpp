#include <vector>
#include <utility>
#include <cmath>
#include <algorithm>
#include <map>

using namespace std;

// Structure to hold animal statistics
struct AnimalStats {
    int count;
    int consumption;
};

// Function to solve the animal counting problem with enhanced logic
pair<int, int> count_animals(int sheep_consumption, int goat_consumption, 
                            int total_animals, int total_consumption) {
    // Check for invalid inputs
    if (total_animals <= 0 || total_consumption <= 0) {
        return make_pair(-1, -1);
    }

    // Special case when both animals consume same amount
    if (sheep_consumption == goat_consumption) {
        if (total_animals == 2 && total_consumption == 2 * sheep_consumption) {
            return make_pair(1, 1);
        }
        return make_pair(-1, -1);
    }

    // Calculate possible solutions
    vector<pair<int, int>> solutions;
    for (int sheep = 1; sheep < total_animals; sheep++) {
        int goats = total_animals - sheep;
        if (sheep * sheep_consumption + goats * goat_consumption == total_consumption) {
            solutions.emplace_back(sheep, goats);
        }
    }

    // Return appropriate result based on number of solutions
    if (solutions.size() == 1) {
        return solutions[0];
    }
    return make_pair(-1, -1);
}

// Enhanced function that returns additional statistics
map<string, int> analyze_animal_counts(int sheep_consumption, int goat_consumption, 
                                      int total_animals, int total_consumption) {
    map<string, int> result;
    auto counts = count_animals(sheep_consumption, goat_consumption, 
                               total_animals, total_consumption);
    
    result["sheep_count"] = counts.first;
    result["goat_count"] = counts.second;
    result["is_valid"] = (counts.first != -1) ? 1 : 0;
    
    if (counts.first != -1) {
        result["total_sheep_consumption"] = counts.first * sheep_consumption;
        result["total_goat_consumption"] = counts.second * goat_consumption;
        result["consumption_ratio"] = (counts.first * sheep_consumption * 100) / total_consumption;
    }
    
    return result;
}
