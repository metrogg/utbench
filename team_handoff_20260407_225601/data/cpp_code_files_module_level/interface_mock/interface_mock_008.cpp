#ifndef TENSORFLOW_CORE_FRAMEWORK_TENSOR_MATCHER_H_
#define TENSORFLOW_CORE_FRAMEWORK_TENSOR_MATCHER_H_
#include <gtest/gtest.h>
#include "tensorflow/core/framework/tensor.h"
namespace tensorflow {
namespace test {
class TensorEq {
 public:
  explicit TensorEq(const tensorflow::Tensor& target) : target_(target) {}
  operator ::testing::Matcher<const tensorflow::Tensor&>() const;  
 private:
  const tensorflow::Tensor& target_;
};
}  
}  
#endif  
#include "tensorflow/core/framework/tensor_matcher.h"
#include <stdint.h>
#include <complex>
#include <ostream>
#include <string>
#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include "absl/log/log.h"
#include "absl/types/span.h"
#include "Eigen/Core"  
#include "tensorflow/core/framework/numeric_types.h"
#include "tensorflow/core/framework/register_types.h"
#include "tensorflow/core/framework/tensor.h"
#include "tensorflow/core/framework/tensor_shape.h"
#include "tensorflow/core/framework/tensor_types.h"
#include "tensorflow/core/framework/types.h"
#include "tensorflow/core/platform/bfloat16.h"
#include "tensorflow/core/platform/tstring.h"
#include "tensorflow/core/platform/types.h"
namespace tensorflow {
namespace test {
namespace {
using tensorflow::Tensor;
template <typename T>
::testing::Matcher<absl::Span<const T>> MakePointwiseMatcher(
    absl::Span<const T> target) {
  return ::testing::MatcherCast<absl::Span<const T>>(
      ::testing::Pointwise(::testing::Eq(), target));
}
template <>
::testing::Matcher<absl::Span<const float>> MakePointwiseMatcher(
    absl::Span<const float> target) {
  return ::testing::MatcherCast<absl::Span<const float>>(
      ::testing::Pointwise(::testing::FloatEq(), target));
}
template <>
::testing::Matcher<absl::Span<const double>> MakePointwiseMatcher(
    absl::Span<const double> target) {
  return ::testing::MatcherCast<absl::Span<const double>>(
      ::testing::Pointwise(::testing::DoubleEq(), target));
}
template <typename T>
bool MatchAndExplainPointwise(absl::Span<const T> value,
                              absl::Span<const T> target,
                              ::testing::MatchResultListener* listener) {
  return MakePointwiseMatcher<T>(target).MatchAndExplain(value, listener);
}
class TensorEqMatcherImpl : public ::testing::MatcherInterface<const Tensor&> {
 public:
  explicit TensorEqMatcherImpl(const Tensor& target) : target_(target) {}
  void DescribeTo(::std::ostream* os) const override {
    *os << "data type is " << tensorflow::DataTypeString(target_.dtype())
        << ", and shape is " << target_.shape();
    switch (target_.dtype()) {
#define CASE_TYPE(T)                                       \
  case tensorflow::DataTypeToEnum<T>::value: {             \
    *os << ", and tensor data ";                           \
    absl::Span<const T> data(target_.unaligned_flat<T>()); \
    MakePointwiseMatcher<T>(data).DescribeTo(os);          \
    break;                                                 \
  }
      TF_CALL_POD_STRING_TYPES(CASE_TYPE);
#undef CASE_TYPE
      default: {
        DLOG(FATAL) << "TensorEq matcher unsupported dtype: "
                    << tensorflow::DataTypeString(target_.dtype());
      }
    }
  }
  void DescribeNegationTo(::std::ostream* os) const override {
    *os << "data type is not " << tensorflow::DataTypeString(target_.dtype())
        << ", or shape is not " << target_.shape();
    switch (target_.dtype()) {
#define CASE_TYPE(T)                                       \
  case tensorflow::DataTypeToEnum<T>::value: {             \
    *os << ", or tensor data ";                            \
    absl::Span<const T> data(target_.unaligned_flat<T>()); \
    MakePointwiseMatcher<T>(data).DescribeNegationTo(os);  \
    break;                                                 \
  }
      TF_CALL_POD_STRING_TYPES(CASE_TYPE);
#undef CASE_TYPE
      default: {
        DLOG(FATAL) << "TensorEq matcher unsupported dtype: "
                    << tensorflow::DataTypeString(target_.dtype());
      }
    }
  }
  bool MatchAndExplain(
      const Tensor& value,
      ::testing::MatchResultListener* listener) const override {
    const bool dtype_compare = value.dtype() == target_.dtype();
    *listener << "whose data type " << tensorflow::DataTypeString(value.dtype())
              << (dtype_compare ? " matches " : " doesn't match ")
              << tensorflow::DataTypeString(target_.dtype());
    const bool shape_compare = value.shape() == target_.shape();
    *listener << ", whose shape " << value.shape()
              << (shape_compare ? " matches " : " doesn't match ")
              << target_.shape();
    if (!dtype_compare || !shape_compare) {
      return false;
    }
    bool result;
    switch (target_.dtype()) {
#define CASE_TYPE(T)                                                       \
  case tensorflow::DataTypeToEnum<T>::value: {                             \
    result = MatchAndExplainPointwise<T>(                                  \
        value.unaligned_flat<T>(), target_.unaligned_flat<T>(), listener); \
    break;                                                                 \
  }
      TF_CALL_POD_STRING_TYPES(CASE_TYPE);
      TF_CALL_QUANTIZED_TYPES(CASE_TYPE);
      TF_CALL_int4(CASE_TYPE);
      TF_CALL_uint4(CASE_TYPE);
#undef CASE_TYPE
      default: {
        DLOG(FATAL) << "TensorEq matcher unsupported dtype: "
                    << tensorflow::DataTypeString(target_.dtype());
        result = false;
      }
    }
    return result;
  }
 private:
  const Tensor target_;
};
}  
TensorEq::operator ::testing::Matcher<const Tensor&>() const {
  return ::testing::MakeMatcher(new TensorEqMatcherImpl(target_));
}
}  
}  