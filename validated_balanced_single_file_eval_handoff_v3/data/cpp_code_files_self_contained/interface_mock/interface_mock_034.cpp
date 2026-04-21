#include <vector>
#include <cmath>
#include <stdexcept>
#include <memory>

using namespace std;

// Simple 3D vector structure
struct Vec3 {
    float x, y, z;
    
    Vec3(float x = 0, float y = 0, float z = 0) : x(x), y(y), z(z) {}
    
    bool operator==(const Vec3& other) const {
        return x == other.x && y == other.y && z == other.z;
    }
};

// Height grid structure
struct HeightGrid {
    int sizeX, sizeZ;
    vector<vector<Vec3>> data;
    
    HeightGrid(int x, int z) : sizeX(x), sizeZ(z), data(x, vector<Vec3>(z)) {}
    
    Vec3& at(int x, int z) { 
        if (x < 0 || x >= sizeX || z < 0 || z >= sizeZ)
            throw out_of_range("Grid index out of range");
        return data[x][z]; 
    }
    
    const Vec3& at(int x, int z) const { 
        if (x < 0 || x >= sizeX || z < 0 || z >= sizeZ)
            throw out_of_range("Grid index out of range");
        return data[x][z]; 
    }
};

// Model structure containing mesh data
struct Model {
    vector<float> vertices;
    vector<float> normals;
    vector<float> texCoords;
    vector<unsigned int> indices;
};

// Helper functions for vector math
Vec3 vectorSub(const Vec3& a, const Vec3& b) {
    return {a.x - b.x, a.y - b.y, a.z - b.z};
}

Vec3 crossProduct(const Vec3& a, const Vec3& b) {
    return {
        a.y * b.z - a.z * b.y,
        a.z * b.x - a.x * b.z,
        a.x * b.y - a.y * b.x
    };
}

Vec3 vectorAdd(const Vec3& a, const Vec3& b) {
    return {a.x + b.x, a.y + b.y, a.z + b.z};
}

Vec3 scalarMult(const Vec3& v, float s) {
    return {v.x * s, v.y * s, v.z * s};
}

float vectorLength(const Vec3& v) {
    return sqrt(v.x*v.x + v.y*v.y + v.z*v.z);
}

Vec3 normalize(const Vec3& v) {
    float len = vectorLength(v);
    if (len == 0) return {0, 0, 0};
    return {v.x/len, v.y/len, v.z/len};
}

// Gets vertex from vertex array with bounds checking
Vec3 getVertex(const vector<float>& vertexArray, int x, int z, const HeightGrid& grid) {
    if (x >= 0 && x < grid.sizeX && z >= 0 && z < grid.sizeZ) {
        int idx = (x + z * grid.sizeX) * 3;
        return {
            vertexArray[idx],
            vertexArray[idx + 1],
            vertexArray[idx + 2]
        };
    }
    return {0, 0, 0};
}

// Calculates normal for a vertex using surrounding vertices
Vec3 calculateNormal(const vector<float>& vertexArray, int x, int z, const HeightGrid& grid) {
    Vec3 center = getVertex(vertexArray, x, z, grid);
    
    // Get surrounding vertices
    Vec3 neighbors[6] = {
        getVertex(vertexArray, x, z-1, grid),
        getVertex(vertexArray, x+1, z-1, grid),
        getVertex(vertexArray, x+1, z, grid),
        getVertex(vertexArray, x, z+1, grid),
        getVertex(vertexArray, x-1, z+1, grid),
        getVertex(vertexArray, x-1, z, grid)
    };
    
    Vec3 normalSum = {0, 0, 0};
    
    // Calculate normals for each adjacent triangle
    for (int i = 0; i < 6; i++) {
        int next = (i + 1) % 6;
        if (!(neighbors[i] == Vec3(0, 0, 0)) && !(neighbors[next] == Vec3(0, 0, 0))) {
            Vec3 v1 = vectorSub(neighbors[i], center);
            Vec3 v2 = vectorSub(neighbors[next], center);
            normalSum = vectorAdd(normalSum, crossProduct(v2, v1));
        }
    }
    
    return normalize(normalSum);
}

// Generates terrain model from height grid
Model generateTerrain(const HeightGrid& heightGrid, bool generateZeroHeight = false, int resolution = 1) {
    if (resolution <= 0) throw invalid_argument("Resolution must be positive");
    
    const int vertexCount = heightGrid.sizeX * heightGrid.sizeZ;
    const int triangleCount = (heightGrid.sizeX - 1) * (heightGrid.sizeZ - 1) * 2;
    
    Model model;
    model.vertices.resize(3 * vertexCount);
    model.normals.resize(3 * vertexCount);
    model.texCoords.resize(2 * vertexCount);
    model.indices.resize(3 * triangleCount);
    
    // Generate vertices and texture coordinates
    for (int x = 0; x < heightGrid.sizeX; x++) {
        for (int z = 0; z < heightGrid.sizeZ; z++) {
            int idx = (x + z * heightGrid.sizeX) * 3;
            model.vertices[idx] = static_cast<float>(x) / resolution;
            model.vertices[idx + 1] = generateZeroHeight ? 0 : heightGrid.at(x, z).y;
            model.vertices[idx + 2] = static_cast<float>(z) / resolution;
            
            int texIdx = (x + z * heightGrid.sizeX) * 2;
            model.texCoords[texIdx] = static_cast<float>(x);
            model.texCoords[texIdx + 1] = static_cast<float>(z);
        }
    }
    
    // Calculate normals
    for (int x = 0; x < heightGrid.sizeX; x++) {
        for (int z = 0; z < heightGrid.sizeZ; z++) {
            Vec3 normal = calculateNormal(model.vertices, x, z, heightGrid);
            int idx = (x + z * heightGrid.sizeX) * 3;
            model.normals[idx] = normal.x;
            model.normals[idx + 1] = normal.y;
            model.normals[idx + 2] = normal.z;
        }
    }
    
    // Generate indices for triangles
    int index = 0;
    for (int x = 0; x < heightGrid.sizeX - 1; x++) {
        for (int z = 0; z < heightGrid.sizeZ - 1; z++) {
            // Triangle 1
            model.indices[index++] = x + z * heightGrid.sizeX;
            model.indices[index++] = x + (z + 1) * heightGrid.sizeX;
            model.indices[index++] = x + 1 + z * heightGrid.sizeX;
            
            // Triangle 2
            model.indices[index++] = x + 1 + z * heightGrid.sizeX;
            model.indices[index++] = x + (z + 1) * heightGrid.sizeX;
            model.indices[index++] = x + 1 + (z + 1) * heightGrid.sizeX;
        }
    }
    
    return model;
}
