#include <vector>
#include <array>
#include <map>
#include <algorithm>
#include <cmath>
#include <iomanip>

using namespace std;

class EnhancedMeshIntersector {
public:
    typedef array<double, 3> Vec3D;
    typedef array<size_t, 3> Face;

    struct Plane {
        Vec3D origin = {0, 0, 0};
        Vec3D normal = {0, 0, 1};
    };

    struct Path3D {
        vector<Vec3D> points;
        bool isClosed = false;
        
        double CalculateLength() const {
            double length = 0.0;
            for (size_t i = 1; i < points.size(); ++i) {
                double dx = points[i][0] - points[i-1][0];
                double dy = points[i][1] - points[i-1][1];
                double dz = points[i][2] - points[i-1][2];
                length += sqrt(dx*dx + dy*dy + dz*dz);
            }
            return length;
        }
        
        Vec3D CalculateCentroid() const {
            Vec3D centroid = {0, 0, 0};
            if (points.empty()) return centroid;
            
            for (const auto& point : points) {
                centroid[0] += point[0];
                centroid[1] += point[1];
                centroid[2] += point[2];
            }
            
            centroid[0] /= points.size();
            centroid[1] /= points.size();
            centroid[2] /= points.size();
            
            return centroid;
        }
    };

    struct Mesh {
        vector<Vec3D> vertices;
        vector<Face> faces;
        
        Mesh(const vector<Vec3D>& v, const vector<Face>& f) : vertices(v), faces(f) {}
        
        vector<Path3D> IntersectWithPlane(const Plane& plane) const {
            vector<Path3D> result;
            auto vertexOffsets = CalculateVertexOffsets(vertices, plane);
            auto edgePaths = FindCrossingEdgePaths(faces, vertexOffsets);
            ChainEdgePaths(edgePaths);
            result = Construct3DPaths(vertices, edgePaths, vertexOffsets);
            
            // Additional analysis
            for (auto& path : result) {
                path.isClosed = (path.points.front()[0] == path.points.back()[0] && 
                                path.points.front()[1] == path.points.back()[1] && 
                                path.points.front()[2] == path.points.back()[2]);
            }
            
            return result;
        }
        
        vector<Path3D> ClipWithPlane(const Plane& plane) const {
            vector<Path3D> result = IntersectWithPlane(plane);
            auto vertexOffsets = CalculateVertexOffsets(vertices, plane);
            auto freeEdges = FindFreeEdges(faces, vertexOffsets);
            auto freeEdgePaths = ProcessFreeEdges(freeEdges, vertexOffsets);
            
            for (const auto& path : freeEdgePaths) {
                Path3D newPath;
                for (const auto& edge : path) {
                    if (edge.first == edge.second) {
                        newPath.points.push_back(vertices[edge.first]);
                    } else {
                        double offset1 = vertexOffsets[edge.first];
                        double offset2 = vertexOffsets[edge.second];
                        double factor = offset1 / (offset1 - offset2);
                        Vec3D intersection;
                        for (int i = 0; i < 3; ++i) {
                            intersection[i] = vertices[edge.first][i] + 
                                            (vertices[edge.second][i] - vertices[edge.first][i]) * factor;
                        }
                        newPath.points.push_back(intersection);
                    }
                }
                result.push_back(newPath);
            }
            
            return result;
        }
        
    private:
        vector<double> CalculateVertexOffsets(const vector<Vec3D>& vertices, const Plane& plane) const {
            vector<double> offsets;
            for (const auto& vertex : vertices) {
                double offset = 0;
                for (int i = 0; i < 3; ++i) {
                    offset += plane.normal[i] * (vertex[i] - plane.origin[i]);
                }
                offsets.push_back(offset);
            }
            return offsets;
        }
        
        typedef pair<size_t, size_t> Edge;
        typedef vector<Edge> EdgePath;
        
        vector<EdgePath> FindCrossingEdgePaths(const vector<Face>& faces, const vector<double>& offsets) const {
            map<Edge, size_t> crossingFaces;
            for (const auto& face : faces) {
                bool edge1 = offsets[face[0]] * offsets[face[1]] < 0;
                bool edge2 = offsets[face[1]] * offsets[face[2]] < 0;
                
                if (edge1 || edge2) {
                    int oddVertex = edge2 - edge1 + 1;
                    bool oddIsHigher = offsets[face[oddVertex]] > 0;
                    size_t v0 = oddVertex + 1 + oddIsHigher;
                    if (v0 > 2) v0 -= 3;
                    size_t v2 = oddVertex + 2 - oddIsHigher;
                    if (v2 > 2) v2 -= 3;
                    
                    Edge key = {face[v0], face[oddVertex]};
                    if (key.first > key.second) swap(key.first, key.second);
                    crossingFaces[key] = face[v2];
                }
            }
            
            vector<EdgePath> edgePaths;
            while (!crossingFaces.empty()) {
                auto current = crossingFaces.begin();
                EdgePath path = {current->first};
                size_t closingVertex = current->second;
                crossingFaces.erase(current);
                
                while (true) {
                    Edge nextKey = {path.back().second, closingVertex};
                    if (nextKey.first > nextKey.second) swap(nextKey.first, nextKey.second);
                    
                    auto next = crossingFaces.find(nextKey);
                    if (next == crossingFaces.end()) {
                        swap(nextKey.first, nextKey.second);
                        next = crossingFaces.find(nextKey);
                        if (next == crossingFaces.end()) break;
                    }
                    
                    path.push_back(next->first);
                    closingVertex = next->second;
                    crossingFaces.erase(next);
                }
                
                path.push_back({path.back().second, closingVertex});
                edgePaths.push_back(path);
            }
            
            return edgePaths;
        }
        
