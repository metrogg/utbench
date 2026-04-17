#ifndef TENSORSTORE_INTERNAL_PARSE_JSON_MATCHES_H_
#define TENSORSTORE_INTERNAL_PARSE_JSON_MATCHES_H_
#include <string>
#include <gtest/gtest.h>
#include <nlohmann/json.hpp>
#include "tensorstore/internal/json_gtest.h"
namespace tensorstore {
namespace internal {
::testing::Matcher<std::string> ParseJsonMatches(
    ::testing::Matcher<::nlohmann::json> json_matcher);
::testing::Matcher<std::string> ParseJsonMatches(::nlohmann::json json);
}  
}  
#endif  
#include "tensorstore/internal/parse_json_matches.h"
#include <ostream>
#include <string>
#include <utility>
#include <gtest/gtest.h>
#include <nlohmann/json.hpp>
#include "tensorstore/internal/json_binding/json_binding.h"
#include "tensorstore/internal/json_gtest.h"
namespace tensorstore {
namespace internal {
namespace {
class Matcher : public ::testing::MatcherInterface<std::string> {
 public:
  Matcher(::testing::Matcher<::nlohmann::json> json_matcher)
      : json_matcher_(std::move(json_matcher)) {}
  bool MatchAndExplain(
      std::string value,
      ::testing::MatchResultListener* listener) const override {
    return json_matcher_.MatchAndExplain(
        tensorstore::internal::ParseJson(value), listener);
  }
  void DescribeTo(std::ostream* os) const override {
    *os << "when parsed as JSON ";
    json_matcher_.DescribeTo(os);
  }
 private:
  ::testing::Matcher<::nlohmann::json> json_matcher_;
};
}  
::testing::Matcher<std::string> ParseJsonMatches(
    ::testing::Matcher<::nlohmann::json> json_matcher) {
  return ::testing::MakeMatcher(new Matcher(std::move(json_matcher)));
}
::testing::Matcher<std::string> ParseJsonMatches(::nlohmann::json json) {
  return ParseJsonMatches(MatchesJson(json));
}
}  
}  