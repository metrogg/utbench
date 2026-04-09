#include <vector>
#include <algorithm>
#include <string>
#include <memory>

class Observer {
public:
    virtual ~Observer() = default;
    virtual void update(const std::string& message) = 0;
    virtual std::string getName() const = 0;
};

class Observable {
protected:
    std::vector<Observer*> observers;
    std::string state;

public:
    virtual ~Observable() = default;
    
    void addObserver(Observer* observer) {
        if (observer && std::find(observers.begin(), observers.end(), observer) == observers.end()) {
            observers.push_back(observer);
        }
    }

    void removeObserver(Observer* observer) {
        if (observer) {
            observers.erase(std::remove(observers.begin(), observers.end(), observer), observers.end());
        }
    }

    void notifyObservers() {
        for (Observer* observer : observers) {
            if (observer) {
                observer->update(state);
            }
        }
    }

    virtual void setState(const std::string& newState) {
        state = newState;
        notifyObservers();
    }

    std::string getState() const {
        return state;
    }

    size_t countObservers() const {
        return observers.size();
    }
};

class NewsAgency : public Observable {
public:
    void publishNews(const std::string& news) {
        setState("BREAKING: " + news);
    }
};

class NewsChannel : public Observer {
    std::string name;
    std::vector<std::string> receivedNews;

public:
    explicit NewsChannel(const std::string& name) : name(name) {}

    void update(const std::string& message) override {
        receivedNews.push_back(message);
    }

    std::string getName() const override {
        return name;
    }

    size_t newsCount() const {
        return receivedNews.size();
    }

    std::string getLatestNews() const {
        return receivedNews.empty() ? "No news received" : receivedNews.back();
    }

    void clearNews() {
        receivedNews.clear();
    }
};