        void ChainEdgePaths(vector<EdgePath>& edgePaths) const {
            if (edgePaths.empty()) return;
            
            vector<bool> used(edgePaths.size(), false);
            vector<EdgePath> chained;
            
            for (size_t i = 0; i < edgePaths.size(); ++i) {
                if (used[i]) continue;
                
                EdgePath chain = edgePaths[i];
                used[i] = true;
                bool changed;
                
                do {
                    changed = false;
                    for (size_t j = 0; j < edgePaths.size(); ++j) {
                        if (used[j]) continue;
                        
                        const EdgePath& path = edgePaths[j];
                        if (path.front() == chain.back()) {
                            chain.insert(chain.end(), path.begin() + 1, path.end());
                            used[j] = true;
                            changed = true;
                        } else if (path.back() == chain.back()) {
                            chain.insert(chain.end(), path.rbegin() + 1, path.rend());
                            used[j] = true;
                            changed = true;
                        } else if (path.back() == chain.front()) {
                            chain.insert(chain.begin(), path.begin(), path.end() - 1);
                            used[j] = true;
                            changed = true;
                        } else if (path.front() == chain.front()) {
                            chain.insert(chain.begin(), path.rbegin(), path.rend() - 1);
                            used[j] = true;
                            changed = true;
                        }
                    }
                } while (changed);
                
                chained.push_back(chain);
            }
            
            edgePaths = chained;
        }
        
        vector<Path3D> Construct3DPaths(const vector<Vec3D>& vertices, 
                                      const vector<EdgePath>& edgePaths,
                                      const vector<double>& offsets) const {
            vector<Path3D> paths;
            
            for (const auto& edgePath : edgePaths) {
                Path3D path;
                bool skipFirst = (edgePath.front() == edgePath.back());
                
                for (const auto& edge : edgePath) {
                    if (skipFirst) {
                        skipFirst = false;
                        continue;
                    }
                    
                    if (edge.first == edge.second) {
                        path.points.push_back(vertices[edge.first]);
                    } else {
                        double offset1 = offsets[edge.first];
                        double offset2 = offsets[edge.second];
                        double factor = offset1 / (offset1 - offset2);
                        
                        Vec3D intersection;
                        for (int i = 0; i < 3; ++i) {
                            intersection[i] = vertices[edge.first][i] + 
                                            (vertices[edge.second][i] - vertices[edge.first][i]) * factor;
                        }
                        path.points.push_back(intersection);
                    }
                }
                
                paths.push_back(path);
            }
            
            return paths;
        }
        
        vector<Edge> FindFreeEdges(const vector<Face>& faces, const vector<double>& offsets) const {
            map<Edge, int> edgeCounts;
            
            for (const auto& face : faces) {
                for (int i = 0; i < 3; ++i) {
                    size_t v0 = face[i];
                    size_t v1 = face[(i + 1) % 3];
                    
                    if (offsets[v0] < 0 && offsets[v1] < 0) continue;
                    
                    Edge edge = {min(v0, v1), max(v0, v1)};
                    edgeCounts[edge]++;
                }
            }
            
            vector<Edge> freeEdges;
            for (const auto& entry : edgeCounts) {
                if (entry.second == 1) {
                    freeEdges.push_back(entry.first);
                }
            }
            
            return freeEdges;
        }
        
        vector<EdgePath> ProcessFreeEdges(const vector<Edge>& freeEdges, 
                                        const vector<double>& offsets) const {
            vector<bool> used(freeEdges.size(), false);
            vector<EdgePath> paths;
            
            for (size_t i = 0; i < freeEdges.size(); ++i) {
                if (used[i]) continue;
                
                EdgePath path;
                Edge current = freeEdges[i];
                used[i] = true;
                
                if (offsets[current.first] > 0) {
                    path.push_back({current.first, current.first});
                }
                if (offsets[current.second] > 0) {
                    path.push_back({current.second, current.second});
                }
                
                bool changed;
                do {
                    changed = false;
                    for (size_t j = 0; j < freeEdges.size(); ++j) {
                        if (used[j]) continue;
                        
                        Edge edge = freeEdges[j];
                        if (offsets[edge.first] > 0 && offsets[edge.second] > 0) {
                            if (edge.first == path.back().first) {
                                path.push_back({edge.second, edge.second});
                                used[j] = true;
                                changed = true;
                            } else if (edge.second == path.back().first) {
                                path.push_back({edge.first, edge.first});
                                used[j] = true;
                                changed = true;
                            }
                        }
                    }
                } while (changed);
                
                if (!path.empty()) {
                    paths.push_back(path);
                }
            }
            
            return paths;
        }
    };
};
