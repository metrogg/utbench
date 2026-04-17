#ifndef TENSORSTORE_INTERNAL_OAUTH2_GOOGLE_AUTH_PROVIDER_H_
#define TENSORSTORE_INTERNAL_OAUTH2_GOOGLE_AUTH_PROVIDER_H_
#include <functional>
#include <memory>
#include "tensorstore/internal/http/curl_transport.h"
#include "tensorstore/internal/http/http_transport.h"
#include "tensorstore/internal/oauth2/auth_provider.h"
#include "tensorstore/util/result.h"
namespace tensorstore {
namespace internal_oauth2 {
Result<std::unique_ptr<AuthProvider>> GetGoogleAuthProvider(
    std::shared_ptr<internal_http::HttpTransport> transport =
        internal_http::GetDefaultHttpTransport());
Result<std::shared_ptr<AuthProvider>> GetSharedGoogleAuthProvider();
void ResetSharedGoogleAuthProvider();
using GoogleAuthProvider =
    std::function<Result<std::unique_ptr<AuthProvider>>()>;
void RegisterGoogleAuthProvider(GoogleAuthProvider provider, int priority);
}  
}  
#endif  
#include "tensorstore/internal/oauth2/google_auth_provider.h"
#include <algorithm>
#include <fstream>
#include <functional>
#include <memory>
#include <optional>
#include <string>
#include <utility>
#include <vector>
#include "absl/base/no_destructor.h"
#include "absl/base/thread_annotations.h"
#include "absl/log/absl_log.h"
#include "absl/status/status.h"
#include "absl/synchronization/mutex.h"
#include <nlohmann/json.hpp>
#include "tensorstore/internal/env.h"
#include "tensorstore/internal/http/http_transport.h"
#include "tensorstore/internal/json_binding/json_binding.h"
#include "tensorstore/internal/oauth2/auth_provider.h"
#include "tensorstore/internal/oauth2/fixed_token_auth_provider.h"
#include "tensorstore/internal/oauth2/gce_auth_provider.h"
#include "tensorstore/internal/oauth2/google_service_account_auth_provider.h"
#include "tensorstore/internal/oauth2/oauth2_auth_provider.h"
#include "tensorstore/internal/oauth2/oauth_utils.h"
#include "tensorstore/internal/path.h"
#include "tensorstore/util/result.h"
#include "tensorstore/util/status.h"
#include "tensorstore/util/str_cat.h"
namespace tensorstore {
namespace internal_oauth2 {
namespace {
using ::tensorstore::internal::GetEnv;
using ::tensorstore::internal::JoinPath;
constexpr char kGoogleAuthTokenForTesting[] = "GOOGLE_AUTH_TOKEN_FOR_TESTING";
constexpr char kGoogleApplicationCredentials[] =
    "GOOGLE_APPLICATION_CREDENTIALS";
constexpr char kCloudSdkConfig[] = "CLOUDSDK_CONFIG";
constexpr char kGCloudConfigFolder[] = ".config/gcloud/";
constexpr char kWellKnownCredentialsFile[] =
    "application_default_credentials.json";
constexpr char kOAuthV3Url[] = "https:
bool IsFile(const std::string& filename) {
  std::ifstream fstream(filename.c_str());
  return fstream.good();
}
Result<std::string> GetEnvironmentVariableFileName() {
  auto env = GetEnv(kGoogleApplicationCredentials);
  if (!env || !IsFile(*env)) {
    return absl::NotFoundError(tensorstore::StrCat(
        "$", kGoogleApplicationCredentials, " is not set or corrupt."));
  }
  return *env;
}
Result<std::string> GetWellKnownFileName() {
  std::string result;
  auto config_dir_override = GetEnv(kCloudSdkConfig);
  if (config_dir_override) {
    result = JoinPath(*config_dir_override, kWellKnownCredentialsFile);
  } else {
    auto home_dir = GetEnv("HOME");
    if (!home_dir) {
      return absl::NotFoundError("Could not read $HOME.");
    }
    result =
        JoinPath(*home_dir, kGCloudConfigFolder, kWellKnownCredentialsFile);
  }
  if (!IsFile(result)) {
    return absl::NotFoundError(
        tensorstore::StrCat("Could not find the credentials file in the "
                            "standard gcloud location [",
                            result, "]"));
  }
  return result;
}
struct AuthProviderRegistry {
  std::vector<std::pair<int, GoogleAuthProvider>> providers;
  absl::Mutex mutex;
};
AuthProviderRegistry& GetGoogleAuthProviderRegistry() {
  static absl::NoDestructor<AuthProviderRegistry> registry;
  return *registry;
}
Result<std::unique_ptr<AuthProvider>> GetDefaultGoogleAuthProvider(
    std::shared_ptr<internal_http::HttpTransport> transport) {
  std::unique_ptr<AuthProvider> result;
  auto var = GetEnv(kGoogleAuthTokenForTesting);
  if (var) {
    ABSL_LOG(INFO) << "Using GOOGLE_AUTH_TOKEN_FOR_TESTING";
    result.reset(new FixedTokenAuthProvider(std::move(*var)));
    return std::move(result);
  }
  absl::Status status;
  auto credentials_filename = GetEnvironmentVariableFileName();
  if (!credentials_filename) {
    credentials_filename = GetWellKnownFileName();
  }
  if (credentials_filename.ok()) {
    ABSL_LOG(INFO) << "Using credentials at " << *credentials_filename;
    std::ifstream credentials_fstream(*credentials_filename);
    auto json = ::nlohmann::json::parse(credentials_fstream, nullptr, false);
    auto refresh_token = internal_oauth2::ParseRefreshToken(json);
    if (refresh_token.ok()) {
      ABSL_LOG(INFO) << "Using OAuth2 AuthProvider";
      result.reset(new OAuth2AuthProvider(*refresh_token, kOAuthV3Url,
                                          std::move(transport)));
      return std::move(result);
    }
    auto service_account =
        internal_oauth2::ParseGoogleServiceAccountCredentials(json);
    if (service_account.ok()) {
      ABSL_LOG(INFO) << "Using ServiceAccount AuthProvider";
      result.reset(new GoogleServiceAccountAuthProvider(*service_account,
                                                        std::move(transport)));
      return std::move(result);
    }
    status = absl::UnknownError(
        tensorstore::StrCat("Unexpected content of the JSON credentials file: ",
                            *credentials_filename));
  }
  if (auto gce_service_account =
          GceAuthProvider::GetDefaultServiceAccountInfoIfRunningOnGce(
              transport.get());
      gce_service_account.ok()) {
    ABSL_LOG(INFO) << "Running on GCE, using service account "
                   << gce_service_account->email;
    result.reset(
        new GceAuthProvider(std::move(transport), *gce_service_account));
    return std::move(result);
  }
  if (!credentials_filename.ok()) {
    ABSL_LOG(ERROR)
        << credentials_filename.status().message()
        << ". You may specify a credentials file using $"
        << kGoogleApplicationCredentials
        << ", or to use Google application default credentials, run: "
           "gcloud auth application-default login";
  }
  TENSORSTORE_RETURN_IF_ERROR(status);
  return absl::NotFoundError(
      "Could not locate the credentials file and not running on GCE.");
}
struct SharedGoogleAuthProviderState {
  absl::Mutex mutex;
  std::optional<Result<std::shared_ptr<AuthProvider>>> auth_provider
      ABSL_GUARDED_BY(mutex);
};
SharedGoogleAuthProviderState& GetSharedGoogleAuthProviderState() {
  static absl::NoDestructor<SharedGoogleAuthProviderState> state;
  return *state;
}
}  
void RegisterGoogleAuthProvider(GoogleAuthProvider provider, int priority) {
  auto& registry = GetGoogleAuthProviderRegistry();
  absl::WriterMutexLock lock(&registry.mutex);
  registry.providers.emplace_back(priority, std::move(provider));
  std::sort(registry.providers.begin(), registry.providers.end(),
            [](const auto& a, const auto& b) { return a.first < b.first; });
}
Result<std::unique_ptr<AuthProvider>> GetGoogleAuthProvider(
    std::shared_ptr<internal_http::HttpTransport> transport) {
  {
    auto& registry = GetGoogleAuthProviderRegistry();
    absl::ReaderMutexLock lock(&registry.mutex);
    for (const auto& provider : registry.providers) {
      auto auth_result = provider.second();
      if (auth_result.ok()) return auth_result;
    }
  }
  return internal_oauth2::GetDefaultGoogleAuthProvider(std::move(transport));
}
Result<std::shared_ptr<AuthProvider>> GetSharedGoogleAuthProvider() {
  auto& state = GetSharedGoogleAuthProviderState();
  absl::MutexLock lock(&state.mutex);
  if (!state.auth_provider) {
    state.auth_provider.emplace(GetGoogleAuthProvider());
  }
  return *state.auth_provider;
}
void ResetSharedGoogleAuthProvider() {
  auto& state = GetSharedGoogleAuthProviderState();
  absl::MutexLock lock(&state.mutex);
  state.auth_provider = std::nullopt;
}
}  
}  