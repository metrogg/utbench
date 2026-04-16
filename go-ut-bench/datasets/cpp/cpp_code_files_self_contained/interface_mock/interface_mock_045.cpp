#include <iomanip>
#include <string>
#include <vector>
#include <cmath>
#include <algorithm>

using namespace std;

class Fahrzeug {
protected:
    static int p_iMaxID;
    int p_iID;
    string p_sName;
    double p_dMaxGeschwindigkeit;
    double p_dGesamtStrecke;
    double p_dGesamtZeit;
    double p_dZeit;
    double p_dFuelLevel;
    double p_dMaxFuel;
    double p_dVerbrauch;
    bool p_bIsParked;
    static double dGlobaleZeit;

public:
    // Constructors
    Fahrzeug();
    Fahrzeug(string name);
    Fahrzeug(string name, double maxSpeed);
    Fahrzeug(string name, double maxSpeed, double fuelCapacity, double consumption);
    Fahrzeug(const Fahrzeug& other);

    // Core methods
    virtual void vInitialisierung();
    virtual void vAusgabe() const;
    virtual void vAbfertigung();
    virtual double dTanken(double dMenge);
    virtual double dGeschwindigkeit() const;
    virtual void vParken();
    virtual void vStart();
    virtual double dVerbrauchBerechnen() const;

    // Operator overloads
    bool operator<(const Fahrzeug& other) const;
    Fahrzeug& operator=(const Fahrzeug& other);
    friend ostream& operator<<(ostream& out, const Fahrzeug& fahrzeug);

    // Getters
    int getID() const { return p_iID; }
    string getName() const { return p_sName; }
    double getMaxSpeed() const { return p_dMaxGeschwindigkeit; }
    double getTotalDistance() const { return p_dGesamtStrecke; }
    double getFuelLevel() const { return p_dFuelLevel; }
    bool isParked() const { return p_bIsParked; }

    // Static methods
    static void vSetzeGlobaleZeit(double zeit) { dGlobaleZeit = zeit; }
    static double dGetGlobaleZeit() { return dGlobaleZeit; }

protected:
    virtual ostream& ostreamAusgabe(ostream& out) const;
};

// Initialize static members
int Fahrzeug::p_iMaxID = 0;
double Fahrzeug::dGlobaleZeit = 0.0;

// Constructors
Fahrzeug::Fahrzeug() {
    vInitialisierung();
    p_sName = "";
}

Fahrzeug::Fahrzeug(string name) {
    vInitialisierung();
    p_sName = name;
}

Fahrzeug::Fahrzeug(string name, double maxSpeed) {
    vInitialisierung();
    p_sName = name;
    p_dMaxGeschwindigkeit = maxSpeed;
}

Fahrzeug::Fahrzeug(string name, double maxSpeed, double fuelCapacity, double consumption) {
    vInitialisierung();
    p_sName = name;
    p_dMaxGeschwindigkeit = maxSpeed;
    p_dMaxFuel = fuelCapacity;
    p_dFuelLevel = fuelCapacity;
    p_dVerbrauch = consumption;
}

Fahrzeug::Fahrzeug(const Fahrzeug& other) {
    vInitialisierung();
    p_sName = other.p_sName + "'";
    p_dMaxGeschwindigkeit = other.p_dMaxGeschwindigkeit;
    p_dMaxFuel = other.p_dMaxFuel;
    p_dFuelLevel = other.p_dFuelLevel;
    p_dVerbrauch = other.p_dVerbrauch;
}

// Core methods
void Fahrzeug::vInitialisierung() {
    p_dMaxGeschwindigkeit = 0;
    p_dGesamtStrecke = 0;
    p_dGesamtZeit = 0;
    p_dZeit = 0;
    p_sName = "";
    p_dMaxFuel = 0;
    p_dFuelLevel = 0;
    p_dVerbrauch = 0;
    p_bIsParked = false;
    p_iMaxID++;
    p_iID = p_iMaxID;
}

void Fahrzeug::vAusgabe() const {
    cout << resetiosflags(ios::adjustfield);
    cout << setiosflags(ios::left) << setw(5) << p_iID
         << setw(15) << p_sName
         << setw(10) << fixed << setprecision(2) << p_dMaxGeschwindigkeit
         << setw(15) << p_dGesamtStrecke
         << setw(10) << p_dFuelLevel << "/" << p_dMaxFuel
         << setw(10) << (p_bIsParked ? "Parked" : "Moving")
         << endl;
}

void Fahrzeug::vAbfertigung() {
    if ((dGlobaleZeit - p_dZeit) > 0 && !p_bIsParked) {
        double timeDiff = dGlobaleZeit - p_dZeit;
        double distance = timeDiff * dGeschwindigkeit();
        
        // Calculate fuel consumption
        double fuelUsed = (distance * p_dVerbrauch) / 100.0;
        if (fuelUsed > p_dFuelLevel) {
            distance = (p_dFuelLevel * 100.0) / p_dVerbrauch;
            p_dFuelLevel = 0;
            p_bIsParked = true; // Automatically park when out of fuel
        } else {
            p_dFuelLevel -= fuelUsed;
        }
        
        p_dGesamtStrecke += distance;
        p_dZeit = dGlobaleZeit;
        p_dGesamtZeit += timeDiff;
    }
}

double Fahrzeug::dTanken(double dMenge) {
    if (dMenge <= 0) return 0;
    
    double fuelNeeded = p_dMaxFuel - p_dFuelLevel;
    double fuelAdded = min(dMenge, fuelNeeded);
    p_dFuelLevel += fuelAdded;
    
    if (p_bIsParked && p_dFuelLevel > 0) {
        p_bIsParked = false; // Auto-start when refueled
    }
    
    return fuelAdded;
}

double Fahrzeug::dGeschwindigkeit() const {
    return p_dMaxGeschwindigkeit;
}

void Fahrzeug::vParken() {
    p_bIsParked = true;
}

void Fahrzeug::vStart() {
    if (p_dFuelLevel > 0) {
        p_bIsParked = false;
    }
}

double Fahrzeug::dVerbrauchBerechnen() const {
    return (p_dGesamtStrecke * p_dVerbrauch) / 100.0;
}

// Operator overloads
bool Fahrzeug::operator<(const Fahrzeug& other) const {
    return this->p_dGesamtStrecke < other.p_dGesamtStrecke;
}

Fahrzeug& Fahrzeug::operator=(const Fahrzeug& other) {
    if (this != &other) {
        vInitialisierung();
        p_sName = other.p_sName + "''";
        p_dMaxGeschwindigkeit = other.p_dMaxGeschwindigkeit;
        p_dMaxFuel = other.p_dMaxFuel;
        p_dFuelLevel = other.p_dFuelLevel;
        p_dVerbrauch = other.p_dVerbrauch;
    }
    return *this;
}

ostream& Fahrzeug::ostreamAusgabe(ostream& out) const {
    out << resetiosflags(ios::adjustfield);
    out << setiosflags(ios::left) << setw(5) << p_iID
        << setw(15) << p_sName
        << setw(10) << fixed << setprecision(2) << p_dMaxGeschwindigkeit
        << setw(15) << p_dGesamtStrecke
        << setw(10) << p_dFuelLevel << "/" << p_dMaxFuel
        << setw(10) << (p_bIsParked ? "Parked" : "Moving");
    return out;
}

ostream& operator<<(ostream& out, const Fahrzeug& fahrzeug) {
    return fahrzeug.ostreamAusgabe(out);
}
