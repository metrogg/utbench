#ifndef TENSORFLOW_TSL_PROFILER_UTILS_PARSE_ANNOTATION_H_
#define TENSORFLOW_TSL_PROFILER_UTILS_PARSE_ANNOTATION_H_
#include <vector>
#include "absl/strings/string_view.h"
namespace tsl {
namespace profiler {
struct Annotation {
  absl::string_view name;
  struct Metadata {
    absl::string_view key;
    absl::string_view value;
  };
  std::vector<Metadata> metadata;
};
Annotation ParseAnnotation(absl::string_view annotation);
inline bool HasMetadata(absl::string_view annotation) {
  constexpr char kUserMetadataMarker = '#';
  return !annotation.empty() && annotation.back() == kUserMetadataMarker;
}
std::vector<Annotation> ParseAnnotationStack(
    absl::string_view annotation_stack);
}  
}  
#endif  
#include "tsl/profiler/utils/parse_annotation.h"
#include <stack>
#include <string>
#include <utility>
#include <vector>
#include "absl/strings/ascii.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"
namespace tsl {
namespace profiler {
namespace {
std::vector<absl::string_view> SplitNameAndMetadata(
    absl::string_view annotation) {
  std::vector<absl::string_view> parts;
  if (!HasMetadata(annotation)) {
    parts.emplace_back(annotation);
  } else {
    annotation.remove_suffix(1);
    parts = absl::StrSplit(annotation, '#');
    if (parts.size() > 2) {
      parts.resize(2);
    }
  }
  while (parts.size() < 2) {
    parts.emplace_back();
  }
  return parts;
}
std::vector<absl::string_view> SplitPairs(absl::string_view metadata) {
  std::vector<absl::string_view> key_value_pairs;
  std::stack<char> quotes;
  size_t start = 0, end = 0;
  for (; end < metadata.size(); ++end) {
    char ch = metadata[end];
    switch (ch) {
      case '\"':
      case '\'':
        if (quotes.empty() || quotes.top() != ch) {
          quotes.push(ch);
        } else {
          quotes.pop();
        }
        break;
      case '{':
      case '(':
      case '[':
        quotes.push(ch);
        break;
      case '}':
        if (!quotes.empty() && quotes.top() == '{') {
          quotes.pop();
        }
        break;
      case ')':
        if (!quotes.empty() && quotes.top() == '(') {
          quotes.pop();
        }
        break;
      case ']':
        if (!quotes.empty() && quotes.top() == '[') {
          quotes.pop();
        }
        break;
      case ',':
        if (quotes.empty()) {
          if (end - start > 1) {
            key_value_pairs.emplace_back(metadata.data() + start, end - start);
          }
          start = end + 1;  
        }
        break;
    }
  }
  if (end - start > 1) {
    key_value_pairs.emplace_back(metadata.data() + start, end - start);
  }
  return key_value_pairs;
}
std::vector<std::pair<absl::string_view, absl::string_view>> ParseMetadata(
    absl::string_view metadata) {
  std::vector<std::pair<absl::string_view, absl::string_view>> key_values;
  for (absl::string_view pair : SplitPairs(metadata)) {
    std::vector<absl::string_view> parts =
        absl::StrSplit(pair, absl::MaxSplits('=', 1));
    if (parts.size() == 2) {
      absl::string_view key = absl::StripAsciiWhitespace(parts[0]);
      absl::string_view value = absl::StripAsciiWhitespace(parts[1]);
      if (!key.empty() && !value.empty()) {
        key_values.push_back({key, value});
      }
    }
  }
  return key_values;
}
}  
Annotation ParseAnnotation(absl::string_view annotation) {
  Annotation result;
  std::vector<absl::string_view> parts = SplitNameAndMetadata(annotation);
  if (!parts.empty()) {
    result.name = absl::StripAsciiWhitespace(parts[0]);
    for (const auto& key_value : ParseMetadata(parts[1])) {
      result.metadata.push_back({key_value.first, key_value.second});
    }
  }
  return result;
}
std::vector<Annotation> ParseAnnotationStack(
    absl::string_view annotation_stack) {
  std::vector<Annotation> annotations;
  const std::string kAnnotationDelimiter = "::";
  for (absl::string_view annotation : absl::StrSplit(
           annotation_stack, kAnnotationDelimiter, absl::SkipEmpty())) {
    annotations.emplace_back(ParseAnnotation(annotation));
  }
  return annotations;
}
}  
}  