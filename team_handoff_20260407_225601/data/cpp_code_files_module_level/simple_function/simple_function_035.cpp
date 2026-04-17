#ifndef ABSL_RANDOM_INTERNAL_CHI_SQUARE_H_
#define ABSL_RANDOM_INTERNAL_CHI_SQUARE_H_
#include <cassert>
#include "absl/base/config.h"
namespace absl {
ABSL_NAMESPACE_BEGIN
namespace random_internal {
constexpr const char kChiSquared[] = "chi-squared";
template <typename Iterator>
double ChiSquareWithExpected(Iterator begin, Iterator end, double expected) {
  assert(expected >= 10);  
  double chi_square = 0;
  for (auto it = begin; it != end; it++) {
    double d = static_cast<double>(*it) - expected;
    chi_square += d * d;
  }
  chi_square = chi_square / expected;
  return chi_square;
}
template <typename Iterator, typename Expected>
double ChiSquare(Iterator it, Iterator end, Expected eit, Expected eend) {
  double chi_square = 0;
  for (; it != end && eit != eend; ++it, ++eit) {
    if (*it > 0) {
      assert(*eit > 0);
    }
    double e = static_cast<double>(*eit);
    double d = static_cast<double>(*it - *eit);
    if (d != 0) {
      assert(e > 0);
      chi_square += (d * d) / e;
    }
  }
  assert(it == end && eit == eend);
  return chi_square;
}
double ChiSquareValue(int dof, double p);
double ChiSquarePValue(double chi_square, int dof);
}  
ABSL_NAMESPACE_END
}  
#endif  
#include "absl/random/internal/chi_square.h"
#include <cmath>
#include "absl/random/internal/distribution_test_util.h"
namespace absl {
ABSL_NAMESPACE_BEGIN
namespace random_internal {
namespace {
#if defined(__EMSCRIPTEN__)
inline double fma(double x, double y, double z) {
  return (x * y) + z;
}
#endif
template <typename T, unsigned N>
inline T EvaluatePolynomial(T x, const T (&poly)[N]) {
#if !defined(__EMSCRIPTEN__)
  using std::fma;
#endif
  T p = poly[N - 1];
  for (unsigned i = 2; i <= N; i++) {
    p = fma(p, x, poly[N - i]);
  }
  return p;
}
static constexpr int kLargeDOF = 150;
double POZ(double z) {
  static constexpr double kP1[] = {
      0.797884560593,  -0.531923007300, 0.319152932694,
      -0.151968751364, 0.059054035642,  -0.019198292004,
      0.005198775019,  -0.001075204047, 0.000124818987,
  };
  static constexpr double kP2[] = {
      0.999936657524,  0.000535310849,  -0.002141268741, 0.005353579108,
      -0.009279453341, 0.011630447319,  -0.010557625006, 0.006549791214,
      -0.002034254874, -0.000794620820, 0.001390604284,  -0.000676904986,
      -0.000019538132, 0.000152529290,  -0.000045255659,
  };
  const double kZMax = 6.0;  
  if (z == 0.0) {
    return 0.5;
  }
  double x;
  double y = 0.5 * std::fabs(z);
  if (y >= (kZMax * 0.5)) {
    x = 1.0;
  } else if (y < 1.0) {
    double w = y * y;
    x = EvaluatePolynomial(w, kP1) * y * 2.0;
  } else {
    y -= 2.0;
    x = EvaluatePolynomial(y, kP2);
  }
  return z > 0.0 ? ((x + 1.0) * 0.5) : ((1.0 - x) * 0.5);
}
double normal_survival(double z) {
  static constexpr double kR[] = {
      1.0, 0.196854, 0.115194, 0.000344, 0.019527,
  };
  double r = EvaluatePolynomial(z, kR);
  r *= r;
  return 0.5 / (r * r);
}
}  
double ChiSquareValue(int dof, double p) {
  static constexpr double kChiEpsilon =
      0.000001;  
  static constexpr double kChiMax =
      99999.0;  
  const double p_value = 1.0 - p;
  if (dof < 1 || p_value > 1.0) {
    return 0.0;
  }
  if (dof > kLargeDOF) {
    const double z = InverseNormalSurvival(p_value);
    const double mean = 1 - 2.0 / (9 * dof);
    const double variance = 2.0 / (9 * dof);
    if (variance != 0) {
      double term = z * std::sqrt(variance) + mean;
      return dof * (term * term * term);
    }
  }
  if (p_value <= 0.0) return kChiMax;
  double min_chisq = 0.0;
  double max_chisq = kChiMax;
  double current = dof / std::sqrt(p_value);
  while ((max_chisq - min_chisq) > kChiEpsilon) {
    if (ChiSquarePValue(current, dof) < p_value) {
      max_chisq = current;
    } else {
      min_chisq = current;
    }
    current = (max_chisq + min_chisq) * 0.5;
  }
  return current;
}
double ChiSquarePValue(double chi_square, int dof) {
  static constexpr double kLogSqrtPi =
      0.5723649429247000870717135;  
  static constexpr double kInverseSqrtPi =
      0.5641895835477562869480795;  
  if (dof > kLargeDOF) {
    const double chi_square_scaled = std::pow(chi_square / dof, 1.0 / 3);
    const double mean = 1 - 2.0 / (9 * dof);
    const double variance = 2.0 / (9 * dof);
    if (variance != 0) {
      const double z = (chi_square_scaled - mean) / std::sqrt(variance);
      if (z > 0) {
        return normal_survival(z);
      } else if (z < 0) {
        return 1.0 - normal_survival(-z);
      } else {
        return 0.5;
      }
    }
  }
  if (chi_square <= 0.0) return 1.0;
  if (dof < 1) return 0;
  auto capped_exp = [](double x) { return x < -20 ? 0.0 : std::exp(x); };
  static constexpr double kBigX = 20;
  double a = 0.5 * chi_square;
  const bool even = !(dof & 1);  
  const double y = capped_exp(-a);
  double s = even ? y : (2.0 * POZ(-std::sqrt(chi_square)));
  if (dof <= 2) {
    return s;
  }
  chi_square = 0.5 * (dof - 1.0);
  double z = (even ? 1.0 : 0.5);
  if (a > kBigX) {
    double e = (even ? 0.0 : kLogSqrtPi);
    double c = std::log(a);
    while (z <= chi_square) {
      e = std::log(z) + e;
      s += capped_exp(c * z - a - e);
      z += 1.0;
    }
    return s;
  }
  double e = (even ? 1.0 : (kInverseSqrtPi / std::sqrt(a)));
  double c = 0.0;
  while (z <= chi_square) {
    e = e * (a / z);
    c = c + e;
    z += 1.0;
  }
  return c * y + s;
}
}  
ABSL_NAMESPACE_END
}  