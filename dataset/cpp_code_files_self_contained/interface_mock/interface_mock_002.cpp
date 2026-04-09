#include <vector>
#include <queue>
#include <stack>
#include <limits>
#include <stdexcept>
#include <iomanip>

using namespace std;

enum ResultCode { Success, Failure, Duplicate, NotPresent };

template <class T>
class ENode {
public:
    ENode() : nextArc(nullptr) {}
    ENode(int vertex, T weight, ENode* next) : adjVex(vertex), w(weight), nextArc(next) {}
    
    int adjVex;
    T w;
    ENode* nextArc;
};

template <class T>
class LGraph {
public:
    LGraph(int mSize);
    ~LGraph();
    
    // Basic graph operations
    ResultCode Insert(int u, int v, T w);
    ResultCode Remove(int u, int v);
    bool Exist(int u, int v) const;
    void Output() const;
    
    // Enhanced functionality
    vector<int> BFS(int start) const;
    vector<int> DFS(int start) const;
    vector<T> Dijkstra(int start) const;
    vector<vector<T>> Floyd() const;
    bool IsConnected() const;
    int GetConnectedComponents() const;
    vector<int> TopologicalSort() const;
    
    // Utility functions
    int VertexCount() const { return n; }
    int EdgeCount() const { return e; }
    vector<pair<int, T>> GetAdjacent(int u) const;
    
private:
    ENode<T>** a;
    int n, e;
    
    void DFSUtil(int v, vector<bool>& visited, vector<int>& result) const;
    void DFSUtil(int v, vector<bool>& visited) const;
};

template <class T>
LGraph<T>::LGraph(int mSize) : n(mSize), e(0) {
    if (n <= 0) throw invalid_argument("Graph size must be positive");
    a = new ENode<T>*[n];
    for (int i = 0; i < n; i++) a[i] = nullptr;
}

template <class T>
LGraph<T>::~LGraph() {
    for (int i = 0; i < n; i++) {
        ENode<T>* p = a[i];
        while (p) {
            ENode<T>* q = p;
            p = p->nextArc;
            delete q;
        }
    }
    delete[] a;
}

template <class T>
bool LGraph<T>::Exist(int u, int v) const {
    if (u < 0 || v < 0 || u >= n || v >= n || u == v) return false;
    ENode<T>* p = a[u];
    while (p && p->adjVex != v) p = p->nextArc;
    return p != nullptr;
}

template <class T>
ResultCode LGraph<T>::Insert(int u, int v, T w) {
    if (u < 0 || v < 0 || u >= n || v >= n || u == v) return Failure;
    if (Exist(u, v)) return Duplicate;
    a[u] = new ENode<T>(v, w, a[u]);
    e++;
    return Success;
}

template <class T>
ResultCode LGraph<T>::Remove(int u, int v) {
    if (u < 0 || v < 0 || u >= n || v >= n || u == v) return Failure;
    
    ENode<T>* p = a[u], *q = nullptr;
    while (p && p->adjVex != v) {
        q = p;
        p = p->nextArc;
    }
    
    if (!p) return NotPresent;
    if (q) q->nextArc = p->nextArc;
    else a[u] = p->nextArc;
    
    delete p;
    e--;
    return Success;
}

template <class T>
void LGraph<T>::Output() const {
    cout << "Adjacency List:" << endl;
    for (int i = 0; i < n; i++) {
        cout << i << ": ";
        for (ENode<T>* p = a[i]; p; p = p->nextArc) {
            cout << "-> [" << p->adjVex << ":" << p->w << "] ";
        }
        cout << endl;
    }
}

template <class T>
vector<int> LGraph<T>::BFS(int start) const {
    if (start < 0 || start >= n) throw out_of_range("Invalid start vertex");
    
    vector<int> result;
    vector<bool> visited(n, false);
    queue<int> q;
    
    q.push(start);
    visited[start] = true;
    
    while (!q.empty()) {
        int u = q.front();
        q.pop();
        result.push_back(u);
        
        for (ENode<T>* p = a[u]; p; p = p->nextArc) {
            int v = p->adjVex;
            if (!visited[v]) {
                visited[v] = true;
                q.push(v);
            }
        }
    }
    return result;
}

template <class T>
vector<int> LGraph<T>::DFS(int start) const {
    if (start < 0 || start >= n) throw out_of_range("Invalid start vertex");
    
    vector<int> result;
    vector<bool> visited(n, false);
    DFSUtil(start, visited, result);
    return result;
}

template <class T>
void LGraph<T>::DFSUtil(int v, vector<bool>& visited, vector<int>& result) const {
    visited[v] = true;
    result.push_back(v);
    
    for (ENode<T>* p = a[v]; p; p = p->nextArc) {
        int u = p->adjVex;
        if (!visited[u]) {
            DFSUtil(u, visited, result);
        }
    }
}

template <class T>
vector<T> LGraph<T>::Dijkstra(int start) const {
    if (start < 0 || start >= n) throw out_of_range("Invalid start vertex");
    
    vector<T> dist(n, numeric_limits<T>::max());
    vector<bool> visited(n, false);
    dist[start] = 0;
    
    priority_queue<pair<T, int>, vector<pair<T, int>>, greater<pair<T, int>>> pq;
    pq.push({0, start});
    
    while (!pq.empty()) {
        int u = pq.top().second;
        pq.pop();
        
        if (visited[u]) continue;
        visited[u] = true;
        
        for (ENode<T>* p = a[u]; p; p = p->nextArc) {
            int v = p->adjVex;
            T weight = p->w;
            
            if (dist[v] > dist[u] + weight) {
                dist[v] = dist[u] + weight;
                pq.push({dist[v], v});
            }
        }
    }
    return dist;
}

template <class T>
vector<pair<int, T>> LGraph<T>::GetAdjacent(int u) const {
    if (u < 0 || u >= n) throw out_of_range("Invalid vertex");
    
    vector<pair<int, T>> result;
    for (ENode<T>* p = a[u]; p; p = p->nextArc) {
        result.emplace_back(p->adjVex, p->w);
    }
    return result;
}
