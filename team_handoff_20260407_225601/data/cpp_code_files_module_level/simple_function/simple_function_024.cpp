#ifndef TENSORFLOW_TSL_PLATFORM_STR_UTIL_H_
#define TENSORFLOW_TSL_PLATFORM_STR_UTIL_H_
#include <cstdint>
#include <string>
#include <vector>
#include "absl/strings/str_join.h"
#include "absl/strings/str_split.h"
#include "tsl/platform/macros.h"
#include "tsl/platform/stringpiece.h"
#include "tsl/platform/types.h"
namespace tsl {
namespace str_util {
std::string CEscape(StringPiece src);
bool CUnescape(StringPiece source, std::string* dest, std::string* error);
void StripTrailingWhitespace(std::string* s);
size_t RemoveLeadingWhitespace(StringPiece* text);
size_t RemoveTrailingWhitespace(StringPiece* text);
size_t RemoveWhitespaceContext(StringPiece* text);
bool ConsumeLeadingDigits(StringPiece* s, uint64_t* val);
bool ConsumeNonWhitespace(StringPiece* s, StringPiece* val);
bool ConsumePrefix(StringPiece* s, StringPiece expected);
bool ConsumeSuffix(StringPiece* s, StringPiece expected);
TF_MUST_USE_RESULT StringPiece StripPrefix(StringPiece s, StringPiece expected);
TF_MUST_USE_RESULT StringPiece StripSuffix(StringPiece s, StringPiece expected);
std::string Lowercase(StringPiece s);
std::string Uppercase(StringPiece s);
void TitlecaseString(std::string* s, StringPiece delimiters);
std::string StringReplace(StringPiece s, StringPiece oldsub, StringPiece newsub,
                          bool replace_all);
template <typename T>
std::string Join(const T& s, const char* sep) {
  return absl::StrJoin(s, sep);
}
template <typename T, typename Formatter>
std::string Join(const T& s, const char* sep, Formatter f) {
  return absl::StrJoin(s, sep, f);
}
struct AllowEmpty {
  bool operator()(StringPiece sp) const { return true; }
};
struct SkipEmpty {
  bool operator()(StringPiece sp) const { return !sp.empty(); }
};
struct SkipWhitespace {
  bool operator()(StringPiece sp) const {
    return !absl::StripTrailingAsciiWhitespace(sp).empty();
  }
};
inline std::vector<string> Split(StringPiece text, StringPiece delims) {
  return text.empty() ? std::vector<string>()
                      : absl::StrSplit(text, absl::ByAnyChar(delims));
}
template <typename Predicate>
std::vector<string> Split(StringPiece text, StringPiece delims, Predicate p) {
  return text.empty() ? std::vector<string>()
                      : absl::StrSplit(text, absl::ByAnyChar(delims), p);
}
inline std::vector<string> Split(StringPiece text, char delim) {
  return text.empty() ? std::vector<string>() : absl::StrSplit(text, delim);
}
template <typename Predicate>
std::vector<string> Split(StringPiece text, char delim, Predicate p) {
  return text.empty() ? std::vector<string>() : absl::StrSplit(text, delim, p);
}
bool StartsWith(StringPiece text, StringPiece prefix);
bool EndsWith(StringPiece text, StringPiece suffix);
bool StrContains(StringPiece haystack, StringPiece needle);
size_t Strnlen(const char* str, const size_t string_max_len);
std::string ArgDefCase(StringPiece s);
}  
}  
#endif  
#include "tsl/platform/str_util.h"
#include <cctype>
#include <cstdint>
#include <string>
#include <vector>
#include "absl/strings/ascii.h"
#include "absl/strings/escaping.h"
#include "absl/strings/match.h"
#include "absl/strings/strip.h"
#include "tsl/platform/logging.h"
#include "tsl/platform/stringpiece.h"
namespace tsl {
namespace str_util {
string CEscape(StringPiece src) { return absl::CEscape(src); }
bool CUnescape(StringPiece source, string* dest, string* error) {
  return absl::CUnescape(source, dest, error);
}
void StripTrailingWhitespace(string* s) {
  absl::StripTrailingAsciiWhitespace(s);
}
size_t RemoveLeadingWhitespace(StringPiece* text) {
  absl::string_view new_text = absl::StripLeadingAsciiWhitespace(*text);
  size_t count = text->size() - new_text.size();
  *text = new_text;
  return count;
}
size_t RemoveTrailingWhitespace(StringPiece* text) {
  absl::string_view new_text = absl::StripTrailingAsciiWhitespace(*text);
  size_t count = text->size() - new_text.size();
  *text = new_text;
  return count;
}
size_t RemoveWhitespaceContext(StringPiece* text) {
  absl::string_view new_text = absl::StripAsciiWhitespace(*text);
  size_t count = text->size() - new_text.size();
  *text = new_text;
  return count;
}
bool ConsumeLeadingDigits(StringPiece* s, uint64_t* val) {
  const char* p = s->data();
  const char* limit = p + s->size();
  uint64_t v = 0;
  while (p < limit) {
    const char c = *p;
    if (c < '0' || c > '9') break;
    uint64_t new_v = (v * 10) + (c - '0');
    if (new_v / 8 < v) {
      return false;
    }
    v = new_v;
    p++;
  }
  if (p > s->data()) {
    s->remove_prefix(p - s->data());
    *val = v;
    return true;
  } else {
    return false;
  }
}
bool ConsumeNonWhitespace(StringPiece* s, StringPiece* val) {
  const char* p = s->data();
  const char* limit = p + s->size();
  while (p < limit) {
    const char c = *p;
    if (isspace(c)) break;
    p++;
  }
  const size_t n = p - s->data();
  if (n > 0) {
    *val = StringPiece(s->data(), n);
    s->remove_prefix(n);
    return true;
  } else {
    *val = StringPiece();
    return false;
  }
}
bool ConsumePrefix(StringPiece* s, StringPiece expected) {
  return absl::ConsumePrefix(s, expected);
}
bool ConsumeSuffix(StringPiece* s, StringPiece expected) {
  return absl::ConsumeSuffix(s, expected);
}
StringPiece StripPrefix(StringPiece s, StringPiece expected) {
  return absl::StripPrefix(s, expected);
}
StringPiece StripSuffix(StringPiece s, StringPiece expected) {
  return absl::StripSuffix(s, expected);
}
string Lowercase(StringPiece s) { return absl::AsciiStrToLower(s); }
string Uppercase(StringPiece s) { return absl::AsciiStrToUpper(s); }
void TitlecaseString(string* s, StringPiece delimiters) {
  bool upper = true;
  for (string::iterator ss = s->begin(); ss != s->end(); ++ss) {
    if (upper) {
      *ss = toupper(*ss);
    }
    upper = (delimiters.find(*ss) != StringPiece::npos);
  }
}
string StringReplace(StringPiece s, StringPiece oldsub, StringPiece newsub,
                     bool replace_all) {
  string res(s);
  size_t pos = 0;
  while ((pos = res.find(oldsub.data(), pos, oldsub.size())) != string::npos) {
    res.replace(pos, oldsub.size(), newsub.data(), newsub.size());
    pos += newsub.size();
    if (oldsub.empty()) {
      pos++;  
    }
    if (!replace_all) {
      break;
    }
  }
  return res;
}
bool StartsWith(StringPiece text, StringPiece prefix) {
  return absl::StartsWith(text, prefix);
}
bool EndsWith(StringPiece text, StringPiece suffix) {
  return absl::EndsWith(text, suffix);
}
bool StrContains(StringPiece haystack, StringPiece needle) {
  return absl::StrContains(haystack, needle);
}
size_t Strnlen(const char* str, const size_t string_max_len) {
  size_t len = 0;
  while (len < string_max_len && str[len] != '\0') {
    ++len;
  }
  return len;
}
string ArgDefCase(StringPiece s) {
  const size_t n = s.size();
  size_t extra_us = 0;
  size_t to_skip = 0;
  for (size_t i = 0; i < n; ++i) {
    if (i == to_skip && !isalpha(s[i])) {
      ++to_skip;
      continue;
    }
    if (isupper(s[i]) && i != to_skip && i > 0 && isalnum(s[i - 1])) {
      ++extra_us;
    }
  }
  string result(n + extra_us - to_skip, '_');
  for (size_t i = to_skip, j = 0; i < n; ++i, ++j) {
    DCHECK_LT(j, result.size());
    char c = s[i];
    if (isalnum(c)) {
      if (isupper(c)) {
        if (i != to_skip) {
          DCHECK_GT(j, 0);
          if (result[j - 1] != '_') ++j;
        }
        result[j] = tolower(c);
      } else {
        result[j] = c;
      }
    }
  }
  return result;
}
}  
}  