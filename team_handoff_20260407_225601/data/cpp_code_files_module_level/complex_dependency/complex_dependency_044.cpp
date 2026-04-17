#ifndef TENSORFLOW_TSL_PLATFORM_CLOUD_GCS_DNS_CACHE_H_
#define TENSORFLOW_TSL_PLATFORM_CLOUD_GCS_DNS_CACHE_H_
#include <random>
#include "tsl/platform/cloud/http_request.h"
#include "tsl/platform/env.h"
namespace tsl {
const int64_t kDefaultRefreshRateSecs = 60;
class GcsDnsCache {
 public:
  GcsDnsCache() : GcsDnsCache(kDefaultRefreshRateSecs) {}
  GcsDnsCache(int64_t refresh_rate_secs)
      : GcsDnsCache(Env::Default(), refresh_rate_secs) {}
  GcsDnsCache(Env* env, int64_t refresh_rate_secs);
  ~GcsDnsCache() {
    mutex_lock l(mu_);
    cancelled_ = true;
    cond_var_.notify_one();
  }
  void AnnotateRequest(HttpRequest* request);
 private:
  static std::vector<string> ResolveName(const string& name);
  static std::vector<std::vector<string>> ResolveNames(
      const std::vector<string>& names);
  void WorkerThread();
  friend class GcsDnsCacheTest;
  mutex mu_;
  Env* env_;
  condition_variable cond_var_;
  std::default_random_engine random_ TF_GUARDED_BY(mu_);
  bool started_ TF_GUARDED_BY(mu_) = false;
  bool cancelled_ TF_GUARDED_BY(mu_) = false;
  std::unique_ptr<Thread> worker_ TF_GUARDED_BY(mu_);  
  const int64_t refresh_rate_secs_;
  std::vector<std::vector<string>> addresses_ TF_GUARDED_BY(mu_);
};
}  
#endif  
#include "tsl/platform/cloud/gcs_dns_cache.h"
#include <cstring>
#include "absl/status/status.h"
#include "absl/strings/str_cat.h"
#include "tsl/platform/errors.h"
#include "tsl/platform/retrying_utils.h"
#include "tsl/platform/status.h"
#ifndef _WIN32
#include <arpa/inet.h>
#include <netdb.h>
#include <netinet/in.h>
#include <sys/socket.h>
#else
#include <Windows.h>
#include <winsock2.h>
#include <ws2tcpip.h>
#endif
#include <sys/types.h>
namespace tsl {
namespace {
const std::vector<string>& kCachedDomainNames =
    *new std::vector<string>{"www.googleapis.com", "storage.googleapis.com"};
inline void print_getaddrinfo_error(const string& name, Status return_status) {
  LOG(ERROR) << "Error resolving " << name << ": " << return_status;
}
template <typename T>
const T& SelectRandomItemUniform(std::default_random_engine* random,
                                 const std::vector<T>& items) {
  CHECK_GT(items.size(), 0);
  std::uniform_int_distribution<size_t> distribution(0u, items.size() - 1u);
  size_t choice_index = distribution(*random);
  return items[choice_index];
}
}  
GcsDnsCache::GcsDnsCache(Env* env, int64_t refresh_rate_secs)
    : env_(env), refresh_rate_secs_(refresh_rate_secs) {}
void GcsDnsCache::AnnotateRequest(HttpRequest* request) {
  mutex_lock l(mu_);
  if (!started_) {
    VLOG(1) << "Starting GCS DNS cache.";
    DCHECK(!worker_) << "Worker thread already exists!";
    addresses_ = ResolveNames(kCachedDomainNames);
    worker_.reset(env_->StartThread({}, "gcs_dns_worker",
                                    [this]() { return WorkerThread(); }));
    started_ = true;
  }
  CHECK_EQ(kCachedDomainNames.size(), addresses_.size());
  for (size_t i = 0; i < kCachedDomainNames.size(); ++i) {
    const string& name = kCachedDomainNames[i];
    const std::vector<string>& addresses = addresses_[i];
    if (!addresses.empty()) {
      const string& chosen_address =
          SelectRandomItemUniform(&random_, addresses);
      request->AddResolveOverride(name, 443, chosen_address);
      VLOG(1) << "Annotated DNS mapping: " << name << " --> " << chosen_address;
    } else {
      LOG(WARNING) << "No IP addresses available for " << name;
    }
  }
}
 std::vector<string> GcsDnsCache::ResolveName(const string& name) {
  VLOG(1) << "Resolving DNS name: " << name;
  addrinfo hints;
  memset(&hints, 0, sizeof(hints));
  hints.ai_family = AF_INET;  
  hints.ai_socktype = SOCK_STREAM;
  addrinfo* result = nullptr;
  RetryConfig retryConfig(
       5000,
       50 * 1000 * 5000,
       5);
  const Status getaddrinfo_status = RetryingUtils::CallWithRetries(
      [&name, &hints, &result]() {
        int return_code = getaddrinfo(name.c_str(), nullptr, &hints, &result);
        absl::Status return_status;
        switch (return_code) {
          case 0:
            return_status = OkStatus();
            break;
#ifndef _WIN32
          case EAI_ADDRFAMILY:
          case EAI_SERVICE:
          case EAI_SOCKTYPE:
          case EAI_NONAME:
            return_status = absl::FailedPreconditionError(
                absl::StrCat("System in invalid state for getaddrinfo call: ",
                             gai_strerror(return_code)));
            break;
          case EAI_AGAIN:
          case EAI_NODATA:  
            return_status = absl::UnavailableError(absl::StrCat(
                "Resolving ", name, " is temporarily unavailable"));
            break;
          case EAI_BADFLAGS:
          case EAI_FAMILY:
            return_status = absl::InvalidArgumentError(absl::StrCat(
                "Bad arguments for getaddrinfo: ", gai_strerror(return_code)));
            break;
          case EAI_FAIL:
            return_status = absl::NotFoundError(
                absl::StrCat("Permanent failure resolving ", name, ": ",
                             gai_strerror(return_code)));
            break;
          case EAI_MEMORY:
            return_status = absl::ResourceExhaustedError("Out of memory");
            break;
          case EAI_SYSTEM:
          default:
            return_status = absl::UnknownError(strerror(return_code));
#else
          case WSATYPE_NOT_FOUND:
          case WSAESOCKTNOSUPPORT:
          case WSAHOST_NOT_FOUND:
            return_status = absl::FailedPreconditionError(
                absl::StrCat("System in invalid state for getaddrinfo call: ",
                             gai_strerror(return_code)));
            break;
          case WSATRY_AGAIN:
            return_status = absl::UnavailableError(absl::StrCat(
                "Resolving ", name, " is temporarily unavailable"));
            break;
          case WSAEINVAL:
          case WSAEAFNOSUPPORT:
            return_status = absl::InvalidArgumentError(absl::StrCat(
                "Bad arguments for getaddrinfo: ", gai_strerror(return_code)));
            break;
          case WSANO_RECOVERY:
            return_status = absl::NotFoundError(
                absl::StrCat("Permanent failure resolving ", name, ": ",
                             gai_strerror(return_code)));
            break;
          case WSA_NOT_ENOUGH_MEMORY:
            return_status = absl::ResourceExhaustedError("Out of memory");
            break;
          default:
            return_status = absl::UnknownError(strerror(return_code));
#endif
        }
        return Status(return_status);
      },
      retryConfig);
  std::vector<string> output;
  if (getaddrinfo_status.ok()) {
    for (const addrinfo* i = result; i != nullptr; i = i->ai_next) {
      if (i->ai_family != AF_INET || i->ai_addr->sa_family != AF_INET) {
        LOG(WARNING) << "Non-IPv4 address returned. ai_family: " << i->ai_family
                     << ". sa_family: " << i->ai_addr->sa_family << ".";
        continue;
      }
      char buf[INET_ADDRSTRLEN];
      void* address_ptr =
          &(reinterpret_cast<sockaddr_in*>(i->ai_addr)->sin_addr);
      const char* formatted = nullptr;
      if ((formatted = inet_ntop(i->ai_addr->sa_family, address_ptr, buf,
                                 INET_ADDRSTRLEN)) == nullptr) {
        LOG(ERROR) << "Error converting response to IP address for " << name
                   << ": " << strerror(errno);
      } else {
        output.emplace_back(buf);
        VLOG(1) << "... address: " << buf;
      }
    }
  } else {
    print_getaddrinfo_error(name, getaddrinfo_status);
  }
  if (result != nullptr) {
    freeaddrinfo(result);
  }
  return output;
}
std::vector<std::vector<string>> GcsDnsCache::ResolveNames(
    const std::vector<string>& names) {
  std::vector<std::vector<string>> all_addresses;
  all_addresses.reserve(names.size());
  for (const string& name : names) {
    all_addresses.push_back(ResolveName(name));
  }
  return all_addresses;
}
void GcsDnsCache::WorkerThread() {
  while (true) {
    {
      mutex_lock l(mu_);
      if (cancelled_) return;
      cond_var_.wait_for(l, std::chrono::seconds(refresh_rate_secs_));
      if (cancelled_) return;
    }
    auto new_addresses = ResolveNames(kCachedDomainNames);
    {
      mutex_lock l(mu_);
      addresses_.swap(new_addresses);
    }
  }
}
}  