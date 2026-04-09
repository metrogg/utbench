#include <vector>
#include <string>
#include <map>
#include <memory>
#include <algorithm>
#include <stdexcept>

using namespace std;

// Enhanced GameObject class with more functionality
class GameObject {
public:
    string name;
    string tag;
    bool isActive;
    map<string, string> properties;

    GameObject(const string& name = "", const string& tag = "") 
        : name(name), tag(tag), isActive(true) {}

    virtual ~GameObject() = default;

    virtual unique_ptr<GameObject> Clone() const {
        auto newObj = make_unique<GameObject>(name, tag);
        newObj->isActive = isActive;
        newObj->properties = properties;
        return newObj;
    }

    virtual void Update(float deltaTime) {
        // Base game object has no update logic
    }

    void SetProperty(const string& key, const string& value) {
        properties[key] = value;
    }

    string GetProperty(const string& key) const {
        auto it = properties.find(key);
        return it != properties.end() ? it->second : "";
    }
};

// Enhanced Level class with object management
class Level {
    vector<unique_ptr<GameObject>> objects;
    string name;
    bool isActive;

public:
    Level(const string& name = "") : name(name), isActive(false) {}

    void AddObject(unique_ptr<GameObject> obj) {
        objects.push_back(move(obj));
    }

    GameObject* FindObject(const string& name) const {
        auto it = find_if(objects.begin(), objects.end(),
            [&name](const unique_ptr<GameObject>& obj) {
                return obj->name == name;
            });
        return it != objects.end() ? it->get() : nullptr;
    }

    vector<GameObject*> FindObjectsByTag(const string& tag) const {
        vector<GameObject*> result;
        for (const auto& obj : objects) {
            if (obj->tag == tag) {
                result.push_back(obj.get());
            }
        }
        return result;
    }

    void DeleteObject(GameObject* obj) {
        objects.erase(remove_if(objects.begin(), objects.end(),
            [obj](const unique_ptr<GameObject>& ptr) {
                return ptr.get() == obj;
            }), objects.end());
    }

    void Update(float deltaTime) {
        for (auto& obj : objects) {
            if (obj->isActive) {
                obj->Update(deltaTime);
            }
        }
    }

    size_t ObjectCount() const { return objects.size(); }
    void SetActive(bool active) { isActive = active; }
    bool IsActive() const { return isActive; }
    const string& GetName() const { return name; }
};

// Enhanced Game Manager class
class GameManager {
    static GameManager* instance;
    vector<unique_ptr<Level>> levels;
    int currentLevelIndex;
    map<string, string> globalProperties;

    GameManager() : currentLevelIndex(-1) {}

public:
    static GameManager& GetInstance() {
        if (!instance) {
            instance = new GameManager();
        }
        return *instance;
    }

    Level* AddLevel(const string& name) {
        levels.push_back(make_unique<Level>(name));
        if (currentLevelIndex == -1) {
            currentLevelIndex = 0;
        }
        return levels.back().get();
    }

    void LoadLevel(int index) {
        if (index >= 0 && index < static_cast<int>(levels.size())) {
            if (currentLevelIndex != -1) {
                levels[currentLevelIndex]->SetActive(false);
            }
            currentLevelIndex = index;
            levels[currentLevelIndex]->SetActive(true);
        }
    }

    GameObject* SpawnCopy(const GameObject* original) {
        if (currentLevelIndex == -1) return nullptr;
        auto newObj = original->Clone();
        auto rawPtr = newObj.get();
        levels[currentLevelIndex]->AddObject(move(newObj));
        return rawPtr;
    }

    void DestroyObject(GameObject* obj) {
        if (currentLevelIndex != -1) {
            levels[currentLevelIndex]->DeleteObject(obj);
        }
    }

    GameObject* FindObject(const string& name) const {
        if (currentLevelIndex == -1) return nullptr;
        return levels[currentLevelIndex]->FindObject(name);
    }

    vector<GameObject*> FindObjectsByTag(const string& tag) const {
        if (currentLevelIndex == -1) return {};
        return levels[currentLevelIndex]->FindObjectsByTag(tag);
    }

    void Update(float deltaTime) {
        if (currentLevelIndex != -1) {
            levels[currentLevelIndex]->Update(deltaTime);
        }
    }

    void SetGlobalProperty(const string& key, const string& value) {
        globalProperties[key] = value;
    }

    string GetGlobalProperty(const string& key) const {
        auto it = globalProperties.find(key);
        return it != globalProperties.end() ? it->second : "";
    }

    size_t LevelCount() const { return levels.size(); }
    int CurrentLevelIndex() const { return currentLevelIndex; }
    Level* CurrentLevel() const {
        return currentLevelIndex != -1 ? levels[currentLevelIndex].get() : nullptr;
    }
};

GameManager* GameManager::instance = nullptr;
