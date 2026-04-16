#include <string>
#include <vector>
#include <sstream>
#include <stdexcept>
#include <iomanip>

using namespace std;

struct WizardingMoney {
    int galleons;
    int sickles;
    int knuts;
    
    WizardingMoney(int g = 0, int s = 0, int k = 0) : galleons(g), sickles(s), knuts(k) {
        normalize();
    }
    
    void normalize() {
        // Convert overflow knuts to sickles
        sickles += knuts / 29;
        knuts %= 29;
        
        // Convert overflow sickles to galleons
        galleons += sickles / 17;
        sickles %= 17;
        
        // Handle negative values
        if (knuts < 0) {
            sickles -= 1;
            knuts += 29;
        }
        if (sickles < 0) {
            galleons -= 1;
            sickles += 17;
        }
    }
    
    int toKnuts() const {
        return galleons * 17 * 29 + sickles * 29 + knuts;
    }
    
    static WizardingMoney fromKnuts(int totalKnuts) {
        bool isNegative = totalKnuts < 0;
        if (isNegative) totalKnuts = -totalKnuts;
        
        int g = totalKnuts / (17 * 29);
        int remaining = totalKnuts % (17 * 29);
        int s = remaining / 29;
        int k = remaining % 29;
        
        if (isNegative) g = -g;
        return WizardingMoney(g, s, k);
    }
    
    WizardingMoney operator+(const WizardingMoney& other) const {
        return WizardingMoney(galleons + other.galleons, 
                            sickles + other.sickles, 
                            knuts + other.knuts);
    }
    
    WizardingMoney operator-(const WizardingMoney& other) const {
        return WizardingMoney(galleons - other.galleons, 
                            sickles - other.sickles, 
                            knuts - other.knuts);
    }
    
    string toString() const {
        stringstream ss;
        ss << galleons << "." << sickles << "." << knuts;
        return ss.str();
    }
    
    static WizardingMoney parse(const string& moneyStr) {
        char dot;
        int g, s, k;
        stringstream ss(moneyStr);
        
        if (!(ss >> g >> dot >> s >> dot >> k) || dot != '.') {
            throw invalid_argument("Invalid money format");
        }
        
        return WizardingMoney(g, s, k);
    }
};

WizardingMoney calculateChange(const WizardingMoney& paid, const WizardingMoney& owed) {
    return paid - owed;
}
