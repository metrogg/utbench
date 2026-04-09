#include <vector>
#include <cmath>
#include <algorithm>
#include <numeric>
#include <iomanip>
#include <map>

using namespace std;

// Function to calculate linear regression (trend line) parameters
map<string, double> calculate_linear_trend(const vector<double>& x, const vector<double>& y) {
    if (x.size() != y.size() || x.empty()) {
        throw invalid_argument("Invalid input data for trend calculation");
    }

    double n = x.size();
    double sum_x = accumulate(x.begin(), x.end(), 0.0);
    double sum_y = accumulate(y.begin(), y.end(), 0.0);
    double sum_xy = inner_product(x.begin(), x.end(), y.begin(), 0.0);
    double sum_x2 = inner_product(x.begin(), x.end(), x.begin(), 0.0);

    double slope = (n * sum_xy - sum_x * sum_y) / (n * sum_x2 - sum_x * sum_x);
    double intercept = (sum_y - slope * sum_x) / n;

    // Calculate R-squared value
    double y_mean = sum_y / n;
    double ss_tot = 0.0, ss_res = 0.0;
    for (size_t i = 0; i < n; ++i) {
        ss_tot += pow(y[i] - y_mean, 2);
        double y_pred = slope * x[i] + intercept;
        ss_res += pow(y[i] - y_pred, 2);
    }
    double r_squared = 1.0 - (ss_res / ss_tot);

    return {{"slope", slope}, {"intercept", intercept}, {"r_squared", r_squared}};
}

// Function to calculate exponential trend line parameters
map<string, double> calculate_exponential_trend(const vector<double>& x, const vector<double>& y) {
    if (x.size() != y.size() || x.empty()) {
        throw invalid_argument("Invalid input data for trend calculation");
    }

    // Transform y values by taking natural log
    vector<double> log_y(y.size());
    transform(y.begin(), y.end(), log_y.begin(), [](double val) { 
        if (val <= 0) throw invalid_argument("y values must be positive for exponential trend");
        return log(val); 
    });

    // Perform linear regression on x and log(y)
    auto lin_result = calculate_linear_trend(x, log_y);
    double a = exp(lin_result["intercept"]);
    double b = lin_result["slope"];

    return {{"a", a}, {"b", b}, {"r_squared", lin_result["r_squared"]}};
}

// Function to calculate logarithmic trend line parameters
map<string, double> calculate_logarithmic_trend(const vector<double>& x, const vector<double>& y) {
    if (x.size() != y.size() || x.empty()) {
        throw invalid_argument("Invalid input data for trend calculation");
    }

    // Transform x values by taking natural log
    vector<double> log_x(x.size());
    transform(x.begin(), x.end(), log_x.begin(), [](double val) { 
        if (val <= 0) throw invalid_argument("x values must be positive for logarithmic trend");
        return log(val); 
    });

    // Perform linear regression on log(x) and y
    auto lin_result = calculate_linear_trend(log_x, y);

    return {{"a", lin_result["intercept"]}, {"b", lin_result["slope"]}, {"r_squared", lin_result["r_squared"]}};
}

// Function to calculate polynomial trend line parameters (2nd order)
map<string, vector<double>> calculate_polynomial_trend(const vector<double>& x, const vector<double>& y, int order = 2) {
    if (x.size() != y.size() || x.empty() || order < 1) {
        throw invalid_argument("Invalid input data for polynomial trend");
    }

    // Create Vandermonde matrix
    size_t n = x.size();
    vector<vector<double>> A(order + 1, vector<double>(order + 1, 0));
    vector<double> B(order + 1, 0);

    for (size_t i = 0; i <= order; ++i) {
        for (size_t j = 0; j <= order; ++j) {
            for (size_t k = 0; k < n; ++k) {
                A[i][j] += pow(x[k], i + j);
            }
        }
        for (size_t k = 0; k < n; ++k) {
            B[i] += pow(x[k], i) * y[k];
        }
    }

    // Solve the system of equations (using simple Gaussian elimination)
    for (size_t i = 0; i <= order; ++i) {
        // Search for maximum in this column
        double maxEl = abs(A[i][i]);
        size_t maxRow = i;
        for (size_t k = i + 1; k <= order; ++k) {
            if (abs(A[k][i]) > maxEl) {
                maxEl = abs(A[k][i]);
                maxRow = k;
            }
        }

        // Swap maximum row with current row
        for (size_t k = i; k <= order; ++k) {
            swap(A[maxRow][k], A[i][k]);
        }
        swap(B[maxRow], B[i]);

        // Make all rows below this one 0 in current column
        for (size_t k = i + 1; k <= order; ++k) {
            double c = -A[k][i] / A[i][i];
            for (size_t j = i; j <= order; ++j) {
                if (i == j) {
                    A[k][j] = 0;
                } else {
                    A[k][j] += c * A[i][j];
                }
            }
            B[k] += c * B[i];
        }
    }

    // Solve equation Ax = B for an upper triangular matrix A
    vector<double> coefficients(order + 1);
    for (int i = order; i >= 0; --i) {
        coefficients[i] = B[i] / A[i][i];
        for (int k = i - 1; k >= 0; --k) {
            B[k] -= A[k][i] * coefficients[i];
        }
    }

    // Calculate R-squared
    double y_mean = accumulate(y.begin(), y.end(), 0.0) / n;
    double ss_tot = 0.0, ss_res = 0.0;
    for (size_t i = 0; i < n; ++i) {
        double y_pred = 0.0;
        for (int j = 0; j <= order; ++j) {
            y_pred += coefficients[j] * pow(x[i], j);
        }
        ss_tot += pow(y[i] - y_mean, 2);
        ss_res += pow(y[i] - y_pred, 2);
    }
    double r_squared = 1.0 - (ss_res / ss_tot);

    return {{"coefficients", coefficients}, {"r_squared", {r_squared}}};
}

// Function to calculate moving average
vector<double> calculate_moving_average(const vector<double>& y, int period) {
    if (y.empty() || period <= 0 || period > y.size()) {
        throw invalid_argument("Invalid period for moving average");
    }

    vector<double> moving_avg;
    for (size_t i = period - 1; i < y.size(); ++i) {
        double sum = 0.0;
        for (int j = 0; j < period; ++j) {
            sum += y[i - j];
        }
        moving_avg.push_back(sum / period);
    }

    return moving_avg;
}
