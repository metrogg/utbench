#include <memory>
#include <vector>
#include <algorithm>
#include <thread>
#include <mutex>
#include <map>
#include <typeinfo>
#include <chrono>

using namespace std;

// Enhanced Resource class with more functionality
class Resource {
public:
    Resource() { cout << "Resource created\n"; }
    ~Resource() { cout << "Resource destroyed\n"; }
    void use() { cout << "Resource used\n"; }
    int getId() const { return id; }
private:
    static int counter;
    int id = counter++;
};
int Resource::counter = 0;

// Enhanced Person class with circular reference handling
class Person {
public:
    string mName;
    weak_ptr<Person> mPartner; // Changed to weak_ptr to avoid circular references
    
    Person(const string& name) : mName(name) {
        cout << "Person " << mName << " created\n";
    }
    
    ~Person() {
        cout << "Person " << mName << " destroyed\n";
    }
    
    shared_ptr<Person> getPartner() const {
        return mPartner.lock();
    }
    
    void setPartner(shared_ptr<Person> partner) {
        mPartner = partner;
    }
};

// Function to demonstrate unique_ptr ownership transfer
void takeOwnership(unique_ptr<Resource> ptr) {
    if (ptr) {
        cout << "Resource ID " << ptr->getId() << " ownership taken\n";
        ptr->use();
    }
}

// Enhanced partner function with thread safety
mutex partnerMutex;
bool partnerUp(shared_ptr<Person> p1, shared_ptr<Person> p2) {
    if (!p1 || !p2) return false;
    
    lock_guard<mutex> lock(partnerMutex);
    p1->setPartner(p2);
    p2->setPartner(p1);
    
    cout << p1->mName << " is now partner with " << p2->mName << endl;
    return true;
}

// Thread-safe resource manager
class ResourceManager {
    vector<shared_ptr<Resource>> resources;
    mutex mtx;
public:
    void addResource(shared_ptr<Resource> res) {
        lock_guard<mutex> lock(mtx);
        resources.push_back(res);
    }
    
    void useAllResources() {
        lock_guard<mutex> lock(mtx);
        for (auto& res : resources) {
            if (res) res->use();
        }
    }
    
    size_t count() const {
        return resources.size();
    }
};

// Template function with perfect forwarding
template<typename T, typename U>
auto multiply(T&& a, U&& b) -> decltype(a * b) {
    return forward<T>(a) * forward<U>(b);
}

// Lambda generator with capture modes
auto makeLambda(int& x, int y) {
    return [&x, y](int z) {
        cout << "x=" << x << ", y=" << y << ", z=" << z << endl;
        x += y + z;
        return x;
    };
}
