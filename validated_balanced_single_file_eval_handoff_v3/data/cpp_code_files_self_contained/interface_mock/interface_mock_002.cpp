#include <vector>
#include <string>
#include <map>
#include <cmath>

using namespace std;

// Simple 2D Point class
class Point {
public:
    float x, y;
    Point(float x = 0, float y = 0) : x(x), y(y) {}
    
    Point operator+(const Point& other) const {
        return Point(x + other.x, y + other.y);
    }
    
    Point operator-(const Point& other) const {
        return Point(x - other.x, y - other.y);
    }
    
    Point operator*(float scalar) const {
        return Point(x * scalar, y * scalar);
    }
    
    float magnitude() const {
        return sqrt(x*x + y*y);
    }
    
    Point normalize() const {
        float mag = magnitude();
        if (mag == 0) return Point(0, 0);
        return Point(x/mag, y/mag);
    }
};

// Axis-Aligned Bounding Box for collision detection
class Hitbox {
    Point position;
    float width, height;
    
public:
    Hitbox(Point pos = Point(), float w = 0, float h = 0) 
        : position(pos), width(w), height(h) {}
    
    bool collides_with(const Hitbox& other) const {
        return (position.x < other.position.x + other.width &&
                position.x + width > other.position.x &&
                position.y < other.position.y + other.height &&
                position.y + height > other.position.y);
    }
    
    void set_position(Point pos) { position = pos; }
    Point get_position() const { return position; }
    float get_width() const { return width; }
    float get_height() const { return height; }
    void set_size(float w, float h) { width = w; height = h; }
};

// Animation frame data
struct Frame {
    string texture_id;
    float duration;
    Frame(string id, float dur) : texture_id(id), duration(dur) {}
};

// Animation sequence
class Animation {
    vector<Frame> frames;
    string name;
    float current_time;
    int current_frame;
    bool looping;
    
public:
    Animation(string n = "", bool loop = true) 
        : name(n), current_time(0), current_frame(0), looping(loop) {}
    
    void add_frame(string texture_id, float duration) {
        frames.emplace_back(texture_id, duration);
    }
    
    void update(float delta_time) {
        if (frames.empty()) return;
        
        current_time += delta_time;
        while (current_time >= frames[current_frame].duration) {
            current_time -= frames[current_frame].duration;
            current_frame++;
            
            if (current_frame >= frames.size()) {
                if (looping) current_frame = 0;
                else current_frame = frames.size() - 1;
            }
        }
    }
    
    string get_current_frame() const {
        if (frames.empty()) return "missing_texture";
        return frames[current_frame].texture_id;
    }
    
    void reset() {
        current_time = 0;
        current_frame = 0;
    }
    
    bool is_finished() const {
        return !looping && current_frame == frames.size() - 1;
    }
};

// Enhanced Entity class with physics and state management
class Entity {
protected:
    Point position;
    Point velocity;
    Point acceleration;
    float mass;
    float friction;
    bool is_kinematic;
    Hitbox hitbox;
    Point hitbox_offset;
    vector<Animation> animations;
    int current_animation;
    map<string, int> animation_map;
    string state;
    
public:
    Entity() : position(0, 0), velocity(0, 0), acceleration(0, 0),
               mass(1.0f), friction(0.1f), is_kinematic(false),
               current_animation(0), state("idle") {}
    
    virtual ~Entity() = default;
    
    // Physics update
    virtual void update(float delta_time) {
        if (!is_kinematic) {
            velocity = velocity + acceleration * delta_time;
            velocity = velocity * (1.0f - friction);
            position = position + velocity * delta_time;
            acceleration = Point(0, 0);
        }
        
        if (current_animation < animations.size()) {
            animations[current_animation].update(delta_time);
        }
    }
    
    // Add force to the entity
    void add_force(Point force) {
        acceleration = acceleration + force * (1.0f / mass);
    }
    
    // Animation management
    void add_animation(string name, Animation anim) {
        animations.push_back(anim);
        animation_map[name] = animations.size() - 1;
    }
    
    void set_animation(string name) {
        auto it = animation_map.find(name);
        if (it != animation_map.end()) {
            current_animation = it->second;
            animations[current_animation].reset();
        }
    }
    
    string get_current_frame() const {
        if (current_animation < animations.size()) {
            return animations[current_animation].get_current_frame();
        }
        return "missing_texture";
    }
    
    // Collision detection
    bool check_collision(const Entity& other) const {
        Hitbox this_hitbox = hitbox;
        this_hitbox.set_position(position + hitbox_offset);
        
        Hitbox other_hitbox = other.hitbox;
        other_hitbox.set_position(other.position + other.hitbox_offset);
        
        return this_hitbox.collides_with(other_hitbox);
    }
    
    // Getters and setters
    void set_position(Point pos) { position = pos; }
    Point get_position() const { return position; }
    
    void set_velocity(Point vel) { velocity = vel; }
    Point get_velocity() const { return velocity; }
    
    void set_hitbox(Hitbox hb) { hitbox = hb; }
    Hitbox get_hitbox() const { return hitbox; }
    
    void set_state(string s) { state = s; }
    string get_state() const { return state; }
    
    void set_kinematic(bool kinematic) { is_kinematic = kinematic; }
    bool get_kinematic() const { return is_kinematic; }
};
