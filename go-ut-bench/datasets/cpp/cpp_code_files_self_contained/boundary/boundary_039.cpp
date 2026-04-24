#include <vector>
#include <cmath>
#include <map>
#include <iomanip>

using namespace std;

struct vec3 {
    float x, y, z;
    vec3(float x = 0, float y = 0, float z = 0) : x(x), y(y), z(z) {}
};

struct vec2 {
    float x, y;
    vec2(float x = 0, float y = 0) : x(x), y(y) {}
};

struct Vertex {
    vec3 position;
    vec3 normal;
    vec2 texcoord;
    vec3 color;
};

struct MeshData {
    vector<Vertex> vertices;
    vector<unsigned int> indices;
    int triangleCount;
    int vertexCount;
};

class SphereGenerator {
public:
    static MeshData GenerateSphere(float radius = 1.0f, 
                                 int latitudeBands = 12, 
                                 int longitudeBands = 12,
                                 vec3 color = vec3(1.0f),
                                 bool smoothNormals = true) {
        MeshData mesh;
        const float PI = 3.14159265358979323846f;

        // Validate parameters
        if (radius <= 0) radius = 1.0f;
        if (latitudeBands < 3) latitudeBands = 3;
        if (longitudeBands < 3) longitudeBands = 3;

        // Generate vertices
        for (int lat = 0; lat <= latitudeBands; lat++) {
            float theta = lat * PI / latitudeBands;
            float sinTheta = sin(theta);
            float cosTheta = cos(theta);

            for (int lon = 0; lon <= longitudeBands; lon++) {
                float phi = lon * 2 * PI / longitudeBands;
                float sinPhi = sin(phi);
                float cosPhi = cos(phi);

                Vertex vertex;
                vertex.position.x = radius * cosPhi * sinTheta;
                vertex.position.y = radius * cosTheta;
                vertex.position.z = radius * sinPhi * sinTheta;

                if (smoothNormals) {
                    vertex.normal = vec3(cosPhi * sinTheta, cosTheta, sinPhi * sinTheta);
                }

                vertex.texcoord.x = 1.0f - (float)lon / longitudeBands;
                vertex.texcoord.y = 1.0f - (float)lat / latitudeBands;
                vertex.color = color;

                mesh.vertices.push_back(vertex);
            }
        }

        // Generate indices
        for (int lat = 0; lat < latitudeBands; lat++) {
            for (int lon = 0; lon < longitudeBands; lon++) {
                int first = lat * (longitudeBands + 1) + lon;
                int second = first + longitudeBands + 1;

                mesh.indices.push_back(first);
                mesh.indices.push_back(second);
                mesh.indices.push_back(first + 1);

                mesh.indices.push_back(second);
                mesh.indices.push_back(second + 1);
                mesh.indices.push_back(first + 1);
            }
        }

        // Calculate flat normals if needed
        if (!smoothNormals) {
            CalculateFlatNormals(mesh);
        }

        mesh.triangleCount = mesh.indices.size() / 3;
        mesh.vertexCount = mesh.vertices.size();

        return mesh;
    }

private:
    static void CalculateFlatNormals(MeshData& mesh) {
        for (size_t i = 0; i < mesh.indices.size(); i += 3) {
            unsigned int i0 = mesh.indices[i];
            unsigned int i1 = mesh.indices[i+1];
            unsigned int i2 = mesh.indices[i+2];

            vec3 v0 = mesh.vertices[i0].position;
            vec3 v1 = mesh.vertices[i1].position;
            vec3 v2 = mesh.vertices[i2].position;

            vec3 edge1 = vec3(v1.x - v0.x, v1.y - v0.y, v1.z - v0.z);
            vec3 edge2 = vec3(v2.x - v0.x, v2.y - v0.y, v2.z - v0.z);

            vec3 normal = vec3(
                edge1.y * edge2.z - edge1.z * edge2.y,
                edge1.z * edge2.x - edge1.x * edge2.z,
                edge1.x * edge2.y - edge1.y * edge2.x
            );

            // Normalize the normal
            float length = sqrt(normal.x * normal.x + normal.y * normal.y + normal.z * normal.z);
            if (length > 0) {
                normal.x /= length;
                normal.y /= length;
                normal.z /= length;
            }

            // Assign the same normal to all three vertices
            mesh.vertices[i0].normal = normal;
            mesh.vertices[i1].normal = normal;
            mesh.vertices[i2].normal = normal;
        }
    }
};
