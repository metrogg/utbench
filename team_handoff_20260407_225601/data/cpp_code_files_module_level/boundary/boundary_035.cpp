#ifndef TENSORSTORE_INTERNAL_HTTP_CURL_WRAPPERS_H_
#define TENSORSTORE_INTERNAL_HTTP_CURL_WRAPPERS_H_
#include <memory>
#include <string>
#include <string_view>
#include "absl/status/status.h"
#include <curl/curl.h>  
#include "tensorstore/internal/source_location.h"
namespace tensorstore {
namespace internal_http {
struct CurlPtrCleanup {
  void operator()(CURL*);
};
struct CurlMultiCleanup {
  void operator()(CURLM*);
};
struct CurlSlistCleanup {
  void operator()(curl_slist*);
};
using CurlPtr = std::unique_ptr<CURL, CurlPtrCleanup>;
using CurlMulti = std::unique_ptr<CURLM, CurlMultiCleanup>;
using CurlHeaders = std::unique_ptr<curl_slist, CurlSlistCleanup>;
std::string GetCurlUserAgentSuffix();
absl::Status CurlCodeToStatus(
    CURLcode code, std::string_view detail,
    SourceLocation loc = tensorstore::SourceLocation::current());
absl::Status CurlMCodeToStatus(
    CURLMcode code, std::string_view,
    SourceLocation loc = tensorstore::SourceLocation::current());
}  
}  
#endif  
#include "tensorstore/internal/http/curl_wrappers.h"
#include <string>
#include <string_view>
#include "absl/status/status.h"
#include "absl/strings/cord.h"
#include <curl/curl.h>
#include "tensorstore/internal/source_location.h"
#include "tensorstore/util/status.h"
#include "tensorstore/util/str_cat.h"
namespace tensorstore {
namespace internal_http {
void CurlPtrCleanup::operator()(CURL* c) { curl_easy_cleanup(c); }
void CurlMultiCleanup::operator()(CURLM* m) { curl_multi_cleanup(m); }
void CurlSlistCleanup::operator()(curl_slist* s) { curl_slist_free_all(s); }
std::string GetCurlUserAgentSuffix() {
  static std::string agent =
      tensorstore::StrCat("tensorstore/0.1 ", curl_version());
  return agent;
}
absl::Status CurlCodeToStatus(CURLcode code, std::string_view detail,
                              SourceLocation loc) {
  auto error_code = absl::StatusCode::kUnknown;
  switch (code) {
    case CURLE_OK:
      return absl::OkStatus();
    case CURLE_COULDNT_RESOLVE_PROXY:
      error_code = absl::StatusCode::kUnavailable;
      if (detail.empty()) detail = "Failed to resolve proxy";
      break;
    case CURLE_OPERATION_TIMEDOUT:
      error_code = absl::StatusCode::kDeadlineExceeded;
      if (detail.empty()) detail = "Timed out";
      break;
    case CURLE_COULDNT_CONNECT:
    case CURLE_COULDNT_RESOLVE_HOST:
    case CURLE_GOT_NOTHING:
    case CURLE_HTTP2:
    case CURLE_HTTP2_STREAM:
    case CURLE_PARTIAL_FILE:
    case CURLE_RECV_ERROR:
    case CURLE_SEND_ERROR:
    case CURLE_SSL_CONNECT_ERROR:
    case CURLE_UNSUPPORTED_PROTOCOL:
      error_code = absl::StatusCode::kUnavailable;
      break;
    case CURLE_URL_MALFORMAT:
      error_code = absl::StatusCode::kInvalidArgument;
      break;
    case CURLE_WRITE_ERROR:
      error_code = absl::StatusCode::kCancelled;
      break;
    case CURLE_ABORTED_BY_CALLBACK:
      error_code = absl::StatusCode::kAborted;
      break;
    case CURLE_REMOTE_ACCESS_DENIED:
      error_code = absl::StatusCode::kPermissionDenied;
      break;
    case CURLE_SEND_FAIL_REWIND:
    case CURLE_RANGE_ERROR:  
      error_code = absl::StatusCode::kInternal;
      break;
    case CURLE_BAD_FUNCTION_ARGUMENT:
    case CURLE_OUT_OF_MEMORY:
    case CURLE_NOT_BUILT_IN:
    case CURLE_UNKNOWN_OPTION:
    case CURLE_BAD_DOWNLOAD_RESUME:
      error_code = absl::StatusCode::kInternal;
      break;
    default:
      break;
  }
  absl::Status status(
      error_code, tensorstore::StrCat("CURL error ", curl_easy_strerror(code),
                                      detail.empty() ? "" : ": ", detail));
  status.SetPayload("curl_code", absl::Cord(tensorstore::StrCat(code)));
  MaybeAddSourceLocation(status, loc);
  return status;
}
absl::Status CurlMCodeToStatus(CURLMcode code, std::string_view detail,
                               SourceLocation loc) {
  if (code == CURLM_OK) {
    return absl::OkStatus();
  }
  absl::Status status(
      absl::StatusCode::kInternal,
      tensorstore::StrCat("CURLM error ", curl_multi_strerror(code),
                          detail.empty() ? "" : ": ", detail));
  status.SetPayload("curlm_code", absl::Cord(tensorstore::StrCat(code)));
  MaybeAddSourceLocation(status, loc);
  return status;
}
}  
}  