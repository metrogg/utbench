#include <string>
#include <vector>
#include <stdexcept>
#include <map>
#include <algorithm>

using namespace std;

enum class Feature { Standard, Toilet, Bar, Restaurant, Couchette, Luggage, FirstClass };

class Wagon {
    int mass;
    vector<Feature> features;
public:
    Wagon(int m, const vector<Feature>& f) : mass(m), features(f) {
        if (m < 10 || m > 50) throw runtime_error("Invalid wagon mass");
    }
    
    int get_mass() const { return mass; }
    bool has_feature(Feature f) const {
        return find(features.begin(), features.end(), f) != features.end();
    }
    vector<Feature> get_features() const { return features; }
    
    friend ostream& operator<<(ostream& o, const Wagon& w) {
        o << "Wagon[" << w.mass << " tons, features:";
        for (auto f : w.features) {
            switch(f) {
                case Feature::Standard: o << " Standard"; break;
                case Feature::Toilet: o << " Toilet"; break;
                case Feature::Bar: o << " Bar"; break;
                case Feature::Restaurant: o << " Restaurant"; break;
                case Feature::Couchette: o << " Couchette"; break;
                case Feature::Luggage: o << " Luggage"; break;
                case Feature::FirstClass: o << " FirstClass"; break;
            }
        }
        return o << "]";
    }
};

class Train {
    int lokmasse;
    int maximallast;
    vector<Wagon> wagonliste;
    
public:
    Train(int a, int b, const vector<Wagon>& c) : lokmasse(a), maximallast(b), wagonliste(c) {
        validate_mass();
    }
    
    Train(int a, int b) : lokmasse(a), maximallast(b) {
        validate_mass();
    }
    
    void validate_mass() const {
        if (lokmasse < 50 || lokmasse > 200 || maximallast < 200 || maximallast > 10000) {
            throw runtime_error("Invalid locomotive mass or max load");
        }
    }
    
    int total_load() const {
        int sum = lokmasse;
        for (const auto& w : wagonliste) sum += w.get_mass();
        return sum;
    }
    
    bool ready() const {
        if (total_load() > maximallast) return false;
        
        // Check toilet spacing
        int since_last_toilet = 0;
        for (const auto& w : wagonliste) {
            if (w.has_feature(Feature::Toilet)) {
                since_last_toilet = 0;
            } else {
                if (++since_last_toilet > 3) return false;
            }
        }
        
        // Check required features
        map<Feature, int> feature_counts;
        for (const auto& w : wagonliste) {
            for (auto f : w.get_features()) {
                feature_counts[f]++;
            }
        }
        
        return feature_counts[Feature::Bar] > 0 &&
               feature_counts[Feature::Restaurant] > 0 &&
               feature_counts[Feature::Couchette] > 0 &&
               feature_counts[Feature::Standard] > 0 &&
               feature_counts[Feature::Toilet] > 0;
    }
    
    void couple(const vector<Wagon>& wagvect) {
        wagonliste.insert(wagonliste.end(), wagvect.begin(), wagvect.end());
    }
    
    vector<Wagon> uncouple(size_t from) {
        if (from >= wagonliste.size()) throw runtime_error("Invalid uncoupling position");
        
        vector<Wagon> disconnected(wagonliste.begin() + from, wagonliste.end());
        wagonliste.erase(wagonliste.begin() + from, wagonliste.end());
        return disconnected;
    }
    
    double efficiency() const {
        if (wagonliste.empty()) return 0.0;
        return static_cast<double>(total_load()) / maximallast;
    }
    
    vector<Wagon> get_wagons_by_feature(Feature f) const {
        vector<Wagon> result;
        for (const auto& w : wagonliste) {
            if (w.has_feature(f)) result.push_back(w);
        }
        return result;
    }
    
    friend ostream& operator<<(ostream& o, const Train& t) {
        o << "Train[" << t.total_load() << "/" << t.maximallast << " tons, ";
        o << (t.ready() ? "READY" : "NOT READY") << ", efficiency: " << t.efficiency() << " {";
        for (size_t i = 0; i < t.wagonliste.size(); ++i) {
            if (i != 0) o << ", ";
            o << t.wagonliste[i];
        }
        return o << "}]";
    }
};
