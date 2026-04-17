#ifndef QUICHE_BLIND_SIGN_AUTH_CACHED_BLIND_SIGN_AUTH_H_
#define QUICHE_BLIND_SIGN_AUTH_CACHED_BLIND_SIGN_AUTH_H_
#include <string>
#include <vector>
#include "absl/status/statusor.h"
#include "absl/types/span.h"
#include "quiche/blind_sign_auth/blind_sign_auth_interface.h"
#include "quiche/common/platform/api/quiche_export.h"
#include "quiche/common/platform/api/quiche_mutex.h"
#include "quiche/common/quiche_circular_deque.h"
namespace quiche {
inline constexpr int kBlindSignAuthRequestMaxTokens = 1024;
class QUICHE_EXPORT CachedBlindSignAuth : public BlindSignAuthInterface {
 public:
  CachedBlindSignAuth(
      BlindSignAuthInterface* blind_sign_auth,
      int max_tokens_per_request = kBlindSignAuthRequestMaxTokens)
      : blind_sign_auth_(blind_sign_auth),
        max_tokens_per_request_(max_tokens_per_request) {}
  void GetTokens(std::optional<std::string> oauth_token, int num_tokens,
                 ProxyLayer proxy_layer, BlindSignAuthServiceType service_type,
                 SignedTokenCallback callback) override;
  void ClearCache() {
    QuicheWriterMutexLock lock(&mutex_);
    cached_tokens_.clear();
  }
 private:
  void HandleGetTokensResponse(
      SignedTokenCallback callback, int num_tokens,
      absl::StatusOr<absl::Span<BlindSignToken>> tokens);
  std::vector<BlindSignToken> CreateOutputTokens(int num_tokens)
      QUICHE_EXCLUSIVE_LOCKS_REQUIRED(mutex_);
  void RemoveExpiredTokens() QUICHE_EXCLUSIVE_LOCKS_REQUIRED(mutex_);
  BlindSignAuthInterface* blind_sign_auth_;
  int max_tokens_per_request_;
  QuicheMutex mutex_;
  QuicheCircularDeque<BlindSignToken> cached_tokens_ QUICHE_GUARDED_BY(mutex_);
};
}  
#endif  
#include "quiche/blind_sign_auth/cached_blind_sign_auth.h"
#include <optional>
#include <string>
#include <utility>
#include <vector>
#include "absl/functional/bind_front.h"
#include "absl/status/status.h"
#include "absl/status/statusor.h"
#include "absl/strings/str_format.h"
#include "absl/time/clock.h"
#include "absl/time/time.h"
#include "absl/types/span.h"
#include "quiche/blind_sign_auth/blind_sign_auth_interface.h"
#include "quiche/common/platform/api/quiche_logging.h"
#include "quiche/common/platform/api/quiche_mutex.h"
namespace quiche {
constexpr absl::Duration kFreshnessConstant = absl::Minutes(5);
void CachedBlindSignAuth::GetTokens(std::optional<std::string> oauth_token,
                                    int num_tokens, ProxyLayer proxy_layer,
                                    BlindSignAuthServiceType service_type,
                                    SignedTokenCallback callback) {
  if (num_tokens > max_tokens_per_request_) {
    std::move(callback)(absl::InvalidArgumentError(
        absl::StrFormat("Number of tokens requested exceeds maximum: %d",
                        kBlindSignAuthRequestMaxTokens)));
    return;
  }
  if (num_tokens < 0) {
    std::move(callback)(absl::InvalidArgumentError(absl::StrFormat(
        "Negative number of tokens requested: %d", num_tokens)));
    return;
  }
  std::vector<BlindSignToken> output_tokens;
  {
    QuicheWriterMutexLock lock(&mutex_);
    RemoveExpiredTokens();
    if (static_cast<size_t>(num_tokens) <= cached_tokens_.size()) {
      output_tokens = CreateOutputTokens(num_tokens);
    }
  }
  if (!output_tokens.empty() || num_tokens == 0) {
    std::move(callback)(absl::MakeSpan(output_tokens));
    return;
  }
  SignedTokenCallback caching_callback =
      absl::bind_front(&CachedBlindSignAuth::HandleGetTokensResponse, this,
                       std::move(callback), num_tokens);
  blind_sign_auth_->GetTokens(oauth_token, kBlindSignAuthRequestMaxTokens,
                              proxy_layer, service_type,
                              std::move(caching_callback));
}
void CachedBlindSignAuth::HandleGetTokensResponse(
    SignedTokenCallback callback, int num_tokens,
    absl::StatusOr<absl::Span<BlindSignToken>> tokens) {
  if (!tokens.ok()) {
    QUICHE_LOG(WARNING) << "BlindSignAuth::GetTokens failed: "
                        << tokens.status();
    std::move(callback)(tokens);
    return;
  }
  if (tokens->size() < static_cast<size_t>(num_tokens) ||
      tokens->size() > kBlindSignAuthRequestMaxTokens) {
    QUICHE_LOG(WARNING) << "Expected " << num_tokens << " tokens, got "
                        << tokens->size();
  }
  std::vector<BlindSignToken> output_tokens;
  size_t cache_size;
  {
    QuicheWriterMutexLock lock(&mutex_);
    for (const BlindSignToken& token : *tokens) {
      cached_tokens_.push_back(token);
    }
    RemoveExpiredTokens();
    cache_size = cached_tokens_.size();
    if (cache_size >= static_cast<size_t>(num_tokens)) {
      output_tokens = CreateOutputTokens(num_tokens);
    }
  }
  if (!output_tokens.empty()) {
    std::move(callback)(absl::MakeSpan(output_tokens));
    return;
  }
  std::move(callback)(absl::ResourceExhaustedError(absl::StrFormat(
      "Requested %d tokens, cache only has %d after GetTokensRequest",
      num_tokens, cache_size)));
}
std::vector<BlindSignToken> CachedBlindSignAuth::CreateOutputTokens(
    int num_tokens) {
  std::vector<BlindSignToken> output_tokens;
  if (cached_tokens_.size() < static_cast<size_t>(num_tokens)) {
    QUICHE_LOG(FATAL) << "Check failed, not enough tokens in cache: "
                      << cached_tokens_.size() << " < " << num_tokens;
  }
  for (int i = 0; i < num_tokens; i++) {
    output_tokens.push_back(std::move(cached_tokens_.front()));
    cached_tokens_.pop_front();
  }
  return output_tokens;
}
void CachedBlindSignAuth::RemoveExpiredTokens() {
  size_t original_size = cached_tokens_.size();
  absl::Time now_plus_five_mins = absl::Now() + kFreshnessConstant;
  for (size_t i = 0; i < original_size; i++) {
    BlindSignToken token = std::move(cached_tokens_.front());
    cached_tokens_.pop_front();
    if (token.expiration > now_plus_five_mins) {
      cached_tokens_.push_back(std::move(token));
    }
  }
}
}  