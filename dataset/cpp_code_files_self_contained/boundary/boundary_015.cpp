#include <vector>
#include <cmath>
#include <stdexcept>
#include <map>
#include <iomanip>

using namespace std;

class Mesh2D {
private:
    string name;
    vector<float> vertices;
    vector<float> textureCoordinates;
    vector<unsigned int> indices;
    map<string, vector<float>> attributes;

public:
    Mesh2D(const string& meshName) : name(meshName) {}

    // Create a right triangle mesh with customizable size and position
    void createRightTriangle(float base, float height, float centerX = 0.0f, float centerY = 0.0f) {
        // Clear existing data
        vertices.clear();
        textureCoordinates.clear();
        indices.clear();

        // Calculate vertex positions relative to center
        float halfBase = base / 2.0f;
        float halfHeight = height / 2.0f;

        // Define vertices (counter-clockwise winding)
        vector<float> newVertices = {
            centerX - halfBase, centerY + halfHeight,  // Top-left
            centerX - halfBase, centerY - halfHeight,  // Bottom-left
            centerX + halfBase, centerY - halfHeight   // Bottom-right
        };

        // Default texture coordinates
        vector<float> newTexCoords = {
            0.0f, 0.0f,  // Top-left
            0.0f, 1.0f,  // Bottom-left
            1.0f, 1.0f   // Bottom-right
        };

        // Define indices (single triangle)
        vector<unsigned int> newIndices = {0, 1, 2};

        vertices = newVertices;
        textureCoordinates = newTexCoords;
        indices = newIndices;

        // Calculate and store mesh properties
        calculateMeshProperties();
    }

    // Calculate various mesh properties
    void calculateMeshProperties() {
        if (vertices.empty()) {
            throw runtime_error("Cannot calculate properties: no vertices defined");
        }

        // Calculate area
        float x1 = vertices[0], y1 = vertices[1];
        float x2 = vertices[2], y2 = vertices[3];
        float x3 = vertices[4], y3 = vertices[5];
        
        float area = abs((x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2)) / 2.0f);
        attributes["area"] = {area};

        // Calculate edge lengths
        float edge1 = sqrt(pow(x2-x1, 2) + pow(y2-y1, 2));
        float edge2 = sqrt(pow(x3-x2, 2) + pow(y3-y2, 2));
        float edge3 = sqrt(pow(x1-x3, 2) + pow(y1-y3, 2));
        attributes["edge_lengths"] = {edge1, edge2, edge3};

        // Calculate centroid
        float centroidX = (x1 + x2 + x3) / 3.0f;
        float centroidY = (y1 + y2 + y3) / 3.0f;
        attributes["centroid"] = {centroidX, centroidY};
    }

    // Print mesh information
    void printInfo() const {
        cout << "Mesh: " << name << endl;
        cout << "Vertices (" << vertices.size()/2 << "):" << endl;
        for (size_t i = 0; i < vertices.size(); i += 2) {
            cout << "  " << vertices[i] << ", " << vertices[i+1] << endl;
        }
        
        cout << "Texture Coordinates:" << endl;
        for (size_t i = 0; i < textureCoordinates.size(); i += 2) {
            cout << "  " << textureCoordinates[i] << ", " << textureCoordinates[i+1] << endl;
        }
        
        cout << "Indices (" << indices.size() << "):" << endl;
        for (size_t i = 0; i < indices.size(); i += 3) {
            if (i+2 < indices.size()) {
                cout << "  " << indices[i] << ", " << indices[i+1] << ", " << indices[i+2] << endl;
            }
        }

        cout << "Properties:" << endl;
        for (const auto& prop : attributes) {
            cout << "  " << prop.first << ": ";
            for (float val : prop.second) {
                cout << val << " ";
            }
            cout << endl;
        }
    }

    // Getter methods
    vector<float> getVertices() const { return vertices; }
    vector<float> getTextureCoordinates() const { return textureCoordinates; }
    vector<unsigned int> getIndices() const { return indices; }
    map<string, vector<float>> getAttributes() const { return attributes; }
};
