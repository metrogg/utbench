#include <vector>
#include <string>
#include <algorithm>
#include <cmath>
#include <sstream>
#include <stdexcept>

using namespace std;

class Pos {
public:
    int i;
    int j;

    // Constructors
    Pos() : i(0), j(0) {}
    Pos(int i, int j) : i(i), j(j) {}
    
    // Operator overloading
    bool operator==(const Pos& other) const {
        return i == other.i && j == other.j;
    }
    
    bool operator!=(const Pos& other) const {
        return !(*this == other);
    }
    
    bool operator<(const Pos& other) const {
        if (i != other.i) return i < other.i;
        return j < other.j;
    }
    
    Pos operator+(const Pos& other) const {
        return Pos(i + other.i, j + other.j);
    }
    
    Pos operator-(const Pos& other) const {
        return Pos(i - other.i, j - other.j);
    }
    
    // Distance calculations
    static int manhattan(const Pos& a, const Pos& b) {
        return abs(a.i - b.i) + abs(a.j - b.j);
    }
    
    static double euclidean(const Pos& a, const Pos& b) {
        return sqrt(pow(a.i - b.i, 2) + pow(a.j - b.j, 2));
    }
    
    static int chebyshev(const Pos& a, const Pos& b) {
        return max(abs(a.i - b.i), abs(a.j - b.j));
    }
    
    // Vector operations
    static int compare_vec(const vector<Pos>& veca, const vector<Pos>& vecb) {
        if (veca.size() != vecb.size()) {
            throw invalid_argument("Vectors must be of same size for comparison");
        }
        
        for (size_t k = 0; k < veca.size(); ++k) {
            if (veca[k] < vecb[k]) return -1;
            if (vecb[k] < veca[k]) return 1;
        }
        return 0;
    }
    
    // Path finding helpers
    static vector<Pos> get_neighbors(const Pos& pos, int grid_width, int grid_height) {
        vector<Pos> neighbors;
        vector<Pos> directions = {{0,1}, {1,0}, {0,-1}, {-1,0}}; // 4-directional
        
        for (const auto& dir : directions) {
            Pos neighbor = pos + dir;
            if (neighbor.i >= 0 && neighbor.i < grid_height && 
                neighbor.j >= 0 && neighbor.j < grid_width) {
                neighbors.push_back(neighbor);
            }
        }
        return neighbors;
    }
    
    // String representation
    string to_str() const {
        ostringstream oss;
        oss << "(" << i << "," << j << ")";
        return oss.str();
    }
    
    // Static factory method
    static Pos from_str(const string& str) {
        int x, y;
        char c;
        istringstream iss(str);
        if (!(iss >> c >> x >> c >> y >> c) || c != ')') {
            throw invalid_argument("Invalid position string format");
        }
        return Pos(x, y);
    }
};
