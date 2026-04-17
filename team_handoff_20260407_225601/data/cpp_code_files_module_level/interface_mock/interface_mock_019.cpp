#ifndef TENSORSTORE_INTERNAL_DECODED_MATCHES_H_
#define TENSORSTORE_INTERNAL_DECODED_MATCHES_H_
#include <functional>
#include <string>
#include <string_view>
#include <gtest/gtest.h>
#include "absl/strings/cord.h"
#include "tensorstore/util/result.h"
namespace tensorstore {
namespace internal {
::testing::Matcher<absl::Cord> DecodedMatches(
    ::testing::Matcher<std::string_view> value_matcher,
    std::function<Result<std::string>(std::string_view)> decoder);
}  
}  
#endif  
#include "tensorstore/internal/decoded_matches.h"
#include <functional>
#include <ostream>
#include <string>
#include <string_view>
#include <utility>
#include <gtest/gtest.h>
#include "absl/status/status.h"
#include "absl/strings/cord.h"
#include "tensorstore/util/result.h"
namespace tensorstore {
namespace internal {
namespace {
using DecodeFunction = std::function<Result<std::string>(std::string_view)>;
class Matcher : public ::testing::MatcherInterface<absl::Cord> {
 public:
  Matcher(::testing::Matcher<std::string_view> value_matcher,
          DecodeFunction decoder)
      : value_matcher_(std::move(value_matcher)),
        decoder_(std::move(decoder)) {}
  bool MatchAndExplain(
      absl::Cord value,
      ::testing::MatchResultListener* listener) const override {
    auto decoded = decoder_(value.Flatten());
    if (!decoded.ok()) {
      *listener << "Failed to decode value: " << decoded.status();
      return false;
    }
    return value_matcher_.MatchAndExplain(*decoded, listener);
  }
  void DescribeTo(std::ostream* os) const override {
    *os << "when decoded ";
    value_matcher_.DescribeTo(os);
  }
 private:
  ::testing::Matcher<std::string_view> value_matcher_;
  DecodeFunction decoder_;
};
}  
::testing::Matcher<absl::Cord> DecodedMatches(
    ::testing::Matcher<std::string_view> value_matcher,
    DecodeFunction decoder) {
  return ::testing::MakeMatcher(
      new Matcher(std::move(value_matcher), std::move(decoder)));
}
}  
}  