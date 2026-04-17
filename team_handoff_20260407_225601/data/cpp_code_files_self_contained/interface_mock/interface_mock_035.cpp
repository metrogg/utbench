#include <vector>
#include <map>
#include <cmath>
#include <algorithm>
#include <tuple>
#include <numeric>
#include <iomanip>

using namespace std;

// Represents a user-item rating matrix
class RatingMatrix {
private:
    map<string, map<string, double>> userItemRatings;
    map<string, map<string, double>> itemUserRatings;
    
public:
    // Add a rating to the matrix
    void addRating(const string& user, const string& item, double rating) {
        userItemRatings[user][item] = rating;
        itemUserRatings[item][user] = rating;
    }
    
    // Get user-item ratings
    const map<string, map<string, double>>& getUserItemRatings() const {
        return userItemRatings;
    }
    
    // Get item-user ratings (transposed matrix)
    const map<string, map<string, double>>& getItemUserRatings() const {
        return itemUserRatings;
    }
    
    // Get all items
    vector<string> getAllItems() const {
        vector<string> items;
        for (const auto& pair : itemUserRatings) {
            items.push_back(pair.first);
        }
        return items;
    }
    
    // Get all users
    vector<string> getAllUsers() const {
        vector<string> users;
        for (const auto& pair : userItemRatings) {
            users.push_back(pair.first);
        }
        return users;
    }
};

// Calculates cosine similarity between two vectors
double cosineSimilarity(const map<string, double>& vec1, const map<string, double>& vec2) {
    double dotProduct = 0.0;
    double norm1 = 0.0;
    double norm2 = 0.0;
    
    // Combine keys from both vectors
    vector<string> allKeys;
    for (const auto& pair : vec1) allKeys.push_back(pair.first);
    for (const auto& pair : vec2) {
        if (find(allKeys.begin(), allKeys.end(), pair.first) == allKeys.end()) {
            allKeys.push_back(pair.first);
        }
    }
    
    // Calculate dot product and norms
    for (const string& key : allKeys) {
        double v1 = vec1.count(key) ? vec1.at(key) : 0.0;
        double v2 = vec2.count(key) ? vec2.at(key) : 0.0;
        dotProduct += v1 * v2;
        norm1 += v1 * v1;
        norm2 += v2 * v2;
    }
    
    if (norm1 == 0 || norm2 == 0) return 0.0;
    return dotProduct / (sqrt(norm1) * sqrt(norm2));
}

// Finds similar items using cosine similarity
map<string, vector<pair<string, double>>> findSimilarItems(
    const RatingMatrix& matrix, 
    int k, 
    double similarityThreshold = 0.0) {
    
    map<string, vector<pair<string, double>>> similarItems;
    const auto& itemRatings = matrix.getItemUserRatings();
    vector<string> items = matrix.getAllItems();
    
    for (size_t i = 0; i < items.size(); ++i) {
        const string& item1 = items[i];
        vector<pair<string, double>> similarities;
        
        for (size_t j = 0; j < items.size(); ++j) {
            if (i == j) continue;
            
            const string& item2 = items[j];
            double sim = cosineSimilarity(itemRatings.at(item1), itemRatings.at(item2));
            
            if (sim > similarityThreshold) {
                similarities.emplace_back(item2, sim);
            }
        }
        
        // Sort by similarity and keep top k
        sort(similarities.begin(), similarities.end(), 
             [](const pair<string, double>& a, const pair<string, double>& b) {
                 return a.second > b.second;
             });
        
        if (k > 0 && similarities.size() > static_cast<size_t>(k)) {
            similarities.resize(k);
        }
        
        similarItems[item1] = similarities;
    }
    
    return similarItems;
}

// Predicts ratings using item-item collaborative filtering
map<string, map<string, double>> predictRatings(
    const RatingMatrix& matrix,
    const map<string, vector<pair<string, double>>>& similarItems,
    int k = 5) {
    
    map<string, map<string, double>> predictions;
    const auto& userRatings = matrix.getUserItemRatings();
    vector<string> users = matrix.getAllUsers();
    vector<string> items = matrix.getAllItems();
    
    for (const string& user : users) {
        for (const string& item : items) {
            // Skip if user already rated this item
            if (userRatings.at(user).count(item)) continue;
            
            double weightedSum = 0.0;
            double similaritySum = 0.0;
            int neighborsUsed = 0;
            
            // Find similar items that the user has rated
            if (!similarItems.count(item)) continue;
            
            for (const auto& simPair : similarItems.at(item)) {
                const string& similarItem = simPair.first;
                double similarity = simPair.second;
                
                if (userRatings.at(user).count(similarItem) && neighborsUsed < k) {
                    weightedSum += similarity * userRatings.at(user).at(similarItem);
                    similaritySum += similarity;
                    neighborsUsed++;
                }
            }
            
            if (similaritySum > 0) {
                predictions[user][item] = weightedSum / similaritySum;
            }
        }
    }
    
    return predictions;
}
