#ifndef TENSORFLOW_TSL_PLATFORM_SCANNER_H_
#define TENSORFLOW_TSL_PLATFORM_SCANNER_H_
#include <string>
#include "tsl/platform/macros.h"
#include "tsl/platform/str_util.h"
#include "tsl/platform/stringpiece.h"
namespace tsl {
namespace strings {
class Scanner {
 public:
  enum CharClass {
    ALL,
    DIGIT,
    LETTER,
    LETTER_DIGIT,
    LETTER_DIGIT_DASH_UNDERSCORE,
    LETTER_DIGIT_DASH_DOT_SLASH,             
    LETTER_DIGIT_DASH_DOT_SLASH_UNDERSCORE,  
    LETTER_DIGIT_DOT,
    LETTER_DIGIT_DOT_PLUS_MINUS,
    LETTER_DIGIT_DOT_UNDERSCORE,
    LETTER_DIGIT_UNDERSCORE,
    LOWERLETTER,
    LOWERLETTER_DIGIT,
    LOWERLETTER_DIGIT_UNDERSCORE,
    NON_ZERO_DIGIT,
    SPACE,
    UPPERLETTER,
    RANGLE,
  };
  explicit Scanner(StringPiece source) : cur_(source) { RestartCapture(); }
  Scanner& One(CharClass clz) {
    if (cur_.empty() || !Matches(clz, cur_[0])) {
      return Error();
    }
    cur_.remove_prefix(1);
    return *this;
  }
  Scanner& ZeroOrOneLiteral(StringPiece s) {
    str_util::ConsumePrefix(&cur_, s);
    return *this;
  }
  Scanner& OneLiteral(StringPiece s) {
    if (!str_util::ConsumePrefix(&cur_, s)) {
      error_ = true;
    }
    return *this;
  }
  Scanner& Any(CharClass clz) {
    while (!cur_.empty() && Matches(clz, cur_[0])) {
      cur_.remove_prefix(1);
    }
    return *this;
  }
  Scanner& Many(CharClass clz) { return One(clz).Any(clz); }
  Scanner& RestartCapture() {
    capture_start_ = cur_.data();
    capture_end_ = nullptr;
    return *this;
  }
  Scanner& StopCapture() {
    capture_end_ = cur_.data();
    return *this;
  }
  Scanner& Eos() {
    if (!cur_.empty()) error_ = true;
    return *this;
  }
  Scanner& AnySpace() { return Any(SPACE); }
  Scanner& ScanUntil(char end_ch) {
    ScanUntilImpl(end_ch, false);
    return *this;
  }
  Scanner& ScanEscapedUntil(char end_ch) {
    ScanUntilImpl(end_ch, true);
    return *this;
  }
  char Peek(char default_value = '\0') const {
    return cur_.empty() ? default_value : cur_[0];
  }
  int empty() const { return cur_.empty(); }
  bool GetResult(StringPiece* remaining = nullptr,
                 StringPiece* capture = nullptr);
 private:
  void ScanUntilImpl(char end_ch, bool escaped);
  Scanner& Error() {
    error_ = true;
    return *this;
  }
  static bool IsLetter(char ch) {
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z');
  }
  static bool IsLowerLetter(char ch) { return ch >= 'a' && ch <= 'z'; }
  static bool IsDigit(char ch) { return ch >= '0' && ch <= '9'; }
  static bool IsSpace(char ch) {
    return (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\v' || ch == '\f' ||
            ch == '\r');
  }
  static bool Matches(CharClass clz, char ch) {
    switch (clz) {
      case ALL:
        return true;
      case DIGIT:
        return IsDigit(ch);
      case LETTER:
        return IsLetter(ch);
      case LETTER_DIGIT:
        return IsLetter(ch) || IsDigit(ch);
      case LETTER_DIGIT_DASH_UNDERSCORE:
        return (IsLetter(ch) || IsDigit(ch) || ch == '-' || ch == '_');
      case LETTER_DIGIT_DASH_DOT_SLASH:
        return IsLetter(ch) || IsDigit(ch) || ch == '-' || ch == '.' ||
               ch == '/';
      case LETTER_DIGIT_DASH_DOT_SLASH_UNDERSCORE:
        return (IsLetter(ch) || IsDigit(ch) || ch == '-' || ch == '.' ||
                ch == '/' || ch == '_');
      case LETTER_DIGIT_DOT:
        return IsLetter(ch) || IsDigit(ch) || ch == '.';
      case LETTER_DIGIT_DOT_PLUS_MINUS:
        return IsLetter(ch) || IsDigit(ch) || ch == '+' || ch == '-' ||
               ch == '.';
      case LETTER_DIGIT_DOT_UNDERSCORE:
        return IsLetter(ch) || IsDigit(ch) || ch == '.' || ch == '_';
      case LETTER_DIGIT_UNDERSCORE:
        return IsLetter(ch) || IsDigit(ch) || ch == '_';
      case LOWERLETTER:
        return ch >= 'a' && ch <= 'z';
      case LOWERLETTER_DIGIT:
        return IsLowerLetter(ch) || IsDigit(ch);
      case LOWERLETTER_DIGIT_UNDERSCORE:
        return IsLowerLetter(ch) || IsDigit(ch) || ch == '_';
      case NON_ZERO_DIGIT:
        return IsDigit(ch) && ch != '0';
      case SPACE:
        return IsSpace(ch);
      case UPPERLETTER:
        return ch >= 'A' && ch <= 'Z';
      case RANGLE:
        return ch == '>';
    }
    return false;
  }
  StringPiece cur_;
  const char* capture_start_ = nullptr;
  const char* capture_end_ = nullptr;
  bool error_ = false;
  friend class ScannerTest;
  Scanner(const Scanner&) = delete;
  void operator=(const Scanner&) = delete;
};
}  
}  
#endif  
#include "tsl/platform/scanner.h"
namespace tsl {
namespace strings {
void Scanner::ScanUntilImpl(char end_ch, bool escaped) {
  for (;;) {
    if (cur_.empty()) {
      Error();
      return;
    }
    const char ch = cur_[0];
    if (ch == end_ch) {
      return;
    }
    cur_.remove_prefix(1);
    if (escaped && ch == '\\') {
      if (cur_.empty()) {
        Error();
        return;
      }
      cur_.remove_prefix(1);
    }
  }
}
bool Scanner::GetResult(StringPiece* remaining, StringPiece* capture) {
  if (error_) {
    return false;
  }
  if (remaining != nullptr) {
    *remaining = cur_;
  }
  if (capture != nullptr) {
    const char* end = capture_end_ == nullptr ? cur_.data() : capture_end_;
    *capture = StringPiece(capture_start_, end - capture_start_);
  }
  return true;
}
}  
}  