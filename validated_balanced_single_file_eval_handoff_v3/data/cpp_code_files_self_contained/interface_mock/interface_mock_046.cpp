#include <vector>
#include <string>
#include <map>
#include <algorithm>

using namespace std;

class Player {
private:
    vector<string> backpack;
    map<string, bool> specialItems;
    map<string, int> itemCounts;
    bool escaped;
    int capacity;
    int health;
    
public:
    // Constructor with initial capacity
    Player(int initialCapacity = 10) : capacity(initialCapacity), health(100), escaped(false) {
        specialItems = {
            {"key", false},
            {"sledgehammer", false},
            {"map", false},
            {"flashlight", false}
        };
    }

    // Add item to backpack with capacity check
    bool addItem(const string& item) {
        if (backpack.size() >= capacity) {
            return false;
        }
        
        backpack.push_back(item);
        
        // Update special items status
        if (specialItems.find(item) != specialItems.end()) {
            specialItems[item] = true;
        }
        
        // Update item counts
        itemCounts[item]++;
        
        return true;
    }

    // Remove item from backpack
    bool removeItem(const string& item) {
        auto it = find(backpack.begin(), backpack.end(), item);
        if (it != backpack.end()) {
            backpack.erase(it);
            
            // Check if this was the last special item
            if (specialItems.find(item) != specialItems.end()) {
                specialItems[item] = (count(backpack.begin(), backpack.end(), item) > 0);
            }
            
            itemCounts[item]--;
            return true;
        }
        return false;
    }

    // Check if player has a specific item
    bool hasItem(const string& item) const {
        if (specialItems.find(item) != specialItems.end()) {
            return specialItems.at(item);
        }
        return (find(backpack.begin(), backpack.end(), item) != backpack.end());
    }

    // Get count of a specific item
    int getItemCount(const string& item) const {
        auto it = itemCounts.find(item);
        return (it != itemCounts.end()) ? it->second : 0;
    }

    // Get current backpack contents
    vector<string> getBackpackContents() const {
        return backpack;
    }

    // Get current backpack capacity
    int getRemainingCapacity() const {
        return capacity - backpack.size();
    }

    // Health management
    void takeDamage(int amount) {
        health = max(0, health - amount);
    }

    void heal(int amount) {
        health = min(100, health + amount);
    }

    int getHealth() const {
        return health;
    }

    // Escape status
    void setEscaped(bool status) {
        escaped = status;
    }

    bool getEscaped() const {
        return escaped;
    }

    // Upgrade backpack capacity
    void upgradeCapacity(int additionalSpace) {
        capacity += additionalSpace;
    }
};
