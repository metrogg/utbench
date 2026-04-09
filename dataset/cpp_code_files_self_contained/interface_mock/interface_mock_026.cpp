#include <vector>
#include <string>
#include <algorithm>
#include <map>
#include <cmath>
#include <limits>

using namespace std;

// Enum for girl types
enum GirlType { CHOOSY, NORMAL, DESPERATE };

// Enum for boy types
enum BoyType { MISER, GENEROUS, GEEK };

// Person base class
class Person {
protected:
    string name;
    int attractiveness;
    int intelligence;
    double budget;
    bool committed;

public:
    Person(string n, int a, int i, double b) 
        : name(n), attractiveness(a), intelligence(i), budget(b), committed(false) {}

    string getName() const { return name; }
    int getAttractiveness() const { return attractiveness; }
    int getIntelligence() const { return intelligence; }
    double getBudget() const { return budget; }
    bool isCommitted() const { return committed; }
    void setCommitted(bool status) { committed = status; }
};

// Girl class
class Girl : public Person {
private:
    GirlType type;
    double maintenanceCost;
    double happiness;
    double compatibility;
    string boyfriendName;

public:
    Girl(string n, int a, int i, double b, GirlType t, double mc) 
        : Person(n, a, i, b), type(t), maintenanceCost(mc), happiness(0), compatibility(0) {}

    GirlType getType() const { return type; }
    double getMaintenanceCost() const { return maintenanceCost; }
    double getHappiness() const { return happiness; }
    double getCompatibility() const { return compatibility; }
    string getBoyfriendName() const { return boyfriendName; }

    void setHappiness(double h) { happiness = h; }
    void setCompatibility(double c) { compatibility = c; }
    void setBoyfriendName(const string& name) { boyfriendName = name; }

    // Calculate happiness based on type and gift value
    void calculateHappiness(double totalGiftValue, double totalGiftCost) {
        switch(type) {
            case CHOOSY:
                happiness = log(totalGiftValue + totalGiftCost);
                break;
            case NORMAL:
                happiness = totalGiftValue + totalGiftCost;
                break;
            case DESPERATE:
                happiness = exp(totalGiftValue / 1000.0);
                break;
        }
    }
};

// Boy class
class Boy : public Person {
private:
    BoyType type;
    int minAttractivenessReq;
    double happiness;

public:
    Boy(string n, int a, int i, double b, BoyType t, int mar) 
        : Person(n, a, i, b), type(t), minAttractivenessReq(mar), happiness(0) {}

    BoyType getType() const { return type; }
    int getMinAttractivenessReq() const { return minAttractivenessReq; }
    double getHappiness() const { return happiness; }

    void setHappiness(double h) { happiness = h; }

    // Calculate happiness based on type and parameters
    void calculateHappiness(const Girl& girlfriend) {
        switch(type) {
            case MISER:
                happiness = budget - girlfriend.getMaintenanceCost();
                break;
            case GENEROUS:
                happiness = girlfriend.getHappiness();
                break;
            case GEEK:
                happiness = girlfriend.getIntelligence();
                break;
        }
    }
};

// Gift structure
struct Gift {
    double price;
    double value;
    bool isLuxury;
};

// Couple allocator class
class CoupleAllocator {
public:
    // Allocate girlfriends to boys based on different criteria
    static void allocateGirlfriends(vector<Boy>& boys, vector<Girl>& girls) {
        for (auto& girl : girls) {
            if (girl.isCommitted()) continue;

            vector<Boy*> eligibleBoys;
            for (auto& boy : boys) {
                if (!boy.isCommitted() && 
                    boy.getBudget() >= girl.getMaintenanceCost() &&
                    girl.getAttractiveness() >= boy.getMinAttractivenessReq()) {
                    eligibleBoys.push_back(&boy);
                }
            }

            if (!eligibleBoys.empty()) {
                // Sort based on different criteria for different girl types
                switch(girl.getType()) {
                    case CHOOSY:
                        sort(eligibleBoys.begin(), eligibleBoys.end(), 
                            [](const Boy* a, const Boy* b) {
                                return a->getAttractiveness() > b->getAttractiveness();
                            });
                        break;
                    case NORMAL:
                        sort(eligibleBoys.begin(), eligibleBoys.end(), 
                            [](const Boy* a, const Boy* b) {
                                return a->getBudget() > b->getBudget();
                            });
                        break;
                    case DESPERATE:
                        sort(eligibleBoys.begin(), eligibleBoys.end(), 
                            [](const Boy* a, const Boy* b) {
                                return a->getIntelligence() > b->getIntelligence();
                            });
                        break;
                }

                // Allocate the best match
                eligibleBoys[0]->setCommitted(true);
                girl.setCommitted(true);
                girl.setBoyfriendName(eligibleBoys[0]->getName());
            }
        }
    }

