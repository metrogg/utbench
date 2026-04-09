#include <string>
#include <vector>
#include <map>
#include <stdexcept>
#include <iomanip>

using namespace std;

enum class ArmorType { LIGHT, MEDIUM, HEAVY };
enum class Material { LEATHER, CHAINMAIL, PLATE, SCALE, PADDED };

class Armor {
private:
    string name;
    int defenseValue;
    int maxDexBonus;
    float weight;
    int id;
    ArmorType type;
    Material material;
    int durability;
    int maxDurability;
    bool isMagic;
    float magicModifier;
    
    static int nextId;
    
    void calculateWorth() {
        float baseWorth = defenseValue * 10.0f;
        if (isMagic) baseWorth *= (1.0f + magicModifier);
        worth = baseWorth * (1.0f + static_cast<float>(material) * 0.3f);
    }

public:
    Armor() : name("Unnamed Armor"), defenseValue(0), maxDexBonus(0), 
              weight(0.0f), id(nextId++), type(ArmorType::LIGHT), 
              material(Material::LEATHER), durability(100), maxDurability(100),
              isMagic(false), magicModifier(0.0f) {
        calculateWorth();
    }
    
    Armor(string name, int defense, int maxDex, float weight, 
          ArmorType type, Material material, bool magical = false, 
          float magicMod = 0.0f) 
        : name(name), defenseValue(defense), maxDexBonus(maxDex), 
          weight(weight), id(nextId++), type(type), material(material),
          durability(100), maxDurability(100), isMagic(magical), 
          magicModifier(magicMod) {
        calculateWorth();
    }
    
    // Getters
    string getName() const { return name; }
    int getDefenseValue() const { return defenseValue; }
    int getMaxDexBonus() const { return maxDexBonus; }
    float getWeight() const { return weight; }
    int getId() const { return id; }
    ArmorType getType() const { return type; }
    Material getMaterial() const { return material; }
    int getDurability() const { return durability; }
    int getMaxDurability() const { return maxDurability; }
    bool getIsMagic() const { return isMagic; }
    float getMagicModifier() const { return magicModifier; }
    float getWorth() const { return worth; }
    
    // Setters with validation
    void setDurability(int newDurability) {
        if (newDurability < 0) durability = 0;
        else if (newDurability > maxDurability) durability = maxDurability;
        else durability = newDurability;
        calculateWorth();
    }
    
    void repair(int amount) {
        setDurability(durability + amount);
    }
    
    void takeDamage(int damage) {
        setDurability(durability - damage);
    }
    
    float getEffectiveDefense() const {
        float effectiveness = durability / static_cast<float>(maxDurability);
        return defenseValue * effectiveness * (isMagic ? (1.0f + magicModifier) : 1.0f);
    }
    
    string getTypeString() const {
        switch(type) {
            case ArmorType::LIGHT: return "Light";
            case ArmorType::MEDIUM: return "Medium";
            case ArmorType::HEAVY: return "Heavy";
            default: return "Unknown";
        }
    }
    
    string getMaterialString() const {
        switch(material) {
            case Material::LEATHER: return "Leather";
            case Material::CHAINMAIL: return "Chainmail";
            case Material::PLATE: return "Plate";
            case Material::SCALE: return "Scale";
            case Material::PADDED: return "Padded";
            default: return "Unknown";
        }
    }
    
    void displayInfo() const {
        cout << "=== Armor Info ===" << endl;
        cout << "ID: " << id << endl;
        cout << "Name: " << name << endl;
        cout << "Type: " << getTypeString() << endl;
        cout << "Material: " << getMaterialString() << endl;
        cout << "Defense: " << defenseValue << endl;
        cout << "Max Dex Bonus: " << maxDexBonus << endl;
        cout << "Weight: " << weight << " lbs" << endl;
        cout << "Durability: " << durability << "/" << maxDurability << endl;
        cout << "Magic: " << (isMagic ? "Yes" : "No") << endl;
        if (isMagic) cout << "Magic Bonus: +" << magicModifier*100 << "%" << endl;
        cout << "Worth: " << fixed << setprecision(2) << worth << " gold" << endl;
        cout << "Effective Defense: " << fixed << setprecision(1) << getEffectiveDefense() << endl;
    }
    
private:
    float worth;
};

int Armor::nextId = 1;
