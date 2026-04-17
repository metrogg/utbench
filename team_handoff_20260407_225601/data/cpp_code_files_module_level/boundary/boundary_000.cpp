#ifndef TENSORSTORE_INTERNAL_HTTP_HTTP_RESPONSE_H_
#define TENSORSTORE_INTERNAL_HTTP_HTTP_RESPONSE_H_
#include <stddef.h>
#include <stdint.h>
#include <string>
#include "absl/container/btree_map.h"
#include "absl/status/status.h"
#include "absl/strings/cord.h"
#include "absl/strings/match.h"
#include "absl/strings/str_format.h"
#include "tensorstore/internal/source_location.h"
#include "tensorstore/util/result.h"
namespace tensorstore {
namespace internal_http {
struct HttpResponse {
  int32_t status_code;
  absl::Cord payload;
  absl::btree_multimap<std::string, std::string> headers;
  template <typename Sink>
  friend void AbslStringify(Sink& sink, const HttpResponse& response) {
    absl::Format(&sink, "HttpResponse{code=%d, headers=<",
                 response.status_code);
    const char* sep = "";
    for (const auto& kv : response.headers) {
      sink.Append(sep);
      sink.Append(kv.first);
      sink.Append(": ");
#ifndef NDEBUG
      if (absl::StrContainsIgnoreCase(kv.first, "auth_token")) {
        sink.Append("#####");
      } else
#endif
      {
        sink.Append(kv.second);
      }
      sep = "  ";
    }
    if (response.payload.size() <= 64) {
      absl::Format(&sink, ">, payload=%v}", response.payload);
    } else {
      absl::Format(&sink, ">, payload.size=%d}", response.payload.size());
    }
  }
};
const char* HttpResponseCodeToMessage(const HttpResponse& response);
absl::StatusCode HttpResponseCodeToStatusCode(const HttpResponse& response);
absl::Status HttpResponseCodeToStatus(
    const HttpResponse& response,
    SourceLocation loc = ::tensorstore::SourceLocation::current());
struct ParsedContentRange {
  int64_t inclusive_min;
  int64_t exclusive_max;
  int64_t total_size;  
};
Result<ParsedContentRange> ParseContentRangeHeader(
    const HttpResponse& response);
}  
}  
#endif  
#include "tensorstore/internal/http/http_response.h"
#include <stddef.h>
#include <stdint.h>
#include <limits>
#include <optional>
#include <string>
#include <utility>
#include "absl/status/status.h"
#include "absl/strings/cord.h"
#include "absl/strings/str_format.h"
#include "re2/re2.h"
#include "tensorstore/internal/source_location.h"
#include "tensorstore/util/quote_string.h"
#include "tensorstore/util/result.h"
#include "tensorstore/util/status.h"
#include "tensorstore/util/str_cat.h"
namespace tensorstore {
namespace internal_http {
const char* HttpResponseCodeToMessage(const HttpResponse& response) {
  switch (response.status_code) {
    case 400:
      return "Bad Request";
    case 401:
      return "Unauthorized";
    case 402:
      return "Payment Required";
    case 403:
      return "Forbidden";
    case 404:
      return "Not Found";
    case 405:
      return "Method Not Allowed";
    case 406:
      return "Not Acceptable";
    case 407:
      return "Proxy Authentication Required";
    case 408:
      return "Request Timeout";
    case 409:
      return "Conflict";
    case 410:
      return "Gone";
    case 411:
      return "Length Required";
    case 412:
      return "Precondition Failed";
    case 413:
      return "Payload Too Large";
    case 414:
      return "URI Too Long";
    case 415:
      return "Unsupported Media Type";
    case 416:
      return "Range Not Satisfiable";
    case 417:
      return "Expectation Failed";
    case 418:
      return "I'm a teapot";
    case 421:
      return "Misdirected Request";
    case 422:
      return "Unprocessable Content";
    case 423:
      return "Locked";
    case 424:
      return "Failed Dependency";
    case 425:
      return "Too Early";
    case 426:
      return "Upgrade Required";
    case 428:
      return "Precondition Required";
    case 429:
      return "Too Many Requests";
    case 431:
      return "Request Header Fields Too Large";
    case 451:
      return "Unavailable For Legal Reasons";
    case 500:
      return "Internal Server Error";
    case 501:
      return "Not Implemented";
    case 502:
      return "Bad Gateway";
    case 503:
      return "Service Unavailable";
    case 504:
      return "Gateway Timeout";
    case 505:
      return "HTTP Version Not Supported";
    case 506:
      return "Variant Also Negotiates";
    case 507:
      return "Insufficient Storage";
    case 508:
      return "Loop Detected";
    case 510:
      return "Not Extended";
    case 511:
      return "Network Authentication Required";
    default:
      return nullptr;
  }
}
absl::StatusCode HttpResponseCodeToStatusCode(const HttpResponse& response) {
  switch (response.status_code) {
    case 200:  
    case 201:  
    case 202:  
    case 204:  
    case 206:  
      return absl::StatusCode::kOk;
    case 400:  
    case 411:  
      return absl::StatusCode::kInvalidArgument;
    case 401:  
    case 403:  
      return absl::StatusCode::kPermissionDenied;
    case 404:  
    case 410:  
      return absl::StatusCode::kNotFound;
    case 302:  
    case 303:  
    case 304:  
    case 307:  
    case 412:  
    case 413:  
      return absl::StatusCode::kFailedPrecondition;
    case 416:  
      return absl::StatusCode::kOutOfRange;
    case 308:  
    case 408:  
    case 409:  
    case 429:  
    case 500:  
    case 502:  
    case 503:  
    case 504:  
      return absl::StatusCode::kUnavailable;
  }
  if (response.status_code < 300) {
    return absl::StatusCode::kOk;
  }
  return absl::StatusCode::kUnknown;
}
absl::Status HttpResponseCodeToStatus(const HttpResponse& response,
                                      SourceLocation loc) {
  auto code = HttpResponseCodeToStatusCode(response);
  if (code == absl::StatusCode::kOk) {
    return absl::OkStatus();
  }
  auto status_message = HttpResponseCodeToMessage(response);
  if (!status_message) status_message = "Unknown";
  absl::Status status(code, status_message);
  if (!response.payload.empty()) {
    status.SetPayload(
        "http_response_body",
        response.payload.Subcord(
            0, response.payload.size() < 256 ? response.payload.size() : 256));
  }
  MaybeAddSourceLocation(status, loc);
  status.SetPayload("http_response_code",
                    absl::Cord(tensorstore::StrCat(response.status_code)));
  return status;
}
Result<ParsedContentRange> ParseContentRangeHeader(
    const HttpResponse& response) {
  auto it = response.headers.find("content-range");
  if (it == response.headers.end()) {
    if (response.status_code != 206) {
      return absl::FailedPreconditionError(
          tensorstore::StrCat("No Content-Range header expected with HTTP ",
                              response.status_code, " response"));
    }
    return absl::FailedPreconditionError(
        "Expected Content-Range header with HTTP 206 response");
  }
  static const RE2 kContentRangeRegex(R"(^bytes (\d+)-(\d+)/(?:(\d+)|\*))");
  int64_t a, b;
  std::optional<int64_t> total_size;
  if (!RE2::FullMatch(it->second, kContentRangeRegex, &a, &b, &total_size) ||
      a > b || (total_size && b >= *total_size) ||
      b == std::numeric_limits<int64_t>::max()) {
    return absl::FailedPreconditionError(tensorstore::StrCat(
        "Unexpected Content-Range header received: ", QuoteString(it->second)));
  }
  return ParsedContentRange{a, b + 1, total_size.value_or(-1)};
}
}  
}  