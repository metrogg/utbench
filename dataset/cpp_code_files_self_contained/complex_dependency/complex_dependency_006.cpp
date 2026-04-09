#include <vector>
#include <queue>
#include <algorithm>
#include <climits>
#include <map>
#include <set>
using namespace std;

struct Edge {
    int u, v, weight;
    bool operator<(const Edge& other) const {
        return weight > other.weight; // For min-heap
    }
};

class Graph {
private:
    int V;
    vector<vector<Edge>> adj;
    vector<Edge> edges;

public:
    Graph(int vertices) : V(vertices), adj(vertices + 1) {}

    void addEdge(int u, int v, int weight) {
        Edge e1{u, v, weight};
        Edge e2{v, u, weight};
        adj[u].push_back(e1);
        adj[v].push_back(e2);
        edges.push_back(e1);
    }

    // Prim's algorithm for MST
    int primMST(int start = 1) {
        if (V == 0) return 0;

        vector<bool> inMST(V + 1, false);
        priority_queue<Edge> pq;
        int mstWeight = 0;

        inMST[start] = true;
        for (const Edge& e : adj[start]) {
            pq.push(e);
        }

        while (!pq.empty()) {
            Edge current = pq.top();
            pq.pop();

            if (inMST[current.v]) continue;

            inMST[current.v] = true;
            mstWeight += current.weight;

            for (const Edge& e : adj[current.v]) {
                if (!inMST[e.v]) {
                    pq.push(e);
                }
            }
        }

        // Check if MST includes all vertices
        for (int i = 1; i <= V; i++) {
            if (!inMST[i]) return -1; // Graph is disconnected
        }

        return mstWeight;
    }

    // Kruskal's algorithm for MST
    int kruskalMST() {
        vector<int> parent(V + 1);
        for (int i = 1; i <= V; i++) {
            parent[i] = i;
        }

        auto find = [&](int u) {
            while (parent[u] != u) {
                parent[u] = parent[parent[u]];
                u = parent[u];
            }
            return u;
        };

        auto unite = [&](int u, int v) {
            u = find(u);
            v = find(v);
            if (u != v) {
                parent[v] = u;
                return true;
            }
            return false;
        };

        sort(edges.begin(), edges.end(), [](const Edge& a, const Edge& b) {
            return a.weight < b.weight;
        });

        int mstWeight = 0, edgeCount = 0;
        for (const Edge& e : edges) {
            if (unite(e.u, e.v)) {
                mstWeight += e.weight;
                if (++edgeCount == V - 1) break;
            }
        }

        return edgeCount == V - 1 ? mstWeight : -1;
    }

    // Check if graph is connected
    bool isConnected() {
        if (V == 0) return true;
        vector<bool> visited(V + 1, false);
        queue<int> q;
        q.push(1);
        visited[1] = true;
        int count = 1;

        while (!q.empty()) {
            int u = q.front();
            q.pop();

            for (const Edge& e : adj[u]) {
                if (!visited[e.v]) {
                    visited[e.v] = true;
                    count++;
                    q.push(e.v);
                }
            }
        }

        return count == V;
    }
};
