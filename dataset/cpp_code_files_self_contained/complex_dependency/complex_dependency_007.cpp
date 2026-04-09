#include <vector>
#include <algorithm>
#include <map>
#include <set>
#include <numeric>
#include <queue>

using namespace std;

typedef pair<int, int> ii;

class TreeAnalyzer {
private:
    int n;
    vector<vector<int>> adj;
    vector<ii> guesses;
    vector<int> parent;
    vector<int> depth;
    vector<int> subtree_size;
    vector<int> correct_guess;
    vector<int> correct_reversed_guess;
    int total_correct;

    void reset() {
        adj.assign(n+1, vector<int>());
        guesses.clear();
        parent.assign(n+1, -1);
        depth.assign(n+1, 0);
        subtree_size.assign(n+1, 1);
        correct_guess.assign(n+1, 0);
        correct_reversed_guess.assign(n+1, 0);
        total_correct = 0;
    }

    void build_tree(const vector<ii>& edges) {
        for (const auto& edge : edges) {
            int u = edge.first;
            int v = edge.second;
            adj[u].push_back(v);
            adj[v].push_back(u);
        }
    }

    void dfs(int u) {
        for (int v : adj[u]) {
            if (v != parent[u]) {
                parent[v] = u;
                depth[v] = depth[u] + 1;
                dfs(v);
                subtree_size[u] += subtree_size[v];
            }
        }
    }

    void process_guesses() {
        sort(guesses.begin(), guesses.end());
        
        queue<int> q;
        vector<bool> visited(n+1, false);
        q.push(1);
        visited[1] = true;

        while (!q.empty()) {
            int u = q.front();
            q.pop();

            for (int v : adj[u]) {
                if (!visited[v]) {
                    correct_guess[v] = correct_guess[u];
                    correct_reversed_guess[v] = correct_reversed_guess[u];

                    if (binary_search(guesses.begin(), guesses.end(), ii(u, v))) {
                        total_correct++;
                        correct_guess[v]++;
                    }

                    if (binary_search(guesses.begin(), guesses.end(), ii(v, u))) {
                        correct_reversed_guess[v]++;
                    }

                    visited[v] = true;
                    q.push(v);
                }
            }
        }
    }

public:
    TreeAnalyzer(int nodes) : n(nodes) {
        reset();
    }

    void initialize_tree(const vector<ii>& edges, const vector<ii>& guess_pairs) {
        build_tree(edges);
        guesses = guess_pairs;
        dfs(1);
        process_guesses();
    }

    map<string, int> analyze_tree(int k) {
        map<string, int> results;
        int valid_roots = (total_correct >= k) ? 1 : 0;

        for (int i = 2; i <= n; i++) {
            if (total_correct - correct_guess[i] + correct_reversed_guess[i] >= k) {
                valid_roots++;
            }
        }

        int gcd_value = gcd(valid_roots, n);
        results["valid_roots"] = valid_roots;
        results["total_nodes"] = n;
        results["reduced_numerator"] = valid_roots / gcd_value;
        results["reduced_denominator"] = n / gcd_value;
        results["depth_root"] = depth[1];
        results["max_subtree"] = *max_element(subtree_size.begin(), subtree_size.end());

        return results;
    }

    vector<int> get_depths() {
        return depth;
    }

    vector<int> get_subtree_sizes() {
        return subtree_size;
    }
};
