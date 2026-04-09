#include <vector>
#include <cmath>
#include <algorithm>
#include <map>

using namespace std;

struct Coord2D {
    float x, y;
    Coord2D(float x = 0, float y = 0) : x(x), y(y) {}
    
    Coord2D operator+(const Coord2D& other) const {
        return Coord2D(x + other.x, y + other.y);
    }
    
    Coord2D operator-(const Coord2D& other) const {
        return Coord2D(x - other.x, y - other.y);
    }
    
    Coord2D operator*(float scalar) const {
        return Coord2D(x * scalar, y * scalar);
    }
};

class EnhancedSprite {
private:
    vector<string> m_sprite;
    Coord2D m_position;
    string m_tag;
    unsigned short m_width, m_height;
    
public:
    EnhancedSprite() : m_width(0), m_height(0) {}
    
    EnhancedSprite(const vector<string>& sprite, const string& tag = "") {
        create(sprite, tag);
    }
    
    void create(const vector<string>& sprite, const string& tag = "") {
        m_sprite = sprite;
        m_tag = tag;
        m_height = sprite.size();
        m_width = 0;
        
        for (const auto& line : sprite) {
            if (line.length() > m_width) {
                m_width = line.length();
            }
        }
    }
    
    // Enhanced collision detection with multiple methods
    bool boxCollision(const EnhancedSprite& other) const {
        Coord2D thisCenter = m_position + Coord2D(m_width/2.0f, m_height/2.0f);
        Coord2D otherCenter = other.m_position + Coord2D(other.m_width/2.0f, other.m_height/2.0f);
        
        float xCollision = abs(thisCenter.x - otherCenter.x) < (m_width + other.m_width)/2.0f;
        float yCollision = abs(thisCenter.y - otherCenter.y) < (m_height + other.m_height)/2.0f;
        
        return xCollision && yCollision;
    }
    
    bool pixelPerfectCollision(const EnhancedSprite& other) const {
        if (!boxCollision(other)) return false;
        
        // Get overlapping area
        int x1 = max((int)m_position.x, (int)other.m_position.x);
        int x2 = min((int)(m_position.x + m_width), (int)(other.m_position.x + other.m_width));
        int y1 = max((int)m_position.y, (int)other.m_position.y);
        int y2 = min((int)(m_position.y + m_height), (int)(other.m_position.y + other.m_height));
        
        for (int y = y1; y < y2; y++) {
            for (int x = x1; x < x2; x++) {
                int thisX = x - (int)m_position.x;
                int thisY = y - (int)m_position.y;
                int otherX = x - (int)other.m_position.x;
                int otherY = y - (int)other.m_position.y;
                
                // Check if both sprites have non-space characters at this position
                if (thisY >= 0 && thisY < m_sprite.size() && 
                    thisX >= 0 && thisX < m_sprite[thisY].length() &&
                    m_sprite[thisY][thisX] != ' ' &&
                    otherY >= 0 && otherY < other.m_sprite.size() && 
                    otherX >= 0 && otherX < other.m_sprite[otherY].length() &&
                    other.m_sprite[otherY][otherX] != ' ') {
                    return true;
                }
            }
        }
        
        return false;
    }
    
    float distanceTo(const EnhancedSprite& other) const {
        Coord2D thisCenter = getCenter();
        Coord2D otherCenter = other.getCenter();
        float dx = thisCenter.x - otherCenter.x;
        float dy = thisCenter.y - otherCenter.y;
        return sqrt(dx*dx + dy*dy);
    }
    
    bool isTouchingEdge(int screenWidth, int screenHeight) const {
        return (m_position.x <= 0) || 
               (m_position.x + m_width >= screenWidth) ||
               (m_position.y <= 0) || 
               (m_position.y + m_height >= screenHeight);
    }
    
    // Getters and setters
    void setPosition(const Coord2D& pos) { m_position = pos; }
    Coord2D getPosition() const { return m_position; }
    Coord2D getCenter() const { return m_position + Coord2D(m_width/2.0f, m_height/2.0f); }
    unsigned short getWidth() const { return m_width; }
    unsigned short getHeight() const { return m_height; }
    string getTag() const { return m_tag; }
    
    void render() const {
        for (const auto& line : m_sprite) {
            cout << line << endl;
        }
    }
};