    // Calculate happiness for all couples
    static void calculateHappiness(vector<Boy>& boys, vector<Girl>& girls, const vector<Gift>& gifts) {
        map<string, double> girlGiftValues;
        map<string, double> girlGiftCosts;

        // Simulate gift giving based on boy type
        for (auto& boy : boys) {
            if (!boy.isCommitted()) continue;

            // Find girlfriend
            auto girlIt = find_if(girls.begin(), girls.end(), 
                [&boy](const Girl& g) { return g.getBoyfriendName() == boy.getName(); });
            
            if (girlIt != girls.end()) {
                double totalValue = 0;
                double totalCost = 0;

                // Different gifting strategies based on boy type
                switch(boy.getType()) {
                    case MISER:
                        // Miser gives cheapest gifts just above maintenance
                        for (const auto& gift : gifts) {
                            if (totalCost < girlIt->getMaintenanceCost()) {
                                totalCost += gift.price;
                                totalValue += gift.value;
                            }
                        }
                        break;
                    case GENEROUS:
                        // Generous gives most valuable gifts within budget
                        for (const auto& gift : gifts) {
                            if (totalCost + gift.price <= boy.getBudget()) {
                                totalCost += gift.price;
                                totalValue += gift.value;
                            }
                        }
                        break;
                    case GEEK:
                        // Geek gives normal gifts plus one luxury
                        bool luxuryGiven = false;
                        for (const auto& gift : gifts) {
                            if (totalCost + gift.price <= boy.getBudget()) {
                                if (!gift.isLuxury || !luxuryGiven) {
                                    totalCost += gift.price;
                                    totalValue += gift.value;
                                    if (gift.isLuxury) luxuryGiven = true;
                                }
                            }
                        }
                        break;
                }

                girlIt->calculateHappiness(totalValue, totalCost);
                boy.calculateHappiness(*girlIt);
            }
        }
    }

    // Calculate compatibility for all couples
    static void calculateCompatibility(vector<Boy>& boys, vector<Girl>& girls) {
        for (auto& girl : girls) {
            if (!girl.isCommitted()) continue;

            // Find boyfriend
            auto boyIt = find_if(boys.begin(), boys.end(), 
                [&girl](const Boy& b) { return b.getName() == girl.getBoyfriendName(); });
            
            if (boyIt != boys.end()) {
                double compatibility = (boyIt->getBudget() - girl.getMaintenanceCost()) + 
                                      abs(boyIt->getAttractiveness() - girl.getAttractiveness()) + 
                                      abs(boyIt->getIntelligence() - girl.getIntelligence());
                girl.setCompatibility(compatibility);
            }
        }
    }

    // Get top K happiest couples
    static vector<pair<string, string>> getTopKHappiest(const vector<Boy>& boys, const vector<Girl>& girls, int k) {
        vector<pair<const Girl*, double>> girlHappiness;
        
        for (const auto& girl : girls) {
            if (girl.isCommitted()) {
                girlHappiness.emplace_back(&girl, girl.getHappiness());
            }
        }

        sort(girlHappiness.begin(), girlHappiness.end(), 
            [](const pair<const Girl*, double>& a, const pair<const Girl*, double>& b) {
                return a.second > b.second;
            });

        vector<pair<string, string>> result;
        for (int i = 0; i < min(k, (int)girlHappiness.size()); ++i) {
            result.emplace_back(girlHappiness[i].first->getBoyfriendName(), 
                               girlHappiness[i].first->getName());
        }

        return result;
    }

    // Get top K most compatible couples
    static vector<pair<string, string>> getTopKCompatible(const vector<Boy>& boys, const vector<Girl>& girls, int k) {
        vector<pair<const Girl*, double>> girlCompatibility;
        
        for (const auto& girl : girls) {
            if (girl.isCommitted()) {
                girlCompatibility.emplace_back(&girl, girl.getCompatibility());
            }
        }

        sort(girlCompatibility.begin(), girlCompatibility.end(), 
            [](const pair<const Girl*, double>& a, const pair<const Girl*, double>& b) {
                return a.second > b.second;
            });

        vector<pair<string, string>> result;
        for (int i = 0; i < min(k, (int)girlCompatibility.size()); ++i) {
            result.emplace_back(girlCompatibility[i].first->getBoyfriendName(), 
                               girlCompatibility[i].first->getName());
        }

        return result;
    }
};
