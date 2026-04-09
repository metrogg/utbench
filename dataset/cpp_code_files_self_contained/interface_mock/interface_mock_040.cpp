#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <algorithm>

using namespace std;

class AdvancedUnionFind {
private:
    vector<int> parent;
    vector<int> rank;
    vector<int> size;
    int count;
    unordered_map<int, unordered_set<int>> groups;

public:
    AdvancedUnionFind(int n) {
        count = n;
        parent.resize(n);
        rank.resize(n, 0);
        size.resize(n, 1);
        
        for (int i = 0; i < n; ++i) {
            parent[i] = i;
            groups[i].insert(i);
        }
    }

    int Find(int p) {
        if (p != parent[p]) {
            parent[p] = Find(parent[p]); // Path compression
        }
        return parent[p];
    }

    void Union(int p, int q) {
        int rootP = Find(p);
        int rootQ = Find(q);
        
        if (rootP == rootQ) return;

        // Union by rank
        if (rank[rootP] > rank[rootQ]) {
            parent[rootQ] = rootP;
            size[rootP] += size[rootQ];
            groups[rootP].insert(groups[rootQ].begin(), groups[rootQ].end());
            groups.erase(rootQ);
        } else {
            parent[rootP] = rootQ;
            size[rootQ] += size[rootP];
            groups[rootQ].insert(groups[rootP].begin(), groups[rootP].end());
            groups.erase(rootP);
            if (rank[rootP] == rank[rootQ]) {
                rank[rootQ]++;
            }
        }
        count--;
    }

    int GetCount() const {
        return count;
    }

    int GetSize(int p) {
        return size[Find(p)];
    }

    const unordered_set<int>& GetGroup(int p) {
        return groups[Find(p)];
    }

    vector<vector<int>> GetAllGroups() {
        vector<vector<int>> result;
        for (const auto& pair : groups) {
            vector<int> group(pair.second.begin(), pair.second.end());
            sort(group.begin(), group.end());
            result.push_back(group);
        }
        sort(result.begin(), result.end());
        return result;
    }

    bool AreConnected(int p, int q) {
        return Find(p) == Find(q);
    }
};

class SocialNetworkAnalyzer {
public:
    static int findFriendCircles(vector<vector<int>>& friendshipMatrix) {
        int n = friendshipMatrix.size();
        AdvancedUnionFind uf(n);

        for (int i = 0; i < n; ++i) {
            for (int j = i + 1; j < n; ++j) {
                if (friendshipMatrix[i][j] == 1) {
                    uf.Union(i, j);
                }
            }
        }

        return uf.GetCount();
    }

    static vector<vector<int>> getFriendCircles(vector<vector<int>>& friendshipMatrix) {
        int n = friendshipMatrix.size();
        AdvancedUnionFind uf(n);

        for (int i = 0; i < n; ++i) {
            for (int j = i + 1; j < n; ++j) {
                if (friendshipMatrix[i][j] == 1) {
                    uf.Union(i, j);
                }
            }
        }

        return uf.GetAllGroups();
    }
};
