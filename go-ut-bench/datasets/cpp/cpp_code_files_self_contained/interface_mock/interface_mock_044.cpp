#include <vector>
#include <queue>
#include <algorithm>
#include <map>
#include <iomanip>

using namespace std;

class GraphAnalyzer {
private:
    vector<vector<int>> adjacencyMatrix;
    int vertexCount;

public:
    GraphAnalyzer(int v) : vertexCount(v), adjacencyMatrix(v, vector<int>(v, 0)) {}

    void addEdge(int f, int s) {
        adjacencyMatrix[f][s] = 1;
        adjacencyMatrix[s][f] = 1;
    }

    bool isConnected() {
        vector<bool> visited(vertexCount, false);
        dfs(0, visited);
        return all_of(visited.begin(), visited.end(), [](bool v) { return v; });
    }

    vector<int> getConnectedComponents() {
        vector<int> components(vertexCount, -1);
        int currentComponent = 0;
        
        for (int i = 0; i < vertexCount; ++i) {
            if (components[i] == -1) {
                queue<int> q;
                q.push(i);
                components[i] = currentComponent;
                
                while (!q.empty()) {
                    int current = q.front();
                    q.pop();
                    
                    for (int neighbor = 0; neighbor < vertexCount; ++neighbor) {
                        if (adjacencyMatrix[current][neighbor] && components[neighbor] == -1) {
                            components[neighbor] = currentComponent;
                            q.push(neighbor);
                        }
                    }
                }
                currentComponent++;
            }
        }
        return components;
    }

    map<int, int> getDegreeDistribution() {
        map<int, int> degreeCount;
        for (int i = 0; i < vertexCount; ++i) {
            int degree = 0;
            for (int j = 0; j < vertexCount; ++j) {
                if (adjacencyMatrix[i][j]) degree++;
            }
            degreeCount[degree]++;
        }
        return degreeCount;
    }

    bool hasCycle() {
        vector<bool> visited(vertexCount, false);
        for (int i = 0; i < vertexCount; ++i) {
            if (!visited[i] && hasCycleUtil(i, visited, -1)) {
                return true;
            }
        }
        return false;
    }

private:
    void dfs(int sv, vector<bool>& visited) {
        visited[sv] = true;
        for (int i = 0; i < vertexCount; ++i) {
            if (adjacencyMatrix[sv][i] && !visited[i]) {
                dfs(i, visited);
            }
        }
    }

    bool hasCycleUtil(int v, vector<bool>& visited, int parent) {
        visited[v] = true;
        for (int i = 0; i < vertexCount; ++i) {
            if (adjacencyMatrix[v][i]) {
                if (!visited[i]) {
                    if (hasCycleUtil(i, visited, v)) {
                        return true;
                    }
                } else if (i != parent) {
                    return true;
                }
            }
        }
        return false;
    }
};
