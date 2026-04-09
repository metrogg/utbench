#include <vector>
#include <string>
#include <memory>
#include <stdexcept>
#include <map>

using namespace std;

class Enchant {
public:
    Enchant(string name, int power, string effect) 
        : m_name(name), m_power(power), m_effect(effect) {}
    
    virtual ~Enchant() {}
    
    string getName() const { return m_name; }
    int getPower() const { return m_power; }
    string getEffect() const { return m_effect; }
    
    virtual void applyEffect() = 0;
    
protected:
    string m_name;
    int m_power;
    string m_effect;
};

class Equipment {
public:
    Equipment(string name, int wgt, int con) 
        : m_name(name), m_weight(wgt), 
          m_condition(con), m_maxCondition(con),
          m_isEnchanted(false) {}
    
    virtual ~Equipment() {
        for (auto enchant : m_enchants) {
            delete enchant;
        }
    }
    
    string getName() const { return m_name; }
    int getWeight() const { return m_weight; }
    int getCondition() const { return m_condition; }
    int getMaxCondition() const { return m_maxCondition; }
    bool isEnchanted() const { return m_isEnchanted; }
    
    void setCondition(int con) { 
        if (con < 0) con = 0;
        if (con > m_maxCondition) con = m_maxCondition;
        m_condition = con; 
    }
    
    void repair(int amount) {
        m_condition += amount;
        if (m_condition > m_maxCondition) {
            m_condition = m_maxCondition;
        }
    }
    
    void addEnchantment(Enchant* enchant) {
        if (enchant == nullptr) {
            throw invalid_argument("Cannot add null enchantment");
        }
        m_enchants.push_back(enchant);
        m_isEnchanted = true;
    }
    
    vector<Enchant*> getEnchantments() const {
        return m_enchants;
    }
    
    virtual string getType() const = 0;
    virtual int calculateValue() const = 0;
    
    void displayInfo() const {
        cout << "Name: " << m_name << endl;
        cout << "Type: " << getType() << endl;
        cout << "Weight: " << m_weight << endl;
        cout << "Condition: " << m_condition << "/" << m_maxCondition << endl;
        cout << "Estimated Value: " << calculateValue() << endl;
        
        if (m_isEnchanted) {
            cout << "Enchantments (" << m_enchants.size() << "):" << endl;
            for (const auto& enchant : m_enchants) {
                cout << " - " << enchant->getName() 
                     << " (Power: " << enchant->getPower() 
                     << ", Effect: " << enchant->getEffect() << ")" << endl;
            }
        }
    }
    
protected:
    string m_name;
    int m_weight;
    int m_condition;
    int m_maxCondition;
    bool m_isEnchanted;
    vector<Enchant*> m_enchants;
};

class Weapon : public Equipment {
public:
    Weapon(string name, int wgt, int con, int dmg, string dmgType)
        : Equipment(name, wgt, con), m_damage(dmg), m_damageType(dmgType) {}
    
    string getType() const override { return "Weapon"; }
    
    int calculateValue() const override {
        int baseValue = m_damage * 10 + (m_condition * 5);
        if (m_isEnchanted) {
            for (const auto& enchant : m_enchants) {
                baseValue += enchant->getPower() * 20;
            }
        }
        return baseValue;
    }
    
    int getDamage() const { return m_damage; }
    string getDamageType() const { return m_damageType; }
    
private:
    int m_damage;
    string m_damageType;
};

class Armor : public Equipment {
public:
    Armor(string name, int wgt, int con, int def, string armorType)
        : Equipment(name, wgt, con), m_defense(def), m_armorType(armorType) {}
    
    string getType() const override { return "Armor"; }
    
    int calculateValue() const override {
        int baseValue = m_defense * 15 + (m_condition * 5);
        if (m_isEnchanted) {
            for (const auto& enchant : m_enchants) {
                baseValue += enchant->getPower() * 25;
            }
        }
        return baseValue;
    }
    
    int getDefense() const { return m_defense; }
    string getArmorType() const { return m_armorType; }
    
private:
    int m_defense;
    string m_armorType;
};

class FireEnchant : public Enchant {
public:
    FireEnchant(int power) 
        : Enchant("Fire Enchant", power, "Adds fire damage to attacks") {}
    
    void applyEffect() override {
        cout << "Applying fire effect with power " << m_power << endl;
    }
};

class FrostEnchant : public Enchant {
public:
    FrostEnchant(int power) 
        : Enchant("Frost Enchant", power, "Slows enemies on hit") {}
    
    void applyEffect() override {
        cout << "Applying frost effect with power " << m_power << endl;
    }
};
