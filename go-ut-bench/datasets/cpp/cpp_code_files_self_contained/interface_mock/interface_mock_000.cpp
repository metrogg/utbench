#include <vector>
#include <string>
#include <algorithm>
#include <map>

using namespace std;

// Enum for unit types (simplified from original)
enum UnitTypes {
    UNIT_WARRIOR,
    UNIT_ARCHER,
    UNIT_SWORDSMAN,
    UNIT_CAVALRY,
    UNIT_CATAPULT,
    UNIT_TRANSPORT,
    NUM_UNIT_TYPES
};

// Unit data structure containing all relevant attributes
struct UnitData {
    string name;
    int cost;
    int strength;
    int move;
    int collateral;
    int range;
    int bombard;
    int cargo;
    int withdrawal;
    int power;
};

// Global unit database (simulating the original game's unit data)
map<UnitTypes, UnitData> unitDatabase = {
    {UNIT_WARRIOR, {"Warrior", 20, 4, 1, 0, 0, 0, 0, 0, 10}},
    {UNIT_ARCHER, {"Archer", 30, 3, 1, 0, 2, 0, 0, 10, 15}},
    {UNIT_SWORDSMAN, {"Swordsman", 40, 6, 1, 0, 0, 0, 0, 0, 20}},
    {UNIT_CAVALRY, {"Cavalry", 60, 10, 3, 0, 0, 0, 0, 25, 30}},
    {UNIT_CATAPULT, {"Catapult", 50, 5, 1, 25, 1, 8, 0, 0, 25}},
    {UNIT_TRANSPORT, {"Transport", 80, 2, 2, 0, 0, 0, 3, 40, 10}}
};

// Enum for sort types (from original)
enum UnitSortTypes {
    NO_UNIT_SORT = -1,
    UNIT_SORT_NAME,
    UNIT_SORT_COST,
    UNIT_SORT_STRENGTH,
    UNIT_SORT_MOVE,
    UNIT_SORT_COLLATERAL,
    UNIT_SORT_RANGE,
    UNIT_SORT_BOMBARD,
    UNIT_SORT_CARGO,
    UNIT_SORT_WITHDRAWAL,
    UNIT_SORT_POWER,
    NUM_UNIT_SORT
};

// Base class for unit sorting (from original, expanded)
class UnitSortBase {
public:
    UnitSortBase(bool bInvert = false) : m_bInvert(bInvert) {};
    virtual ~UnitSortBase() {};
    
    // Returns the value of eUnit in the sorting category
    virtual int getUnitValue(UnitTypes eUnit) = 0;
    
    // Returns if eUnit1 is before eUnit2
    virtual bool isLesserUnit(UnitTypes eUnit1, UnitTypes eUnit2) {
        return m_bInvert ? 
            (getUnitValue(eUnit1) < getUnitValue(eUnit2)) : 
            (getUnitValue(eUnit1) > getUnitValue(eUnit2));
    }
    
    bool isInverse() { return m_bInvert; }
    bool setInverse(bool bInvert) { m_bInvert = bInvert; return true; }
    
protected:
    bool m_bInvert;
};

// Derived sorting classes for each attribute (expanded with additional logic)
class UnitSortStrength : public UnitSortBase {
public:
    UnitSortStrength(bool bInvert = false) : UnitSortBase(bInvert) {};
    int getUnitValue(UnitTypes eUnit) override {
        return unitDatabase[eUnit].strength;
    }
};

class UnitSortMove : public UnitSortBase {
public:
    UnitSortMove(bool bInvert = false) : UnitSortBase(bInvert) {};
    int getUnitValue(UnitTypes eUnit) override {
        return unitDatabase[eUnit].move;
    }
};

class UnitSortCollateral : public UnitSortBase {
public:
    UnitSortCollateral(bool bInvert = false) : UnitSortBase(bInvert) {};
    int getUnitValue(UnitTypes eUnit) override {
        return unitDatabase[eUnit].collateral;
    }
};

class UnitSortRange : public UnitSortBase {
public:
    UnitSortRange(bool bInvert = false) : UnitSortBase(bInvert) {};
    int getUnitValue(UnitTypes eUnit) override {
        return unitDatabase[eUnit].range;
    }
};

class UnitSortBombard : public UnitSortBase {
public:
    UnitSortBombard(bool bInvert = false) : UnitSortBase(bInvert) {};
    int getUnitValue(UnitTypes eUnit) override {
        return unitDatabase[eUnit].bombard;
    }
};

