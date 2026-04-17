#include <vector>
#include <algorithm>
#include <iomanip>
#include <map>
#include <stdexcept>

using namespace std;

const int MAXN = 100001;

struct Person {
    int type;
    int ID;
    int v_grade;
    int t_grade;
    int tot;
    string name;  // Added name field
    int age;      // Added age field
};

enum PersonType {
    SAGE = 0,
    NOBLEMAN = 1,
    FOOLMAN = 2,
    SMALLMAN = 3,
    FAILED = 4
};

// Enhanced comparison function with more sorting criteria
bool comparison(const Person &A, const Person &B) {
    if (A.type == B.type) {
        if (A.tot == B.tot) {
            if (A.v_grade == B.v_grade) {
                if (A.age == B.age) {
                    return A.ID < B.ID;
                }
                return A.age < B.age;  // Younger first if same grades
            }
            return A.v_grade > B.v_grade;
        }
        return A.tot > B.tot;
    }
    return A.type < B.type;
}

// Enhanced classification function with more categories
void classifyPerson(Person &p, int H, int L) {
    if (p.v_grade < L || p.t_grade < L) {
        p.type = FAILED;
        return;
    }
    
    if (p.v_grade >= H && p.t_grade >= H) {
        p.type = SAGE;
    } else if (p.v_grade >= H && p.t_grade < H) {
        p.type = NOBLEMAN;
    } else if (p.v_grade < H && p.t_grade < H && p.v_grade >= p.t_grade) {
        p.type = FOOLMAN;
    } else {
        p.type = SMALLMAN;
    }
    p.tot = p.v_grade + p.t_grade;
}

// Function to process and sort persons
vector<Person> processPersons(const vector<Person>& input, int L, int H) {
    vector<Person> qualified;
    
    for (const auto& p : input) {
        Person temp = p;
        classifyPerson(temp, H, L);
        if (temp.type != FAILED) {
            qualified.push_back(temp);
        }
    }
    
    sort(qualified.begin(), qualified.end(), comparison);
    return qualified;
}

// Function to generate statistics about the results
map<string, int> generateStats(const vector<Person>& persons) {
    map<string, int> stats;
    stats["total"] = persons.size();
    
    for (const auto& p : persons) {
        switch(p.type) {
            case SAGE: stats["sages"]++; break;
            case NOBLEMAN: stats["noblemen"]++; break;
            case FOOLMAN: stats["foolmen"]++; break;
            case SMALLMAN: stats["smallmen"]++; break;
        }
    }
    
    if (!persons.empty()) {
        int min_age = persons[0].age;
        int max_age = persons[0].age;
        double avg_age = 0;
        
        for (const auto& p : persons) {
            if (p.age < min_age) min_age = p.age;
            if (p.age > max_age) max_age = p.age;
            avg_age += p.age;
        }
        avg_age /= persons.size();
        
        stats["min_age"] = min_age;
        stats["max_age"] = max_age;
        stats["avg_age"] = static_cast<int>(avg_age);
    }
    
    return stats;
}
