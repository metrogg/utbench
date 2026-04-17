#include <vector>
#include <mutex>
#include <condition_variable>
#include <thread>
#include <numeric>
#include <algorithm>
#include <cmath>
#include <map>

class DataProcessor {
private:
    std::vector<int> data;
    std::mutex mtx;
    std::condition_variable cv;
    bool data_ready = false;
    bool processing_complete = false;

    // Calculate statistics on the data
    std::map<std::string, double> calculate_stats(const std::vector<int>& data) {
        std::map<std::string, double> stats;
        if (data.empty()) return stats;

        // Basic statistics
        int sum = std::accumulate(data.begin(), data.end(), 0);
        double mean = static_cast<double>(sum) / data.size();
        
        // Variance and standard deviation
        double variance = 0.0;
        for (int num : data) {
            variance += std::pow(num - mean, 2);
        }
        variance /= data.size();
        double std_dev = std::sqrt(variance);

        // Median
        std::vector<int> sorted_data = data;
        std::sort(sorted_data.begin(), sorted_data.end());
        double median;
        if (sorted_data.size() % 2 == 0) {
            median = (sorted_data[sorted_data.size()/2 - 1] + sorted_data[sorted_data.size()/2]) / 2.0;
        } else {
            median = sorted_data[sorted_data.size()/2];
        }

        // Min and max
        auto minmax = std::minmax_element(data.begin(), data.end());

        stats["sum"] = sum;
        stats["mean"] = mean;
        stats["variance"] = variance;
        stats["std_dev"] = std_dev;
        stats["median"] = median;
        stats["min"] = *minmax.first;
        stats["max"] = *minmax.second;

        return stats;
    }

public:
    // Producer function to add data
    void add_data(const std::vector<int>& new_data) {
        std::unique_lock<std::mutex> lock(mtx);
        data = new_data;
        data_ready = true;
        cv.notify_one();
    }

    // Consumer function to process data
    std::map<std::string, double> process_data() {
        std::unique_lock<std::mutex> lock(mtx);
        cv.wait(lock, [this]{ return data_ready; });
        
        auto stats = calculate_stats(data);
        
        data_ready = false;
        processing_complete = true;
        cv.notify_one();
        
        return stats;
    }

    // Wait for processing to complete
    void wait_for_completion() {
        std::unique_lock<std::mutex> lock(mtx);
        cv.wait(lock, [this]{ return processing_complete; });
        processing_complete = false;
    }
};
