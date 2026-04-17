#ifndef QUICHE_QUIC_CORE_QUIC_CONNECTION_CONTEXT_H_
#define QUICHE_QUIC_CORE_QUIC_CONNECTION_CONTEXT_H_
#include <memory>
#include "absl/strings/str_format.h"
#include "absl/strings/string_view.h"
#include "quiche/quic/platform/api/quic_export.h"
#include "quiche/common/platform/api/quiche_logging.h"
namespace quic {
class QUICHE_EXPORT QuicConnectionTracer {
 public:
  virtual ~QuicConnectionTracer() = default;
  virtual void PrintLiteral(const char* literal) = 0;
  virtual void PrintString(absl::string_view s) = 0;
  template <typename... Args>
  void Printf(const absl::FormatSpec<Args...>& format, const Args&... args) {
    std::string s = absl::StrFormat(format, args...);
    PrintString(s);
  }
};
class QUICHE_EXPORT QuicBugListener {
 public:
  virtual ~QuicBugListener() = default;
  virtual void OnQuicBug(const char* bug_id, const char* file, int line,
                         absl::string_view bug_message) = 0;
};
class QUICHE_EXPORT QuicConnectionContextListener {
 public:
  virtual ~QuicConnectionContextListener() = default;
 private:
  friend class QuicConnectionContextSwitcher;
  virtual void Activate() = 0;
  virtual void Deactivate() = 0;
};
struct QUICHE_EXPORT QuicConnectionContext final {
  static QuicConnectionContext* Current();
  std::unique_ptr<QuicConnectionContextListener> listener;
  std::unique_ptr<QuicConnectionTracer> tracer;
  std::unique_ptr<QuicBugListener> bug_listener;
};
class QUICHE_EXPORT QuicConnectionContextSwitcher final {
 public:
  explicit QuicConnectionContextSwitcher(QuicConnectionContext* new_context);
  ~QuicConnectionContextSwitcher();
 private:
  QuicConnectionContext* old_context_;
};
inline void QUIC_TRACELITERAL(const char* literal) {
  QuicConnectionContext* current = QuicConnectionContext::Current();
  if (current && current->tracer) {
    current->tracer->PrintLiteral(literal);
  }
}
inline void QUIC_TRACESTRING(absl::string_view s) {
  QuicConnectionContext* current = QuicConnectionContext::Current();
  if (current && current->tracer) {
    current->tracer->PrintString(s);
  }
}
template <typename... Args>
void QUIC_TRACEPRINTF(const absl::FormatSpec<Args...>& format,
                      const Args&... args) {
  QuicConnectionContext* current = QuicConnectionContext::Current();
  if (current && current->tracer) {
    current->tracer->Printf(format, args...);
  }
}
inline QuicBugListener* CurrentBugListener() {
  QuicConnectionContext* current = QuicConnectionContext::Current();
  return (current != nullptr) ? current->bug_listener.get() : nullptr;
}
}  
#endif  
#include "quiche/quic/core/quic_connection_context.h"
#include "absl/base/attributes.h"
namespace quic {
namespace {
ABSL_CONST_INIT thread_local QuicConnectionContext* current_context = nullptr;
}  
QuicConnectionContext* QuicConnectionContext::Current() {
  return current_context;
}
QuicConnectionContextSwitcher::QuicConnectionContextSwitcher(
    QuicConnectionContext* new_context)
    : old_context_(QuicConnectionContext::Current()) {
  current_context = new_context;
  if (new_context && new_context->listener) {
    new_context->listener->Activate();
  }
}
QuicConnectionContextSwitcher::~QuicConnectionContextSwitcher() {
  QuicConnectionContext* current = QuicConnectionContext::Current();
  if (current && current->listener) {
    current->listener->Deactivate();
  }
  current_context = old_context_;
}
}  