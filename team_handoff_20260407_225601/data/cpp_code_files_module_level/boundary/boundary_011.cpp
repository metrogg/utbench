#ifndef TENSORSTORE_KVSTORE_S3_OBJECT_METADATA_H_
#define TENSORSTORE_KVSTORE_S3_OBJECT_METADATA_H_
#include <stddef.h>
#include <stdint.h>
#include <optional>
#include <string>
#include "absl/container/btree_map.h"
#include "absl/status/status.h"
#include "absl/time/time.h"
#include "tensorstore/internal/http/http_response.h"
#include "tensorstore/internal/source_location.h"
#include "tensorstore/kvstore/generation.h"
#include "tensorstore/util/result.h"
namespace tinyxml2 {
class XMLNode;
}
namespace tensorstore {
namespace internal_kvstore_s3 {
std::string GetNodeText(tinyxml2::XMLNode* node);
std::optional<int64_t> GetNodeInt(tinyxml2::XMLNode* node);
std::optional<absl::Time> GetNodeTimestamp(tinyxml2::XMLNode* node);
Result<StorageGeneration> StorageGenerationFromHeaders(
    const absl::btree_multimap<std::string, std::string>& headers);
absl::Status AwsHttpResponseToStatus(
    const internal_http::HttpResponse& response, bool& retryable,
    SourceLocation loc = ::tensorstore::SourceLocation::current());
}  
}  
#endif  
#include "tensorstore/kvstore/s3/s3_metadata.h"
#include <stddef.h>
#include <stdint.h>
#include <cassert>
#include <initializer_list>
#include <optional>
#include <string>
#include <string_view>
#include <utility>
#include "absl/base/no_destructor.h"
#include "absl/container/btree_map.h"
#include "absl/container/flat_hash_set.h"
#include "absl/status/status.h"
#include "absl/strings/numbers.h"
#include "absl/strings/str_format.h"
#include "absl/time/time.h"
#include "re2/re2.h"
#include "tensorstore/internal/http/http_response.h"
#include "tensorstore/internal/source_location.h"
#include "tensorstore/kvstore/generation.h"
#include "tensorstore/util/result.h"
#include "tensorstore/util/status.h"
#include "tinyxml2.h"
using ::tensorstore::internal_http::HttpResponse;
namespace tensorstore {
namespace internal_kvstore_s3 {
namespace {
static constexpr char kEtag[] = "etag";
static constexpr char kLt[] = "&lt;";
static constexpr char kGt[] = "&gt;";
static constexpr char kQuot[] = "&quot;";
static constexpr char kApos[] = "&apos;";
static constexpr char kAmp[] = "&amp;";
std::string UnescapeXml(std::string_view data) {
  static LazyRE2 kSpecialXmlSymbols = {"(&gt;|&lt;|&quot;|&apos;|&amp;)"};
  std::string_view search = data;
  std::string_view symbol;
  size_t result_len = data.length();
  while (RE2::FindAndConsume(&search, *kSpecialXmlSymbols, &symbol)) {
    result_len -= symbol.length() - 1;
  }
  if (result_len == data.length()) {
    return std::string(data);
  }
  search = data;
  size_t pos = 0;
  size_t res_pos = 0;
  auto result = std::string(result_len, '0');
  while (RE2::FindAndConsume(&search, *kSpecialXmlSymbols, &symbol)) {
    size_t next = data.length() - search.length();
    for (size_t i = pos; i < next - symbol.length(); ++i, ++res_pos) {
      result[res_pos] = data[i];
    }
    if (symbol == kGt) {
      result[res_pos++] = '>';
    } else if (symbol == kLt) {
      result[res_pos++] = '<';
    } else if (symbol == kQuot) {
      result[res_pos++] = '"';
    } else if (symbol == kApos) {
      result[res_pos++] = '`';
    } else if (symbol == kAmp) {
      result[res_pos++] = '&';
    } else {
      assert(false);
    }
    pos = next;
  }
  for (size_t i = pos; i < data.length(); ++i, ++res_pos) {
    result[res_pos] = data[i];
  }
  return result;
}
bool IsRetryableAwsStatusCode(int32_t status_code) {
  switch (status_code) {
    case 408:  
    case 419:  
    case 429:  
    case 440:  
    case 500:  
    case 502:  
    case 503:  
    case 504:  
    case 509:  
    case 598:  
    case 599:  
      return true;
    default:
      return false;
  }
}
bool IsRetryableAwsMessageCode(std::string_view code) {
  static const absl::NoDestructor<absl::flat_hash_set<std::string_view>>
      kRetryableMessages(absl::flat_hash_set<std::string_view>({
          "InternalFailureException",
          "InternalFailure",
          "InternalServerError",
          "InternalError",
          "RequestExpiredException",
          "RequestExpired",
          "ServiceUnavailableException",
          "ServiceUnavailableError",
          "ServiceUnavailable",
          "RequestThrottledException",
          "RequestThrottled",
          "ThrottlingException",
          "ThrottledException",
          "Throttling",
          "SlowDownException",
          "SlowDown",
          "RequestTimeTooSkewedException",
          "RequestTimeTooSkewed",
          "RequestTimeoutException",
          "RequestTimeout",
      }));
  return kRetryableMessages->contains(code);
}
}  
std::optional<int64_t> GetNodeInt(tinyxml2::XMLNode* node) {
  if (!node) {
    return std::nullopt;
  }
  tinyxml2::XMLPrinter printer;
  for (auto* child = node->FirstChild(); child != nullptr;
       child = child->NextSibling()) {
    child->Accept(&printer);
  }
  int64_t result;
  if (absl::SimpleAtoi(printer.CStr(), &result)) {
    return result;
  }
  return std::nullopt;
}
std::optional<absl::Time> GetNodeTimestamp(tinyxml2::XMLNode* node) {
  if (!node) {
    return std::nullopt;
  }
  tinyxml2::XMLPrinter printer;
  for (auto* child = node->FirstChild(); child != nullptr;
       child = child->NextSibling()) {
    child->Accept(&printer);
  }
  absl::Time result;
  if (absl::ParseTime(absl::RFC3339_full, printer.CStr(), absl::UTCTimeZone(),
                      &result, nullptr)) {
    return result;
  }
  return std::nullopt;
}
std::string GetNodeText(tinyxml2::XMLNode* node) {
  if (!node) {
    return "";
  }
  tinyxml2::XMLPrinter printer;
  for (auto* child = node->FirstChild(); child != nullptr;
       child = child->NextSibling()) {
    child->Accept(&printer);
  }
  return UnescapeXml(printer.CStr());
}
Result<StorageGeneration> StorageGenerationFromHeaders(
    const absl::btree_multimap<std::string, std::string>& headers) {
  if (auto it = headers.find(kEtag); it != headers.end()) {
    return StorageGeneration::FromString(it->second);
  }
  return absl::NotFoundError("etag not found in response headers");
}
absl::Status AwsHttpResponseToStatus(const HttpResponse& response,
                                     bool& retryable, SourceLocation loc) {
  auto absl_status_code = internal_http::HttpResponseCodeToStatusCode(response);
  if (absl_status_code == absl::StatusCode::kOk) {
    return absl::OkStatus();
  }
  std::string error_type;
  if (auto error_header = response.headers.find("x-amzn-errortype");
      error_header != response.headers.end()) {
    error_type = error_header->second;
  }
  absl::Cord request_id;
  if (auto request_id_header = response.headers.find("x-amzn-requestid");
      request_id_header != response.headers.end()) {
    request_id = request_id_header->second;
  }
  std::string message;
  auto payload = response.payload;
  auto payload_str = payload.Flatten();
  [&]() {
    if (payload.empty()) return;
    tinyxml2::XMLDocument xmlDocument;
    if (int xmlcode = xmlDocument.Parse(payload_str.data(), payload_str.size());
        xmlcode != tinyxml2::XML_SUCCESS) {
      return;
    }
    auto* root_node = xmlDocument.FirstChildElement("Error");
    if (root_node == nullptr) return;
    if (error_type.empty()) {
      error_type = GetNodeText(root_node->FirstChildElement("Code"));
    }
    if (request_id.empty()) {
      request_id = GetNodeText(root_node->FirstChildElement("RequestId"));
    }
    message = GetNodeText(root_node->FirstChildElement("Message"));
  }();
  retryable = error_type.empty()
                  ? IsRetryableAwsStatusCode(response.status_code)
                  : IsRetryableAwsMessageCode(error_type);
  if (error_type.empty()) {
    error_type = "Unknown";
  }
  absl::Status status(absl_status_code,
                      absl::StrFormat("%s%s%s", error_type,
                                      message.empty() ? "" : ": ", message));
  status.SetPayload("http_response_code",
                    absl::Cord(absl::StrFormat("%d", response.status_code)));
  if (!payload_str.empty()) {
    status.SetPayload(
        "http_response_body",
        payload.Subcord(0,
                        payload_str.size() < 256 ? payload_str.size() : 256));
  }
  if (!request_id.empty()) {
    status.SetPayload("x-amzn-requestid", request_id);
  }
  MaybeAddSourceLocation(status, loc);
  return status;
}
}  
}  