#ifndef AROLLA_UTIL_BINARY_SEARCH_H_
#define AROLLA_UTIL_BINARY_SEARCH_H_
#include <cstddef>
#include <cstdint>
#include <optional>
#include "absl/base/attributes.h"
#include "absl/types/span.h"
namespace arolla {
size_t LowerBound(float value, absl::Span<const float> array);
size_t LowerBound(double value, absl::Span<const double> array);
size_t LowerBound(int32_t value, absl::Span<const int32_t> array);
size_t LowerBound(int64_t value, absl::Span<const int64_t> array);
size_t UpperBound(float value, absl::Span<const float> array);
size_t UpperBound(double value, absl::Span<const double> array);
size_t UpperBound(int32_t value, absl::Span<const int32_t> array);
size_t UpperBound(int64_t value, absl::Span<const int64_t> array);
template <typename T, typename Iter>
Iter GallopingLowerBound(Iter begin, Iter end, const T& value);
}  
namespace arolla::binary_search_details {
constexpr size_t kSupremacySizeThreshold = 1'000'000;
template <typename T>
size_t LowerBound(T value, absl::Span<const T> array);
template <typename T>
size_t UpperBound(T value, absl::Span<const T> array);
template <typename T, typename Predicate>
inline ABSL_ATTRIBUTE_ALWAYS_INLINE std::optional<size_t> SmallLinearSearch(
    absl::Span<const T> array, Predicate predicate) {
  if (array.size() <= 2) {
    if (array.empty() || predicate(array[0])) {
      return 0;
    } else if (array.size() == 1 || predicate(array[1])) {
      return 1;
    }
    return 2;
  }
  return std::nullopt;
}
size_t UpperBoundImpl(float value, absl::Span<const float> array);
size_t UpperBoundImpl(double value, absl::Span<const double> array);
size_t UpperBoundImpl(int32_t value, absl::Span<const int32_t> array);
size_t UpperBoundImpl(int64_t value, absl::Span<const int64_t> array);
size_t LowerBoundImpl(float value, absl::Span<const float> array);
size_t LowerBoundImpl(double value, absl::Span<const double> array);
size_t LowerBoundImpl(int32_t value, absl::Span<const int32_t> array);
size_t LowerBoundImpl(int64_t value, absl::Span<const int64_t> array);
template <typename T>
inline ABSL_ATTRIBUTE_ALWAYS_INLINE size_t
LowerBound(T value, absl::Span<const T> array) {
  if (auto result =
          SmallLinearSearch(array, [value](T arg) { return !(arg < value); })) {
    return *result;
  }
  return LowerBoundImpl(value, array);
}
template <typename T>
inline ABSL_ATTRIBUTE_ALWAYS_INLINE size_t
UpperBound(T value, absl::Span<const T> array) {
  if (auto result =
          SmallLinearSearch(array, [value](T arg) { return value < arg; })) {
    return *result;
  }
  return UpperBoundImpl(value, array);
}
}  
namespace arolla {
inline size_t LowerBound(float value, absl::Span<const float> array) {
  return binary_search_details::LowerBound<float>(value, array);
}
inline size_t LowerBound(double value, absl::Span<const double> array) {
  return binary_search_details::LowerBound<double>(value, array);
}
inline size_t LowerBound(int32_t value, absl::Span<const int32_t> array) {
  return binary_search_details::LowerBound<int32_t>(value, array);
}
inline size_t LowerBound(int64_t value, absl::Span<const int64_t> array) {
  return binary_search_details::LowerBound<int64_t>(value, array);
}
inline size_t UpperBound(float value, absl::Span<const float> array) {
  return binary_search_details::UpperBound<float>(value, array);
}
inline size_t UpperBound(double value, absl::Span<const double> array) {
  return binary_search_details::UpperBound<double>(value, array);
}
inline size_t UpperBound(int32_t value, absl::Span<const int32_t> array) {
  return binary_search_details::UpperBound<int32_t>(value, array);
}
inline size_t UpperBound(int64_t value, absl::Span<const int64_t> array) {
  return binary_search_details::UpperBound<int64_t>(value, array);
}
template <typename T, typename Iter>
Iter GallopingLowerBound(Iter begin, Iter end, const T& value) {
  size_t i = 0;
  size_t size = end - begin;
  if (begin >= end || !(*begin < value)) {
    return std::min<Iter>(begin, end);
  }
  size_t d = 1;
  while (i + d < size && begin[i + d] < value) {
    i += d;
    d <<= 1;
  }
  while (d > 1) {
    d >>= 1;
    if (i + d < size && begin[i + d] < value) {
      i += d;
    }
  }
  return begin + i + 1;
}
}  
#endif  
#include "arolla/util/binary_search.h"
#include <cassert>
#include <cmath>
#include <cstddef>
#include <cstdint>
#include "absl/types/span.h"
#include "arolla/util/bits.h"
#include "arolla/util/switch_index.h"
namespace arolla::binary_search_details {
namespace {
template <size_t kArraySize, typename T, class Predicate>
size_t FastBinarySearchT(const T* const array, Predicate predicate) {
  static_assert((kArraySize & (kArraySize + 1)) == 0);
  size_t offset = 0;
  for (size_t k = kArraySize; k > 0;) {
    k >>= 1;
    offset = (!predicate(array[offset + k]) ? offset + k + 1 : offset);
  }
  return offset;
}
template <typename T, typename Predicate>
size_t BinarySearchT(absl::Span<const T> array, Predicate predicate) {
  assert(!array.empty());
  const int log2_size = BitScanReverse(array.size());
  return switch_index<8 * sizeof(size_t)>(
      log2_size, [array, predicate](auto constexpr_log2_size) {
        constexpr size_t size =
            (1ULL << static_cast<int>(constexpr_log2_size)) - 1;
        size_t offset = 0;
#if !defined(__clang__) && defined(__GNUC__)
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Warray-bounds"
#endif
        offset = (!predicate(array[size]) ? array.size() - size : offset);
#if !defined(__clang__) && defined(__GNUC__)
#pragma GCC diagnostic pop
#endif
        return offset +
               FastBinarySearchT<size>(array.begin() + offset, predicate);
      });
}
}  
size_t LowerBoundImpl(float value, absl::Span<const float> array) {
  return BinarySearchT(array, [value](auto arg) { return !(arg < value); });
}
size_t LowerBoundImpl(double value, absl::Span<const double> array) {
  return BinarySearchT(array, [value](auto arg) { return !(arg < value); });
}
size_t LowerBoundImpl(int32_t value, absl::Span<const int32_t> array) {
  return BinarySearchT(array, [value](auto arg) { return arg >= value; });
}
size_t LowerBoundImpl(int64_t value, absl::Span<const int64_t> array) {
  return BinarySearchT(array, [value](auto arg) { return arg >= value; });
}
size_t UpperBoundImpl(float value, absl::Span<const float> array) {
  if (std::isnan(value)) {
    return array.size();
  }
  return BinarySearchT(array, [value](auto arg) { return !(arg <= value); });
}
size_t UpperBoundImpl(double value, absl::Span<const double> array) {
  if (std::isnan(value)) {
    return array.size();
  }
  return BinarySearchT(array, [value](auto arg) { return !(arg <= value); });
}
size_t UpperBoundImpl(int32_t value, absl::Span<const int32_t> array) {
  return BinarySearchT(array, [value](auto arg) { return arg > value; });
}
size_t UpperBoundImpl(int64_t value, absl::Span<const int64_t> array) {
  return BinarySearchT(array, [value](auto arg) { return arg > value; });
}
}  