#include <vector>
#include <map>
#include <string>
#include <cmath>
#include <algorithm>

using namespace std;

class GameEntity {
private:
    int id;
    int teamId;
    int maxHealth;
    int currentHealth;
    int playerId;
    bool isControlled;
    string entityType;
    map<string, int> attributes;
    vector<string> abilities;

    void checkDestroyed() {
        if (currentHealth <= 0) {
            currentHealth = 0;
            isControlled = false;
            playerId = -1;
            cout << "Entity " << id << " has been destroyed!" << endl;
        }
    }

public:
    GameEntity(int _id, int _teamId, int _maxHealth, string _type, 
               const map<string, int>& _attrs = {}, 
               const vector<string>& _abilities = {}) 
        : id(_id), teamId(_teamId), maxHealth(_maxHealth), 
          currentHealth(_maxHealth), entityType(_type), 
          attributes(_attrs), abilities(_abilities) {
        playerId = -1;
        isControlled = false;
    }

    void takeDamage(int damage, string damageType = "normal") {
        if (damageType == "normal") {
            currentHealth -= damage;
        } else if (damageType == "piercing") {
            currentHealth -= (damage * 1.5);
        } else if (damageType == "explosive") {
            currentHealth -= (damage * 0.8);
        }
        checkDestroyed();
    }

    void heal(int amount) {
        currentHealth = min(currentHealth + amount, maxHealth);
    }

    void upgrade(string attribute, int value) {
        if (attribute == "maxHealth") {
            maxHealth += value;
            currentHealth += value;
        } else if (attributes.count(attribute)) {
            attributes[attribute] += value;
        }
    }

    bool useAbility(string ability) {
        if (find(abilities.begin(), abilities.end(), ability) != abilities.end()) {
            cout << "Using ability: " << ability << endl;
            return true;
        }
        return false;
    }

    void changeControl(int newPlayerId) {
        playerId = newPlayerId;
        isControlled = (newPlayerId > 0);
    }

    void displayStatus() const {
        cout << "Entity ID: " << id << " | Type: " << entityType 
             << " | Health: " << currentHealth << "/" << maxHealth 
             << " | Controlled: " << (isControlled ? "Yes" : "No") << endl;
        cout << "Attributes: ";
        for (const auto& attr : attributes) {
            cout << attr.first << ": " << attr.second << " ";
        }
        cout << "\nAbilities: ";
        for (const auto& ab : abilities) {
            cout << ab << " ";
        }
        cout << "\n" << endl;
    }

    int getId() const { return id; }
    int getTeamId() const { return teamId; }
    int getCurrentHealth() const { return currentHealth; }
    int getPlayerId() const { return playerId; }
    bool getIsControlled() const { return isControlled; }
    string getEntityType() const { return entityType; }
};