class UnitSortCargo : public UnitSortBase {
public:
    UnitSortCargo(bool bInvert = false) : UnitSortBase(bInvert) {};
    int getUnitValue(UnitTypes eUnit) override {
        return unitDatabase[eUnit].cargo;
    }
};

class UnitSortWithdrawal : public UnitSortBase {
public:
    UnitSortWithdrawal(bool bInvert = false) : UnitSortBase(bInvert) {};
    int getUnitValue(UnitTypes eUnit) override {
        return unitDatabase[eUnit].withdrawal;
    }
};

class UnitSortPower : public UnitSortBase {
public:
    UnitSortPower(bool bInvert = false) : UnitSortBase(bInvert) {};
    int getUnitValue(UnitTypes eUnit) override {
        return unitDatabase[eUnit].power;
    }
};

class UnitSortCost : public UnitSortBase {
public:
    UnitSortCost(bool bInvert = false) : UnitSortBase(bInvert) {};
    int getUnitValue(UnitTypes eUnit) override {
        return unitDatabase[eUnit].cost;
    }
};

class UnitSortName : public UnitSortBase {
public:
    UnitSortName(bool bInvert = false) : UnitSortBase(bInvert) {};
    bool isLesserUnit(UnitTypes eUnit1, UnitTypes eUnit2) override {
        string name1 = unitDatabase[eUnit1].name;
        string name2 = unitDatabase[eUnit2].name;
        return m_bInvert ? (name1 > name2) : (name1 < name2);
    }
    int getUnitValue(UnitTypes eUnit) override {
        return 0; // Not used for name sorting
    }
};

// UnitSortList class to manage sorting (expanded with additional functionality)
class UnitSortList {
public:
    UnitSortList() {
        // Initialize all sort types
        m_apUnitSort[UNIT_SORT_NAME] = new UnitSortName();
        m_apUnitSort[UNIT_SORT_COST] = new UnitSortCost();
        m_apUnitSort[UNIT_SORT_STRENGTH] = new UnitSortStrength();
        m_apUnitSort[UNIT_SORT_MOVE] = new UnitSortMove();
        m_apUnitSort[UNIT_SORT_COLLATERAL] = new UnitSortCollateral();
        m_apUnitSort[UNIT_SORT_RANGE] = new UnitSortRange();
        m_apUnitSort[UNIT_SORT_BOMBARD] = new UnitSortBombard();
        m_apUnitSort[UNIT_SORT_CARGO] = new UnitSortCargo();
        m_apUnitSort[UNIT_SORT_WITHDRAWAL] = new UnitSortWithdrawal();
        m_apUnitSort[UNIT_SORT_POWER] = new UnitSortPower();
        
        m_eActiveSort = UNIT_SORT_NAME; // Default sort
    }
    
    ~UnitSortList() {
        for (int i = 0; i < NUM_UNIT_SORT; i++) {
            delete m_apUnitSort[i];
        }
    }
    
    UnitSortTypes getActiveSort() { return m_eActiveSort; }
    
    bool setActiveSort(UnitSortTypes eActiveSort) {
        if (eActiveSort >= 0 && eActiveSort < NUM_UNIT_SORT) {
            m_eActiveSort = eActiveSort;
            return true;
        }
        return false;
    }
    
    int getNumSort() { return NUM_UNIT_SORT; }
    
    bool operator()(UnitTypes eUnit1, UnitTypes eUnit2) {
        return m_apUnitSort[m_eActiveSort]->isLesserUnit(eUnit1, eUnit2);
    }
    
    // New method to sort a vector of units
    void sortUnits(vector<UnitTypes>& units) {
        sort(units.begin(), units.end(), [this](UnitTypes a, UnitTypes b) {
            return this->operator()(a, b);
        });
    }
    
    // New method to set sort direction
    void setSortDirection(bool bInvert) {
        m_apUnitSort[m_eActiveSort]->setInverse(bInvert);
    }
    
protected:
    UnitSortBase* m_apUnitSort[NUM_UNIT_SORT];
    UnitSortTypes m_eActiveSort;
};

// Helper function to print unit list
void printUnits(const vector<UnitTypes>& units) {
    for (UnitTypes unit : units) {
        cout << unitDatabase[unit].name;
        cout << " (Cost: " << unitDatabase[unit].cost;
        cout << ", Strength: " << unitDatabase[unit].strength;
        cout << ", Move: " << unitDatabase[unit].move << ")";
        cout << endl;
    }
    cout << "----------------------" << endl;
}
