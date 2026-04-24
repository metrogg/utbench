#include <queue>
#include <tuple>
#include <vector>
#include <functional>
#include <limits>
#include <map>
#include <set>
#include <algorithm>

using namespace std;

const int INF = numeric_limits<int>::max();

class GraphAnalyzer {
private:
    int node_count;
    vector<vector<pair<int, int>>> adjacency_list;
    vector<int> hotel_nodes;
    vector<vector<int>> distance_matrix;
    int max_stay_duration;

    void compute_shortest_paths(int source) {
        distance_matrix[source].assign(node_count, INF);
        priority_queue<pair<int, int>, vector<pair<int, int>>, greater<pair<int, int>>> pq;
        pq.emplace(0, source);
        distance_matrix[source][source] = 0;

        while (!pq.empty()) {
            auto [current_dist, u] = pq.top();
            pq.pop();

            if (current_dist > distance_matrix[source][u]) continue;

            for (const auto& [v, weight] : adjacency_list[u]) {
                if (distance_matrix[source][v] > current_dist + weight) {
                    distance_matrix[source][v] = current_dist + weight;
                    pq.emplace(distance_matrix[source][v], v);
                }
            }
        }
    }

public:
    GraphAnalyzer(int nodes, const vector<vector<pair<int, int>>>& graph, 
                 const vector<int>& hotels, int stay_duration = 600)
        : node_count(nodes), adjacency_list(graph), hotel_nodes(hotels), 
          max_stay_duration(stay_duration) {
        
        // Ensure start (0) and end (n-1) are always in hotel list
        hotel_nodes.push_back(0);
        hotel_nodes.push_back(node_count - 1);
        sort(hotel_nodes.begin(), hotel_nodes.end());
        hotel_nodes.erase(unique(hotel_nodes.begin(), hotel_nodes.end()), hotel_nodes.end());

        distance_matrix.resize(node_count);
        for (int hotel : hotel_nodes) {
            compute_shortest_paths(hotel);
        }
    }

    int find_min_hotel_stays() {
        queue<pair<int, int>> q;
        vector<int> visited(node_count, -1);
        q.emplace(0, 0);
        visited[0] = 0;

        while (!q.empty()) {
            auto [u, stays] = q.front();
            q.pop();

            if (u == node_count - 1) {
                return stays - 1;
            }

            for (int hotel : hotel_nodes) {
                if (distance_matrix[u][hotel] <= max_stay_duration && visited[hotel] == -1) {
                    visited[hotel] = stays + 1;
                    q.emplace(hotel, stays + 1);
                }
            }
        }

        return -1;
    }

    vector<int> get_reachable_hotels(int node) const {
        vector<int> result;
        for (int hotel : hotel_nodes) {
            if (distance_matrix[node][hotel] <= max_stay_duration) {
                result.push_back(hotel);
            }
        }
        return result;
    }

    int get_distance(int from, int to) const {
        if (from < 0 || from >= node_count || to < 0 || to >= node_count) {
            return -1;
        }
        return distance_matrix[from][to];
    }
};
