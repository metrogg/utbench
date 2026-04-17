#include <chrono>
#include <vector>
#include <string>
#include <iomanip>
#include <map>
#include <numeric>
#include <cmath>
#include <cstdarg>
#include <algorithm>

using namespace std;
using namespace chrono;

class Benchmark {
private:
    time_point<system_clock> start_time;
    string delimiter;
    vector<pair<string, duration<double>>> records;
    map<string, vector<duration<double>>> grouped_metrics;
    bool scientific_notation;
    int precision;

    void update_grouped_metrics(const string& label, duration<double> elapsed) {
        grouped_metrics[label].push_back(elapsed);
    }

public:
    explicit Benchmark(const string& delim = ",", bool scientific = true, int prec = 6)
        : delimiter(delim), scientific_notation(scientific), precision(prec) {
        start_time = system_clock::now();
        if (scientific_notation) {
            cout << scientific << setprecision(precision);
        }
    }

    void reset_timer() {
        start_time = system_clock::now();
    }

    void set_delimiter(const string& delim) {
        delimiter = delim;
    }

    void set_output_format(bool scientific, int prec) {
        scientific_notation = scientific;
        precision = prec;
        if (scientific_notation) {
            cout << scientific << setprecision(precision);
        } else {
            cout << fixed << setprecision(precision);
        }
    }

    duration<double> lap(const string& label = "") {
        auto end = system_clock::now();
        auto elapsed = end - start_time;
        if (!label.empty()) {
            records.emplace_back(label, elapsed);
            update_grouped_metrics(label, elapsed);
        }
        start_time = end;
        return elapsed;
    }

    void print_stats(const string& label = "") const {
        if (grouped_metrics.empty()) return;

        auto metrics = label.empty() ? grouped_metrics.begin()->second 
                                   : grouped_metrics.at(label);

        double sum = 0.0;
        for (const auto& d : metrics) {
            sum += d.count();
        }
        double mean = sum / metrics.size();

        double variance = 0.0;
        for (const auto& d : metrics) {
            variance += pow(d.count() - mean, 2);
        }
        variance /= metrics.size();
        double stddev = sqrt(variance);

        auto minmax = minmax_element(metrics.begin(), metrics.end());
        double min_val = minmax.first->count();
        double max_val = minmax.second->count();

        cout << "Benchmark stats" << (label.empty() ? "" : " for " + label) << ":\n";
        cout << "Count: " << metrics.size() << "\n";
        cout << "Mean: " << mean << " s\n";
        cout << "StdDev: " << stddev << " s\n";
        cout << "Min: " << min_val << " s\n";
        cout << "Max: " << max_val << " s\n";
    }

    template<typename... Args>
    void log(Args&&... args) {
        (cout << ... << args) << delimiter;
    }

    template<typename... Args>
    void log_time(Args&&... args) {
        auto now = system_clock::now();
        duration<double> elapsed = now - start_time;
        cout << elapsed.count() << delimiter;
        (cout << ... << args) << endl;
    }

    void save_to_csv(const string& filename) const {
        // Implementation would write records to CSV file
        // Not implemented to keep code self-contained
    }
};
