#include <vector>
#include <string>
#include <map>
#include <algorithm>
#include <stdexcept>

using namespace std;

class UnionFind {
private:
    struct Node {
        string name;
        Node* parent;
        int rank;
        int size;
        
        Node(string n) : name(n), parent(nullptr), rank(0), size(1) {}
    };

    map<string, Node*> nodes;

public:
    // Insert a new node into the structure
    bool insertNode(const string& name) {
        if (nodes.find(name) != nodes.end()) {
            return false;
        }
        nodes[name] = new Node(name);
        return true;
    }

    // Find the root of a node with path compression
    Node* findSet(const string& name) {
        if (nodes.find(name) == nodes.end()) {
            return nullptr;
        }

        vector<Node*> path;
        Node* current = nodes[name];

        // Find the root
        while (current->parent != nullptr) {
            path.push_back(current);
            current = current->parent;
        }

        // Path compression
        for (Node* node : path) {
            node->parent = current;
        }

        return current;
    }

    // Union two sets using union by rank
    bool unionSets(const string& nameA, const string& nameB) {
        Node* rootA = findSet(nameA);
        Node* rootB = findSet(nameB);

        if (!rootA || !rootB) {
            return false;
        }

        if (rootA == rootB) {
            return true; // Already in same set
        }

        // Union by rank
        if (rootA->rank > rootB->rank) {
            rootB->parent = rootA;
            rootA->size += rootB->size;
        } else if (rootA->rank < rootB->rank) {
            rootA->parent = rootB;
            rootB->size += rootA->size;
        } else {
            rootB->parent = rootA;
            rootA->rank++;
            rootA->size += rootB->size;
        }

        return true;
    }

    // Check if two elements are connected
    bool areConnected(const string& nameA, const string& nameB) {
        Node* rootA = findSet(nameA);
        Node* rootB = findSet(nameB);
        
        if (!rootA || !rootB) {
            return false;
        }
        return rootA == rootB;
    }

    // Get all connected components
    map<string, vector<string>> getConnectedComponents() {
        map<Node*, vector<string>> componentMap;
        
        for (const auto& pair : nodes) {
            Node* root = findSet(pair.first);
            componentMap[root].push_back(pair.first);
        }

        map<string, vector<string>> result;
        for (const auto& pair : componentMap) {
            result[pair.second[0]] = pair.second; // Use first element as key
        }

        return result;
    }

    // Get the size of the component containing name
    int getComponentSize(const string& name) {
        Node* root = findSet(name);
        if (!root) return 0;
        return root->size;
    }

    ~UnionFind() {
        for (auto& pair : nodes) {
            delete pair.second;
        }
    }
};
